package operator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	msli "github.com/slok/service-level-operator/mocks/service/sli"
	mslo "github.com/slok/service-level-operator/mocks/service/slo"
	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/operator"
	"github.com/slok/service-level-operator/pkg/service/sli"
	"github.com/slok/service-level-operator/pkg/service/slo"
)

var (
	sl0 = &measurev1alpha1.ServiceLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-service0",
			Namespace: "fake",
		},
		Spec: measurev1alpha1.ServiceLevelSpec{
			ServiceLevelObjectives: []measurev1alpha1.SLO{
				{
					Name:                         "slo0",
					AvailabilityObjectivePercent: 99.99,
					Disable:                      true,
					ServiceLevelIndicator: measurev1alpha1.SLI{
						SLISource: measurev1alpha1.SLISource{
							Prometheus: &measurev1alpha1.PrometheusSLISource{
								Address:    "http://127.0.0.1:9090",
								TotalQuery: `sum(increase(skipper_serve_host_duration_seconds_count{host="www_spotahome_com"}[5m]))`,
								ErrorQuery: `sum(increase(skipper_serve_host_duration_seconds_count{host="www_spotahome_com", code=~"5.."}[5m]))`,
							},
						},
					},
					Output: measurev1alpha1.Output{
						Prometheus: &measurev1alpha1.PrometheusOutputSource{},
					},
				},
			},
		},
	}

	sl1 = &measurev1alpha1.ServiceLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-service0",
			Namespace: "fake",
		},
		Spec: measurev1alpha1.ServiceLevelSpec{
			ServiceLevelObjectives: []measurev1alpha1.SLO{
				{
					Name:                         "slo0",
					AvailabilityObjectivePercent: 99.95,
					ServiceLevelIndicator: measurev1alpha1.SLI{
						SLISource: measurev1alpha1.SLISource{
							Prometheus: &measurev1alpha1.PrometheusSLISource{},
						},
					},
					Output: measurev1alpha1.Output{
						Prometheus: &measurev1alpha1.PrometheusOutputSource{},
					},
				},
				{
					Name:                         "slo1",
					AvailabilityObjectivePercent: 99.99,
					ServiceLevelIndicator: measurev1alpha1.SLI{
						SLISource: measurev1alpha1.SLISource{
							Prometheus: &measurev1alpha1.PrometheusSLISource{},
						},
					},
					Output: measurev1alpha1.Output{
						Prometheus: &measurev1alpha1.PrometheusOutputSource{},
					},
				},
				{
					Name:                         "slo2",
					AvailabilityObjectivePercent: 99.9,
					ServiceLevelIndicator: measurev1alpha1.SLI{
						SLISource: measurev1alpha1.SLISource{
							Prometheus: &measurev1alpha1.PrometheusSLISource{},
						},
					},
					Output: measurev1alpha1.Output{
						Prometheus: &measurev1alpha1.PrometheusOutputSource{},
					},
				},
				{
					Name:                         "slo3",
					AvailabilityObjectivePercent: 99.9999,
					Disable:                      true,
					ServiceLevelIndicator: measurev1alpha1.SLI{
						SLISource: measurev1alpha1.SLISource{
							Prometheus: &measurev1alpha1.PrometheusSLISource{},
						},
					},
					Output: measurev1alpha1.Output{
						Prometheus: &measurev1alpha1.PrometheusOutputSource{},
					},
				},
			},
		},
	}
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name         string
		serviceLevel *measurev1alpha1.ServiceLevel
		processTimes int
		expErr       bool
	}{
		{
			name:         "With disabled SLO should not process anything.",
			serviceLevel: sl0,
			processTimes: 0,
			expErr:       false,
		},
		{
			name:         "A service level with multiple slos should process all slos.",
			serviceLevel: sl1,
			processTimes: 3,
			expErr:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mout := &mslo.Output{}
			moutf := slo.MockOutputFactory{Mock: mout}
			mret := &msli.Retriever{}
			mretf := sli.MockRetrieverFactory{Mock: mret}

			if test.processTimes > 0 {
				mout.On("Create", mock.Anything, mock.Anything, mock.Anything).Times(test.processTimes).Return(nil)
				mret.On("Retrieve", mock.Anything).Times(test.processTimes).Return(sli.Result{}, nil)
			}

			h := operator.NewHandler(moutf, mretf, log.Dummy)
			err := h.Add(context.Background(), test.serviceLevel)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				mout.AssertExpectations(t)
				mret.AssertExpectations(t)
			}
		})
	}
}
