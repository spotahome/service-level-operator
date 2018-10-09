package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceLevel represents a service level policy to measure the service level
// of an application.
type ServiceLevel struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the ddesired behaviour of the pod terminator.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Spec ServiceLevelSpec `json:"spec,omitempty"`
}

// ServiceLevelSpec is the spec for a ServiceLevel resource.
type ServiceLevelSpec struct {
	// ServiceLevelObjectives is the list of SLOs of a service/app.
	// +optional
	ServiceLevelObjectives []SLO `json:"serviceLevelObjectives,omitempty"`
}

// SLO represents a SLO.
type SLO struct {
	// Name of the SLO, must be made of [a-zA-z0-9] and '_'(underscore) characters.
	Name string `json:"name"`
	// Description is a description of the SLO.
	// +optional
	Description string `json:"description,omitempty"`
	// Disable will disable the SLO.
	Disable bool `json:"disable,omitempty"`
	// Availability is the percentage of availability target for the SLO.
	Availability float64 `json:"availability"`
	// ServiceLevelIndicator is the SLI associated with the SLO.
	ServiceLevelIndicator SLI `json:"serviceLevelIndicator"`
	// Output is the output backedn of the SLO.
	Output Output `json:"output"`
}

// SLI is the SLI to get for the SLO.
type SLI struct {
	SLISource `json:",inline"`
}

// SLISource is where the SLI will get from.
type SLISource struct {
	// Prometheus is the prometheus SLI source.
	// +optional
	Prometheus *PrometheusSLISource `json:"prometheus,omitempty"`
}

// PrometheusSLISource is the source to get SLIs from a Prometheus backend.
type PrometheusSLISource struct {
	// Address is the address of the Prometheus.
	Address string `json:"address"`
	// TotalQuery is the query that gets the total that will be the base to get the unavailability
	// of the SLO based on the errorQuery (errorQuery / totalQuery).
	TotalQuery string `json:"totalQuery"`
	// ErrorQuery is the query that gets the total errors that then will be divided against the total.
	ErrorQuery string `json:"errorQuery"`
}

// Output is how the SLO will expose the generated SLO.
type Output struct {
	//Prometheus is the prometheus format for the SLO output.
	// +optional
	Prometheus *PrometheusOutputSource `json:"prometheus,omitempty"`
}

// PrometheusOutputSource  is the source of the output in prometheus format.
type PrometheusOutputSource struct {
	// Prefix is a prefix to set to the generated output metrics.
	// +optional
	Prefix string `json:"prefix,omitempty"`
	// Labels are the labels that will be set to the output metrics of this SLO.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceLevelList is a list of ServiceLevel resources
type ServiceLevelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ServiceLevel `json:"items"`
}
