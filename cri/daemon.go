package cri

import (
	"fmt"

	"github.com/lxc/lxd/shared/logger"
)

// Domain of the daemon
const Domain = "lxe"

// A Daemon can respond to requests from a shared client.
type Daemon struct {
	setupChan    chan struct{} // Closed when basic Daemon setup is completed
	shutdownChan chan struct{}

	daemonConfig *DaemonConfig
	cri          *Server
	criConfig    *Config
}

// DaemonConfig holds configuration values for Daemon.
type DaemonConfig struct {
	Group string   // Group name the local unix socket should be chown'ed to
	Trace []string // List of sub-systems to trace
}

// NewDaemon returns a new Daemon object with the given configuration.
func NewDaemon(daemonConfig *DaemonConfig, criConfig *Config) *Daemon {
	return &Daemon{
		daemonConfig: daemonConfig,
		criConfig:    criConfig,
		setupChan:    make(chan struct{}),
		shutdownChan: make(chan struct{}),
	}
}

// NewDaemonConfig returns a DaemonConfig object
func NewDaemonConfig(group string, trace []string) *DaemonConfig {
	return &DaemonConfig{
		Group: group,
		Trace: trace,
	}
}

// Init the daemon
func (d *Daemon) Init() error {
	err := d.init()

	// If an error occurred synchronously while starting up, let's try to
	// cleanup any state we produced so far. Errors happening here will be
	// ignored.
	if err != nil {
		err2 := d.Stop()
		if err2 != nil {
			logger.Errorf("Init errored and also errored during stop: %v", err2)
		}
	}

	return err
}

func (d *Daemon) init() error {
	// Initialize CRI server
	d.cri = NewServer(d.criConfig)
	go d.cri.Serve() // nolint

	return nil
}

// Kill signals the daemon that we want to shutdown, and that any work
// initiated from this point (e.g. database queries over gRPC) should not be
// retried in case of failure.
func (d *Daemon) Kill() {
	//d.cri.Kill()
}

// Stop stops the shared daemon.
func (d *Daemon) Stop() error {
	errs := []error{}
	trackError := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	trackError(d.cri.Stop())

	var err error
	if n := len(errs); n > 0 {
		format := "%v"
		if n > 1 {
			format += fmt.Sprintf(" (and %d more errors)", n)
		}
		err = fmt.Errorf(format, errs[0])
	}
	return err
}
