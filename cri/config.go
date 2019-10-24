package cri

// Config options that LXE will need to interface with LXD
type Config struct {
	// UnixSocket this LXE will be reachable under
	UnixSocket string
	// LXDSocket where LXD is reachable under
	LXDSocket string
	// LXDRemoteConfig file path where lxd remote settings are stored
	LXDRemoteConfig string
	// LXDImageRemote to use by default when ImageSpec doesn't provide an explicit remote
	LXDImageRemote string
	// LXDProfiles which all cri containers inherit
	LXDProfiles []string
	// LXEStreamingServerEndpoint contains the listen address for the streaming server
	LXEStreamingServerEndpoint string
	// LXEStreamingPort is the port for the streaming server
	LXEStreamingPort int
	// LXEHostnetworkFile file path to use for lxc's raw.include
	LXEHostnetworkFile string
	// Which LXENetworkPlugin to use
	LXENetworkPlugin string
	// LXEBridgeDHCPRange to configure for lxebr0 if NetworkPlugin is default
	LXEBridgeDHCPRange string
}
