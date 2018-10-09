package sli

import (
	"fmt"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
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

	dw, err := r.ErrorRatio()
	if err != nil {
		return 0, err
	}

	return 1 - dw, nil
}

// ErrorRatio returns the error of an SLI result in.
// ratio unit (0-1).
func (r *Result) ErrorRatio() (float64, error) {
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
