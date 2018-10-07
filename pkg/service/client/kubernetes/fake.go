package kubernetes

import (
	crdcli "github.com/slok/service-level-operator/pkg/k8sautogen/client/clientset/versioned"
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsclifake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	crdclifake "github.com/slok/service-level-operator/pkg/k8sautogen/client/clientset/versioned/fake"
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
		&measurev1alpha1.ServiceLevel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service0",
				Namespace: "fake",
			},
			Spec: measurev1alpha1.ServiceLevelSpec{},
		},
		&measurev1alpha1.ServiceLevel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service1",
				Namespace: "fake",
			},
			Spec: measurev1alpha1.ServiceLevelSpec{},
		},
		&measurev1alpha1.ServiceLevel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service2",
				Namespace: "fake",
			},
			Spec: measurev1alpha1.ServiceLevelSpec{},
		},
	}

	aexObjs = []runtime.Object{}
)
