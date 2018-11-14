package sli

import (
	"errors"

	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
)

// RetrieverFactory is a factory that knows how to get the correct
// Retriever based on the SLI source.
type RetrieverFactory interface {
	// GetRetriever returns a retriever based on the SLI source.
	GetStrategy(*monitoringv1alpha1.SLI) (Retriever, error)
}

// retrieverFactory doesn't create objects per se, it only knows
// what strategy to return based on the passed SLI.
type retrieverFactory struct {
	promRetriever Retriever
}

// NewRetrieverFactory returns a new retriever factory.
func NewRetrieverFactory(promRetriever Retriever) RetrieverFactory {
	return &retrieverFactory{
		promRetriever: promRetriever,
	}
}

// GetRetriever satsifies RetrieverFactory interface.
func (r retrieverFactory) GetStrategy(s *monitoringv1alpha1.SLI) (Retriever, error) {
	if s.Prometheus != nil {
		return r.promRetriever, nil
	}

	return nil, errors.New("unsupported retriever kind")
}

// MockRetrieverFactory returns the mocked retriever strategy.
type MockRetrieverFactory struct {
	Mock Retriever
}

// GetStrategy satisfies RetrieverFactory interface.
func (m MockRetrieverFactory) GetStrategy(_ *monitoringv1alpha1.SLI) (Retriever, error) {
	return m.Mock, nil
}
