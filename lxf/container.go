package lxf

import (
	"crypto/md5" // nolint: gosec #nosec (no sensitive data)
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxe/lxf/device"
	"github.com/lxc/lxe/lxf/lxo"
	"github.com/lxc/lxe/network"
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
)

// ContainerState defines the state of the container, default is ContainerStateCreated
type ContainerState string

const (
	// ContainerStateCreated it's there but not started yet
	ContainerStateCreated ContainerState = "created"
	// ContainerStateRunning it's there and running
	ContainerStateRunning ContainerState = "running"
	// ContainerStateExited it's there but terminated
	ContainerStateExited ContainerState = "exited"
	// ContainerStateUnknown it's there but we don't know what it's doing
	ContainerStateUnknown ContainerState = "unknown"
)

func (s ContainerState) String() string {
	return string(s)
}

var (
	containerConfigStore = NewConfigStore().WithReserved(cfgSchema, cfgLogPath, cfgIsCRI,
		cfgSecurityPrivileged, cfgSecurityNesting, cfgState, cfgMetaName, cfgMetaAttempt, cfgCreatedAt, cfgStartedAt, cfgCloudInitUserData, cfgCloudInitMetaData,
		cfgCloudInitNetworkConfig).
		WithReservedPrefixes(cfgLabels, cfgAnnotations, "volatile")
)

