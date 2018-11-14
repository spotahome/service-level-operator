package kubernetes

import (
	koopercrd "github.com/spotahome/kooper/client/crd"
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"github.com/spotahome/service-level-operator/pkg/log"
)

// CRDConf is the configuration of the crd.
type CRDConf = koopercrd.Conf

// CRD is the CRD service that knows how to interact with k8s to manage them.
type CRD interface {
	// EnsurePresentCRD will create the custom resource and wait to be ready
	// if there is not already present.
	EnsurePresentCRD(conf CRDConf) error
}

// crdService is the CRD service implementation using API calls to kubernetes.
type crd struct {
	crdCli koopercrd.Interface
	logger log.Logger
}

// NewCRD returns a new CRD KubeService.
func NewCRD(aeClient apiextensionscli.Interface, logger log.Logger) CRD {
	logger = logger.With("service", "k8s.crd")
	crdCli := koopercrd.NewClient(aeClient, logger)

	return &crd{
		crdCli: crdCli,
		logger: logger,
	}
}

// EnsurePresentCRD satisfies workspace.Service interface.
func (c *crd) EnsurePresentCRD(conf CRDConf) error {
	return c.crdCli.EnsurePresent(conf)
}
