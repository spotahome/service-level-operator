package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	crdcli "github.com/slok/service-level-operator/pkg/k8sautogen/client/clientset/versioned"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/operator"
	kubernetesclifactory "github.com/slok/service-level-operator/pkg/service/client/kubernetes"
	kubernetesservice "github.com/slok/service-level-operator/pkg/service/kubernetes"
)

const (
	kubeCliQPS   = 100
	kubeCliBurst = 100
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

	// Operator.
	{
		cfg := m.flags.toOperatorConfig()
		op, err := operator.New(cfg, k8ssvc, m.logger)
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
