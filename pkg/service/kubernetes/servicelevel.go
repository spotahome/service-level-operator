package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	crdcli "github.com/slok/service-level-operator/pkg/k8sautogen/client/clientset/versioned"
	"github.com/slok/service-level-operator/pkg/log"
)

// ServiceLevel knows how to interact with Kubernetes on the
// ServiceLevel CRs
type ServiceLevel interface {
	// ListServiceLevels will list the service levels.
	ListServiceLevels(namespace string, opts metav1.ListOptions) (*measurev1alpha1.ServiceLevelList, error)
	// ListServiceLevels will list the service levels.
	WatchServiceLevels(namespace string, opt metav1.ListOptions) (watch.Interface, error)
}

type serviceLevel struct {
	cli    crdcli.Interface
	logger log.Logger
}

// NewServiceLevel returns a new service level service.
func NewServiceLevel(crdcli crdcli.Interface, logger log.Logger) ServiceLevel {
	return &serviceLevel{
		cli:    crdcli,
		logger: logger,
	}
}

func (s *serviceLevel) ListServiceLevels(namespace string, opts metav1.ListOptions) (*measurev1alpha1.ServiceLevelList, error) {
	return s.cli.MeasureV1alpha1().ServiceLevels(namespace).List(opts)
}
func (s *serviceLevel) WatchServiceLevels(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return s.cli.MeasureV1alpha1().ServiceLevels(namespace).Watch(opts)
}
