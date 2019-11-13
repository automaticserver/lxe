package network

import (
	"context"
	"net"

	"github.com/automaticserver/lxe/lxf/device"
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
}

// PodNetwork is the interface for a pod network environment.
type PodNetwork interface {
	DeviceHandler
	PidHandler
}

// With DeviceHandler the result are lxf.Network devices. The calls for this behavior are made before the resources are
// effectively started and deleted respectively.
type DeviceHandler interface {
	// SetupDevice creates a lxd device for that pod.
	SetupDevice(ctx context.Context) (device.Nic, error)
	// TeardownDevice removes the network of that pod. Must tear down networking as good as possible, an error will only be
	// logged and doesn't stop execution of further statements.
	TeardownDevice(ctx context.Context) error
	// Attach a container to the pod network.
	AttachDevice(ctx context.Context) (device.Nic, error)
	// Detach a container from the pod network. Must detach networking as good as possible, an error will only be logged
	// and doesn't stop execution of further statements.
	DetachDevice(ctx context.Context) error
	// StatusDevice reports IP and any error with the network of that pod.
	StatusDevice(ctx context.Context) (*PodNetworkStatus, error)
}

// With PidHandler the resulting interfaces are directly setup in the process. The calls for this behavior are made
// after the resources are effectively started and stopped respectively.
type PidHandler interface {
	// SetupPid creates the network for that pod. pid is the process id of the pod. The retured result bytes are provided for
	// the other calls of this PodNetwork.
	SetupPid(ctx context.Context, pid int64) ([]byte, error)
	// TeardownPid removes the network of that pod. Must tear down networking as good as possible, an error will only be
	// logged and doesn't stop execution of further statements.
	TeardownPid(ctx context.Context, result []byte) error
	// AttachPid attaches a container to the pod network. pid is the process id of the container. Can return arbitrary
	// bytes or nic devices or both. The retured result bytes replace the existing one if provided.
	AttachPid(ctx context.Context, result []byte, pid int64) ([]byte, error)
	// DetachPid detaches a container from the pod network. Must detach networking as good as possible, an error will only
	// be logged and doesn't stop execution of further statements.
	DetachPid(ctx context.Context, result []byte) error
	// StatusPid reports IP and any error with the network of that pod. Bytes can be nil if LXE thinks it never ran Setup
	// and thus also pid is not set or weren't returned yet.
	StatusPid(ctx context.Context, result []byte, pid int64) (*PodNetworkStatus, error)
}

// PodNetworkStatus contains status and addresses of that pod network
type PodNetworkStatus struct {
	// The IP of the pod network
	IPs []net.IP
}
