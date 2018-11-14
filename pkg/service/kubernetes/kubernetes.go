package kubernetes

import (
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"

	crdcli "github.com/spotahome/service-level-operator/pkg/k8sautogen/client/clientset/versioned"
	"github.com/spotahome/service-level-operator/pkg/log"
)

// Service is the service used to interact with the Kubernetes
// objects.
type Service interface {
	ServiceLevel
	CRD
}

type service struct {
	ServiceLevel
	CRD
}

// New returns a new Kubernetes service.
func New(stdcli kubernetes.Interface, crdcli crdcli.Interface, apiextcli apiextensionscli.Interface, logger log.Logger) Service {
	return &service{
		ServiceLevel: NewServiceLevel(crdcli, logger),
		CRD:          NewCRD(apiextcli, logger),
	}
}
