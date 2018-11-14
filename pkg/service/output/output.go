package output

import (
	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
	"github.com/spotahome/service-level-operator/pkg/log"
	"github.com/spotahome/service-level-operator/pkg/service/sli"
)

// Output knows how expose/send/create the output of a SLO and SLI result.
type Output interface {
	// Create will create the SLI result and the SLO on the specific format.
	// It receives the SLI's SLO and it's result.
	Create(serviceLevel *monitoringv1alpha1.ServiceLevel, slo *monitoringv1alpha1.SLO, result *sli.Result) error
}

type logger struct {
	logger log.Logger
}

// NewLogger returns a new output logger service that will output the SLOs on
// the specified logger.
func NewLogger(l log.Logger) Output {
	return &logger{
		logger: l,
	}
}

// Create will log the result on the console.
func (l *logger) Create(serviceLevel *monitoringv1alpha1.ServiceLevel, slo *monitoringv1alpha1.SLO, result *sli.Result) error {
	errorRat, err := result.ErrorRatio()
	if err != nil {
		return err
	}
	l.logger.With("id", serviceLevel.Name).
		With("slo", slo.Name).
		With("availability-target", slo.AvailabilityObjectivePercent).
		Infof("SLI error ratio: %f", errorRat)
	return nil
}
