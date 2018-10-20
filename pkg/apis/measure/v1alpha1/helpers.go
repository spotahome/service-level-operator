package v1alpha1

import "fmt"

// Validate validates and sets defaults on the ServiceLevel
// Kubernetes resource object.
func (s *ServiceLevel) Validate() error {

	if len(s.Spec.ServiceLevelObjectives) == 0 {
		return fmt.Errorf("the number of SLOs on a service level must be more than 0")
	}

	// Check if there is an input.
	for _, slo := range s.Spec.ServiceLevelObjectives {
		err := s.validateSLO(&slo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceLevel) validateSLO(slo *SLO) error {
	if slo.Name == "" {
		return fmt.Errorf("a SLO must have a name")
	}

	if slo.AvailabilityObjectivePercent == 0 {
		return fmt.Errorf("the %s SLO must have a availability objective percent", slo.Name)
	}

	// Check inputs.
	if slo.ServiceLevelIndicator.Prometheus == nil {
		return fmt.Errorf("the %s SLO must have at least one input source", slo.Name)
	}

	// Check outputs.
	if slo.Output.Prometheus == nil {
		return fmt.Errorf("the %s SLO must have at least one output source", slo.Name)
	}

	return nil
}
