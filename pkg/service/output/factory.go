package output

import (
	"fmt"

	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
)

// Factory is a factory that knows how to get the correct
// Output strategy based on the SLO output source.
type Factory interface {
	// GetGetStrategy returns a output based on the SLO source.
	GetStrategy(*monitoringv1alpha1.SLO) (Output, error)
}

// factory doesn't create objects per se, it only knows
// what strategy to return based on the passed SLI.
type factory struct {
	promOutput Output
}

// NewFactory returns a new output factory.
func NewFactory(promOutput Output) Factory {
	return &factory{
		promOutput: promOutput,
	}
}

// GetStrategy satsifies OutputFactory interface.
func (f factory) GetStrategy(s *monitoringv1alpha1.SLO) (Output, error) {
	if s.Output.Prometheus != nil {
		return f.promOutput, nil
	}

	return nil, fmt.Errorf("%s unsupported output kind", s.Name)
}

// MockFactory returns the mocked output strategy.
type MockFactory struct {
	Mock Output
}

// GetStrategy satisfies Factory interface.
func (m MockFactory) GetStrategy(_ *monitoringv1alpha1.SLO) (Output, error) {
	return m.Mock, nil
}
