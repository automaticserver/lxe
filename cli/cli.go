package cli // import "github.com/automaticserver/lxe/cli"

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// TODO apply as argument?
	log   = logrus.StandardLogger()
	venom *viper.Viper
)

type Type string

const (
	AnnType          = "type"
	TypeTool    Type = "tool"
	TypeService Type = "service"
)

type Options struct {
	Type         Type
	KeyDelimiter string
	EnvReplace   []string
}

// New instantiates a new root command with defaults based on submitted type
func New(o Options) (*viper.Viper, *cobra.Command) {
	if len(o.EnvReplace) > 0 {
		envReplacer = strings.NewReplacer(o.EnvReplace...)
	}

	if o.KeyDelimiter != "" {
		keyDelimiter = o.KeyDelimiter
	}

	venom = viper.NewWithOptions(viper.EnvKeyReplacer(envReplacer), viper.KeyDelimiter(keyDelimiter))

	initRoot()
	initCmpl()
	initConf()
	initLog()
	initVersion()

	rootCmd.Annotations[AnnType] = string(o.Type)

	if o.Type == TypeService {
		ToDaemon(rootCmd)
	}

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

func Is(cmd *cobra.Command, t Type) bool {
	return cmd.Root().Annotations[AnnType] == string(t)
}

func IsOperational(cmd *cobra.Command) bool {
	return cmd.Root().Annotations[AnnIsNonoperational] == ""
}
