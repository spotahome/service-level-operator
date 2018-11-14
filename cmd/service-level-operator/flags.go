package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/util/homedir"

	"github.com/spotahome/service-level-operator/pkg/operator"
)

// defaults
const (
	defMetricsPath   = "/metrics"
	defListenAddress = ":8080"
	defResyncSeconds = 5
	defWorkers       = 10
)

type cmdFlags struct {
	fs *flag.FlagSet

	kubeConfig    string
	resyncSeconds int
	workers       int
	metricsPath   string
	listenAddress string
	labelSelector string
	namespace     string
	debug         bool
	development   bool
	fake          bool
}

func newCmdFlags() *cmdFlags {
	c := &cmdFlags{
		fs: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
	c.init()

	return c
}

func (c *cmdFlags) init() {

	kubehome := filepath.Join(homedir.HomeDir(), ".kube", "config")
	// register flags
	c.fs.StringVar(&c.kubeConfig, "kubeconfig", kubehome, "kubernetes configuration path, only used when development mode enabled")
	c.fs.StringVar(&c.metricsPath, "metrics-path", defMetricsPath, "the path where the metrics will be served")
	c.fs.StringVar(&c.listenAddress, "listen-addr", defListenAddress, "the address where the metrics will be exposed")
	c.fs.StringVar(&c.labelSelector, "selector", "", "selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	c.fs.StringVar(&c.namespace, "namespace", "", "the namespace to filter on, by default all")
	c.fs.IntVar(&c.resyncSeconds, "resync-seconds", defResyncSeconds, "the number of seconds for the SLO calculation interval")
	c.fs.IntVar(&c.workers, "workers", defWorkers, "the number of concurrent workers per controller handling events")
	c.fs.BoolVar(&c.development, "development", false, "development flag will allow to run the operator outside a kubernetes cluster")
	c.fs.BoolVar(&c.debug, "debug", false, "enable debug mode")
	c.fs.BoolVar(&c.fake, "fake", false, "enable faked mode, in faked node external services/dependencies are not needed")

	// Parse flags
	c.fs.Parse(os.Args[1:])
}

func (c *cmdFlags) toOperatorConfig() operator.Config {
	return operator.Config{
		ResyncPeriod:     time.Duration(c.resyncSeconds) * time.Second,
		ConcurretWorkers: c.workers,
		LabelSelector:    c.labelSelector,
		Namespace:        c.namespace,
	}
}
