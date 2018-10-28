package metrics

import (
	"time"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
)

// Service knows how to measure the different parts, flows and processes
// of the application to give more insights and improve the observability
// of the application.
type Service interface {
	// ObserveSLIRetrieveDuration will measure the duration of the process of gathering the group of
	// SLIs for a SLO.
	ObserveSLIRetrieveDuration(sli *measurev1alpha1.SLI, kind string, startTime time.Time)
	// IncSLIRetrieveError will increment the number of errors on the retrieval of the SLIs.
	IncSLIRetrieveError(sli *measurev1alpha1.SLI, kind string)
	// ObserveSLOOuputCreateDuration measures the duration of the process of creating the output for the SLO
	ObserveSLOOuputCreateDuration(slo *measurev1alpha1.SLO, kind string, startTime time.Time)
	// IncSSLOOuputCreateError will increment the number of errors on the SLO output creation.
	IncSSLOOuputCreateError(slo *measurev1alpha1.SLO, kind string)
}
