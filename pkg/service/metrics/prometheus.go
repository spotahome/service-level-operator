package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
)

const (
	promNamespace = "service_level"
	promSubsystem = "processing"
)

var (
	buckets = prometheus.DefBuckets
)

type prometheusService struct {
	sliHistogram  *prometheus.HistogramVec
	sliErrCounter *prometheus.CounterVec
	sloHistogram  *prometheus.HistogramVec
	sloErrCounter *prometheus.CounterVec

	reg prometheus.Registerer
}

// NewPrometheus returns a new metrics.Service implementation that
// knows how to measureusing Prometheus as backend.
func NewPrometheus(reg prometheus.Registerer) Service {
	p := &prometheusService{
		sliHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "sli_retrieve_duration_seconds",
			Help:      "The duration seconds to retrieve the SLIs.",
			Buckets:   buckets,
		}, []string{"kind"}),

		sliErrCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "sli_retrieve_failures_total",
			Help:      "Total number sli retrieval failures.",
		}, []string{"kind"}),

		sloHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "slo_output_create_duration_seconds",
			Help:      "The duration seconds to create the output of the SLO results.",
			Buckets:   buckets,
		}, []string{"kind"}),

		sloErrCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: promSubsystem,
			Name:      "slo_output_create_failures_total",
			Help:      "Total number slo output creation failures.",
		}, []string{"kind"}),

		reg: reg,
	}

	p.registerMetrics()

	return p
}

func (p prometheusService) registerMetrics() {
	p.reg.MustRegister(
		p.sliHistogram,
		p.sliErrCounter,
		p.sloHistogram,
		p.sloErrCounter,
	)
}

// ObserveSLIRetrieveDuration satisfies metrics.Service interface.
func (p prometheusService) ObserveSLIRetrieveDuration(_ *measurev1alpha1.SLI, kind string, startTime time.Time) {
	p.sliHistogram.WithLabelValues(kind).Observe(time.Since(startTime).Seconds())
}

// IncSLIRetrieveError satisfies metrics.Service interface.
func (p prometheusService) IncSLIRetrieveError(_ *measurev1alpha1.SLI, kind string) {
	p.sliErrCounter.WithLabelValues(kind).Inc()
}

// ObserveSLOOuputCreateDuration satisfies metrics.Service interface.
func (p prometheusService) ObserveSLOOuputCreateDuration(_ *measurev1alpha1.SLO, kind string, startTime time.Time) {
	p.sloHistogram.WithLabelValues(kind).Observe(time.Since(startTime).Seconds())
}

// IncSSLOOuputCreateError satisfies metrics.Service interface.
func (p prometheusService) IncSSLOOuputCreateError(_ *measurev1alpha1.SLO, kind string) {
	p.sloErrCounter.WithLabelValues(kind).Inc()
}
