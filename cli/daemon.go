package cli // import "github.com/automaticserver/lxe/cli"

import (
	"github.com/automaticserver/lxe/cli/version"
	"github.com/spf13/cobra"
)

func daemonPreRunE(cmd *cobra.Command, args []string) error {
	log.WithFields(version.MapInf()).Warnf("starting %s...", cmd.Name())

	return nil
}

func daemonPostRunE(cmd *cobra.Command, args []string) error {
	log.WithFields(version.MapInf()).Warnf("stopping %s...", cmd.Name())

	return nil
}

func ToDaemon(c *cobra.Command) {
	c.PreRunE = daemonPreRunE
	c.PostRunE = daemonPostRunE
}
