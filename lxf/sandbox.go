package lxf

import (
	"crypto/md5" // nolint: gosec #nosec (no sensitive data)
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxe/lxf/device"
	"github.com/lxc/lxe/network"
)

const (
	cfgHostname     = "user.host_name"
	cfgLogDirectory = "user.log_directory"
	cfgCreatedAt    = "user.created_at"
	cfgState        = "user.state"

	cfgNetworkConfigNameservers = "user.networkconfig.nameservers"
	cfgNetworkConfigSearches    = "user.networkconfig.searches"
	cfgNetworkConfigMode        = "user.networkconfig.mode"
	cfgNetworkConfigModeData    = "user.networkconfig.modedata"
	cfgCloudInitNetworkConfig   = "user.network-config" // write-only field
	cfgCloudInitVendorData      = "user.vendor-data"    // write-only field

	cfgRawLXC = "raw.lxc"

	// CfgRawLXCInclude is the raw lxc config field name for the include file
	CfgRawLXCInclude = "lxc.include"
	// CfgRawLXCNamespaces is the lxc config field name for what namespaces to keep
	CfgRawLXCNamespaces = "lxc.namespace.keep"
	// CfgRawLXCKernelModules is the lxc config field name for what kernel modules to load
	CfgRawLXCKernelModules = "linux.kernel_modules"
	// CfgRawLXCMounts is the raw key to add lxc mounts. Useful for mounting proc for example,
	// for nested containers
	CfgRawLXCMounts = "lxc.mount.auto"
)

var (
	sandboxConfigStore = NewConfigStore().WithReserved(cfgSchema, cfgHostname, cfgLogDirectory, cfgCreatedAt,
		cfgState, cfgIsCRI, cfgMetaAttempt, cfgMetaName, cfgMetaNamespace, cfgMetaUID, cfgCloudInitNetworkConfig,
		cfgCloudInitVendorData, cfgNetworkConfigModeData, cfgRawLXC).
		WithReservedPrefixes(cfgLabels, cfgAnnotations)
)

// RawLXCOption contains a single option plus its value
type RawLXCOption struct {
	Option string
	Value  string
}

