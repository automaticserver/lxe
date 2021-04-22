package cli // import "github.com/automaticserver/lxe/cli"

import (
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"text/template"

	"github.com/automaticserver/lxe/cli/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AnnIsNonoperational = "nonoperational"
)

var (
	envReplacer  = strings.NewReplacer("-", "_")
	keyDelimiter = "-"
)

var rootCmd = &cobra.Command{
	DisableAutoGenTag: true,
	SilenceUsage:      false,
	SilenceErrors:     true,
	Version:           strings.TrimPrefix(version.String(), "version "),
	Args:              cobra.NoArgs,
	Annotations:       map[string]string{},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		// Apply basic logging settings
		err := setLoggingBasic()
		if err != nil {
			return err
		}

		// Several subcommands should not fail if logging hook can be misconfigured
		markIsNonoperational(cmd)
		cmd.VisitParents(markIsNonoperational)

		if IsOperational(cmd) {
			err := setLoggingHook()
			if err != nil {
				return err
			}

			go handleSignals(cmd)
		}

		log.Trace("all settings: ", venom.AllSettings())

		return nil
	},
}

func initRoot() {
	lflags := rootCmd.LocalFlags()
	lflags.Bool("version", false, "Print version information")

	pflags := rootCmd.PersistentFlags()
	pflags.SortFlags = true

	pflags.StringP("config", "c", "", "Load configuration from this file. The path may be absolute or relative. Supported extensions: "+strings.Join(viper.SupportedExts, ", "))

	cobra.OnInitialize(readConfig)
	venom.AllowEmptyEnv(true)
	venom.AutomaticEnv()
}

func readConfig() {
	// Try to find config argument in parameter and env
	c := rootCmd.PersistentFlags().Lookup("config").Value.String()
	if c == "" {
		c = os.Getenv("CONFIG")
	}

	if c != "" {
		venom.SetConfigFile(c)
	} else { // or try to find it in these explicit paths
		venom.SetConfigName(rootCmd.Name())
		venom.AddConfigPath(path.Join("$HOME/.local", version.PackageName))
		venom.AddConfigPath(path.Join("/etc", version.PackageName))
	}

	err := venom.ReadInConfig()
	if err != nil {
		// if it was explicitly set from a parameter, show any error, otherwise show error only when it's not the not found error
		_, is := err.(viper.ConfigFileNotFoundError)
		if c != "" || !is {
			log.WithError(err).Fatal("unable to read config")
		}

		// else show a warning that no config file was loaded
		//log.WithError(err).Warn("no config file was loaded")
	} // nolint: wsl

	// Reset config variable from config file, as this makes no sense to be able to override that from config files
	venom.Set("config", c)
}

func initVersion() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := cmd.Root()
		t := template.New("top")
		t, err := t.Parse(root.VersionTemplate())
		if err != nil {
			return err
		}
		err = t.Execute(root.OutOrStdout(), root)
		if err != nil {
			return err
		}

		return nil
	},
}

func handleSignals(c *cobra.Command) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGUSR2)

	autostartPProf()

	for sig := range ch {
		log.WithField("sig", sig).Info("received signal")

		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			err := gracefulShutdown(c)
			if err != nil {
				log.WithError(err).Fatalf("unable to stop %v", c.Name())
			}

			stopPProf()

			os.Exit(0)
		case syscall.SIGUSR2:
			togglePProf()
		}
	}
}

// gracefulShutdown is a copy of *cobra.Command.execute() with only the relevant post run functions
func gracefulShutdown(c *cobra.Command) error {
	argWoFlags := c.Flags().Args()

	if c.PostRunE != nil {
		if err := c.PostRunE(c, argWoFlags); err != nil {
			return err
		}
	} else if c.PostRun != nil {
		c.PostRun(c, argWoFlags)
	}

	for p := c; p != nil; p = p.Parent() {
		if p.PersistentPostRunE != nil {
			if err := p.PersistentPostRunE(c, argWoFlags); err != nil {
				return err
			}

			break
		} else if p.PersistentPostRun != nil {
			p.PersistentPostRun(c, argWoFlags)

			break
		}
	}

	return nil
}

func markIsNonoperational(c *cobra.Command) {
	if c == confCmd || c == cmplCmd || c == versionCmd {
		c.Root().Annotations[AnnIsNonoperational] = strconv.FormatBool(true)
	}
}
