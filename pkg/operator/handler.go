package operator

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
)

// handler is the Operator handler.
type handler struct {
	logger log.Logger
}

// NewHandler returns a new project handler
func newHandler(logger log.Logger) *handler {
	return &handler{
		logger: logger,
	}
}

// Add will ensure the the ci builds and jobs are persisted.
func (h *handler) Add(_ context.Context, obj runtime.Object) error {
	sl, ok := obj.(*measurev1alpha1.ServiceLevel)
	if !ok {
		return fmt.Errorf("can't handle received object, it's not a service level object")
	}

	slc := sl.DeepCopy()
	h.logger.Infof("serviceLevel: %s/%s", slc.Namespace, slc.Name)

	return nil
}

// Delete handles the deletion of a release.
func (h *handler) Delete(_ context.Context, name string) error {
	h.logger.Debugf("delete received")
	return nil
}
