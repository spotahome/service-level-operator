package slo

import (
	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/service/sli"
)

// Output knows how expose/send/create the output of a SLO.
type Output interface {
	// Create will create the SLO result on the specific format.
	// It receives the SLO processed and it's result.
	Create(slo *measurev1alpha1.SLO, result *sli.Result) error
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
func (l *logger) Create(slo *measurev1alpha1.SLO, result *sli.Result) error {
	down, err := result.DowntimeRatio()
	if err != nil {
		return err
	}
	l.logger.With("slo", slo.Name).
		With("availability-target", slo.Availability).
		Infof("SLI downtime result: %s", down)
	return nil
}
