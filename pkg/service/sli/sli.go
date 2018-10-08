package sli

import (
	"errors"
	"fmt"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	promcli "github.com/slok/service-level-operator/pkg/service/client/prometheus"
)

// Result is the result of getting a SLI from a backend.
type Result struct {
	// TotalQ is the result of applying the total query.
	TotalQ float64
	// ErrorQ is the result of applying  the error query.
	ErrorQ float64
}

// AvailabilityRatio returns the availability of an SLI result in
// ratio unit (0-1).
func (r *Result) AvailabilityRatio() (float64, error) {
	if r.TotalQ < r.ErrorQ {
		return 0, fmt.Errorf("%f can't be higher than %f", r.ErrorQ, r.TotalQ)
	}

	// If no total then everything ok.
	if r.TotalQ == 0 {
		return 1, nil
	}

	dw, err := r.DowntimeRatio()
	if err != nil {
		return 0, err
	}

	return 1 - dw, nil
}

// DowntimeRatio returns the downtime of an SLI result in.
// ratio unit (0-1).
func (r *Result) DowntimeRatio() (float64, error) {
	if r.TotalQ < r.ErrorQ {
		return 0, fmt.Errorf("%f can't be higher than %f", r.ErrorQ, r.TotalQ)
	}

	// If no total then everything ok.
	if r.TotalQ == 0 {
		return 1, nil
	}

	return r.ErrorQ / r.TotalQ, nil
}

// Retriever knows how to get SLIs from different backends.
type Retriever interface {
	// Retrieve returns the result of a SLI retrieved from the implemented backend.
	Retrieve(*measurev1alpha1.SLI) (Result, error)
}

// prometheus knows how to get SLIs from a prometheus backend.
type prometheus struct {
	cliFactory promcli.ClientFactory
	logger     log.Logger
}

// NewPrometheus returns a new prometheus SLI service.
func NewPrometheus(promCliFactory promcli.ClientFactory, logger log.Logger) Retriever {
	return &prometheus{
		cliFactory: promCliFactory,
		logger:     logger,
	}
}

// Retrieve satisfies Service interface..
func (p *prometheus) Retrieve(sli *measurev1alpha1.SLI) (Result, error) {
	// TODO
	return Result{}, errors.New("not implemented")
}