// Sandbox is an abstraction of a CRI PodSandbox saved as a LXD profile
type Sandbox struct {
	CRIObject
	Hostname      string
	LogDirectory  string
	Metadata      SandboxMetadata
	NetworkConfig NetworkConfig
	// State is readonly
	State     SandboxState
	CreatedAt time.Time
	// RawLXCOptions are additional raw.lxc fields
	RawLXCOptions []RawLXCOption
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

func (s SandboxState) toString() string {
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
	Version int                  `json:"version"`
	Config  []NetworkConfigEntry `json:"config"`
}

// NetworkConfigEntry is an entry in the v1 network config, currently limited to nameserver
type NetworkConfigEntry struct {
	Type    string   `json:"type"`
	Address []string `json:"address"`
	Search  []string `json:"search"`
}

// CreateSandbox will create the provided sandbox and put it into state ready
func (l *LXF) CreateSandbox(s *Sandbox) error {
	s.State = SandboxReady
	s.CreatedAt = time.Now()

	// Apply defined network mode
	switch s.NetworkConfig.Mode {
	case NetworkBridged:
		s.Nics = append(s.Nics, device.Nic{
			Name:    network.DefaultInterface,
			NicType: "bridged",
			Parent:  s.NetworkConfig.ModeData["bridge"],
		})
	case NetworkHost:
		fallthrough
	case NetworkCNI:
		fallthrough
	case NetworkNone:
		fallthrough
	default:
		// do nothing
	}

	return l.saveSandbox(s)
}

// StopSandbox will find a sandbox by id and set it's state to "not ready".
func (l *LXF) StopSandbox(id string) error {
	p, _, err := l.server.GetProfile(id)
	if err != nil {
		return err
	}

	p.Config[cfgState] = string(SandboxNotReady)
	err = l.server.UpdateProfile(id, p.Writable(), "")
	return err
}

// GetSandbox will find a sandbox by id and return it.
func (l *LXF) GetSandbox(id string) (*Sandbox, error) {
	p, _, err := l.server.GetProfile(id)
	if err != nil {
		return nil, err
	}

	if !IsCRI(p) {
		return nil, fmt.Errorf(ErrorNotFound)
	}
	return l.toSandbox(p)
}

// DeleteSandbox will delete the given sandbox
func (l *LXF) DeleteSandbox(name string) error {
	return l.server.DeleteProfile(name)
}

// ListSandboxes will return a list with all the available sandboxes
func (l *LXF) ListSandboxes() ([]*Sandbox, error) { // nolint:dupl
	ps, err := l.server.GetProfiles()
	if err != nil {
		return nil, err
	}

	sandboxes := []*Sandbox{}
	for _, p := range ps {
		if !IsCRI(p) {
			continue
		}
		sb, err2 := l.toSandbox(&p)
		if err2 != nil {
			return nil, err2
		}
		sandboxes = append(sandboxes, sb)
	}

	return sandboxes, nil
}

// saveSandbox will take a sandbox and saves it as a profile
// if the profile already exists it will be created, otherwise updated
func (l *LXF) saveSandbox(s *Sandbox) error {
	config := map[string]string{
		cfgState:                    s.State.toString(),
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
	highestSearch := ""
	if len(s.NetworkConfig.Nameservers) > 0 &&
		len(s.NetworkConfig.Searches) > 0 {
		data := networkConfigData{
			Version: 1,
			Config: []NetworkConfigEntry{
				NetworkConfigEntry{
					Type:    "nameserver",
					Address: s.NetworkConfig.Nameservers,
					Search:  s.NetworkConfig.Searches,
				},
			},
		}

		yml, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		config[cfgCloudInitNetworkConfig] = string(yml)
		highestSearch = "." + s.NetworkConfig.Searches[0]
	}

	// write cloud-init vendor data if we have hostname and search
	if s.Hostname != "" && highestSearch != "" {
		config[cfgCloudInitVendorData] = fmt.Sprintf(`#cloud-config
hostname: %s
fqdn: %s
manage_etc_hosts: true`, s.Hostname, s.Hostname+highestSearch)
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

	// Apply raw lxc options
	re := regexp.MustCompile(`\r?\n`)
	for _, lxcOption := range s.RawLXCOptions {
		// TODO: these variables are probably not sanitized enough! So far we don't allow newlines
		if !re.Match([]byte(lxcOption.Option)) && !re.Match([]byte(lxcOption.Value)) {
			config[cfgRawLXC] += fmt.Sprintf("%s = %s\n", lxcOption.Option, lxcOption.Value)
		}
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
	return l.server.UpdateProfile(s.ID, profile, "")
}

// toSandbox will take a profile and convert it to a sandbox.
func (l *LXF) toSandbox(p *api.Profile) (*Sandbox, error) {
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
	s.RawLXCOptions = make([]RawLXCOption, 0)
	s.State = getSandboxState(p.Config[cfgState])
	s.CreatedAt = time.Unix(0, createdAt)

	err = yaml.Unmarshal([]byte(p.Config[cfgNetworkConfigModeData]), &s.NetworkConfig.ModeData)
	if err != nil {
		return nil, err
	}
	if len(s.NetworkConfig.ModeData) == 0 {
		s.NetworkConfig.ModeData = make(map[string]string)
	}

	// Hint: cloud-init network config & vendor-data are write-only so not readed

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

	// get raw lxc options
	for _, line := range strings.Split(strings.TrimSuffix(p.Config[cfgRawLXC], "\n"), "\n") {
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		s.RawLXCOptions = append(s.RawLXCOptions, RawLXCOption{
			Option: strings.TrimSpace(parts[0]),
			Value:  strings.TrimSpace(parts[1]),
		})
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

// CreateID creates the unique sandbox id based on Kubernetes sandbox values
// This is currently not expected to be a long term stable hashing for these informations
func (s *Sandbox) CreateID() string {
	var parts []string
	parts = append(parts, "k8s")
	parts = append(parts, s.Metadata.Name)
	parts = append(parts, s.Metadata.Namespace)
	parts = append(parts, strconv.FormatUint(uint64(s.Metadata.Attempt), 10))
	parts = append(parts, s.Metadata.UID)
	name := strings.Join(parts, "-")

	bin := md5.Sum([]byte(name)) // nolint: gosec #nosec
	return string(s.Metadata.Name[0]) + b32lowerEncoder.EncodeToString(bin[:])[:15]
}
