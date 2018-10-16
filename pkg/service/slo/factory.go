package slo

import (
	"fmt"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
)

// OutputFactory is a factory that knows how to get the correct
// Output strategy based on the SLO output source.
type OutputFactory interface {
	// GetGetStrategy returns a output based on the SLO source.
	GetStrategy(*measurev1alpha1.SLO) (Output, error)
}

// retrieverFactory doesn't create objects per se, it only knows
// what strategy to return based on the passed SLI.
type outputFactory struct {
	promOutput Output
}

// NewOutputFactory returns a new output factory.
func NewOutputFactory(promOutput Output) OutputFactory {
	return &outputFactory{
		promOutput: promOutput,
	}
}

// GetStrategy satsifies OutputFactory interface.
func (o outputFactory) GetStrategy(s *measurev1alpha1.SLO) (Output, error) {
	if s.Output.Prometheus != nil {
		return o.promOutput, nil
	}

	return nil, fmt.Errorf("%s unsupported output kind", s.Name)
}

// MockOutputFactory returns the mocked output strategy.
type MockOutputFactory struct {
	Mock Output
}

// GetStrategy satisfies OutputFactory interface.
func (m MockOutputFactory) GetStrategy(_ *measurev1alpha1.SLO) (Output, error) {
	return m.Mock, nil
}
