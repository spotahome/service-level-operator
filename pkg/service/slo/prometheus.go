package slo

import (
	"github.com/prometheus/client_golang/prometheus"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/service/sli"
)

const (
	promNS        = "service_level"
	promSubsystem = "slo"
)

// Prometheus knows how to set the output of the SLO on a Prometheus backend.
// The way it works this output is creating two main counters, one that increments
// the error and other that increments the full ratio.
// Example:
// error ratio:    0 +  0  + 0.001 +  0.1  +  0.01  = 0.111
// full ratio:     1 +  1  +     1 +    1  +     1  = 5
//
// You could get the total availability ratio with 1-(0.111/5) = 0.9778
// In other words the availability of all this time is: 97.78%
type prometheusOutput struct {
	sloErrorCounter *prometheus.CounterVec
	sloFullCounter  *prometheus.CounterVec
	sloTargetGauge  *prometheus.GaugeVec

	reg prometheus.Registerer
}

// NewPrometheus returns a new Prometheus output.
func NewPrometheus(reg prometheus.Registerer) Output {
	p := &prometheusOutput{
		sloErrorCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNS,
			Subsystem: promSubsystem,
			Name:      "error_ratio_total",
			Help:      "Is the total error ratio counter for SLOs.",
		}, []string{"service_level", "slo"}),
		sloFullCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNS,
			Subsystem: promSubsystem,
			Name:      "full_ratio_total",
			Help:      "Is the full SLOs ratio counter in time.",
		}, []string{"service_level", "slo"}),
		sloTargetGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: promNS,
			Subsystem: promSubsystem,
			Name:      "objective_ratio",
			Help:      "Is the objective of the SLO in ratio unit.",
		}, []string{"service_level", "slo"}),

		reg: reg,
	}

	p.registerMetrics()

	return p
}

func (p *prometheusOutput) registerMetrics() {
	p.reg.MustRegister(
		p.sloErrorCounter,
		p.sloFullCounter,
		p.sloTargetGauge,
	)
}

// Create satisfies output interface. By setting the correct values on the different
// metrics of the SLO.
func (p *prometheusOutput) Create(serviceLevel *measurev1alpha1.ServiceLevel, slo *measurev1alpha1.SLO, result *sli.Result) error {
	errRat, err := result.ErrorRatio()
	if err != nil {
		return err
	}

	p.sloErrorCounter.WithLabelValues(serviceLevel.Name, slo.Name).Add(errRat)
	p.sloFullCounter.WithLabelValues(serviceLevel.Name, slo.Name).Add(1)

	// Objective is in %  so we convert to ratio (0-1)
	objRat := slo.AvailabilityObjectivePercent / 100
	p.sloTargetGauge.WithLabelValues(serviceLevel.Name, slo.Name).Set(objRat)
	return nil
}
