package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/automaticserver/lxe/cri"
	log "github.com/lxc/lxd/shared/log15"
	"github.com/lxc/lxd/shared/logger"
	"github.com/spf13/cobra"
)

type cmdDaemon struct {
	cmd    *cobra.Command
	global *cmdGlobal

	// Common options
	flagGroup string

	// Debug options
	flagCPUProfile      string
	flagMemoryProfile   string
	flagPrintGoroutines int
}

func (c *cmdDaemon) Command() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = cri.Domain
	cmd.Short = "LXE is a shim of the Kubernetes Container Runtime Interface for LXD"

	cmd.RunE = c.Run
	cmd.Flags().StringVar(&c.flagCPUProfile, "cpu-profile", "", "Enable CPU profiling, writing into the specified file"+"``")
	cmd.Flags().StringVar(&c.flagMemoryProfile, "memory-profile", "", "Enable memory profiling, writing into the specified file"+"``")
	cmd.Flags().IntVar(&c.flagPrintGoroutines, "print-goroutines", 0, "How often to print all the goroutines"+"``")

	c.cmd = cmd

	return cmd
}

func (c *cmdDaemon) Run(cmd *cobra.Command, args []string) error {
	daemonConf := cri.NewDaemonConfig(c.flagGroup, c.global.flagLogTrace)
	d := cri.NewDaemon(daemonConf, &c.global.cri)

	err := d.Init()
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGPWR)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGQUIT)
	signal.Notify(ch, syscall.SIGTERM)

	sig := <-ch
	if sig == syscall.SIGPWR {
		logger.Infof("Received '%s signal', shutting down.", sig)
	} else {
		logger.Infof("Received '%s signal', exiting.", sig)
	}

	return d.Stop()
}

type noHandler struct {
}

// Log 's nothing
func (h noHandler) Log(r *log.Record) error {
	return nil
}
