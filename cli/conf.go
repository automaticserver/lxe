package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initConf() {
	rootCmd.AddCommand(confCmd)
	confCmd.AddCommand(confShowCmd)
}

var confCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Manage configuration options",
	Hidden:  false,
	Args:    cobra.NoArgs,
}

var confShowCmd = &cobra.Command{
	Use:       "show <extension>",
	Short:     "Display the currently loaded configuration in specified format",
	Args:      cobra.ExactArgs(1),
	ValidArgs: viper.SupportedExts,
	RunE:      confShowCmdRun,
}

func confShowCmdRun(cmd *cobra.Command, args []string) error {
	filename := filepath.Join(os.TempDir(), fmt.Sprintf("%s-config.%s", cmd.Root().Name(), args[0]))

	err := venom.WriteConfigAs(filename)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = cmd.Root().OutOrStdout().Write(b)

	return err
}
