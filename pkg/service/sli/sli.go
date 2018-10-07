package sli

import (
	"errors"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	promcli "github.com/slok/service-level-operator/pkg/service/client/prometheus"
)

// Result is the result of getting a SLI from a backend.
type Result struct {
	TotalQ float64
	ErrorQ float64
}

// Service knows how to get .
type Service interface {
	// GetSLI returns the result of a ServiceLevel SLI
	GetSLI(*measurev1alpha1.ServiceLevel) (Result, error)
}

// Prometheus knows how to get SLIs from a prometheus backend.
type prometheus struct {
	cliFactory promcli.ClientFactory
	logger     log.Logger
}

// NewPrometheus returns a new prometheus SLI service.
func NewPrometheus(promCliFactory promcli.ClientFactory, logger log.Logger) Service {
	return &prometheus{
		cliFactory: promCliFactory,
		logger:     logger,
	}
}

// GetSLI satisfies Service interface..
func (p *prometheus) GetSLI(sl *measurev1alpha1.ServiceLevel) (Result, error) {
	// TODO
	return Result{}, errors.New("not implemented")
}
