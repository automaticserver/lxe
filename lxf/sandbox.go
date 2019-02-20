package lxf

import (
	"crypto/md5" // nolint: gosec #nosec (no sensitive data)
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxe/lxf/device"
	"github.com/lxc/lxe/network"
	utilNet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	cfgHostname     = "user.host_name"
	cfgLogDirectory = "user.log_directory"
	cfgCreatedAt    = "user.created_at"

	cfgNetworkConfigNameservers = "user.networkconfig.nameservers"
	cfgNetworkConfigSearches    = "user.networkconfig.searches"
	cfgNetworkConfigMode        = "user.networkconfig.mode"
	cfgNetworkConfigModeData    = "user.networkconfig.modedata"
	cfgCloudInitNetworkConfig   = "user.network-config" // write-only field
	cfgCloudInitVendorData      = "user.vendor-data"    // write-only field
)

var (
	sandboxConfigStore = NewConfigStore().WithReserved(cfgSchema, cfgHostname, cfgLogDirectory, cfgCreatedAt,
		cfgState, cfgIsCRI, cfgMetaAttempt, cfgMetaName, cfgMetaNamespace, cfgMetaUID, cfgCloudInitNetworkConfig,
		cfgCloudInitVendorData, cfgNetworkConfigModeData).
		WithReservedPrefixes(cfgLabels, cfgAnnotations)
)

// RawLXCOption contains a single option plus its value
type RawLXCOption struct {
	Option string
	Value  string
}

// Sandbox is an abstraction of a CRI PodSandbox saved as a LXD profile
type Sandbox struct {
	LXDObject
	CRIObject
	Hostname      string
	LogDirectory  string
	Metadata      SandboxMetadata
	NetworkConfig NetworkConfig
	// State is readonly
	State     SandboxState
	CreatedAt time.Time
	// UsedBy contains the names of the containers using this profile
	// It is read only.
	UsedBy []string
}

// SandboxMetadata contains common metadata values
type SandboxMetadata struct {
	Attempt   uint32
	Name      string
	Namespace string
	UID       string
}

// NetworkConfig contains the network config
// searches and nameservers must not be empty to be valid
type NetworkConfig struct {
	Nameservers []string
	Searches    []string
	// Mode describes the type of networking
	Mode NetworkMode
	// ModeData allows Mode-specific data to be persisted
	ModeData map[string]string
}

// NetworkMode defines the type of the container network
type NetworkMode string

// These are valid network modes. NetworkHost means the container to share the host's network namespace
// NetworkCNI means the CNI handles the interface, NetworkBridged means the container gets a interface
// from a predefined bridge, NetworkNone is used when the requested mode can't be used
const (
	NetworkHost    NetworkMode = "node"
	NetworkCNI     NetworkMode = "cni"
	NetworkBridged NetworkMode = "bridged"
	NetworkNone    NetworkMode = "none"
)

func (s NetworkMode) toString() string {
	return string(s)
}

func getNetworkMode(str string) NetworkMode {
	for _, v := range []NetworkMode{NetworkHost, NetworkCNI, NetworkBridged, NetworkNone} {
		if str == string(v) {
			return v
		}
	}
	return NetworkNone
}

// SandboxState defines the state of the sandbox, default is SandboxReady
type SandboxState string

// These are valid sandbox statuses. SandboxReady means a resource is in the condition.
// SandboxNotReady means a resource is not in the condition.
const (
	SandboxNotReady SandboxState = "notready"
	SandboxReady    SandboxState = "ready"
)

func (s SandboxState) String() string {
	return string(s)
}

func getSandboxState(str string) SandboxState {
	if str == string(SandboxNotReady) {
		return SandboxNotReady
	}
	return SandboxReady
}

// networkConfigData is used as root element to serialize to cloud config
type networkConfigData struct {
	Version int           `json:"version"`
	Config  []interface{} `json:"config"`
}

// NetworkConfigEntry is an entry in the v1 network config
type NetworkConfigEntry struct {
	Type string `json:"type"`
}

// NetworkConfigEntryNameserver is a nameserver entry
type NetworkConfigEntryNameserver struct {
	NetworkConfigEntry
	Address []string `json:"address"`
	Search  []string `json:"search"`
}

// NetworkConfigEntryPhysical is a nameserver entry
type NetworkConfigEntryPhysical struct {
	NetworkConfigEntry
	Name    string                             `json:"name"`
	Subnets []NetworkConfigEntryPhysicalSubnet `json:"subnets"`
}

// NetworkConfigEntryPhysicalSubnet is a subnet entry in the v1 network config of a physical device
type NetworkConfigEntryPhysicalSubnet struct {
	Type string `json:"type"`
}

