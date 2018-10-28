package operator

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	kmetrics "github.com/spotahome/kooper/monitoring/metrics"
	"github.com/spotahome/kooper/operator"
	"github.com/spotahome/kooper/operator/controller"

	"github.com/slok/service-level-operator/pkg/log"
	promcli "github.com/slok/service-level-operator/pkg/service/client/prometheus"
	"github.com/slok/service-level-operator/pkg/service/kubernetes"
	"github.com/slok/service-level-operator/pkg/service/metrics"
	"github.com/slok/service-level-operator/pkg/service/sli"
	"github.com/slok/service-level-operator/pkg/service/slo"
)

const (
	operatorName = "service-level-operator"
	jobRetries   = 3
)

// Config is the configuration for the ci operator.
type Config struct {
	// ResyncPeriod is the resync period of the controllers.
	ResyncPeriod time.Duration
	// ConcurretWorkers are number of workers to handle the events.
	ConcurretWorkers int
	// LabelSelector is the label selector to filter Kubernetes resources by labels.
	LabelSelector string
	// Namespace is the namespace to filter Kubernetes resources by a single namespace.
	Namespace string
}

// New returns pod terminator operator.
func New(cfg Config, promreg *prometheus.Registry, promCliFactory promcli.ClientFactory, k8ssvc kubernetes.Service, metricssvc metrics.Service, logger log.Logger) (operator.Operator, error) {

	// Create crd.
	ptCRD := newServiceLevelCRD(cfg, k8ssvc, logger)

	// Create services.
	promRetriever := sli.NewPrometheus(promCliFactory, logger.WithField("sli-retriever", "prometheus"))
	retrieverFact := sli.NewRetrieverFactory(
		sli.NewMetricsMiddleware(metricssvc, "prometheus", promRetriever),
	)

	promOutput := slo.NewPrometheus(slo.PrometheusCfg{}, promreg, logger.WithField("slo-output", "prometheus"))
	outputFact := slo.NewOutputFactory(
		slo.NewMetricsMiddleware(metricssvc, "prometheus", promOutput),
	)

	// Create handler.
	handler := NewHandler(outputFact, retrieverFact, logger)

	// Create controller.
	ctrlCfg := &controller.Config{
		Name:                 operatorName,
		ConcurrentWorkers:    cfg.ConcurretWorkers,
		ResyncInterval:       cfg.ResyncPeriod,
		ProcessingJobRetries: jobRetries,
	}

	ctrl := controller.New(
		ctrlCfg,
		handler,
		ptCRD,
		nil,
		nil,
		kmetrics.NewPrometheus(promreg),
		logger)

	// Assemble CRD and controller to create the operator.
	return operator.NewOperator(ptCRD, ctrl, logger), nil
}
