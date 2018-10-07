package kubernetes

import (
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	crdcli "github.com/slok/service-level-operator/pkg/k8sautogen/client/clientset/versioned"
)

// ClientFactory knows how to get Kubernetes clients.
type ClientFactory interface {
	// GetSTDClient gets the Kubernetes standard client (pods, services...).
	GetSTDClient() (kubernetes.Interface, error)
	// GetCRDClient gets the Kubernetes client for the CRDs described in this application.
	GetCRDClient() (crdcli.Interface, error)
	// GetAPIExtensionClient gets the Kubernetes api extensions client (crds...).
	GetAPIExtensionClient() (apiextensionscli.Interface, error)
}

type factory struct {
	restCfg *rest.Config

	stdcli kubernetes.Interface
	crdcli crdcli.Interface
	aexcli apiextensionscli.Interface
}

// New returns a new kubernetes client factory.
func NewFactory(config *rest.Config) ClientFactory {
	return &factory{
		restCfg: config,
	}
}

func (f *factory) GetSTDClient() (kubernetes.Interface, error) {
	if f.stdcli == nil {
		cli, err := kubernetes.NewForConfig(f.restCfg)
		if err != nil {
			return nil, err
		}
		f.stdcli = cli
	}
	return f.stdcli, nil
}
func (f *factory) GetCRDClient() (crdcli.Interface, error) {
	if f.crdcli == nil {
		cli, err := crdcli.NewForConfig(f.restCfg)
		if err != nil {
			return nil, err
		}
		f.crdcli = cli
	}
	return f.crdcli, nil
}
func (f *factory) GetAPIExtensionClient() (apiextensionscli.Interface, error) {
	if f.aexcli == nil {
		cli, err := apiextensionscli.NewForConfig(f.restCfg)
		if err != nil {
			return nil, err
		}
		f.aexcli = cli
	}
	return f.aexcli, nil
}
