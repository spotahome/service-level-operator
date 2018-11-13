package output_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/service/output"
	"github.com/slok/service-level-operator/pkg/service/sli"
)

var (
	sl0 = &measurev1alpha1.ServiceLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sl0-test",
			Namespace: "ns0",
		},
	}
	sl1 = &measurev1alpha1.ServiceLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sl1-test",
			Namespace: "ns1",
		},
	}
	slo00 = &measurev1alpha1.SLO{
		Name:                         "slo00-test",
		AvailabilityObjectivePercent: 99.999,
		Output: measurev1alpha1.Output{
			Prometheus: &measurev1alpha1.PrometheusOutputSource{},
		},
		BurnRates: []measurev1alpha1.BurnRate{
			{
				ErrorBudgetDays: 30,
				Thresholds: []measurev1alpha1.BurnRateThreshold{
					{TimeRangeHours: 1, ErrorBudgetPercent: 2},
					{TimeRangeHours: 6, ErrorBudgetPercent: 5},
					{TimeRangeHours: 72, ErrorBudgetPercent: 10},
				},
			},
		},
	}
	slo01 = &measurev1alpha1.SLO{
		Name:                         "slo01-test",
		AvailabilityObjectivePercent: 99.98,
		Output: measurev1alpha1.Output{
			Prometheus: &measurev1alpha1.PrometheusOutputSource{},
		},
	}
	slo10 = &measurev1alpha1.SLO{
		Name:                         "slo10-test",
		AvailabilityObjectivePercent: 99.99978,
		Output: measurev1alpha1.Output{
			Prometheus: &measurev1alpha1.PrometheusOutputSource{},
		},
	}
	slo11 = &measurev1alpha1.SLO{
		Name:                         "slo11-test",
		AvailabilityObjectivePercent: 95.9981,
		Output: measurev1alpha1.Output{
			Prometheus: &measurev1alpha1.PrometheusOutputSource{
				Labels: map[string]string{
					"env":  "test",
					"team": "team1",
				},
			},
		},
	}
)

func TestPrometheusOutput(t *testing.T) {
	tests := []struct {
		name              string
		cfg               output.PrometheusCfg
		createResults     func(output output.Output)
		expMetrics        []string
		expMissingMetrics []string
	}{
		{
			name: "Creating a output result should expose all the required metrics",
			createResults: func(output output.Output) {
				output.Create(sl0, slo00, &sli.Result{
					TotalQ: 1000000,
					ErrorQ: 122,
				})
			},
			expMetrics: []string{
				`service_level_sli_result_error_ratio_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 0.000122`,
				`service_level_sli_result_count_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 1`,
				`service_level_slo_objective_ratio{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 0.9999899999999999`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="1h",error_budget_spent="2%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 14.4`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="6h",error_budget_spent="5%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 6`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="72h",error_budget_spent="10%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 1`,
			},
		},
		{
			name: "Expired metrics shouldn't be exposed",
			cfg: output.PrometheusCfg{
				ExpireDuration: 500 * time.Microsecond,
			},
			createResults: func(output output.Output) {
				output.Create(sl0, slo00, &sli.Result{
					TotalQ: 1000000,
					ErrorQ: 122,
				})
				time.Sleep(1 * time.Millisecond)
			},
			expMissingMetrics: []string{
				`service_level_sli_result_error_ratio_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 0.000122`,
				`service_level_sli_result_count_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 1`,
				`service_level_slo_objective_ratio{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 0.9999899999999999`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="1h",error_budget_spent="2%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 14.4`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="6h",error_budget_spent="5%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 6`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="72h",error_budget_spent="10%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 1`,
			},
		},
		{
			name: "Creating a output result should expose all the required metrics (multiple adds on same SLO).",
			createResults: func(output output.Output) {
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
				`service_level_sli_result_error_ratio_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 2.112137520556389`,
				`service_level_sli_result_count_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 8`,
				`service_level_slo_objective_ratio{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 0.9999899999999999`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="1h",error_budget_spent="2%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 14.4`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="6h",error_budget_spent="5%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 6`,
				`service_level_slo_burn_rate_threshold{burn_rate_range="72h",error_budget_spent="10%",namespace="ns0",service_level="sl0-test",slo="slo00-test",total_error_budget_range="30d"} 1`,
			},
		},
		{
			name: "Creating a output result should expose all the required metrics (multiple SLOs).",
			createResults: func(output output.Output) {
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
				output.Create(sl1, slo11, &sli.Result{
					TotalQ: 998,
					ErrorQ: 7,
				})
			},
			expMetrics: []string{
				`service_level_sli_result_error_ratio_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 0.000122`,
				`service_level_sli_result_count_total{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 1`,
				`service_level_slo_objective_ratio{namespace="ns0",service_level="sl0-test",slo="slo00-test"} 0.9999899999999999`,

				`service_level_sli_result_error_ratio_total{namespace="ns0",service_level="sl0-test",slo="slo01-test"} 0.3363006923837784`,
				`service_level_sli_result_count_total{namespace="ns0",service_level="sl0-test",slo="slo01-test"} 1`,
				`service_level_slo_objective_ratio{namespace="ns0",service_level="sl0-test",slo="slo01-test"} 0.9998`,

				`service_level_sli_result_error_ratio_total{namespace="ns1",service_level="sl1-test",slo="slo10-test"} 0.0009766096154773965`,
				`service_level_sli_result_count_total{namespace="ns1",service_level="sl1-test",slo="slo10-test"} 2`,
				`service_level_slo_objective_ratio{namespace="ns1",service_level="sl1-test",slo="slo10-test"} 0.9999978`,

				`service_level_sli_result_error_ratio_total{env="test",namespace="ns1",service_level="sl1-test",slo="slo11-test",team="team1"} 0.0070140280561122245`,
				`service_level_sli_result_count_total{env="test",namespace="ns1",service_level="sl1-test",slo="slo11-test",team="team1"} 1`,
				`service_level_slo_objective_ratio{env="test",namespace="ns1",service_level="sl1-test",slo="slo11-test",team="team1"} 0.959981`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			promReg := prometheus.NewRegistry()

			output := output.NewPrometheus(test.cfg, promReg, log.Dummy)
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
			for _, expMissingMetric := range test.expMissingMetrics {
				assert.NotContains(string(metrics), expMissingMetric)
			}
		})
	}
}
