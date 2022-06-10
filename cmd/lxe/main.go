package main

import (
	"github.com/automaticserver/lxe/cli"
	"github.com/automaticserver/lxe/cri"
	"github.com/automaticserver/lxe/network"
	"github.com/dionysius/errand"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	venom, rootCmd = cli.New()
	log            = logrus.StandardLogger()
)

func main() {
	cli.Run()
}

func init() {
	rootCmd.Use = "lxe"
	rootCmd.Short = "LXE is a shim of the Kubernetes Container Runtime Interface for LXD"
	rootCmd.Long = "LXE implements the Kubernetes Container Runtime Interface and creates LXD containers from Pods. Many options in the PodSpec and ContainerSpec are honored but a fundamental perception of containers (application containers like docker does vs system containers which is what LXD does) lead to slight implementation differences. Read the documentation for the full list of caveats."
	rootCmd.Example = "lxe"

	pflags := rootCmd.PersistentFlags()

	// application flags
	pflags.StringP("socket", "s", "/run/lxe.sock", "Path of the socket where it should provide the runtime and image service to kubelet.")
	pflags.StringP("lxd-socket", "l", "", "Path of the socket where LXD provides it's API. (guessed by default)")
	pflags.StringP("lxd-remote-config", "r", "", "Path to the LXD remote config. (guessed by default)")
	pflags.StringP("lxd-image-remote", "", "local", "Use this remote if ImageSpec doesn't provide an explicit remote.")
	pflags.StringSliceP("lxd-profiles", "p", []string{"default"}, "Set these additional profiles when creating containers.")
	pflags.StringP("streaming-bindaddr", "", "localhost:44124", "Listen address for the streaming service. Be careful from where this service can be accessed from as it allows to run exec commands on the containers! Format: [IP]:Port.")
	pflags.StringP("streaming-baseurl", "", "", "Define which base address to use for constructing streaming URLs for a client to connect to. If this is set to empty, it will use the same host address and port from --streaming-bindaddr. If that has an empty host address, it will obtain the address of the interface to the default gateway. Format: [IP][:Port].")
	// TODO: I was thinking, can't we just create a tmpfile with those contents when running lxe and remember that? Maybe, but it must be a persistent location, otherwise containers won't be able to start without that file existing.
	pflags.StringP("hostnetwork-file", "", "", "EXPERIMENTAL! If host networking is defined in the PodSpec, this persisting file will be set as include in raw.lxc container config. (This process is required to workaround LXD, since it doesn't offer such option in the container or device config out of the box). The file must contain: 'lxc.net.0.type=none'.")
	pflags.StringP("network-plugin", "n", "bridge", "The network plugin to use. 'bridge' manages the lxd bridge defined in --bridge-name. 'cni' uses container network interface to attach interfaces using a configuration defined in --cni-conf-dir.")
	pflags.StringP("bridge-name", "", network.DefaultLXDBridge, "Which bridge to create and use when using --network-plugin 'bridge'.")
	pflags.StringP("bridge-dhcp-range", "", "", "Which DHCP range to configure the lxd bridge when using --network-plugin 'bridge'. If empty, uses random range provided by lxd. Not needed, if kubernetes will publish the range using CRI UpdateRuntimeconfig.")
	pflags.StringP("cni-conf-dir", "", network.DefaultCNIconfPath, "Dir in which to search for CNI configuration files when using --network-plugin 'cni'.")
	pflags.StringP("cni-bin-dir", "", network.DefaultCNIbinPath, "Dir in which to search for CNI plugin binaries when using --network-plugin 'cni'.")
	pflags.StringP("cni-output-target", "", "stderr", "Where to forward the cni command output, one of: stdout, stderr, file.")
	pflags.StringP("cni-output-file-path", "", "stderr", "Path to output file. Only required if --cni-output-target is set to file.")
	pflags.BoolP("critest", "", false, "Enable critest mode to be used with cri-tools' critest since LXE cannot use OCI images. Automatically creates some images and changes a few image names before handling the request. Read the log for arguments to use cri-tools' critest.")

	rootCmd.RunE = rootCmdRunE
}

func rootCmdRunE(cmd *cobra.Command, args []string) error {
	conf := &cri.Config{
		UnixSocket:           venom.GetString("socket"),
		LXDSocket:            venom.GetString("lxd-socket"),
		LXDRemoteConfig:      venom.GetString("lxd-remote-config"),
		LXDImageRemote:       venom.GetString("lxd-image-remote"),
		LXDProfiles:          venom.GetStringSlice("lxd-profiles"),
		LXEStreamingBindAddr: venom.GetString("streaming-bindaddr"),
		LXEStreamingBaseURL:  venom.GetString("streaming-baseurl"),
		LXEHostnetworkFile:   venom.GetString("hostnetwork-file"),
		LXENetworkPlugin:     venom.GetString("network-plugin"),
		LXDBridgeName:        venom.GetString("bridge-name"),
		LXDBridgeDHCPRange:   venom.GetString("bridge-dhcp-range"),
		CNIConfDir:           venom.GetString("cni-conf-dir"),
		CNIBinDir:            venom.GetString("cni-bin-dir"),
		CNIOutputTarget:      venom.GetString("cni-output-target"),
		CNIOutputFile:        venom.GetString("cni-output-file-path"),
		CRITest:              venom.GetBool("critest"),
	}

	criServer := cri.NewServer(conf)

	go func() {
		err := errand.Append(nil, criServer.Serve())
		if err != nil {
			err = errand.Append(err, criServer.Stop())
			log.WithError(err).Fatal("unable to start CRI server")
		}
	}()

	// run forever
	select {}
}
