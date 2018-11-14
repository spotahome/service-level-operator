package sli_test

import (
	"errors"
	"testing"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mpromv1 "github.com/spotahome/service-level-operator/mocks/github.com/prometheus/client_golang/api/prometheus/v1"
	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
	"github.com/spotahome/service-level-operator/pkg/log"
	prometheusvc "github.com/spotahome/service-level-operator/pkg/service/client/prometheus"
	"github.com/spotahome/service-level-operator/pkg/service/sli"
)

func TestPrometheusRetrieve(t *testing.T) {
	sli0 := &monitoringv1alpha1.SLI{
		SLISource: monitoringv1alpha1.SLISource{
			Prometheus: &monitoringv1alpha1.PrometheusSLISource{
				TotalQuery: "test_total_query",
				ErrorQuery: "test_error_query",
			},
		},
	}
	vector2 := model.Vector{
		&model.Sample{
			Metric: model.Metric{},
			Value:  model.SampleValue(2),
		},
	}
	vector100 := model.Vector{
		&model.Sample{
			Metric: model.Metric{},
			Value:  model.SampleValue(100),
		},
	}

	tests := []struct {
		name string
		sli  *monitoringv1alpha1.SLI

		totalQueryResult model.Value
		totalQueryErr    error
		errorQueryResult model.Value
		errorQueryErr    error

		expResult sli.Result
		expErr    bool
	}{
		{
			name:   "If no result from prometheus it should fail.",
			sli:    sli0,
			expErr: true,
		},
		{
			name:             "Failing total query should make the retrieval fail.",
			sli:              sli0,
			totalQueryResult: vector100,
			totalQueryErr:    errors.New("wanted error"),
			errorQueryResult: vector2,
			expErr:           true,
		},
		{
			name:             "Failing error query should make the retrieval fail.",
			sli:              sli0,
			totalQueryResult: vector100,
			errorQueryResult: vector2,
			errorQueryErr:    errors.New("wanted error"),
			expErr:           true,
		},
		{
			name: "If the query doesn't return a vector it should fail.",
			sli:  sli0,
			totalQueryResult: &model.Scalar{
				Value: model.SampleValue(2),
			},
			errorQueryResult: vector2,
			expErr:           true,
		},
		{
			name: "If the query returns more than one metric it should fail.",
			sli:  sli0,
			totalQueryResult: model.Vector{
				&model.Sample{
					Value: model.SampleValue(1),
				},
				&model.Sample{
					Value: model.SampleValue(2),
				},
			},
			errorQueryResult: vector2,
			expErr:           true,
		},
		{
			name:             "If the query returns 0 metrics it should treat as a 0 value.",
			sli:              sli0,
			totalQueryResult: vector2,
			errorQueryResult: model.Vector{},
			expErr:           false,
			expResult: sli.Result{
				TotalQ: 2,
				ErrorQ: 0,
			},
		},
		{
			name:             "Quering prometheus for total and error metrics should return a correct result",
			sli:              sli0,
			totalQueryResult: vector100,
			errorQueryResult: vector2,
			expResult: sli.Result{
				TotalQ: 100,
				ErrorQ: 2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mapi := &mpromv1.API{}
			mpromfactory := &prometheusvc.MockFactory{Cli: mapi}
			mapi.On("Query", mock.Anything, test.sli.Prometheus.TotalQuery, mock.Anything).Return(test.totalQueryResult, test.errorQueryErr)
			mapi.On("Query", mock.Anything, test.sli.Prometheus.ErrorQuery, mock.Anything).Return(test.errorQueryResult, test.totalQueryErr)

			retriever := sli.NewPrometheus(mpromfactory, log.Dummy)
			res, err := retriever.Retrieve(test.sli)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expResult, res)
			}
		})
	}
}
