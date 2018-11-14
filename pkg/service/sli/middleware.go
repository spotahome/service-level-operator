package sli

import (
	"time"

	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
	"github.com/spotahome/service-level-operator/pkg/service/metrics"
)

// metricsMiddleware will monitoring the calls to the SLI Retriever.
type metricsMiddleware struct {
	kind       string
	metricssvc metrics.Service
	next       Retriever
}

// NewMetricsMiddleware returns a new metrics middleware that wraps a Retriever SLI
// service and monitorings with metrics.
func NewMetricsMiddleware(metricssvc metrics.Service, kind string, next Retriever) Retriever {
	return metricsMiddleware{
		kind:       kind,
		metricssvc: metricssvc,
		next:       next,
	}
}

// Retrieve satisfies sli.Retriever interface.
func (m metricsMiddleware) Retrieve(sli *monitoringv1alpha1.SLI) (result Result, err error) {
	defer func(t time.Time) {
		m.metricssvc.ObserveSLIRetrieveDuration(sli, m.kind, t)
		if err != nil {
			m.metricssvc.IncSLIRetrieveError(sli, m.kind)
		}
	}(time.Now())
	return m.next.Retrieve(sli)
}
