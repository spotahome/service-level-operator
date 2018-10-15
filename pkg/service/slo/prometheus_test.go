package slo_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/service/sli"
	"github.com/slok/service-level-operator/pkg/service/slo"
)

var (
	sl0 = &measurev1alpha1.ServiceLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sl0-test",
		},
	}
	sl1 = &measurev1alpha1.ServiceLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sl1-test",
		},
	}
	slo00 = &measurev1alpha1.SLO{
		Name:                         "slo00-test",
		AvailabilityObjectivePercent: 99.999,
	}
	slo01 = &measurev1alpha1.SLO{
		Name:                         "slo01-test",
		AvailabilityObjectivePercent: 99.98,
	}
	slo10 = &measurev1alpha1.SLO{
		Name:                         "slo10-test",
		AvailabilityObjectivePercent: 99.99978,
	}
)

func TestPrometheusOutput(t *testing.T) {
	tests := []struct {
		name          string
		createResults func(output slo.Output)
		expMetrics    []string
	}{
		{
			name: "Creating a output result should expose all the required metrics",
			createResults: func(output slo.Output) {
				output.Create(sl0, slo00, &sli.Result{
					TotalQ: 1000000,
					ErrorQ: 122,
				})
			},
			expMetrics: []string{
				`service_level_slo_error_ratio_total{service_level="sl0-test",slo="slo00-test"} 0.000122`,
				`service_level_slo_full_ratio_total{service_level="sl0-test",slo="slo00-test"} 1`,
				`service_level_slo_objective_ratio{service_level="sl0-test",slo="slo00-test"} 0.9999899999999999`,
			},
		},
		{
			name: "Creating a output result should expose all the required metrics (multiple adds on same SLO).",
			createResults: func(output slo.Output) {
				slis := []*sli.Result{
					&sli.Result{TotalQ: 1000000, ErrorQ: 122},
					&sli.Result{TotalQ: 999, ErrorQ: 1},
					&sli.Result{TotalQ: 812392, ErrorQ: 94},
					&sli.Result{TotalQ: 83, ErrorQ: 83},
					&sli.Result{TotalQ: 11223, ErrorQ: 11222},
					&sli.Result{TotalQ: 9999999999, ErrorQ: 2},
					&sli.Result{TotalQ: 1245, ErrorQ: 0},
					&sli.Result{TotalQ: 9019, ErrorQ: 1001},
				}
				for _, sli := range slis {
					output.Create(sl0, slo00, sli)
				}
			},
			expMetrics: []string{
				`service_level_slo_error_ratio_total{service_level="sl0-test",slo="slo00-test"} 2.1121375205563884`,
				`service_level_slo_full_ratio_total{service_level="sl0-test",slo="slo00-test"} 8`,
				`service_level_slo_objective_ratio{service_level="sl0-test",slo="slo00-test"} 0.9999899999999999`,
			},
		},
		{
			name: "Creating a output result should expose all the required metrics (multiple SLOs).",
			createResults: func(output slo.Output) {
				output.Create(sl0, slo00, &sli.Result{
					TotalQ: 1000000,
					ErrorQ: 122,
				})
				output.Create(sl0, slo01, &sli.Result{
					TotalQ: 1011,
					ErrorQ: 340,
				})
				output.Create(sl1, slo10, &sli.Result{
					TotalQ: 9212,
					ErrorQ: 1,
				})
				output.Create(sl1, slo10, &sli.Result{
					TotalQ: 3456,
					ErrorQ: 3,
				})
			},
			expMetrics: []string{
				`service_level_slo_error_ratio_total{service_level="sl0-test",slo="slo00-test"} 0.000122`,
				`service_level_slo_full_ratio_total{service_level="sl0-test",slo="slo00-test"} 1`,
				`service_level_slo_objective_ratio{service_level="sl0-test",slo="slo00-test"} 0.9999899999999999`,

				`service_level_slo_error_ratio_total{service_level="sl0-test",slo="slo01-test"} 0.3363006923837784`,
				`service_level_slo_full_ratio_total{service_level="sl0-test",slo="slo01-test"} 1`,
				`service_level_slo_objective_ratio{service_level="sl0-test",slo="slo01-test"} 0.9998`,

				`service_level_slo_error_ratio_total{service_level="sl1-test",slo="slo10-test"} 0.0009766096154773965`,
				`service_level_slo_full_ratio_total{service_level="sl1-test",slo="slo10-test"} 2`,
				`service_level_slo_objective_ratio{service_level="sl1-test",slo="slo10-test"} 0.9999978`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			promReg := prometheus.NewRegistry()

			output := slo.NewPrometheus(promReg)
			test.createResults(output)

			// Check metrics
			h := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/metrics", nil)
			h.ServeHTTP(w, req)

			metrics, _ := ioutil.ReadAll(w.Result().Body)
			for _, expMetric := range test.expMetrics {
				assert.Contains(string(metrics), expMetric)
			}

		})
	}
}
