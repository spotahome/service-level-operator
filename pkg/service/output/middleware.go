package output

import (
	"time"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/service/metrics"
	"github.com/slok/service-level-operator/pkg/service/sli"
)

// metricsMiddleware will measure the calls to the SLO output.
type metricsMiddleware struct {
	kind       string
	metricssvc metrics.Service
	next       Output
}

// NewMetricsMiddleware returns a new metrics middleware that wraps a Output SLO
// service and measures with metrics.
func NewMetricsMiddleware(metricssvc metrics.Service, kind string, next Output) Output {
	return metricsMiddleware{
		kind:       kind,
		metricssvc: metricssvc,
		next:       next,
	}
}

// Create satisfies slo.Output interface.
func (m metricsMiddleware) Create(serviceLevel *measurev1alpha1.ServiceLevel, slo *measurev1alpha1.SLO, result *sli.Result) (err error) {
	defer func(t time.Time) {
		m.metricssvc.ObserveSLOOuputCreateDuration(slo, m.kind, t)
		if err != nil {
			m.metricssvc.IncSSLOOuputCreateError(slo, m.kind)
		}
	}(time.Now())
	return m.next.Create(serviceLevel, slo, result)
}
