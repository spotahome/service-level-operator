package output

import (
	"time"

	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
	"github.com/spotahome/service-level-operator/pkg/service/metrics"
	"github.com/spotahome/service-level-operator/pkg/service/sli"
)

// metricsMiddleware will monitoring the calls to the SLO output.
type metricsMiddleware struct {
	kind       string
	metricssvc metrics.Service
	next       Output
}

// NewMetricsMiddleware returns a new metrics middleware that wraps a Output SLO
// service and monitorings with metrics.
func NewMetricsMiddleware(metricssvc metrics.Service, kind string, next Output) Output {
	return metricsMiddleware{
		kind:       kind,
		metricssvc: metricssvc,
		next:       next,
	}
}

// Create satisfies slo.Output interface.
func (m metricsMiddleware) Create(serviceLevel *monitoringv1alpha1.ServiceLevel, slo *monitoringv1alpha1.SLO, result *sli.Result) (err error) {
	defer func(t time.Time) {
		m.metricssvc.ObserveOuputCreateDuration(slo, m.kind, t)
		if err != nil {
			m.metricssvc.IncOuputCreateError(slo, m.kind)
		}
	}(time.Now())
	return m.next.Create(serviceLevel, slo, result)
}
