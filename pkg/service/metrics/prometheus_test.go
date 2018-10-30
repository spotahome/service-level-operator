package metrics_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	"github.com/slok/service-level-operator/pkg/service/metrics"
)

func TestPrometheusMetrics(t *testing.T) {
	kind := "test"

	tests := []struct {
		name       string
		addMetrics func(metrics.Service)
		expMetrics []string
		expCode    int
	}{
		{
			name: "Measuring SLO realted metrics should expose SLO processing metrics on the prometheus endpoint.",
			addMetrics: func(s metrics.Service) {
				now := time.Now()
				s.IncOuputCreateError(nil, kind)
				s.IncOuputCreateError(nil, kind)
				s.ObserveOuputCreateDuration(nil, kind, now.Add(-6*time.Second))
				s.ObserveOuputCreateDuration(nil, kind, now.Add(-27*time.Millisecond))
			},
			expMetrics: []string{
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="0.005"} 0`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="0.01"} 0`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="0.025"} 0`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="0.05"} 1`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="0.1"} 1`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="0.25"} 1`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="0.5"} 1`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="1"} 1`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="2.5"} 1`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="5"} 1`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="10"} 2`,
				`service_level_processing_output_create_duration_seconds_bucket{kind="test",le="+Inf"} 2`,
				`service_level_processing_output_create_duration_seconds_count{kind="test"} 2`,

				`service_level_processing_output_create_failures_total{kind="test"} 2`,
			},
			expCode: 200,
		},
		{
			name: "Measuring SLI realted metrics should expose SLI processing metrics on the prometheus endpoint.",
			addMetrics: func(s metrics.Service) {
				now := time.Now()
				s.IncSLIRetrieveError(nil, kind)
				s.IncSLIRetrieveError(nil, kind)
				s.IncSLIRetrieveError(nil, kind)
				s.ObserveSLIRetrieveDuration(nil, kind, now.Add(-3*time.Second))
				s.ObserveSLIRetrieveDuration(nil, kind, now.Add(-15*time.Millisecond))
				s.ObserveSLIRetrieveDuration(nil, kind, now.Add(-567*time.Millisecond))
			},
			expMetrics: []string{
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="0.005"} 0`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="0.01"} 0`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="0.025"} 1`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="0.05"} 1`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="0.1"} 1`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="0.25"} 1`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="0.5"} 1`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="1"} 2`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="2.5"} 2`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="5"} 3`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="10"} 3`,
				`service_level_processing_sli_retrieve_duration_seconds_bucket{kind="test",le="+Inf"} 3`,
				`service_level_processing_sli_retrieve_duration_seconds_count{kind="test"} 3`,

				`service_level_processing_sli_retrieve_failures_total{kind="test"} 3`,
			},
			expCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			reg := prometheus.NewRegistry()
			m := metrics.NewPrometheus(reg)

			// Add desired metrics
			test.addMetrics(m)

			// Ask prometheus for the metrics
			h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
			r := httptest.NewRequest("GET", "/metrics", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			resp := w.Result()

			// Check all metrics are present.
			if assert.Equal(test.expCode, resp.StatusCode) {
				body, _ := ioutil.ReadAll(resp.Body)
				for _, expMetric := range test.expMetrics {
					assert.Contains(string(body), expMetric, "metric not present on the result of metrics service")
				}
			}
		})
	}
}