// CreateSandbox will create the provided sandbox and put it into state ready
func (l *Client) CreateSandbox(s *Sandbox) error {
	s.State = SandboxReady
	s.CreatedAt = time.Now()

	// Apply defined network mode
	switch s.NetworkConfig.Mode {
	case NetworkBridged:
		s.Nics = append(s.Nics, device.Nic{
			Name:        network.DefaultInterface,
			NicType:     "bridged",
			Parent:      s.NetworkConfig.ModeData["bridge"],
			IPv4Address: s.NetworkConfig.ModeData["interface-address"],
		})
	default:
		// do nothing
	}

	return l.saveSandbox(s)
}

// StopSandbox will find a sandbox by id and set it's state to "not ready".
func (l *Client) StopSandbox(id string) error {
	p, ETag, err := l.server.GetProfile(id)
	if err != nil {
		return err
	}

	p.Config[cfgState] = SandboxNotReady.String()
	err = l.server.UpdateProfile(id, p.Writable(), ETag)
	return err
}

// DeleteSandbox will delete the given sandbox
func (l *Client) DeleteSandbox(name string) error {
	return l.server.DeleteProfile(name)
}

// saveSandbox will take a sandbox and saves it as a profile
// if the profile already exists it will be created, otherwise updated
func (l *Client) saveSandbox(s *Sandbox) error {
	config := map[string]string{
		cfgState:                    s.State.String(),
		cfgIsCRI:                    strconv.FormatBool(true),
		cfgCreatedAt:                strconv.FormatInt(s.CreatedAt.UnixNano(), 10),
		cfgMetaAttempt:              strconv.FormatUint(uint64(s.Metadata.Attempt), 10),
		cfgMetaName:                 s.Metadata.Name,
		cfgMetaNamespace:            s.Metadata.Namespace,
		cfgMetaUID:                  s.Metadata.UID,
		cfgHostname:                 s.Hostname,
		cfgLogDirectory:             s.LogDirectory,
		cfgNetworkConfigNameservers: strings.Join(s.NetworkConfig.Nameservers, ","),
		cfgNetworkConfigSearches:    strings.Join(s.NetworkConfig.Searches, ","),
		cfgNetworkConfigMode:        s.NetworkConfig.Mode.toString(),
	}

	// write NetworkConfigData as yaml
	yml, err := yaml.Marshal(s.NetworkConfig.ModeData)
	if err != nil {
		return err
	}
	config[cfgNetworkConfigModeData] = string(yml)

	// write labels
	for key, val := range s.Labels {
		config[cfgLabels+"."+key] = val
	}
	// and annotations
	for key, val := range s.Annotations {
		config[cfgAnnotations+"."+key] = val
	}
	// and config keys
	for key, val := range s.Config {
		if sandboxConfigStore.IsReserved(key) {
			logger.Warnf("config key '%v' is reserved and can not be used", key)
		} else {
			config[key] = val
		}
	}

	// write cloud-init network config
	if len(s.NetworkConfig.Nameservers) > 0 &&
		len(s.NetworkConfig.Searches) > 0 {
		data := networkConfigData{
			Version: 1,
			Config: []interface{}{
				NetworkConfigEntryNameserver{
					NetworkConfigEntry: NetworkConfigEntry{
						Type: "nameserver",
					},
					Address: s.NetworkConfig.Nameservers,
					Search:  s.NetworkConfig.Searches,
				},
			},
		}
		if s.NetworkConfig.Mode == NetworkBridged && s.NetworkConfig.ModeData["interface-address"] != "" {
			data.Config = append(data.Config, NetworkConfigEntryPhysical{
				NetworkConfigEntry: NetworkConfigEntry{
					Type: "physical",
				},
				Name: network.DefaultInterface,
				Subnets: []NetworkConfigEntryPhysicalSubnet{
					NetworkConfigEntryPhysicalSubnet{
						Type: "dhcp",
					},
				},
			})
		}

		yml, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		config[cfgCloudInitNetworkConfig] = string(yml)
	}

	// write cloud-init vendor data if we have hostname and search
	if s.Hostname != "" {
		config[cfgCloudInitVendorData] = fmt.Sprintf(`#cloud-config
hostname: %s
manage_etc_hosts: true`, s.Hostname)
	}

	devices := map[string]map[string]string{}
	err = device.AddProxiesToMap(devices, s.Proxies...)
	if err != nil {
		return err
	}
	err = device.AddDisksToMap(devices, s.Disks...)
	if err != nil {
		return err
	}
	err = device.AddBlocksToMap(devices, s.Blocks...)
	if err != nil {
		return err
	}
	err = device.AddNicsToMap(devices, s.Nics...)
	if err != nil {
		return err
	}

	config[cfgSchema] = SchemaVersionProfile
	profile := api.ProfilePut{
		Config:  config,
		Devices: devices,
	}

	if s.ID == "" { // profile has to be created
		s.ID = s.CreateID()
		return l.server.CreateProfile(api.ProfilesPost{
			Name:       s.ID,
			ProfilePut: profile,
		})
	}
	// else profile has to be updated
	if s.ETag == "" {
		return fmt.Errorf("Update profile not allowed with empty ETag")
	}
	return l.server.UpdateProfile(s.ID, profile, s.ETag)
}

