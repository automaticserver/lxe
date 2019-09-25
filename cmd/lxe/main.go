package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/automaticserver/lxe/cri"
	"github.com/automaticserver/lxe/shared"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxd/shared/logging"
	"github.com/spf13/cobra"
)

// Global variables
var debug bool
var verbose bool

// Initialize the random number generator
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type cmdGlobal struct {
	flagHelp    bool
	flagVersion bool

	flagLogFile    string
	flagLogDebug   bool
	flagLogSyslog  bool
	flagLogTrace   []string
	flagLogVerbose bool

	cri cri.Config
}

func (c *cmdGlobal) Run(cmd *cobra.Command, args []string) error {
	// TODO: maybe consider using k8slogs?
	//import k8slogs "github.com/kubernetes/kubernetes/pkg/kubectl/util/logs"
	//k8slogs.InitLogs()
	//#defer k8slogs.FlushLogs()

	// Set logging global variables
	debug = c.flagLogVerbose
	verbose = c.flagLogDebug

	// Setup logger
	syslog := ""
	if c.flagLogSyslog {
		syslog = cri.Domain
	}

	handler := noHandler{}
	log, err := logging.GetLogger(syslog, c.flagLogFile, c.flagLogVerbose, c.flagLogDebug, handler)
	if err != nil {
		return err
	}
	logger.Log = log

	return nil
}

func main() {
	// daemon command (main)
	daemonCmd := cmdDaemon{}
	app := daemonCmd.Command()

	// Workaround for main command
	app.Args = cobra.ArbitraryArgs
	app.Version = cri.Version

	// Global flags
	globalCmd := cmdGlobal{}
	daemonCmd.global = &globalCmd
	app.PersistentPreRunE = globalCmd.Run
	app.PersistentFlags().BoolVar(&globalCmd.flagVersion, "version", false, "Print version number.")
	app.PersistentFlags().BoolVarP(&globalCmd.flagHelp, "help", "h", false, "Print help.")
	app.PersistentFlags().StringVar(&globalCmd.flagLogFile, "logfile", "/var/log/lxe.log", "Path to the log file."+"``")
	app.PersistentFlags().StringArrayVar(&globalCmd.flagLogTrace, "trace", []string{}, "Log tracing targets."+"``")
	app.PersistentFlags().BoolVarP(&globalCmd.flagLogDebug, "debug", "d", false, "Show all debug messages.")
	app.PersistentFlags().BoolVarP(&globalCmd.flagLogVerbose, "verbose", "v", false, "Show all information messages.")

	// lxd / lxe specific flags
	app.PersistentFlags().StringVar(&globalCmd.cri.UnixSocket, "socket",
		"/var/run/lxe.sock", "The unix socket under which LXE will expose its service to Kubernetes.")
	app.PersistentFlags().StringVar(&globalCmd.cri.LXDSocket, "lxd-socket",
		"/var/lib/lxd/unix.socket", "LXD's unix socket.")
	app.PersistentFlags().StringVar(&globalCmd.cri.LXDRemoteConfig, "lxd-remote-config",
		"", "Path to the LXD remote config. (guessed by default)")
	app.PersistentFlags().StringVar(&globalCmd.cri.LXDImageRemote, "lxd-image-remote",
		"local", "Use this remote when ImageSpec doesn't provide an explicit remote.")
	app.PersistentFlags().StringSliceVar(&globalCmd.cri.LXDProfiles, "lxd-profiles",
		[]string{"default"}, "Set these additional profiles when creating containers.")
	app.PersistentFlags().StringVar(&globalCmd.cri.LXEStreamingServerEndpoint, "streaming-endpoint",
		"", "IP or Interface for Streaming Server. (guessed by default)")
	app.PersistentFlags().IntVar(&globalCmd.cri.LXEStreamingPort, "streaming-port",
		44124, "Port where LXE's Streaming HTTP Server will listen.")
	app.PersistentFlags().StringVar(&globalCmd.cri.LXEHostnetworkFile, "hostnetwork-file",
		"/var/lib/lxe/hostnetwork.conf", "Path to the hostnetwork file for lxc raw include")
	app.PersistentFlags().StringVar(&globalCmd.cri.LXENetworkPlugin, "network-plugin",
		"", "The network plugin to use. '' is the standard network plugin and manages a lxd bridge 'lxebr0'. 'cni' uses kubernetes cni tools to attach interfaces.")
	app.PersistentFlags().StringVar(&globalCmd.cri.LXEBridgeDHCPRange, "bridge-dhcp-range",
		"", "Which DHCP range to configure in the 'lxebr0' lxd bridge. Only applies if network-plugin is set to '' and does not exist yet. If empty, uses random range provided by lxd. Not needed, if kubernetes will publish the range using CRI UpdateRuntimeconfig")

	// Run the main command and handle errors
	err := app.Execute()
	if err != nil {
		os.Exit(shared.ExitCodeUnspecified)
	}
}
