package lxf

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxe/network"
)

// NewContainer creates a local representation of a container
func (l *Client) NewContainer(sandboxID string) (c *Container, err error) {
	c.client = l
	c.sandbox, err = l.GetSandbox(sandboxID)
	if err != nil {
		return nil, err
	}
	c.Profiles = append(c.Profiles, sandboxID)
	return c, nil
}

// ListContainers returns a list of all available containers
func (l *Client) ListContainers() ([]*Container, error) {
	var err error
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
		res, err := toContainer(&ct, ETag)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

// GetContainer returns the container identified by id
func (l *Client) GetContainer(id string) (*Container, error) {
	ct, ETag, err := l.server.GetContainer(id)
	if err != nil {
		return nil, err
	}

	if !IsCRI(ct) {
		return nil, fmt.Errorf("TODO: container not found")
	}

	return toContainer(ct, ETag)
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
		if err.Error() == ErrorLXDNotFound {
			// The started container is not a cri container, so we get the error not found
			// TODO: we might get a Container Not Found error (once implemented)
			return
		}
		logger.Errorf("Unable to GetContainer %v: %v", containerID, err)
		return
	}

	// add container to queue in order to recheck if mounts are okay
	//l.AddMonitorTask(c, "volumes", 0, true)

	s, err := c.Sandbox()
	if err != nil {
		return
	}

	switch s.NetworkConfig.Mode {
	case NetworkCNI:
		// attach interface using CNI
		result, err := network.AttachCNIInterface(s.Metadata.Namespace, s.Metadata.Name, c.ID, c.State.Pid)
		if err != nil {
			logger.Errorf("unable to attach CNI interface to container (%v): %v", c.ID, err)
		}
		s.NetworkConfig.ModeData["result"] = string(result)
		err = l.saveSandbox(s)
		if err != nil {
			logger.Errorf("unable to save sandbox after attaching CNI interface to container (%v): %v", c.ID, err)
		}
	default:
		// nothing to do, all other modes need no help after starting
	}
}
