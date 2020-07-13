package network

import (
	"context"
	"net"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/network/cloudinit"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	// DefaultInterface for containers is always eth0
	DefaultInterface = "eth0"
)

// NetworkPlugin is the interface for lxe network plugins
type Plugin interface {
	// PodNetwork enters a pod network environment context
	PodNetwork(id string, annotations map[string]string) (PodNetwork, error)
	// Status returns error if the plugin is in error state
	Status() error
	// UpdateRuntimeConfig is called when there are updates to the configuration which the plugin might need to apply
	UpdateRuntimeConfig(conf *rtApi.RuntimeConfig) error
}

// PodNetwork is the interface for a pod network environment.
type PodNetwork interface {
	// ContainerNetwork enters a container network environment context
	ContainerNetwork(id string, annotations map[string]string) (ContainerNetwork, error)
	// Status reports IP and any error with the network of that pod
	Status(ctx context.Context, prop *PropertiesRunning) (*Status, error)
	// WhenCreated is called when the pod is created
	WhenCreated(ctx context.Context, prop *Properties) (*Result, error)
	// WhenStarted is called when the pod is started.
	WhenStarted(ctx context.Context, prop *PropertiesRunning) (*Result, error)
	// WhenStopped is called when the pod is stopped. If tearing down here, must tear down as good as possible. Must tear
	// down here if not implemented for WhenDeleted. If an error is returned it will only be logged
	WhenStopped(ctx context.Context, prop *Properties) error
	// WhenDeleted is called when the pod is deleted. If tearing down here, must tear down as good as possible. Must tear
	// down here if not implemented for WhenStopped. If an error is returned it will only be logged
	WhenDeleted(ctx context.Context, prop *Properties) error
}

// ContainerNetwork is the interface for a container network environment context
type ContainerNetwork interface {
	// WhenCreated is called when the container is created
	WhenCreated(ctx context.Context, prop *Properties) (*Result, error)
	// WhenStarted is called when the container is started
	WhenStarted(ctx context.Context, prop *PropertiesRunning) (*Result, error)
	// WhenStopped is called when the container is stopped. If tearing down here, must tear down as good as possible. Must tear
	// down here if not implemented for WhenDeleted. If an error is returned it will only be logged
	WhenStopped(ctx context.Context, prop *Properties) error
	// WhenDeleted is called when the container is deleted. If tearing down here, must tear down as good as possible. Must tear
	// down here if not implemented for WhenStopped. If an error is returned it will only be logged
	WhenDeleted(ctx context.Context, prop *Properties) error
}

// Properties of the resource at the time of the call
type Properties struct {
	// Arbitrary Data are provided if a previous call on this PodNetwork returned them
	Data map[string]string
}

// PropertiesRunning contains additionally running info
type PropertiesRunning struct {
	Properties
	// Pid of the resource. This value is set for calls where the applicable resource is running
	Pid int64
}

// Result contains additionally info which can only be set on creation
type Result struct {
	// Arbitrary Data to keep related to the pod. If non-nil they will overwrite the previous saved data.
	Data map[string]string
	// List of Nics to add to the resource
	Nics []device.Nic
	// NetworkConfigEntries of cloudinit to be set. Keep in mind cloudinit runs only when the container starts
	NetworkConfigEntries []cloudinit.NetworkConfigEntryPhysical
}

// Contains Status and addresses of that pod network
type Status struct {
	// The IP of the pod network
	IPs []net.IP
}
