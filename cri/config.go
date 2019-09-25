package cri

// Config options that LXE will need to interface with LXD
type Config struct {
	UnixSocket                 string // Unix socket this LXE will be reachable under
	LXDSocket                  string // Unix socket target LXD is reachable under
	LXDRemoteConfig            string // Path where the lxd remote config can be found
	LXDImageRemote             string // Remote to use when ImageSpec doesn't provide an explicit remote
	LXEStreamingServerEndpoint string // IP or Interface for Streaming Server. Guessed by default if not present
	LXEStreamingPort           int    // Port where LXE's Http Server will listen
	LXEHostnetworkFile         string // Path to the hostnetwork file for lxc raw include
	LXENetworkPlugin           string // The network plugin to use as described above
	LXEBrDHCPRange             string // Which DHCP Range to configure to lxebr0
}
