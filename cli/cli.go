package cli // import "github.com/automaticserver/lxe/cli"

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// TODO apply as argument?
	log   = logrus.StandardLogger()
	venom *viper.Viper
)

// New instantiates a new root command with defaults based on submitted type
func New() (*viper.Viper, *cobra.Command) {
	venom = viper.NewWithOptions(viper.EnvKeyReplacer(envReplacer), viper.KeyDelimiter(keyDelimiter))

	initRoot()
	initCmpl()
	initConf()
	initLog()
	initVersion()

	ToDaemon(rootCmd)

	return venom, rootCmd
}

func Run() {
	// Bind pflags as late as possible so all imports were able to set their flags
	pflags := rootCmd.PersistentFlags()

	err := venom.BindPFlags(pflags)
	if err != nil {
		log.Fatal(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func IsOperational(cmd *cobra.Command) bool {
	return cmd.Root().Annotations[AnnIsNonoperational] == ""
}