// Container is a unified interface to LXDs container methodes
type Container struct {
	CRIObject
	LogPath  string
	Metadata ContainerMetadata
	// State is read only
	State ContainerState
	// Pid is readonly
	Pid int64
	// StartedAt is read only, if not started it will be the zero time
	StartedAt              time.Time
	CreatedAt              time.Time
	FinishedAt             time.Time
	Privileged             bool
	CloudInitUserData      string
	CloudInitMetaData      string
	CloudInitNetworkConfig string
	// Network is readonly
	Network map[string]api.ContainerStateNetwork
	// Implements spec.Env
	EnvironmentVars map[string]string

	Stats ContainerStats

	Sandbox *Sandbox
	Image   string // can be hash or local alias
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

// CreateContainer will instantiate a container but not start it
func (l *LXF) CreateContainer(c *Container) error {
	if c.Sandbox == nil {
		return fmt.Errorf("container needs a sandbox")
	}

	c.State = ContainerStateCreated
	c.CreatedAt = time.Now()
	// c.Config[cfgAutoStartOnBoot] = strconv.FormatBool(true)

	switch c.Sandbox.NetworkConfig.Mode {
	case NetworkHost:
		if !c.Privileged {
			return fmt.Errorf("'podSpec.hostNetwork: true' can only be used together with 'containerSpec.securityContext.privileged: true'")
		}
	default:
		// do nothing
	}

	return l.saveContainer(c)
}

// UpdateContainer will update an existing container
func (l *LXF) UpdateContainer(c *Container) error {
	if c.Sandbox == nil {
		return fmt.Errorf("container needs a sandbox")
	}
	return l.saveContainer(c)
}

// StartContainer starts an existing container
func (l *LXF) StartContainer(id string) error {
	err := lxo.StartContainer(l.server, id)
	if err != nil {
		return err
	}

	// TODO: Since we now need the full lxe.Container we could ensure the
	// following steps over that, now it's raw-ish lxd
	ct, ETag, err := l.server.GetContainer(id)
	if err != nil {
		return err
	}

	// custom state created is removed
	delete(ct.Config, cfgState)

	// set started at date
	if ct.Config[cfgStartedAt] == "" {
		ct.Config[cfgStartedAt] = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	// c, err := l.GetContainer(id)
	// if err != nil {
	// 	return err
	// }
	// go l.remountMissingVolumes(c)

	return lxo.UpdateContainer(l.server, id, ct.Writable(), ETag)
}

// StopContainer will try to stop the container, returns nil when container is already deleted or
// got deleted in the meantime, otherwise it will return an error.
// If it's not deleted after 30 seconds it will return an error (might be way to low).
func (l *LXF) StopContainer(id string, timeout int) error {
	err := lxo.StopContainer(l.server, id, timeout, 2)
	if err != nil {
		return err
	}

	// TODO: Since we now need the full lxe.Container we could ensure the
	// following steps over that, now it's raw-ish lxd
	ct, ETag, err := l.server.GetContainer(id)
	if err != nil {
		return err
	}

	// TODO: probably Exit Code 143? With dockershim this happens when forcing containers to stop

	// set finished at date
	if ct.Config[cfgFinishedAt] == "" {
		ct.Config[cfgFinishedAt] = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	return lxo.UpdateContainer(l.server, id, ct.Writable(), ETag)
}

// DeleteContainer will delete the container
func (l *LXF) DeleteContainer(id string) error {
	return lxo.DeleteContainer(l.server, id)
}

// ListContainers returns a list of all available containers
func (l *LXF) ListContainers() ([]*Container, error) { // nolint:dupl
	ETag := ""
	cts, err := l.server.GetContainers()
	if err != nil {
		return nil, err
	}
	result := []*Container{}
	for _, ct := range cts {
		if !IsCRI(ct) {
			continue
		}
		res, err := l.toContainer(&ct, ETag)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

// GetContainer returns the container identified by id
func (l *LXF) GetContainer(id string) (*Container, error) {
	ct, ETag, err := l.server.GetContainer(id)
	if err != nil {
		return nil, err
	}

	if !IsCRI(ct) {
		return nil, fmt.Errorf(ErrorNotFound)
	}

	return l.toContainer(ct, ETag)
}

// saveContainer
func (l *LXF) saveContainer(c *Container) error {
	// TODO: can't this be done easier?
	imageID, err := l.parseImage(c.Image)
	if err != nil {
		return err
	}
	hash, found, err := imageID.Hash(l)
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
		Profiles: []string{c.Sandbox.ID, "default"},
		Config:   config,
		Devices:  devices,
	}
	if c.ID == "" { // container has to be created
		c.ID = c.CreateID()
		return lxo.CreateContainer(l.server, api.ContainersPost{
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
	return lxo.UpdateContainer(l.server, c.ID, contPut, c.ETag)
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

	if c.State == ContainerStateCreated {
		config[cfgState] = ContainerStateCreated.String()
	}

	for k, v := range c.EnvironmentVars {
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

	for _, containerName := range c.Sandbox.SecurityNesting {
		if containerName == c.Metadata.Name {
			config[cfgSecurityNesting] = strconv.FormatBool(true)
		}
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
func (l *LXF) toContainer(ct *api.Container, ETag string) (*Container, error) {
	state, _, err := l.server.GetContainerState(ct.Name)
	if err != nil {
		return nil, err
	}
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
	c.Pid = state.Pid
	c.CreatedAt = time.Unix(0, createdAt)
	c.StartedAt = time.Unix(0, startedAt)
	c.FinishedAt = time.Unix(0, finishedAt)
	c.Stats = ContainerStats{
		CPUUsage:        uint64(state.CPU.Usage),
		MemoryUsage:     uint64(state.Memory.Usage),
		FilesystemUsage: uint64(state.Disk["root"].Usage),
	}
	c.Network = state.Network
	c.EnvironmentVars = extractEnvVars(ct.Config)
	c.Privileged = privileged
	c.CloudInitUserData = ct.Config[cfgCloudInitUserData]
	c.CloudInitMetaData = ct.Config[cfgCloudInitMetaData]
	c.CloudInitNetworkConfig = ct.Config[cfgCloudInitNetworkConfig]

	// get status and map it
	switch state.StatusCode {
	case api.Running:
		c.State = ContainerStateRunning
	case api.Stopped, api.Aborting, api.Stopping:
		// we have to differentiate between stopped and created using the "user.state" config value
		if state, has := ct.Config[cfgState]; has && state == string(ContainerStateCreated) {
			c.State = ContainerStateCreated
		} else {
			c.State = ContainerStateExited
		}
	default:
		c.State = ContainerStateUnknown
	}

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

	// get sandbox
	if len(ct.Profiles) > 0 {
		c.Sandbox, err = l.GetSandbox(ct.Profiles[0])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Container '%v' must have at least one profile", ct.Name)
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

func (l *LXF) lifecycleEventHandler(message interface{}) {
	msg, err := json.Marshal(&message)
	if err != nil {
		logger.Errorf("unable to marshal json event: %v", message)
		return
	}

	event := api.Event{}
	err = json.Unmarshal(msg, &event)
	if err != nil {
		logger.Errorf("unable to unmarshal to json event: %v", message)
		return
	}

	// we should always only get lifecycle events due to the handler setup
	// but just in case ...
	if event.Type != "lifecycle" {
		return
	}

	eventLifecycle := api.EventLifecycle{}
	err = json.Unmarshal(event.Metadata, &eventLifecycle)
	if err != nil {
		logger.Errorf("unable to unmarshal to json lifecycle event: %v", message)
		return
	}

	// we are only interested in container started events
	if eventLifecycle.Action != "container-started" {
		return
	}

	containerID := strings.TrimPrefix(eventLifecycle.Source, "/1.0/containers/")
	logger.Infof("DEBUG: Container start event %v", containerID)
	time.Sleep(10 * time.Second)
	cnt, err := l.GetContainer(containerID)
	if err != nil {
		if IsErrorNotFound(err) {
			// The started container is not a cri container, so we get the error not found
			return
		}
		logger.Errorf("Unable to GetContainer %v: %v", containerID, err)
		return
	}

	// add container to queue in order to recheck if mounts are okay
	//l.AddMonitorTask(cnt, "volumes", 0, true)

	switch cnt.Sandbox.NetworkConfig.Mode {
	case NetworkCNI:
		if len(cnt.Sandbox.NetworkConfig.ModeData) == 0 {
			// new container, attach cni
			result, err := network.AttachCNIInterface(cnt.Sandbox.Metadata.Namespace, cnt.Sandbox.Metadata.Name, cnt.ID, cnt.Pid)
			if err != nil {
				logger.Errorf("unable to attach CNI interface to container (%v): %v", cnt.ID, err)
			}
			cnt.Sandbox.NetworkConfig.ModeData["result"] = string(result)
			err = l.saveSandbox(cnt.Sandbox)
			if err != nil {
				logger.Errorf("unable to save sandbox after attaching CNI interface to container (%v): %v", cnt.ID, err)
			}
		} else {
			// existing container, reattach cni
			err = network.ReattachCNIInterface(
				cnt.Sandbox.Metadata.Namespace,
				cnt.Sandbox.Metadata.Name,
				cnt.ID,
				cnt.Pid,
				cnt.Sandbox.NetworkConfig.ModeData["result"])
			if err != nil {
				logger.Errorf("unable to reattach CNI settings to container (%v): %v", cnt.ID, err)
			}
		}
	default:
		// nothing to do, all other modes need no help after starting
	}
}

// AddMonitorTask adds 'task' to be executed once or everytime for a given interval
func (l *LXF) AddMonitorTask(c *Container, task string, interval time.Duration, once bool) {
	l.cntMonitorChan <- ContainerMonitorChan{
		container:   c,
		task:        task,
		intervalSec: interval,
		once:        once,
	}
}

func (l *LXF) containerMonitor(cntMonitorChan chan ContainerMonitorChan) {
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			for i := range cntMonitorChan {
				if i.lastCheck.Add(i.intervalSec).Sub(time.Now()) <= 0 {
					switch i.task {
					case "volumes":
						// go l.remountMissingVolumes(i.container)
						i.lastCheck = time.Now()
					default:
						logger.Debugf("containerMonitor: unknown task: %v for object: %+v", i.task, i)
					}
				}
				// requeue item
				if !i.once {
					cntMonitorChan <- i
				}
			}
		}
	}
}

// TODO: Why exactly do we remove disks and add them again?
// BTW there is a bug: when LXE is started it tries to mount only the then known list of devices
// If this function somehow didn't save full list before getting restarted, devices are lost!
func (l *LXF) remountMissingVolumes(container *Container) {
	logger.Debugf("remountMissingVolumes triggered: %v", container.ID)

	allDisks := container.Disks
	mountedDisks := []device.Disk{}
	for {
		current, err := l.GetContainer(container.ID)
		if err != nil {
			logger.Errorf("remountMissingVolumes failed to update container %s info: %v", container.ID, err)
			return
		}
		if current == nil {
			logger.Errorf("current container %s was nil", container.ID)
			return
		}
		if (current.State == ContainerStateExited) || (current.State == ContainerStateUnknown) {
			logger.Debugf("state remountMissingVolumes: stale container")
			return
		}
		for _, disk := range current.Disks {
			_, _, err := l.server.GetContainerFile(current.ID, disk.Path)
			if err != nil {
				logger.Debugf("remountMissingVolumes Container(%s) '%s' path: %s: %v. - attempting remounting",
					current.ID, disk.GetName(), disk.Path, err)
			} else {
				mountedDisks = append(mountedDisks, disk)
			}
		}
		if len(mountedDisks) == len(allDisks) {
			return
		}

		// TODO: Can we remove the sleep since we redo this repeatedly in containerMonitor()?
		time.Sleep(time.Second * 1)

		// remove failed devices, to retry later (with all)
		current.Disks = mountedDisks
		err = l.UpdateContainer(current)
		if err != nil {
			logger.Debugf("Failed to update container without failed disks, %v", err)
		}

		// mount with all devices
		current.Disks = allDisks
		err = l.UpdateContainer(current)
		if err != nil {
			logger.Debugf("Failed to update container with all disks, %v", err)
		}
	}
}

// CreateID creates the unique container id based on Kubernetes container and sandbox values
// This is currently not expected to be a long term stable hashing for these informations
func (c *Container) CreateID() string {
	uuid.GetUUID()
	var parts []string
	parts = append(parts, "k8s")
	parts = append(parts, c.Metadata.Name)
	parts = append(parts, c.Sandbox.Metadata.Name)
	parts = append(parts, c.Sandbox.Metadata.Namespace)
	parts = append(parts, strconv.FormatUint(uint64(c.Sandbox.Metadata.Attempt), 10))
	parts = append(parts, c.Sandbox.Metadata.UID)
	name := strings.Join(parts, "-")

	bin := md5.Sum([]byte(name)) // nolint: gosec #nosec
	return string(c.Metadata.Name[0]) + b32lowerEncoder.EncodeToString(bin[:])[:15]
}

// GetContainerIPv4Address returns the IPv4 address of the first matching interface in the parameter list
// empty string if nothing was found
func (c *Container) GetContainerIPv4Address(ifs []string) string {
	for _, i := range ifs {
		if netif, ok := c.Network[i]; ok {
			for _, addr := range netif.Addresses {
				if addr.Family == "inet" {
					return addr.Address
				}
			}
		}
	}
	return ""
}
