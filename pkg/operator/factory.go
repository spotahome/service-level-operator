package operator

import (
	"time"

	"github.com/spotahome/kooper/operator"
	"github.com/spotahome/kooper/operator/controller"

	"github.com/slok/service-level-operator/pkg/log"
	"github.com/slok/service-level-operator/pkg/service/kubernetes"
)

// Config is the configuration for the ci operator.
type Config struct {
	// ResyncPeriod is the resync period of the controllers.
	ResyncPeriod time.Duration
	// ConcurretWorkers are number of workers to handle the events.
	ConcurretWorkers int
}

// New returns pod terminator operator.
func New(cfg Config, k8ssvc kubernetes.Service, logger log.Logger) (operator.Operator, error) {

	// Create crd.
	ptCRD := newServiceLevelCRD(k8ssvc, logger)

	// Create handler.
	handler := newHandler(logger)

	// Create controller.
	ctrl := controller.NewSequential(cfg.ResyncPeriod, handler, ptCRD, nil, logger)

	// Assemble CRD and controller to create the operator.
	return operator.NewOperator(ptCRD, ctrl, logger), nil
}
