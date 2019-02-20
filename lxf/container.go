package lxf

import (
	"crypto/md5" // nolint: gosec #nosec (no sensitive data)
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxe/lxf/device"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	cfgLogPath            = "user.log_path"
	cfgSecurityPrivileged = "security.privileged"
	cfgSecurityNesting    = "security.nesting"
	cfgVolatileBaseImage  = "volatile.base_image"
	cfgStartedAt          = "user.started_at"
	cfgFinishedAt         = "user.finished_at"
	cfgCloudInitUserData  = "user.user-data"
	cfgCloudInitMetaData  = "user.meta-data"
	cfgEnvironmentPrefix  = "environment"
	cfgAutoStartOnBoot    = "boot.autostart"
	cfgVolatile           = "volatile"
)

var (
	containerConfigStore = NewConfigStore().WithReserved(
		cfgSchema,
		cfgLogPath,
		cfgIsCRI,
		cfgSecurityPrivileged,
		cfgSecurityNesting,
		cfgMetaName,
		cfgMetaAttempt,
		cfgCreatedAt,
		cfgStartedAt,
		cfgCloudInitUserData,
		cfgCloudInitMetaData,
		cfgCloudInitNetworkConfig,
	).WithReservedPrefixes(
		cfgLabels,
		cfgAnnotations,
		cfgVolatile,
	)
)

// Container represents a LXD container including CRI specific configuration
type Container struct {
	// LXDObject inherits common CRI fields
	LXDObject
	// Profiles of the container. First entry is always the sandbox profile
	// The default profile is always excluded and managed according to the settings automatically
	Profiles []string
	// Image defines the image to be used, can be the hash or local alias
	Image string
	// Privileged defines if the container is run privileged
	Privileged bool
	// Environment specifies to the container exported environment variables
	Environment map[string]string

	// CRIObject inherits common CRI fields
	CRIObject
	// Metadata contains user defined data
	Metadata ContainerMetadata
	// CreatedAt is when the container was created
	CreatedAt time.Time
	// StartedAt is when the container was started
	StartedAt time.Time
	// FinishedAt is when the container was exited
	FinishedAt time.Time
	// State contains the current state of this container
	State *ContainerState
	// LogPath TODO, not really in use, can it be removed?
	LogPath                string
	CloudInitUserData      string
	CloudInitMetaData      string
	CloudInitNetworkConfig string

	// sandbox is the parent sandbox of this container
	sandbox *Sandbox
}

// ContainerState holds information about the container state
type ContainerState struct {
	// Pid of the container
	// +readonly
	Pid int64
	// State of the current container
	// +readonly
	Name ContainerStateName
	// Stats usage of the current container
	// +readonly
	Stats ContainerStats
	// Network represents the network information section of a LXD container's state
	// +readonly
	Network map[string]api.ContainerStateNetwork
}

// ContainerStateName represents the state name of the container
type ContainerStateName string

const (
	// ContainerStateCreated it's there but not started yet
	ContainerStateCreated ContainerStateName = "created"
	// ContainerStateRunning it's there and running
	ContainerStateRunning ContainerStateName = "running"
	// ContainerStateExited it's there but terminated
	ContainerStateExited ContainerStateName = "exited"
	// ContainerStateUnknown it's there but we don't know what it's doing
	ContainerStateUnknown ContainerStateName = "unknown"
)

func (s ContainerStateName) String() string {
	return string(s)
}

// ContainerStats relevant for cri
type ContainerStats struct {
	MemoryUsage     uint64
	CPUUsage        uint64
	FilesystemUsage uint64
}

// ContainerMetadata has the metadata neede by a container
type ContainerMetadata struct {
	Name    string
	Attempt uint32
}

// Sandbox looks up the parent sandbox
// Implemented as lazy loading, and returns same result if already looked up
// Not thread safe! But it's expected the pointers stay in the same routine
func (c *Container) Sandbox() (*Sandbox, error) {
	var err error
	if c.sandbox == nil {
		c.sandbox, err = c.getSandbox()
		if err != nil {
			return nil, err
		}
	}
	return c.sandbox, nil
}

func (c *Container) getSandbox() (*Sandbox, error) {
	if len(c.Profiles) > 0 {
		profile := c.Profiles[0]
		if profile == "default" {
			return nil, fmt.Errorf("Container '%v' cannot have default as the only profile", c.ID)
		}
		sandbox, err := c.client.GetSandbox(profile)
		if err != nil {
			return nil, err
		}
		return sandbox, nil
	}
	return nil, fmt.Errorf("Container '%v' must have at least one profile", c.ID)
}

