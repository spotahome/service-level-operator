package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
)

func TestServiceLevelValidation(t *testing.T) {
	// Setup the different combinations of service level to validate.
	goodSL := &measurev1alpha1.ServiceLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fake-service0",
		},
		Spec: measurev1alpha1.ServiceLevelSpec{
			ServiceLevelObjectives: []measurev1alpha1.SLO{
				{
					Name:                         "fake_slo0",
					Description:                  "fake slo 0.",
					AvailabilityObjectivePercent: 99.99,
					ServiceLevelIndicator: measurev1alpha1.SLI{
						SLISource: measurev1alpha1.SLISource{
							Prometheus: &measurev1alpha1.PrometheusSLISource{
								Address:    "http://fake:9090",
								TotalQuery: `slo0_total`,
								ErrorQuery: `slo0_error`,
							},
						},
					},
					Output: measurev1alpha1.Output{
						Prometheus: &measurev1alpha1.PrometheusOutputSource{},
					},
					BurnRates: []measurev1alpha1.BurnRate{
						{
							ErrorBudgetDays: 30,
							Thresholds: []measurev1alpha1.BurnRateThreshold{
								{
									TimeRangeHours:     1,
									ErrorBudgetPercent: 2,
								},
							},
						},
					},
				},
			},
		},
	}
	slWithoutSLO := goodSL.DeepCopy()
	slWithoutSLO.Spec.ServiceLevelObjectives = []measurev1alpha1.SLO{}
	slSLOWithoutName := goodSL.DeepCopy()
	slSLOWithoutName.Spec.ServiceLevelObjectives[0].Name = ""
	slSLOWithoutObjective := goodSL.DeepCopy()
	slSLOWithoutObjective.Spec.ServiceLevelObjectives[0].AvailabilityObjectivePercent = 0
	slSLOWithoutSLI := goodSL.DeepCopy()
	slSLOWithoutSLI.Spec.ServiceLevelObjectives[0].ServiceLevelIndicator.Prometheus = nil
	slSLOWithoutOutput := goodSL.DeepCopy()
	slSLOWithoutOutput.Spec.ServiceLevelObjectives[0].Output.Prometheus = nil
	slSLOSameMultipleBurnRates := goodSL.DeepCopy()
	slSLOSameMultipleBurnRates.Spec.ServiceLevelObjectives[0].BurnRates = append(slSLOSameMultipleBurnRates.Spec.ServiceLevelObjectives[0].BurnRates, measurev1alpha1.BurnRate{
		ErrorBudgetDays: 30,
		Thresholds: []measurev1alpha1.BurnRateThreshold{
			{
				TimeRangeHours:     1,
				ErrorBudgetPercent: 2,
			},
		},
	})

	tests := []struct {
		name         string
		serviceLevel *measurev1alpha1.ServiceLevel
		expErr       bool
	}{
		{
			name:         "A valid ServiceLevel should be valid.",
			serviceLevel: goodSL,
			expErr:       false,
		},
		{
			name:         "A ServiceLevel without SLOs houldn't be valid.",
			serviceLevel: slWithoutSLO,
			expErr:       true,
		},
		{
			name:         "A ServiceLevel with an SLO without name shouldn't be valid.",
			serviceLevel: slSLOWithoutName,
			expErr:       true,
		},
		{
			name:         "A ServiceLevel with an SLO without objective shouldn't be valid.",
			serviceLevel: slSLOWithoutObjective,
			expErr:       true,
		},
		{
			name:         "A ServiceLevel with an SLO without SLI shouldn't be valid.",
			serviceLevel: slSLOWithoutSLI,
			expErr:       true,
		},
		{
			name:         "A ServiceLevel with an SLO without output shouldn't be valid.",
			serviceLevel: slSLOWithoutOutput,
			expErr:       true,
		},
		{
			name:         "A ServiceLevel with an SLO With repeated burn rates shouldn't be valid.",
			serviceLevel: slSLOSameMultipleBurnRates,
			expErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			err := test.serviceLevel.Validate()

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
