package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	crdcli "github.com/spotahome/service-level-operator/pkg/k8sautogen/client/clientset/versioned"
	"github.com/spotahome/service-level-operator/pkg/log"
	"github.com/spotahome/service-level-operator/pkg/operator"
	kubernetesclifactory "github.com/spotahome/service-level-operator/pkg/service/client/kubernetes"
	promclifactory "github.com/spotahome/service-level-operator/pkg/service/client/prometheus"
	"github.com/spotahome/service-level-operator/pkg/service/configuration"
	kubernetesservice "github.com/spotahome/service-level-operator/pkg/service/kubernetes"
	"github.com/spotahome/service-level-operator/pkg/service/metrics"
)

const (
	kubeCliQPS   = 100
	kubeCliBurst = 100
	gracePeriod  = 2 * time.Second
)

// Main has the main logic of the app.
type Main struct {
	flags  *cmdFlags
	logger log.Logger
}

// Run runs the main program.
func (m *Main) Run() error {
	// Prepare the logger with the correct settings.
	jsonLog := true
	if m.flags.development {
		jsonLog = false
	}
	m.logger = log.Base(jsonLog)
	if m.flags.debug {
		m.logger.Set("debug")
	}

	if m.flags.fake {
		m.logger = m.logger.With("mode", "fake")
		m.logger.Warnf("running in faked mode, any external service will be faked")
	}

	// Create prometheus registry and metrics service to expose and measure with metrics.
	promReg := prometheus.NewRegistry()
	metricssvc := metrics.NewPrometheus(promReg)

	// Create services
	k8sstdcli, k8scrdcli, k8saexcli, err := m.createKubernetesClients()
	if err != nil {
		return err
	}
	k8ssvc := kubernetesservice.New(k8sstdcli, k8scrdcli, k8saexcli, m.logger)

	// Prepare our run entrypoints.
	var g run.Group

	// OS signals.
	{
		sigC := make(chan os.Signal, 1)
		exitC := make(chan struct{})
		signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

		g.Add(
			func() error {
				select {
				case s := <-sigC:
					m.logger.Infof("signal %s received", s)
					return nil
				case <-exitC:
					return nil
				}
			},
			func(_ error) {
				close(exitC)
			},
		)
	}

	// Metrics.
	{
		s := m.createHTTPServer(promReg)
		g.Add(
			func() error {
				m.logger.Infof("metrics server listening on %s", m.flags.listenAddress)
				return s.ListenAndServe()
			},
			func(_ error) {
				m.logger.Infof("draining metrics server connections")
				ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
				defer cancel()
				err := s.Shutdown(ctx)
				if err != nil {
					m.logger.Errorf("error while drainning connections on metrics sever")
				}
			},
		)
	}

	// Operator.
	{

		// Load configuration.
		var cfgSLISrc *configuration.DefaultSLISource
		if m.flags.defSLISourcePath != "" {
			f, err := os.Open(m.flags.defSLISourcePath)
			if err != nil {
				return err
			}
			defer f.Close()
			cfgSLISrc, err = configuration.JSONLoader{}.LoadDefaultSLISource(context.Background(), f)
		}

		// Create SLI source client factories.
		promCliFactory, err := m.createPrometheusCliFactory(cfgSLISrc)
		if err != nil {
			return err
		}

		cfg := m.flags.toOperatorConfig()
		op, err := operator.New(cfg, promReg, promCliFactory, k8ssvc, metricssvc, m.logger)
		if err != nil {
			return err
		}
		closeC := make(chan struct{})

		g.Add(
			func() error {
				return op.Run(closeC)
			},
			func(_ error) {
				close(closeC)
			},
		)
	}

	// Run everything
	return g.Run()
}

// loadKubernetesConfig loads kubernetes configuration based on flags.
func (m *Main) loadKubernetesConfig() (*rest.Config, error) {
	var cfg *rest.Config
	// If devel mode then use configuration flag path.
	if m.flags.development {
		config, err := clientcmd.BuildConfigFromFlags("", m.flags.kubeConfig)
		if err != nil {
			return nil, fmt.Errorf("could not load configuration: %s", err)
		}
		cfg = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading kubernetes configuration inside cluster, check app is running outside kubernetes cluster or run in development mode: %s", err)
		}
		cfg = config
	}

	// Set better cli rate limiter.
	cfg.QPS = kubeCliQPS
	cfg.Burst = kubeCliBurst

	return cfg, nil
}

func (m *Main) createKubernetesClients() (kubernetes.Interface, crdcli.Interface, apiextensionscli.Interface, error) {
	var factory kubernetesclifactory.ClientFactory

	if m.flags.fake {
		factory = kubernetesclifactory.NewFake()
	} else {
		config, err := m.loadKubernetesConfig()
		if err != nil {
			return nil, nil, nil, err
		}
		factory = kubernetesclifactory.NewFactory(config)
	}

	stdcli, err := factory.GetSTDClient()
	if err != nil {
		return nil, nil, nil, err
	}

	crdcli, err := factory.GetCRDClient()
	if err != nil {
		return nil, nil, nil, err
	}

	aexcli, err := factory.GetAPIExtensionClient()
	if err != nil {
		return nil, nil, nil, err
	}

	return stdcli, crdcli, aexcli, nil
}

func (m *Main) createPrometheusCliFactory(cfg *configuration.DefaultSLISource) (promclifactory.ClientFactory, error) {
	if m.flags.fake {
		return promclifactory.NewFakeFactory(), nil
	}

	f := promclifactory.NewBaseFactory()
	if cfg != nil && cfg.Prometheus.Address != "" {
		err := f.WithDefaultV1APIClient(cfg.Prometheus.Address)
		if err != nil {
			return nil, err
		}
		m.logger.Infof("prometheus default SLI source set to: %s", cfg.Prometheus.Address)
	}

	return f, nil
}

// createHTTPServer creates the http server that serves prometheus metrics and healthchecks.
func (m *Main) createHTTPServer(promReg *prometheus.Registry) http.Server {
	h := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
	mux := http.NewServeMux()
	mux.Handle(m.flags.metricsPath, h)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Service level operator</title></head>
			<body>
			<h1>Service level operator</h1>
			<p><a href="` + m.flags.metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	mux.HandleFunc("/healthz/ready", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`ready`)) })
	mux.HandleFunc("/healthz/live", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`live`)) })

	return http.Server{
		Handler: mux,
		Addr:    m.flags.listenAddress,
	}
}

func main() {
	m := &Main{flags: newCmdFlags()}

	// Party time!
	err := m.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running app: %s", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "see you soon, good bye!")
	os.Exit(0)
}
