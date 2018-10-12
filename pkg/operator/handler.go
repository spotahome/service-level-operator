package operator

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/service/sli"
	"github.com/slok/service-level-operator/pkg/service/slo"
)

// Handler is the Operator handler.
type Handler struct {
	outputer  slo.Output
	retriever sli.Retriever
	logger    log.Logger
}

// NewHandler returns a new project handler
func NewHandler(output slo.Output, retriever sli.Retriever, logger log.Logger) *Handler {
	return &Handler{
		outputer:  output,
		retriever: retriever,
		logger:    logger,
	}
}

// Add will ensure the the ci builds and jobs are persisted.
func (h *Handler) Add(_ context.Context, obj runtime.Object) error {
	sl, ok := obj.(*measurev1alpha1.ServiceLevel)
	if !ok {
		return fmt.Errorf("can't handle received object, it's not a service level object")
	}

	slc := sl.DeepCopy()
	// TODO Check the service level is correct.

	// Retrieve the SLIs.
	// TODO: Concurrency and don't stop if one of the SLOs fails.
	for _, slo := range slc.Spec.ServiceLevelObjectives {
		err := h.processSLO(slc, &slo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) processSLO(sl *measurev1alpha1.ServiceLevel, slo *measurev1alpha1.SLO) error {
	if slo.Disable {
		h.logger.Debugf("ignoring SLO %s", slo.Name)
		return nil
	}
	res, err := h.retriever.Retrieve(&slo.ServiceLevelIndicator)
	if err != nil {
		return err
	}

	err = h.outputer.Create(sl, slo, &res)
	if err != nil {
		return err
	}

	return nil
}

// Delete handles the deletion of a release.
func (h *Handler) Delete(_ context.Context, name string) error {
	h.logger.Debugf("delete received")
	return nil
}
