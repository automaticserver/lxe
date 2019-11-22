package lxf

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	opencontainers "github.com/opencontainers/runtime-spec/specs-go"
)

// NewContainer creates a local representation of a container
func (l *Client) NewContainer(sandboxID string, additionalProfiles ...string) *Container {
	c := &Container{}
	c.client = l
	c.Profiles = append(c.Profiles, additionalProfiles...)
	c.Profiles = append(c.Profiles, sandboxID)
	c.Config = make(map[string]string)
	c.Environment = make(map[string]string)

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
	var (
		err  error
		etag string
	)

	cts, err := l.server.GetContainers()
	if err != nil {
		return nil, NewContainerError("lxdApi", err)
	}

	var cl = []*Container{}

	for _, ct := range cts {
		ct := ct // pin!
		if !IsCRI(ct) {
			continue
		}

		c, err := l.toContainer(&ct, etag)
		if err != nil {
			return nil, err
		}

		cl = append(cl, c)
	}

	return cl, nil
}

// toContainer will convert an lxd container to lxf format
func (l *Client) toContainer(ct *api.Container, etag string) (*Container, error) { // nolint: gocognit
	var err error

	var attempt uint64
	if attemptS, is := ct.Config[cfgMetaAttempt]; is {
		attempt, err = strconv.ParseUint(attemptS, 10, 32)
		if err != nil {
			return nil, err
		}
	}

	var privileged bool
	if privilegedS, is := ct.Config[cfgSecurityPrivileged]; is {
		privileged, err = strconv.ParseBool(privilegedS)
		if err != nil {
			return nil, err
		}
	}

	createdAt := time.Time{}.UnixNano()
	if createdAtS, is := ct.Config[cfgCreatedAt]; is {
		createdAt, err = strconv.ParseInt(createdAtS, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	startedAt := time.Time{}.UnixNano()
	if startedAtS, is := ct.Config[cfgStartedAt]; is {
		startedAt, err = strconv.ParseInt(startedAtS, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	finishedAt := time.Time{}.UnixNano()
	if finishedAtS, is := ct.Config[cfgFinishedAt]; is {
		finishedAt, err = strconv.ParseInt(finishedAtS, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	c := &Container{}
	c.client = l

	c.ID = ct.Name
	c.ETag = etag
	c.Image = ct.Config[cfgVolatileBaseImage]
	c.Metadata = ContainerMetadata{
		Name:    ct.Config[cfgMetaName],
		Attempt: uint32(attempt),
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

	// get devices
	for name, options := range ct.Devices {
		d, err := device.Detect(name, options)
		if err != nil {
			return nil, err
		}

		c.Devices.Upsert(d)
	}

	c.Resources = &opencontainers.LinuxResources{}
	c.Resources.CPU = &opencontainers.LinuxCPU{}

	if sharesS := ct.Config[cfgResourcesCPUShares]; sharesS != "" {
		shares, err := strconv.ParseUint(sharesS, 10, 64)
		if err != nil {
			return nil, err
		}

		c.Resources.CPU.Shares = &shares
	}

	if quotaS := ct.Config[cfgResourcesCPUQuota]; quotaS != "" {
		quota, err := strconv.ParseInt(quotaS, 10, 64)
		if err != nil {
			return nil, err
		}

		c.Resources.CPU.Quota = &quota
	}

	if periodS := ct.Config[cfgResourcesCPUPeriod]; periodS != "" {
		period, err := strconv.ParseUint(periodS, 10, 64)
		if err != nil {
			return nil, err
		}

		c.Resources.CPU.Period = &period
	}

	c.Resources.Memory = &opencontainers.LinuxMemory{}

	if memoryS := ct.Config[cfgResourcesMemoryLimit]; memoryS != "" {
		memory, err := strconv.ParseInt(memoryS, 10, 64)
		if err != nil {
			return nil, err
		}

		c.Resources.Memory.Limit = &memory
	}

	c.Profiles = ct.Profiles
	if len(c.Profiles) == 0 {
		return nil, fmt.Errorf("Container '%v' has no sandbox", c.ID)
	}

	// Map status code of LXD to CRI
	switch ct.StatusCode {
	case api.Running:
		c.StateName = ContainerStateRunning
	case api.Stopped, api.Aborting, api.Stopping:
		// we have to differentiate between stopped and created. If "user.state" exists, then it must be created, otherwise
		// its exited
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

type EventHandler interface {
	ContainerStarted(ctx context.Context, c *Container) error
	ContainerStopped(ctx context.Context, c *Container) error
}

// lifecycleEventHandler is registered to the lxd event handler for listening to container start events
func (l *Client) lifecycleEventHandler(event api.Event) {
	// we should always only get lifecycle events due to the handler setup but just in case ...
	if event.Type != "lifecycle" {
		return
	}

	eventLifecycle := api.EventLifecycle{}

	err := json.Unmarshal(event.Metadata, &eventLifecycle)
	if err != nil {
		logger.Errorf("unable to unmarshal to json lifecycle event: %v", event.Metadata)
		return
	}

	// Early exit. We are only interested in container started and stopped events
	if eventLifecycle.Action != "container-started" && eventLifecycle.Action != "container-stopped" {
		return
	}

	containerID := strings.TrimPrefix(eventLifecycle.Source, "/1.0/containers/")

	c, err := l.GetContainer(containerID)
	if err != nil {
		if IsContainerNotFound(err) {
			// If the started container is not a cri container, we also get the error not found. So this container can be
			// ignored
			return
		}

		// still return immediately since we can't do anything when we get an error here
		logger.Errorf("lifecycle: ContainerID %v trying to get container: %v", containerID, err)

		return
	}

	switch eventLifecycle.Action {
	case "container-started":
		err := l.eventHandler.ContainerStarted(context.TODO(), c)
		if err != nil {
			logger.Errorf("lifecycle: handling event %v for container %v failed: %v", eventLifecycle.Action, containerID, err)
			return
		}
	case "container-stopped":
		err := l.eventHandler.ContainerStopped(context.TODO(), c)
		if err != nil {
			logger.Errorf("lifecycle: handling event %v for container %v failed: %v", eventLifecycle.Action, containerID, err)
			return
		}
	}
}
