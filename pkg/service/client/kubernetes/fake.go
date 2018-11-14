package kubernetes

import (
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsclifake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	monitoringv1alpha1 "github.com/spotahome/service-level-operator/pkg/apis/monitoring/v1alpha1"
	crdcli "github.com/spotahome/service-level-operator/pkg/k8sautogen/client/clientset/versioned"
	crdclifake "github.com/spotahome/service-level-operator/pkg/k8sautogen/client/clientset/versioned/fake"
)

// fakeFactory is a fake factory that has already loaded faked objects on the Kubernetes clients.
type fakeFactory struct{}

// NewFake returns the faked Kubernetes clients factory.
func NewFake() ClientFactory {
	return &fakeFactory{}
}

func (f *fakeFactory) GetSTDClient() (kubernetes.Interface, error) {
	return kubernetesfake.NewSimpleClientset(stdObjs...), nil
}
func (f *fakeFactory) GetCRDClient() (crdcli.Interface, error) {
	return crdclifake.NewSimpleClientset(crdObjs...), nil
}
func (f *fakeFactory) GetAPIExtensionClient() (apiextensionscli.Interface, error) {
	cli := apiextensionsclifake.NewSimpleClientset(aexObjs...)

	// Fake cluster version (Required for CRD version checks).
	fakeDiscovery, _ := cli.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v1.10.5",
	}

	return cli, nil
}

var (
	stdObjs = []runtime.Object{}

	// The field selector doesn't work with a fake K8s client: https://github.com/kubernetes/client-go/issues/326
	crdObjs = []runtime.Object{
		&monitoringv1alpha1.ServiceLevel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service0",
				Namespace: "ns0",
				Labels: map[string]string{
					"wrong": "false",
				},
			},
			Spec: monitoringv1alpha1.ServiceLevelSpec{
				ServiceLevelObjectives: []monitoringv1alpha1.SLO{
					{
						Name:                         "fake_slo0",
						Description:                  "fake slo 0.",
						AvailabilityObjectivePercent: 99.99,
						ServiceLevelIndicator: monitoringv1alpha1.SLI{
							SLISource: monitoringv1alpha1.SLISource{
								Prometheus: &monitoringv1alpha1.PrometheusSLISource{
									Address:    "http://fake:9090",
									TotalQuery: `slo0_total`,
									ErrorQuery: `slo0_error`,
								},
							},
						},
						Output: monitoringv1alpha1.Output{
							Prometheus: &monitoringv1alpha1.PrometheusOutputSource{
								Labels: map[string]string{
									"fake": "true",
									"team": "fake-team0",
								},
							},
						},
					},
					{
						Name:                         "fake_slo1",
						Description:                  "fake slo 1.",
						AvailabilityObjectivePercent: 99.9,
						ServiceLevelIndicator: monitoringv1alpha1.SLI{
							SLISource: monitoringv1alpha1.SLISource{
								Prometheus: &monitoringv1alpha1.PrometheusSLISource{
									Address:    "http://fake:9090",
									TotalQuery: `slo1_total`,
									ErrorQuery: `slo1_error`,
								},
							},
						},
						Output: monitoringv1alpha1.Output{
							Prometheus: &monitoringv1alpha1.PrometheusOutputSource{
								Labels: map[string]string{
									"fake": "true",
									"team": "fake-team1",
								},
							},
						},
					},
					{
						Name:                         "fake_slo2",
						Description:                  "fake slo 2.",
						AvailabilityObjectivePercent: 99.998,
						ServiceLevelIndicator: monitoringv1alpha1.SLI{
							SLISource: monitoringv1alpha1.SLISource{
								Prometheus: &monitoringv1alpha1.PrometheusSLISource{
									Address:    "http://fake:9090",
									TotalQuery: `slo2_total`,
									ErrorQuery: `slo2_error`,
								},
							},
						},
						Output: monitoringv1alpha1.Output{
							Prometheus: &monitoringv1alpha1.PrometheusOutputSource{
								Labels: map[string]string{
									"fake": "true",
									"team": "fake-team2",
								},
							},
						},
					},
				},
			},
		},
		&monitoringv1alpha1.ServiceLevel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service1",
				Namespace: "ns1",
				Labels: map[string]string{
					"wrong": "false",
				},
			},
			Spec: monitoringv1alpha1.ServiceLevelSpec{
				ServiceLevelObjectives: []monitoringv1alpha1.SLO{
					{
						Name:                         "fake_slo3",
						Description:                  "fake slo 3.",
						AvailabilityObjectivePercent: 99,
						ServiceLevelIndicator: monitoringv1alpha1.SLI{
							SLISource: monitoringv1alpha1.SLISource{
								Prometheus: &monitoringv1alpha1.PrometheusSLISource{
									Address:    "http://fake:9090",
									TotalQuery: `slo3_total`,
									ErrorQuery: `slo3_error`,
								},
							},
						},
						Output: monitoringv1alpha1.Output{
							Prometheus: &monitoringv1alpha1.PrometheusOutputSource{
								Labels: map[string]string{
									"fake": "true",
									"team": "fake-team3",
								},
							},
						},
					},
				},
			},
		},
		&monitoringv1alpha1.ServiceLevel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service2-no-output",
				Namespace: "ns0",
				Labels: map[string]string{
					"wrong": "true",
				},
			},
			Spec: monitoringv1alpha1.ServiceLevelSpec{
				ServiceLevelObjectives: []monitoringv1alpha1.SLO{
					{
						Name:                         "fake_slo4",
						Description:                  "fake slo 4.",
						AvailabilityObjectivePercent: 99,
						ServiceLevelIndicator: monitoringv1alpha1.SLI{
							SLISource: monitoringv1alpha1.SLISource{
								Prometheus: &monitoringv1alpha1.PrometheusSLISource{
									Address:    "http://fake:9090",
									TotalQuery: `slo3_total`,
									ErrorQuery: `slo3_error`,
								},
							},
						},
						Output: monitoringv1alpha1.Output{},
					},
				},
			},
		},

		&monitoringv1alpha1.ServiceLevel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service3-no-input",
				Namespace: "ns1",
				Labels: map[string]string{
					"wrong": "true",
				},
			},
			Spec: monitoringv1alpha1.ServiceLevelSpec{
				ServiceLevelObjectives: []monitoringv1alpha1.SLO{
					{
						Name:                         "fake_slo5",
						Description:                  "fake slo 5.",
						AvailabilityObjectivePercent: 99,
						ServiceLevelIndicator:        monitoringv1alpha1.SLI{},
						Output: monitoringv1alpha1.Output{
							Prometheus: &monitoringv1alpha1.PrometheusOutputSource{
								Labels: map[string]string{
									"wrong": "true",
								},
							},
						},
					},
				},
			},
		},
	}

	aexObjs = []runtime.Object{}
)
