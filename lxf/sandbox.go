package lxf

import (
	"crypto/md5" // nolint: gosec
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/network"
	"github.com/ghodss/yaml"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	utilNet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	// Default device name of the root disk when initializing lxd
	lxdInitDefaultDiskName = "root"
	// Default device name of the nic interface when initializing lxd
	lxdInitDefaultNicName = "eth0"

	cfgHostname                 = "user.host_name"
	cfgLogDirectory             = "user.log_directory"
	cfgCreatedAt                = "user.created_at"
	cfgNetworkConfig            = "user.networkconfig"
	cfgNetworkConfigNameservers = cfgNetworkConfig + ".nameservers"
	cfgNetworkConfigSearches    = cfgNetworkConfig + ".searches"
	cfgNetworkConfigMode        = cfgNetworkConfig + ".mode"
	cfgNetworkConfigModeData    = cfgNetworkConfig + ".modedata"
	cfgCloudInitNetworkConfig   = "user.network-config" // write-only field
	cfgCloudInitVendorData      = "user.vendor-data"    // write-only field
)

var (
	sandboxConfigStore = NewConfigStore().WithReserved(
		append([]string{
			cfgLogDirectory,
			cfgState,
			cfgHostname,
			cfgCloudInitNetworkConfig,
			cfgCloudInitVendorData,
			cfgNetworkConfigModeData,
		}, reservedConfigCRI...,
		)...,
	).WithReservedPrefixes(
		reservedConfigPrefixesCRI...,
	)
)

// Sandbox is an abstraction of a CRI PodSandbox saved as a LXD profile
type Sandbox struct {
	// LXDObject inherits common CRI fields
	LXDObject
	// UsedBy contains the names of the containers using this profile
	// It is read only.
	UsedBy []string

	// CRIObject inherits common CRI fields
	CRIObject
	// Metadata contains user defined data
	Metadata SandboxMetadata
	// Hostname to be set for containers if defined
	Hostname string
	// NetworkConfig to be applied for the sandbox and it's containers
	NetworkConfig NetworkConfig
	// State contains the current state of this sandbox
	State SandboxState
	// LogDirectory TODO, to be implemented?
	LogDirectory string

	// sandbox is the parent sandbox of this container
	containers []*Container
}

// SandboxState defines the state of the sandbox
type SandboxState string

// These are valid sandbox statuses. SandboxReady means a resource is in the condition.
// SandboxNotReady means a resource is not in the condition.
const (
	SandboxNotReady SandboxState = "notready"
	SandboxReady    SandboxState = "ready"
)

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

func (s NetworkMode) String() string {
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

// Containers looks up all assigned containers
// Implemented as lazy loading, and returns same result if already looked up
// Not thread safe! But it's expected the pointers stay in the same routine
func (s *Sandbox) Containers() ([]*Container, error) {
	var err error
	if s.containers == nil {
		s.containers, err = s.getContainers()
		if err != nil {
			return nil, err
		}
	}

	return s.containers, nil
}

func (s *Sandbox) getContainers() ([]*Container, error) {
	cl := []*Container{}

	for _, cntName := range s.UsedBy {
		c, err := s.client.GetContainer(cntName)
		if err != nil {
			return nil, err
		}

		cl = append(cl, c)
	}

	return cl, nil
}

// Apply will save the changes of a sandbox
func (s *Sandbox) Apply() error {
	// A new sandbox gets also some default values
	// except ID, which is generated inline in unexported method apply()
	if s.ID == "" {
		s.State = SandboxReady
		s.CreatedAt = time.Now()
	}

	// Apply defined network mode
	switch s.NetworkConfig.Mode { // nolint: gocritic
	case NetworkBridged:
		s.Devices.Upsert(&device.Nic{
			Name:        network.DefaultInterface,
			NicType:     "bridged",
			Parent:      s.NetworkConfig.ModeData["bridge"],
			IPv4Address: s.NetworkConfig.ModeData["interface-address"],
		})
	}

	// Always stop inheriting default eth0 device
	s.Devices.Upsert(&device.None{
		KeyName: lxdInitDefaultNicName,
	})

	return s.apply()
}

// Stop set the sandbox state to SandboxNotReady
func (s *Sandbox) Stop() error {
	s.State = SandboxNotReady

	return s.apply()
}

// Delete will delete the given sandbox, returns nil when sandbox is already deleted
func (s *Sandbox) Delete() error {
	err := s.client.server.DeleteProfile(s.ID)
	if err != nil {
		if err.Error() == ErrorLXDNotFound {
			return nil
		}

		return err
	}

	return nil
}

// apply saves the changes to LXD
func (s *Sandbox) apply() error {
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
		cfgNetworkConfigMode:        s.NetworkConfig.Mode.String(),
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

		if s.NetworkConfig.Mode == NetworkBridged {
			// because added later, if physical-type is empty, it is dhcp
			if s.NetworkConfig.ModeData["physical-type"] == "" {
				s.NetworkConfig.ModeData["physical-type"] = "dhcp"
			}

			entry := NetworkConfigEntryPhysical{
				NetworkConfigEntry: NetworkConfigEntry{
					Type: "physical",
				},
				Name: network.DefaultInterface,
				Subnets: []NetworkConfigEntryPhysicalSubnet{
					{
						Type: s.NetworkConfig.ModeData["physical-type"],
					},
				},
			}
			data.Config = append(data.Config, entry)
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

	devices := make(map[string]map[string]string)

	for _, d := range s.Devices {
		name, options := d.ToMap()
		devices[name] = options
	}

	config[cfgSchema] = SchemaVersionProfile
	profile := api.ProfilePut{
		Config:  config,
		Devices: devices,
	}

	if s.ID == "" { // profile has to be created
		s.ID = s.CreateID()

		return s.client.server.CreateProfile(api.ProfilesPost{
			Name:       s.ID,
			ProfilePut: profile,
		})
	}
	// else profile has to be updated
	if s.ETag == "" {
		return fmt.Errorf("update profile not allowed with empty ETag")
	}

	err = s.client.server.UpdateProfile(s.ID, profile, s.ETag)
	if err != nil {
		if err.Error() == ErrorLXDNotFound {
			return NewSandboxError(s.ID, err)
		}

		return err
	}

	return nil
}

// CreateID creates a unique profile id
func (s *Sandbox) CreateID() string {
	bin := md5.Sum([]byte(uuid.NewUUID())) // nolint: gosec
	return string(s.Metadata.Name[0]) + b32lowerEncoder.EncodeToString(bin[:])[:15]
}

// GetInetAddress returns the ip address of the sandbox
// empty string if nothing was found
func (s *Sandbox) GetInetAddress() string {
	switch s.NetworkConfig.Mode {
	case NetworkHost:
		ip, err := utilNet.ChooseHostInterface()
		if err != nil {
			// TODO: additional debug output
			return ""
		}

		return ip.String()
	case NetworkNone:
		return ""
	case NetworkBridged:
		// if statically assigned ip exists, return that
		if s.NetworkConfig.ModeData["interface-address"] != "" {
			return s.NetworkConfig.ModeData["interface-address"]
		}
		// otherwise look into containers (= fallthrough)
		fallthrough
	case NetworkCNI:
		fallthrough
	default:
		cl, err := s.Containers()
		if err != nil {
			// TODO: additional debug output
			return ""
		}

		for _, c := range cl {
			// ignore any non-running containers
			if c.StateName != ContainerStateRunning {
				continue
			}

			// get the ipv4 address of eth0
			ip := c.GetInetAddress([]string{network.DefaultInterface})
			if ip != "" {
				return ip
			}
		}
	}

	return ""
}