func (c *Container) getState() (*ContainerState, error) {
	cs := &ContainerState{}

	state, _, err := c.client.server.GetContainerState(c.ID)
	if err != nil {
		return nil, err
	}

	cs.Pid = state.Pid
	cs.Network = state.Network
	cs.Stats = ContainerStats{
		CPUUsage:    uint64(state.CPU.Usage),
		MemoryUsage: uint64(state.Memory.Usage),
		// TODO: "root" is a magic constant
		FilesystemUsage: uint64(state.Disk["root"].Usage),
	}

	// Map status code of LXD to CRI
	switch state.StatusCode {
	case api.Running:
		cs.Name = ContainerStateRunning
	case api.Stopped, api.Aborting, api.Stopping:
		// we have to differentiate between stopped and created. If "user.state" exists, then it must be
		// created, otherwise its exited
		if state, has := c.Config[cfgState]; has && state == string(ContainerStateCreated) {
			cs.Name = ContainerStateCreated
		} else {
			cs.Name = ContainerStateExited
		}
	default:
		cs.Name = ContainerStateUnknown
	}

	return cs, nil
}

// Apply will save the changes of an existing container
func (c *Container) Apply() error {
	s, err := c.Sandbox()
	if err != nil {
		return err
	}
	switch s.NetworkConfig.Mode {
	case NetworkHost:
		if !c.Privileged {
			return fmt.Errorf("`podSpec.hostNetwork: true` can only be used together with `containerSpec.securityContext.privileged: true`")
		}
	default:
		// do nothing
	}

	return c.apply()
}

// Start the container
func (c *Container) Start() error {
	err := c.client.opwait.StartContainer(c.ID)
	if err != nil {
		return err
	}

	// delete created mark if exists, so next stopping state can be exited
	delete(c.Config, cfgState)
	// and add started time
	c.StartedAt = time.Now()

	return c.apply()
}

// Stop will try to stop the container, returns nil when container is already deleted or
// got deleted in the meantime, otherwise it will return an error.
// If it's not deleted after 30 seconds it will return an error (might be way to low).
func (c *Container) Stop(timeout int) error {
	err := c.client.opwait.StopContainer(c.ID, timeout, 2)
	if err != nil {
		return err
	}

	c.FinishedAt = time.Now()

	return c.apply()
}

// Delete the container
func (c *Container) Delete() error {
	return c.client.opwait.DeleteContainer(c.ID)
}

// apply saves the changes to LXD
func (c *Container) apply() error {
	// TODO: can't this be done easier?
	imageID, err := c.client.parseImage(c.Image)
	if err != nil {
		return err
	}
	hash, found, err := imageID.Hash(c.client)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("image '%v' not found on local remote", c.Image)
	}

	config := makeContainerConfig(c)
	devices, err := makeContainerDevices(c)
	if err != nil {
		return err
	}

	for key, val := range c.Config {
		if containerConfigStore.IsReserved(key) {
			logger.Warnf("config key '%v' is reserved and can not be used", key)
		} else {
			config[key] = val
		}
	}

	config[cfgSchema] = SchemaVersionContainer
	contPut := api.ContainerPut{
		// TODO: magic constant
		Profiles: []string{c.Profiles[0], "default"},
		Config:   config,
		Devices:  devices,
	}

	if c.ID == "" {
		// container has to be created
		c.ID = c.CreateID()
		c.Config[cfgState] = ContainerStateCreated.String()
		c.CreatedAt = time.Now()
		return c.client.opwait.CreateContainer(api.ContainersPost{
			Name:         c.ID,
			ContainerPut: contPut,
			Source: api.ContainerSource{
				Fingerprint: hash,
				Type:        "image",
			},
		})
	}
	// else container has to be updated
	if c.ETag == "" {
		return fmt.Errorf("Update container not allowed with empty ETag")
	}
	return c.client.opwait.UpdateContainer(c.ID, contPut, c.ETag)
}

func makeContainerConfig(c *Container) map[string]string {
	config := map[string]string{}

	// write labels
	for key, val := range c.Labels {
		config[cfgLabels+"."+key] = val
	}
	// and annotations
	for key, val := range c.Annotations {
		config[cfgAnnotations+"."+key] = val
	}

	config[cfgCreatedAt] = strconv.FormatInt(c.CreatedAt.UnixNano(), 10)
	config[cfgStartedAt] = strconv.FormatInt(c.StartedAt.UnixNano(), 10)
	config[cfgFinishedAt] = strconv.FormatInt(c.FinishedAt.UnixNano(), 10)
	config[cfgSecurityPrivileged] = strconv.FormatBool(c.Privileged)
	config[cfgLogPath] = c.LogPath
	config[cfgIsCRI] = strconv.FormatBool(true)
	config[cfgMetaName] = c.Metadata.Name
	config[cfgMetaAttempt] = strconv.FormatUint(uint64(c.Metadata.Attempt), 10)
	config[cfgVolatileBaseImage] = c.Image

	for k, v := range c.Environment {
		config[cfgEnvironmentPrefix+"."+k] = v
	}

	// and meta-data & cloud-init
	// fields should not exist when there's nothing
	if c.CloudInitMetaData != "" {
		config[cfgCloudInitMetaData] = c.CloudInitMetaData
	}
	if c.CloudInitUserData != "" {
		config[cfgCloudInitUserData] = c.CloudInitUserData
	}
	if c.CloudInitNetworkConfig != "" {
		config[cfgCloudInitNetworkConfig] = c.CloudInitNetworkConfig
	}

	return config
}