// toSandbox will take a profile and convert it to a sandbox.
func (l *Client) toSandbox(p *api.Profile, ETag string) (*Sandbox, error) {
	attempts, err := strconv.ParseUint(p.Config[cfgMetaAttempt], 10, 32)
	if err != nil {
		return nil, err
	}
	createdAt, err := strconv.ParseInt(p.Config[cfgCreatedAt], 10, 64)
	if err != nil {
		return nil, err
	}

	s := &Sandbox{}
	s.ID = p.Name
	s.ETag = ETag
	s.Hostname = p.Config[cfgHostname]
	s.LogDirectory = p.Config[cfgLogDirectory]
	s.Metadata = SandboxMetadata{
		Attempt:   uint32(attempts),
		Name:      p.Config[cfgMetaName],
		Namespace: p.Config[cfgMetaNamespace],
		UID:       p.Config[cfgMetaUID],
	}
	s.NetworkConfig = NetworkConfig{
		Nameservers: strings.Split(p.Config[cfgNetworkConfigNameservers], ","),
		Searches:    strings.Split(p.Config[cfgNetworkConfigSearches], ","),
		Mode:        getNetworkMode(p.Config[cfgNetworkConfigMode]),
		// ModeData:    make(map[string]string),
	}
	s.Labels = sandboxConfigStore.StripedPrefixMap(p.Config, cfgLabels)
	s.Annotations = sandboxConfigStore.StripedPrefixMap(p.Config, cfgAnnotations)
	s.Config = sandboxConfigStore.UnreservedMap(p.Config)
	s.State = getSandboxState(p.Config[cfgState])
	s.CreatedAt = time.Unix(0, createdAt)

	err = yaml.Unmarshal([]byte(p.Config[cfgNetworkConfigModeData]), &s.NetworkConfig.ModeData)
	if err != nil {
		return nil, err
	}
	if len(s.NetworkConfig.ModeData) == 0 {
		s.NetworkConfig.ModeData = make(map[string]string)
	}

	// Hint: cloud-init network config & vendor-data are write-only so not read

	// get devices
	s.Proxies, err = device.GetProxiesFromMap(p.Devices)
	if err != nil {
		return nil, err
	}
	s.Disks, err = device.GetDisksFromMap(p.Devices)
	if err != nil {
		return nil, err
	}
	s.Blocks, err = device.GetBlocksFromMap(p.Devices)
	if err != nil {
		return nil, err
	}
	s.Nics, err = device.GetNicsFromMap(p.Devices)
	if err != nil {
		return nil, err
	}

	// get containers using this sandbox
	for _, name := range p.UsedBy {
		name = strings.TrimPrefix(name, "/1.0/containers/")
		name = strings.TrimSuffix(name, "?project=default")
		if strings.Contains(name, shared.SnapshotDelimiter) {
			// this is a snapshot so dont parse this entry
			continue
		}
		c, _, err := l.server.GetContainer(name)
		if err != nil {
			return nil, err
		}
		s.UsedBy = append(s.UsedBy, c.Name)
	}

	return s, nil
}

// GetSandboxIP returns the ip address of the sandbox
func (l *Client) GetSandboxIP(s *Sandbox) (string, error) {
	switch s.NetworkConfig.Mode {
	case NetworkHost:
		ip, err := utilNet.ChooseHostInterface()
		if err != nil {
			return "", fmt.Errorf("looking up host interface failed: %v", err)
		}
		return ip.String(), nil
	case NetworkNone:
		return "", nil
	case NetworkBridged:
		// if statically assigned ip exists, return that
		if s.NetworkConfig.ModeData["interface-address"] != "" {
			return s.NetworkConfig.ModeData["interface-address"], nil
		}
		fallthrough
	case NetworkCNI:
		fallthrough
	default:
		for _, cntName := range s.UsedBy {
			c, err := l.GetContainer(cntName)
			if err != nil {
				if err.Error() == ErrorLXDNotFound {
					// TODO: we might get a Container Not Found error (once implemented)
					continue
				}
				return "", fmt.Errorf("looking up container failed: %v", err)
			}

			// ignore any non-running containers
			if c.State.Name != ContainerStateRunning {
				continue
			}

			// get the ipv4 address of eth0
			ip := c.GetInetAddress([]string{network.DefaultInterface})
			if ip != "" {
				return ip, nil
			}
		}
	}
	return "", nil
}

// CreateID creates a unique profile id
func (s *Sandbox) CreateID() string {
	bin := md5.Sum([]byte(uuid.NewUUID())) // nolint: gosec #nosec
	return string(s.Metadata.Name[0]) + b32lowerEncoder.EncodeToString(bin[:])[:15]
}
