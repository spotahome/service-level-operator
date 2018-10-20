package operator

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"

	measurev1alpha1 "github.com/slok/service-level-operator/pkg/apis/measure/v1alpha1"
	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/service/sli"
	"github.com/slok/service-level-operator/pkg/service/slo"
)

// Handler is the Operator handler.
type Handler struct {
	outputerFact  slo.OutputFactory
	retrieverFact sli.RetrieverFactory
	logger        log.Logger
}

// NewHandler returns a new project handler
func NewHandler(outputerFact slo.OutputFactory, retrieverFact sli.RetrieverFactory, logger log.Logger) *Handler {
	return &Handler{
		outputerFact:  outputerFact,
		retrieverFact: retrieverFact,
		logger:        logger,
	}
}

// Add will ensure the the ci builds and jobs are persisted.
func (h *Handler) Add(_ context.Context, obj runtime.Object) error {
	sl, ok := obj.(*measurev1alpha1.ServiceLevel)
	if !ok {
		return fmt.Errorf("can't handle received object, it's not a service level object")
	}

	slc := sl.DeepCopy()

	err := slc.Validate()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(slc.Spec.ServiceLevelObjectives))

	// Retrieve the SLIs.
	for _, slo := range slc.Spec.ServiceLevelObjectives {
		slo := slo

		go func() {
			defer wg.Done()
			err := h.processSLO(slc, &slo)
			// Don't stop if one of the SLOs errors, the rest should
			// be processed independently.
			if err != nil {
				h.logger.With("sl", sl.Name).With("slo", slo.Name).Errorf("error processing SLO: %s", err)
			}
		}()
	}

	wg.Wait()
	return nil
}

func (h *Handler) processSLO(sl *measurev1alpha1.ServiceLevel, slo *measurev1alpha1.SLO) error {
	if slo.Disable {
		h.logger.Debugf("ignoring SLO %s", slo.Name)
		return nil
	}

	retriever, err := h.retrieverFact.GetStrategy(&slo.ServiceLevelIndicator)
	if err != nil {
		return err
	}

	res, err := retriever.Retrieve(&slo.ServiceLevelIndicator)
	if err != nil {
		return err
	}

	outputer, err := h.outputerFact.GetStrategy(slo)
	if err != nil {
		return err
	}

	err = outputer.Create(sl, slo, &res)
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
