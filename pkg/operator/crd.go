package operator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/service/kubernetes"
)

// serviceLevelCRD is the crd release.
type serviceLevelCRD struct {
	cfg     Config
	service kubernetes.Service
	logger  log.Logger
}

func newServiceLevelCRD(cfg Config, service kubernetes.Service, logger log.Logger) *serviceLevelCRD {
	logger = logger.With("crd", "servicelevel")
	return &serviceLevelCRD{
		cfg:     cfg,
		service: service,
		logger:  logger,
	}
}

// Initialize satisfies resource.crd interface.
func (s *serviceLevelCRD) Initialize() error {
	crd := kubernetes.CRDConf{
		Kind:                    measurev1alpha1.ServiceLevelKind,
		NamePlural:              measurev1alpha1.ServiceLevelNamePlural,
		Group:                   measurev1alpha1.SchemeGroupVersion.Group,
		Version:                 measurev1alpha1.SchemeGroupVersion.Version,
		Scope:                   measurev1alpha1.ServiceLevelScope,
		Categories:              []string{"measure", "slo"},
		EnableStatusSubresource: true,
	}

	return s.service.EnsurePresentCRD(crd)
}

// GetListerWatcher satisfies resource.crd interface (and retrieve.Retriever).
func (s *serviceLevelCRD) GetListerWatcher() cache.ListerWatcher {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.LabelSelector = s.cfg.LabelSelector
			return s.service.ListServiceLevels(s.cfg.Namespace, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.LabelSelector = s.cfg.LabelSelector
			return s.service.WatchServiceLevels(s.cfg.Namespace, options)
		},
	}
}

// GetObject satisfies resource.crd interface (and retrieve.Retriever).
func (s *serviceLevelCRD) GetObject() runtime.Object {
	return &measurev1alpha1.ServiceLevel{}
}
