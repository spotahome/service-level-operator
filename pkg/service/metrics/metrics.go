package metrics

import (
	"time"

	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
)

// Service knows how to monitoring the different parts, flows and processes
// of the application to give more insights and improve the observability
// of the application.
type Service interface {
	// ObserveSLIRetrieveDuration will monitoring the duration of the process of gathering the group of
	// SLIs for a SLO.
	ObserveSLIRetrieveDuration(sli *monitoringv1alpha1.SLI, kind string, startTime time.Time)
	// IncSLIRetrieveError will increment the number of errors on the retrieval of the SLIs.
	IncSLIRetrieveError(sli *monitoringv1alpha1.SLI, kind string)
	// ObserveOuputCreateDuration monitorings the duration of the process of creating the output for the SLO
	ObserveOuputCreateDuration(slo *monitoringv1alpha1.SLO, kind string, startTime time.Time)
	// IncOuputCreateError will increment the number of errors on the SLO output creation.
	IncOuputCreateError(slo *monitoringv1alpha1.SLO, kind string)
}
