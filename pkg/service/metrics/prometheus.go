package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
)

const (
	promNamespace = "service_level"
	promSubsystem = "processing"
)

var (
	buckets = prometheus.DefBuckets
)

type prometheusService struct {
	sliRetrieveHistogram   *prometheus.HistogramVec
	sliRetrieveErrCounter  *prometheus.CounterVec
	outputCreateHistogram  *prometheus.HistogramVec
	outputCreateErrCounter *prometheus.CounterVec

	reg prometheus.Registerer
}

// NewPrometheus returns a new metrics.Service implementation that
// knows how to monitor gusing Prometheus as backend.
func NewPrometheus(reg prometheus.Registerer) Service {
	p := &prometheusService{
		sliRetrieveHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "sli_retrieve_duration_seconds",
			Help:      "The duration seconds to retrieve the SLIs.",
			Buckets:   buckets,
		}, []string{"kind"}),

		sliRetrieveErrCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "sli_retrieve_failures_total",
			Help:      "Total number sli retrieval failures.",
		}, []string{"kind"}),

		outputCreateHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "output_create_duration_seconds",
			Help:      "The duration seconds to create the output of the SLI and SLO results.",
			Buckets:   buckets,
		}, []string{"kind"}),

		outputCreateErrCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "output_create_failures_total",
			Help:      "Total number SLI and SLO output creation failures.",
		}, []string{"kind"}),

		reg: reg,
	}

	p.registerMetrics()

	return p
}

func (p prometheusService) registerMetrics() {
	p.reg.MustRegister(
		p.sliRetrieveHistogram,
		p.sliRetrieveErrCounter,
		p.outputCreateHistogram,
		p.outputCreateErrCounter,
	)
}

// ObserveSLIRetrieveDuration satisfies metrics.Service interface.
func (p prometheusService) ObserveSLIRetrieveDuration(_ *monitoringv1alpha1.SLI, kind string, startTime time.Time) {
	p.sliRetrieveHistogram.WithLabelValues(kind).Observe(time.Since(startTime).Seconds())
}

// IncSLIRetrieveError satisfies metrics.Service interface.
func (p prometheusService) IncSLIRetrieveError(_ *monitoringv1alpha1.SLI, kind string) {
	p.sliRetrieveErrCounter.WithLabelValues(kind).Inc()
}

// ObserveOuputCreateDuration satisfies metrics.Service interface.
func (p prometheusService) ObserveOuputCreateDuration(_ *monitoringv1alpha1.SLO, kind string, startTime time.Time) {
	p.outputCreateHistogram.WithLabelValues(kind).Observe(time.Since(startTime).Seconds())
}

// IncOuputCreateError satisfies metrics.Service interface.
func (p prometheusService) IncOuputCreateError(_ *monitoringv1alpha1.SLO, kind string) {
	p.outputCreateErrCounter.WithLabelValues(kind).Inc()
}
