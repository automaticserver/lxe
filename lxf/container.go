package lxf // import "github.com/automaticserver/lxe/lxf"

import (
	"crypto/md5" // nolint: gosec
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/automaticserver/lxe/shared"
	"github.com/lxc/lxd/shared/api"
	opencontainers "github.com/opencontainers/runtime-spec/specs-go"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	cfgLogPath              = "user.log_path"
	cfgSecurityPrivileged   = "security.privileged"
	cfgVolatileBaseImage    = cfgVolatile + ".base_image"
	cfgStartedAt            = "user.started_at"
	cfgFinishedAt           = "user.finished_at"
	cfgCloudInitUserData    = "user.user-data"
	cfgCloudInitMetaData    = "user.meta-data"
	cfgEnvironmentPrefix    = "environment"
	cfgResourcesPrefix      = "user.resources"
	cfgResourcesCPUPrefix   = cfgResourcesPrefix + ".cpu"
	cfgResourcesCPUShares   = cfgResourcesCPUPrefix + ".shares"
	cfgResourcesCPUQuota    = cfgResourcesCPUPrefix + ".quota"
	cfgResourcesCPUPeriod   = cfgResourcesCPUPrefix + ".period"
	cfgResourcesMemoryLimit = cfgResourcesPrefix + ".memory.limit"
	cfgLimitCPUAllowance    = "limits.cpu.allowance"
	cfgLimitMemory          = "limits.memory"
)

var (
	containerConfigStore = NewConfigStore().WithReserved(
		append([]string{
			cfgLogPath,
			cfgSecurityPrivileged,
			cfgStartedAt,
			cfgFinishedAt,
			cfgCloudInitUserData,
			cfgCloudInitMetaData,
			cfgCloudInitNetworkConfig,
			cfgVolatileBaseImage,
		}, reservedConfigCRI...,
		)...,
	).WithReservedPrefixes(
		append([]string{
			cfgEnvironmentPrefix,
			cfgResourcesPrefix,
		}, reservedConfigPrefixesCRI...,
		)...,
	)
)

// Container represents a LXD container including CRI specific configuration
type Container struct {
	// LXDObject inherits common CRI fields
	LXDObject
	// Profiles of the container. First entry is always the sandbox profile
	// The default profile is always excluded and managed according to the settings automatically
	Profiles []string
	// Image defines the image to use, can be the hash or local alias
	Image string
	// Privileged defines if the container is run privileged
	Privileged bool
	// Environment specifies to the container exported environment variables
	Environment map[string]string

	// CRIObject inherits common CRI fields
	CRIObject
	// Metadata contains user defined data
	Metadata ContainerMetadata
	// StartedAt is when the container was started
	StartedAt time.Time
	// FinishedAt is when the container was exited
	FinishedAt time.Time
	// StateName of the current container
	StateName ContainerStateName
	// LogPath TODO, to be implemented?
	LogPath string
	// CloudInit fields
	CloudInitUserData      string
	CloudInitMetaData      string
	CloudInitNetworkConfig string
	// Resources contain cgroup information for handling resource constraints for the container
	Resources *opencontainers.LinuxResources

	// sandbox is the parent sandbox of this container
	sandbox *Sandbox
	// State contains the current additional state info of this container
	state *ContainerState
}

// ContainerState holds information about the container state
type ContainerState struct {
	// Pid of the container
	// +readonly
	Pid int64
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

// SandboxID returns the last profile name which is the sandbox profile name
func (c *Container) SandboxID() string {
	return c.Profiles[len(c.Profiles)-1]
}

func (c *Container) getSandbox() (*Sandbox, error) {
	if len(c.Profiles) > 0 {
		sandbox, err := c.client.GetSandbox(c.SandboxID())
		if err != nil {
			return nil, err
		}

		return sandbox, nil
	}

	return nil, fmt.Errorf("%w: container '%v' must have at least one profile", ErrConvert, c.ID)
}

// State looks up additional state info
// Implemented as lazy loading, and returns same result if already looked up
// Not thread safe! But it's expected the pointers stay in the same routine
func (c *Container) State() (*ContainerState, error) {
	var err error
	if c.state == nil {
		c.state, err = c.getState()
		if err != nil {
			return nil, err
		}
	}

	return c.state, nil
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
		CPUUsage:        uint64(state.CPU.Usage),
		MemoryUsage:     uint64(state.Memory.Usage),
		FilesystemUsage: uint64(state.Disk[lxdInitDefaultDiskName].Usage),
	}

	return cs, nil
}

// refresh loads the container again from LXD with data and ETag
func (c *Container) refresh() error {
	r, err := c.client.GetContainer(c.ID)
	if err != nil {
		return err
	}

	*c = *r

	return nil
}

// Apply will save the changes of a container if validation was successful, refreshes ETag after save
func (c *Container) Apply() error {
	err := c.validate()
	if err != nil {
		return err
	}

	err = c.apply()
	if err != nil {
		return err
	}

	return c.refresh()
}

// Start the container
func (c *Container) Start() error {
	err := c.client.opwait.StartContainer(c.ID)
	if err != nil {
		if shared.IsErrNotFound(err) {
			return fmt.Errorf("container %w: %s", shared.NewErrNotFound(), c.ID)
		}

		return err
	}

	// when changing state of container, need to refresh ETag
	err = c.refresh()
	if err != nil {
		return err
	}

	// delete created mark if exists, so next stopping state can be exited
	delete(c.Config, cfgState)
	c.StartedAt = time.Now()

	return c.Apply()
}

