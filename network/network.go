package network

import (
	"context"
	"net"
)

const (
	// DefaultInterface for containers is always eth0
	DefaultInterface = "eth0"
)

// NetworkPlugin is the interface for lxe network plugins
type Plugin interface {
	// PodNetwork enters a pod network environment context
	PodNetwork(namespace, name, id string, annotations map[string]string) (PodNetwork, error)
	// Status returns error if the plugin is in error state
	Status() error
}

// PodNetwork is the interface for a pod network environment
type PodNetwork interface {
	// Setup creates the network for that pod and the result is saved. pid is the process id of the pod.
	Setup(ctx context.Context, pid int64) ([]byte, error)
	// Teardown removes the network of that pod. pid is the process id of the pod, but might be missing. Must tear down
	// networking as good as possible, an error will only be logged and doesn't stop execution of further statements
	Teardown(ctx context.Context, result []byte, pid int64) error
	// Attach a container to the pod network. pid is the process id of the container.
	Attach(ctx context.Context, result []byte, pid int64) error
	// Detach a container from the pod network. pid is the process id of the container, but might be missing. Must detach
	// networking as good as possible, an error will only be logged and doesn't stop execution of further statements
	Detach(ctx context.Context, result []byte, pid int64) error
	// Status reports IP and any error with the network of that pod. Bytes can be nil if LXE thinks it never ran Setup and
	// thus also pid is not set
	Status(ctx context.Context, result []byte, pid int64) (*PodNetworkStatus, error)
}

// PodNetworkStatus contains status and addresses of that pod network
type PodNetworkStatus struct {
	// The IP of the pod network
	IPs []net.IP
}
