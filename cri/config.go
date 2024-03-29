package cri

// Domain of the daemon
const Domain = "lxe"

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
	// LXDBridgeName is the name of the bridge to create and use
	LXDBridgeName string
	// LXDBridgeDHCPRange to configure for bridge if NetworkPlugin is default
	LXDBridgeDHCPRange string
	// LXEStreamingBindAddr contains the listen address for the streaming server
	LXEStreamingBindAddr string
	// LXEStreamingBaseURL is the base address for constructing streaming URLs
	LXEStreamingBaseURL string
	// LXEHostnetworkFile file path to use for lxc's raw.include
	LXEHostnetworkFile string
	// Which LXENetworkPlugin to use
	LXENetworkPlugin string
	// CNIConfDir is the path where the cni configuration files are
	CNIConfDir string
	// CNIBinDir is the path where the cni plugins are
	CNIBinDir string
	// CNIOutputWriter is the writer for CNI call outputs
	CNIOutputTarget string
	// CNIOutputFile is the path to a file
	CNIOutputFile string
	// CRITest mode that enables rewriting of requested images as cri-tools only refer to OCI-images
	CRITest bool
}
