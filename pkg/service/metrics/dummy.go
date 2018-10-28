package metrics

import (
	"time"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
)

// Dummy is a Dummy implementation of the metrics service.
var Dummy = &dummy{}

type dummy struct{}

func (dummy) ObserveSLIRetrieveDuration(_ *measurev1alpha1.SLI, _ string, startTime time.Time)    {}
func (dummy) IncSLIRetrieveError(_ *measurev1alpha1.SLI, _ string)                                {}
func (dummy) ObserveSLOOuputCreateDuration(_ *measurev1alpha1.SLO, _ string, startTime time.Time) {}
func (dummy) IncSSLOOuputCreateError(_ *measurev1alpha1.SLO, _ string)                            {}