// Stop will try to stop the container, returns nil when container is already stopped or
// got stopped in the meantime, otherwise it will return an error.
func (c *Container) Stop(timeout int) error {
	err := c.client.opwait.StopContainer(c.ID, timeout, 1)
	if err != nil {
		if shared.IsErrNotFound(err) {
			return nil
		}

		return err
	}

	// when changing state of container, need to refresh ETag
	err = c.refresh()
	if err != nil {
		return err
	}

	c.FinishedAt = time.Now()

	return c.Apply()
}

// Delete the container, returns nil when container is already deleted or
// got deleted in the meantime, otherwise it will return an error.
func (c *Container) Delete() error {
	err := c.client.opwait.DeleteContainer(c.ID)
	if err != nil {
		if shared.IsErrNotFound(err) {
			return nil
		}

		return err
	}

	return nil
}

// validate checks for misconfigurations
func (c *Container) validate() error {
	s, err := c.Sandbox()
	if err != nil {
		return err
	}

	switch s.NetworkConfig.Mode { // nolint: exhaustive, gocritic
	case NetworkHost:
		if !c.Privileged {
			return fmt.Errorf("%w: `podSpec.hostNetwork: true` can only be used together with `containerSpec.securityContext.privileged: true`", ErrUsage)
		}
	}

	return nil
}

// apply saves the changes to LXD
// Will not obtain the new ETag!
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
		return fmt.Errorf("image %w on local remote: %s", shared.NewErrNotFound(), c.Image)
	}

	config := makeContainerConfig(c)

	devices := make(map[string]map[string]string)

	for _, d := range c.Devices {
		name, options := d.ToMap()
		devices[name] = options
	}

	for key, val := range c.Config {
		if containerConfigStore.IsReserved(key) {
			log.Warnf("config key '%v' is reserved and can't be used", key)
		} else {
			config[key] = val
		}
	}

	config[cfgSchema] = SchemaVersionContainer
	contPut := api.ContainerPut{
		Profiles: c.Profiles,
		Config:   config,
		Devices:  devices,
	}

	if c.ID == "" {
		// container has to be created
		c.ID = c.CreateID()

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
		return fmt.Errorf("update container not allowed: %w", ErrMissingETag)
	}

	err = c.client.opwait.UpdateContainer(c.ID, contPut, c.ETag)
	if err != nil {
		if shared.IsErrNotFound(err) {
			return fmt.Errorf("container %w: %s", shared.NewErrNotFound(), c.ID)
		}

		return err
	}

	return nil
}

// CreateID creates a unique container id
func (c *Container) CreateID() string {
	bin := md5.Sum([]byte(uuid.NewUUID())) // nolint: gosec
	return string(c.Metadata.Name[0]) + b32lowerEncoder.EncodeToString(bin[:])[:15]
}

// GetInetAddress returns the IPv4 address of the first matching interface in the parameter list
// empty string if nothing was found
func (c *Container) GetInetAddress(ifs []string) string {
	st, err := c.State()
	if err != nil {
		return ""
	}

	for _, i := range ifs {
		if netif, ok := st.Network[i]; ok {
			for _, addr := range netif.Addresses {
				if addr.Family == "inet" {
					return addr.Address
				}
			}
		}
	}

	return ""
}

func makeContainerConfig(c *Container) map[string]string { // nolint: gocognit
	// default values for new containers
	if c.ID == "" {
		c.Config[cfgState] = ContainerStateCreated.String()
		c.CreatedAt = time.Now()
	}

	config := make(map[string]string)

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

	if c.Resources != nil { // nolint: nestif
		if c.Resources.CPU != nil {
			if c.Resources.CPU.Shares != nil {
				config[cfgResourcesCPUShares] = strconv.FormatUint(*c.Resources.CPU.Shares, 10)
			}

			if c.Resources.CPU.Quota != nil {
				config[cfgResourcesCPUQuota] = strconv.FormatInt(*c.Resources.CPU.Quota, 10)
			}

			if c.Resources.CPU.Period != nil {
				config[cfgResourcesCPUPeriod] = strconv.FormatUint(*c.Resources.CPU.Period, 10)
			}

			if c.Resources.CPU.Quota != nil && *c.Resources.CPU.Quota > 0 && c.Resources.CPU.Period != nil && *c.Resources.CPU.Period > 0 {
				// nolint:gomnd
				config[cfgLimitCPUAllowance] = fmt.Sprintf("%dms/%dms",
					int(math.Ceil(float64(*c.Resources.CPU.Quota)/1000)),
					int(math.Ceil(float64(*c.Resources.CPU.Period)/1000)),
				)
			}
		}

		if c.Resources.Memory != nil {
			if c.Resources.Memory.Limit != nil && *c.Resources.Memory.Limit > 0 {
				config[cfgLimitMemory] = strconv.FormatInt(*c.Resources.Memory.Limit, 10)
			}
		}
	}

	return config
}

// extractEnvVars extracts all the config options that start with "environment."
// and returns the environment variables + values
func extractEnvVars(config map[string]string) map[string]string {
	envVars := make(map[string]string)

	for k, v := range config {
		if strings.HasPrefix(k, cfgEnvironmentPrefix+".") {
			varName := strings.TrimPrefix(k, cfgEnvironmentPrefix+".")
			varValue := v
			envVars[varName] = varValue
		}
	}

	return envVars
}
