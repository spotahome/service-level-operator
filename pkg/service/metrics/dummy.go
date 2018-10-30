package metrics

import (
	"time"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
)

// Dummy is a Dummy implementation of the metrics service.
var Dummy = &dummy{}

type dummy struct{}

func (dummy) ObserveSLIRetrieveDuration(_ *measurev1alpha1.SLI, _ string, startTime time.Time) {}
func (dummy) IncSLIRetrieveError(_ *measurev1alpha1.SLI, _ string)                             {}
func (dummy) ObserveOuputCreateDuration(_ *measurev1alpha1.SLO, _ string, startTime time.Time) {}
func (dummy) IncOuputCreateError(_ *measurev1alpha1.SLO, _ string)                             {}
