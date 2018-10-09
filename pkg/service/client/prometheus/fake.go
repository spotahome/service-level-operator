package prometheus

import (
	"context"
	"fmt"
	"time"

	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var (
	slo3CallCount int
)

type fakeFactory struct {
}

// NewFakeFactory returns a new fake factory.
func NewFakeFactory() ClientFactory {
	return &fakeFactory{}
}

// GetV1APIClient satisfies ClientFactory interface.
func (f *fakeFactory) GetV1APIClient(_ string) (promv1.API, error) {
	return &fakeAPICli{
		queryFuncs: map[string]func() float64{
			"slo0_total": func() float64 { return 100 },
			"slo0_error": func() float64 { return 1 },
			"slo1_total": func() float64 { return 1000 },
			"slo1_error": func() float64 { return 1 },
			"slo2_total": func() float64 { return 100000 },
			"slo2_error": func() float64 { return 12 },
			"slo3_total": func() float64 { return 10000 },
			"slo3_error": func() float64 {
				// Every 2 calls return error.
				slo3CallCount++
				if slo3CallCount%2 == 0 {
					return 1
				}
				return 0
			},
		},
	}, nil
}

// fakeAPICli is a faked http client.
type fakeAPICli struct {
	queryFuncs map[string]func() float64
}

func (f *fakeAPICli) Query(_ context.Context, query string, ts time.Time) (model.Value, error) {

	fn, ok := f.queryFuncs[query]
	if !ok {
		return nil, fmt.Errorf("not faked result")
	}

	return model.Vector{
		&model.Sample{
			Metric:    model.Metric{},
			Timestamp: model.Time(time.Now().UTC().Nanosecond()),
			Value:     model.SampleValue(fn()),
		},
	}, nil
}

func (f *fakeAPICli) AlertManagers(_ context.Context) (promv1.AlertManagersResult, error) {
	return promv1.AlertManagersResult{}, nil
}
func (f *fakeAPICli) CleanTombstones(_ context.Context) error {
	return nil
}
func (f *fakeAPICli) Config(_ context.Context) (promv1.ConfigResult, error) {
	return promv1.ConfigResult{}, nil
}
func (f *fakeAPICli) DeleteSeries(_ context.Context, matches []string, startTime time.Time, endTime time.Time) error {
	return nil
}
func (f *fakeAPICli) Flags(_ context.Context) (promv1.FlagsResult, error) {
	return promv1.FlagsResult{}, nil
}
func (f *fakeAPICli) LabelValues(_ context.Context, label string) (model.LabelValues, error) {
	return model.LabelValues{}, nil
}
func (f *fakeAPICli) QueryRange(_ context.Context, query string, r promv1.Range) (model.Value, error) {
	return nil, nil
}
func (f *fakeAPICli) Series(_ context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, error) {
	return []model.LabelSet{}, nil
}
func (f *fakeAPICli) Snapshot(_ context.Context, skipHead bool) (promv1.SnapshotResult, error) {
	return promv1.SnapshotResult{}, nil
}
func (f *fakeAPICli) Targets(_ context.Context) (promv1.TargetsResult, error) {
	return promv1.TargetsResult{}, nil
}