func makeContainerDevices(c *Container) (map[string]map[string]string, error) {
	devices := map[string]map[string]string{}
	err := device.AddDisksToMap(devices, c.Disks...)
	if err != nil {
		return devices, err
	}
	err = device.AddProxiesToMap(devices, c.Proxies...)
	if err != nil {
		return devices, err
	}
	err = device.AddBlocksToMap(devices, c.Blocks...)
	if err != nil {
		return devices, err
	}
	return devices, device.AddNicsToMap(devices, c.Nics...)
}

// toContainer will convert an lxd container to lxf format
func toContainer(ct *api.Container, ETag string) (*Container, error) {
	attempts, err := strconv.ParseUint(ct.Config[cfgMetaAttempt], 10, 32)
	if err != nil {
		return nil, err
	}
	privileged, err := strconv.ParseBool(ct.Config[cfgSecurityPrivileged])
	if err != nil {
		return nil, err
	}
	createdAt, err := strconv.ParseInt(ct.Config[cfgCreatedAt], 10, 64)
	if err != nil {
		return nil, err
	}
	startedAt, err := strconv.ParseInt(ct.Config[cfgStartedAt], 10, 64)
	if err != nil {
		return nil, err
	}
	finishedAt, err := strconv.ParseInt(ct.Config[cfgFinishedAt], 10, 64)
	if err != nil {
		return nil, err
	}

	c := &Container{}
	c.ID = ct.Name
	c.ETag = ETag
	c.Metadata = ContainerMetadata{
		Name:    ct.Config[cfgMetaName],
		Attempt: uint32(attempts),
	}
	c.LogPath = ct.Config[cfgLogPath]
	c.Image = ct.Config[cfgVolatileBaseImage]
	c.Annotations = containerConfigStore.StripedPrefixMap(ct.Config, cfgAnnotations)
	c.Labels = containerConfigStore.StripedPrefixMap(ct.Config, cfgLabels)
	c.Config = containerConfigStore.UnreservedMap(ct.Config)

	c.CreatedAt = time.Unix(0, createdAt)
	c.StartedAt = time.Unix(0, startedAt)
	c.FinishedAt = time.Unix(0, finishedAt)

	c.Environment = extractEnvVars(ct.Config)
	c.Privileged = privileged
	c.CloudInitUserData = ct.Config[cfgCloudInitUserData]
	c.CloudInitMetaData = ct.Config[cfgCloudInitMetaData]
	c.CloudInitNetworkConfig = ct.Config[cfgCloudInitNetworkConfig]

	c.Proxies, err = device.GetProxiesFromMap(ct.Devices)
	if err != nil {
		return nil, err
	}
	c.Disks, err = device.GetDisksFromMap(ct.Devices)
	if err != nil {
		return nil, err
	}
	c.Blocks, err = device.GetBlocksFromMap(ct.Devices)
	if err != nil {
		return nil, err
	}
	c.Nics, err = device.GetNicsFromMap(ct.Devices)
	if err != nil {
		return nil, err
	}

	for _, v := range ct.Profiles {
		// TODO: magic constant
		if v != "default" {
			c.Profiles = append(c.Profiles, v)
		}
	}
	if len(c.Profiles) == 0 {
		return nil, fmt.Errorf("Container '%v' has no sandbox", c.ID)
	}

	c.State, err = c.getState()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// extractEnvVars extracts all the config options that start with "environment."
// and returns the environment variables + values
func extractEnvVars(config map[string]string) map[string]string {
	envVars := make(map[string]string)
	for k, v := range config {
		if strings.HasPrefix(k, cfgEnvironmentPrefix+".") {
			varName := strings.TrimLeft(k, cfgEnvironmentPrefix+".")
			varValue := v
			envVars[varName] = varValue
		}
	}
	return envVars
}

// CreateID creates a unique container id
func (c *Container) CreateID() string {
	bin := md5.Sum([]byte(uuid.NewUUID())) // nolint: gosec #nosec
	return string(c.Metadata.Name[0]) + b32lowerEncoder.EncodeToString(bin[:])[:15]
}

// GetInetAddress returns the IPv4 address of the first matching interface in the parameter list
// empty string if nothing was found
func (c *Container) GetInetAddress(ifs []string) string {
	for _, i := range ifs {
		if netif, ok := c.State.Network[i]; ok {
			for _, addr := range netif.Addresses {
				if addr.Family == "inet" {
					return addr.Address
				}
			}
		}
	}
	return ""
}
