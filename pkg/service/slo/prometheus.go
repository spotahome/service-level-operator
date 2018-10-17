package slo

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/service/sli"
)

const (
	promNS            = "service_level"
	promSubsystem     = "slo"
	defExpireDuration = 90 * time.Second
)

// metricValue is an internal type to store the counters
// of the metrics so when the collector is called it creates
// the metrics based on this values.
type metricValue struct {
	serviceLevel *measurev1alpha1.ServiceLevel
	slo          *measurev1alpha1.SLO
	errorSum     float64
	totalSum     float64
	objective    float64
	expire       time.Time // expire is the time where this metric will expire unless it's refreshed.
}

// PrometheusCfg is the configuration of the Prometheus Output.
type PrometheusCfg struct {
	// ExpireDuration is the time a metric will expire if is not refreshed.
	ExpireDuration time.Duration
}

// Validate will validate the cfg setting safe defaults.
func (p *PrometheusCfg) Validate() {
	if p.ExpireDuration == 0 {
		p.ExpireDuration = defExpireDuration
	}
}

// Prometheus knows how to set the output of the SLO on a Prometheus backend.
// The way it works this output is creating two main counters, one that increments
// the error and other that increments the full ratio.
// Example:
// error ratio:    0 +  0  + 0.001 +  0.1  +  0.01  = 0.111
// full ratio:     1 +  1  +     1 +    1  +     1  = 5
//
// You could get the total availability ratio with 1-(0.111/5) = 0.9778
// In other words the availability of all this time is: 97.78%
//
// Under the hood this service is a prometheus collector, it will send to
// prometheus dynamic metrics (because of dynamic labels) when the collect
// process is called. This is made by storing the internal counters and
// generating the metrics when the collect process is callend on each scrape.
type prometheusOutput struct {
	cfg            PrometheusCfg
	metricValuesMu sync.Mutex
	metricValues   map[string]*metricValue
	reg            prometheus.Registerer
	logger         log.Logger
}

// NewPrometheus returns a new Prometheus output.
func NewPrometheus(cfg PrometheusCfg, reg prometheus.Registerer, logger log.Logger) Output {
	cfg.Validate()

	p := &prometheusOutput{
		cfg:          cfg,
		metricValues: map[string]*metricValue{},
		reg:          reg,
		logger:       logger,
	}

	// Autoregister as collector of SLO metrics for prometheus.
	p.reg.MustRegister(p)

	return p
}

// Create satisfies output interface. By setting the correct values on the different
// metrics of the SLO.
func (p *prometheusOutput) Create(serviceLevel *measurev1alpha1.ServiceLevel, slo *measurev1alpha1.SLO, result *sli.Result) error {
	p.metricValuesMu.Lock()
	defer p.metricValuesMu.Unlock()

	// Get the current metrics for the SLO.
	sloID := fmt.Sprintf("%s-%s-%s", serviceLevel.Namespace, serviceLevel.Name, slo.Name)
	if _, ok := p.metricValues[sloID]; !ok {
		p.metricValues[sloID] = &metricValue{}
	}

	// Add metric values.
	errRat, err := result.ErrorRatio()
	if err != nil {
		return err
	}
	metric := p.metricValues[sloID]
	metric.serviceLevel = serviceLevel
	metric.slo = slo
	metric.errorSum += errRat
	metric.totalSum++
	// Objective is in %  so we convert to ratio (0-1).
	metric.objective = slo.AvailabilityObjectivePercent / 100
	// Refresh the metric expiration.
	metric.expire = time.Now().Add(p.cfg.ExpireDuration)

	return nil
}

// Describe satisfies prometheus.Collector interface.
func (p *prometheusOutput) Describe(chan<- *prometheus.Desc) {}

// Collect satisfies prometheus.Collector interface.
func (p *prometheusOutput) Collect(ch chan<- prometheus.Metric) {
	p.metricValuesMu.Lock()
	defer p.metricValuesMu.Unlock()
	p.logger.Debugf("start collecting all SLOs")

	for id, metric := range p.metricValues {
		metric := metric

		// If metric has expired then remove from the map.
		if time.Now().After(metric.expire) {
			p.logger.With("slo", metric.slo.Name).With("service-level", metric.serviceLevel.Name).Infof("metric expired, removing")
			delete(p.metricValues, id)
			continue
		}

		ns := metric.serviceLevel.Namespace
		slName := metric.serviceLevel.Name
		sloName := metric.slo.Name
		var labels map[string]string
		// Check just in case.
		if metric.slo.Output.Prometheus != nil && metric.slo.Output.Prometheus.Labels != nil {
			labels = metric.slo.Output.Prometheus.Labels
		}

		ch <- p.getSLOErrorMetric(ns, slName, sloName, labels, metric.errorSum)
		ch <- p.getSLOFullMetric(ns, slName, sloName, labels, metric.totalSum)
		ch <- p.getSLOObjectiveMetric(ns, slName, sloName, labels, metric.objective)
	}

	// Collect all SLOs metric.
	p.logger.Debugf("finished collecting all SLOs")
}

func (p *prometheusOutput) getSLOErrorMetric(ns, serviceLevel, slo string, constLabels prometheus.Labels, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(promNS, promSubsystem, "error_ratio_total"),
			"Is the total error ratio counter for SLOs.",
			[]string{"namespace", "service_level", "slo"},
			constLabels,
		),
		prometheus.CounterValue,
		value,
		ns, serviceLevel, slo,
	)
}

func (p *prometheusOutput) getSLOFullMetric(ns, serviceLevel, slo string, constLabels prometheus.Labels, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(promNS, promSubsystem, "full_ratio_total"),
			"Is the full SLOs ratio counter in time.",
			[]string{"namespace", "service_level", "slo"},
			constLabels,
		),
		prometheus.CounterValue,
		value,
		ns, serviceLevel, slo,
	)
}

func (p *prometheusOutput) getSLOObjectiveMetric(ns, serviceLevel, slo string, constLabels prometheus.Labels, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(promNS, promSubsystem, "objective_ratio"),
			"Is the objective of the SLO in ratio unit.",
			[]string{"namespace", "service_level", "slo"},
			constLabels,
		),
		prometheus.GaugeValue,
		value,
		ns, serviceLevel, slo,
	)
}
