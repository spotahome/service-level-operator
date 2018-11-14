package metrics

import (
	"time"

	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
)

// Dummy is a Dummy implementation of the metrics service.
var Dummy = &dummy{}

type dummy struct{}

func (dummy) ObserveSLIRetrieveDuration(_ *monitoringv1alpha1.SLI, _ string, startTime time.Time) {}
func (dummy) IncSLIRetrieveError(_ *monitoringv1alpha1.SLI, _ string)                             {}
func (dummy) ObserveOuputCreateDuration(_ *monitoringv1alpha1.SLO, _ string, startTime time.Time) {}
func (dummy) IncOuputCreateError(_ *monitoringv1alpha1.SLO, _ string)                             {}
