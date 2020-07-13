package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/automaticserver/lxe/cri"
	log "github.com/lxc/lxd/shared/log15"
	"github.com/lxc/lxd/shared/logger"
	"github.com/spf13/cobra"
)

type cmdDaemon struct {
	cmd    *cobra.Command
	global *cmdGlobal
}

func (c *cmdDaemon) Command() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = cri.Domain
	cmd.Short = "LXE is a shim of the Kubernetes Container Runtime Interface for LXD"

	cmd.RunE = c.Run

	c.cmd = cmd

	return cmd
}

func (c *cmdDaemon) Run(cmd *cobra.Command, args []string) error {
	d := cri.NewDaemon(&c.global.cri)

	err := d.Init()
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	// let go handle SIGQUIT itself, so not added here
	signal.Notify(ch, syscall.SIGPWR)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGUSR2)

	for sig := range ch {
		logger.Infof("received signal %v", sig)

		switch sig {
		case syscall.SIGPWR, syscall.SIGINT, syscall.SIGTERM:
			logger.Warn("shutting down")
			return d.Stop()
		case syscall.SIGUSR2:
			// Allow manual dump of goroutines until pprof is implemented
			err := dumpGoroutines()
			if err != nil {
				logger.Errorf("Unable to dump goroutines: %v", err)
			}
		}
	}

	return nil
}

type noHandler struct {
}

// Log 's nothing
func (h noHandler) Log(r *log.Record) error {
	return nil
}

func dumpGoroutines() error {
	file, err := ioutil.TempFile(os.TempDir(), "lxe-routines-")
	if err != nil {
		return err
	}
	defer file.Close()

	// Based on answers to this stackoverflow question:
	// https://stackoverflow.com/questions/19094099/how-to-dump-goroutine-stacktraces
	var buf []byte
	var bufsize int
	var stacklen int

	// Create a stack buffer of 1MB and grow it to at most 100MB if necessary
	for bufsize = 1e6; bufsize < 100e6; bufsize *= 2 {
		buf = make([]byte, bufsize)
		stacklen = runtime.Stack(buf, true)

		if stacklen < bufsize {
			break
		}
	}

	_, err = file.Write(buf[:stacklen])
	if err != nil {
		return err
	}

	logger.Errorf("Wrote goroutine dump to %s", file.Name())

	return nil
}
