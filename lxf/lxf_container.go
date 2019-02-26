package lxf

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxe/lxf/device"
	"github.com/lxc/lxe/network"
)

// NewContainer creates a local representation of a container
func (l *Client) NewContainer(sandboxID string) *Container {
	c := &Container{}
	c.client = l
	c.Profiles = append(c.Profiles, sandboxID)
	return c
}

// GetContainer returns the container identified by id
func (l *Client) GetContainer(id string) (*Container, error) {
	ct, ETag, err := l.server.GetContainer(id)
	if err != nil {
		return nil, NewContainerError(id, err)
	}

	if !IsCRI(ct) {
		return nil, NewContainerError(id, fmt.Errorf(ErrorLXDNotFound))
	}

	return l.toContainer(ct, ETag)
}

// ListContainers returns a list of all available containers
func (l *Client) ListContainers() ([]*Container, error) {
	var err error
	ETag := ""
	cts, err := l.server.GetContainers()
	if err != nil {
		return nil, NewContainerError("lxdApi", err)
	}

	cl := []*Container{}
	for _, ct := range cts {
		if !IsCRI(ct) {
			continue
		}
		c, err := l.toContainer(&ct, ETag)
		if err != nil {
			return nil, err
		}
		cl = append(cl, c)
	}

	return cl, nil
}

// toContainer will convert an lxd container to lxf format
func (l *Client) toContainer(ct *api.Container, ETag string) (*Container, error) {
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
	c.client = l

	c.ID = ct.Name
	c.ETag = ETag
	c.Image = ct.Config[cfgVolatileBaseImage]
	c.Metadata = ContainerMetadata{
		Name:    ct.Config[cfgMetaName],
		Attempt: uint32(attempts),
	}
	c.Annotations = containerConfigStore.StripedPrefixMap(ct.Config, cfgAnnotations)
	c.Labels = containerConfigStore.StripedPrefixMap(ct.Config, cfgLabels)
	c.Config = containerConfigStore.UnreservedMap(ct.Config)
	c.LogPath = ct.Config[cfgLogPath]

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
		if v != defaultProfile {
			c.Profiles = append(c.Profiles, v)
		}
	}
	if len(c.Profiles) == 0 {
		return nil, fmt.Errorf("Container '%v' has no sandbox", c.ID)
	}

	// Map status code of LXD to CRI
	switch ct.StatusCode {
	case api.Running:
		c.StateName = ContainerStateRunning
	case api.Stopped, api.Aborting, api.Stopping:
		// we have to differentiate between stopped and created. If "user.state" exists, then it must be
		// created, otherwise its exited
		if state, has := c.Config[cfgState]; has && state == string(ContainerStateCreated) {
			c.StateName = ContainerStateCreated
		} else {
			c.StateName = ContainerStateExited
		}
	default:
		c.StateName = ContainerStateUnknown
	}

	return c, nil
}

// lifecycleEventHandler is registered to the lxd event handler for listening to container start events
func (l *Client) lifecycleEventHandler(event api.Event) {
	// we should always only get lifecycle events due to the handler setup
	// but just in case ...
	if event.Type != "lifecycle" {
		return
	}

	eventLifecycle := api.EventLifecycle{}
	err := json.Unmarshal(event.Metadata, &eventLifecycle)
	if err != nil {
		logger.Errorf("unable to unmarshal to json lifecycle event: %v", event.Metadata)
		return
	}

	// we are only interested in container started events
	// TODO: Unregister IP address when container is stopping if network-plugin is CNI
	if eventLifecycle.Action != "container-started" {
		return
	}

	containerID := strings.TrimPrefix(eventLifecycle.Source, "/1.0/containers/")
	c, err := l.GetContainer(containerID)
	if err != nil {
		if IsContainerNotFound(err) {
			// The started container is not a cri container, we also get the error not found
			// So this container can be ignored
			return
		}
		logger.Errorf("lifecycle: ContainerID %v trying to get container: %v", containerID, err)
		return
	}

	s, err := c.Sandbox()
	if err != nil {
		logger.Errorf("lifecycle: ContainerID %v trying to get sandbox: %v", containerID, err)
		return
	}
	st, err := c.State()
	if err != nil {
		logger.Errorf("lifecycle: ContainerID %v trying to get state: %v", containerID, err)
		return
	}

	switch s.NetworkConfig.Mode {
	case NetworkCNI:
		// attach interface using CNI
		result, err := network.AttachCNIInterface(s.Metadata.Namespace, s.Metadata.Name, c.ID, st.Pid)
		if err != nil {
			logger.Errorf("unable to attach CNI interface to container (%v): %v", c.ID, err)
		}
		s.NetworkConfig.ModeData["result"] = string(result)
		err = s.apply()
		if err != nil {
			logger.Errorf("unable to save sandbox after attaching CNI interface to container (%v): %v", c.ID, err)
		}
	default:
		// nothing to do, all other modes need no help after starting
	}
}
