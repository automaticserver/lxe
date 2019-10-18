package network

import (
	"context"
	"net"
)

const (
	DefaultInterface = "eth0"
)

// NetworkPlugin is the interface for lxe network plugins
type NetworkPlugin interface {
	// PodNetwork enters a pod network environment context
	PodNetwork(namespace, name, id string, annotations map[string]string) (PodNetwork, error)
	// Status returns error if the plugin is in error state
	Status() error
}

// PodNetwork is the interface for a pod network environment
type PodNetwork interface {
	// Setup creates the network for that pod and the result is saved
	Setup(ctx context.Context) ([]byte, error)
	// Teardown removes the network of that pod
	Teardown(ctx context.Context) error
	// Status reports IP and any error with the network of that pod
	Status(ctx context.Context) (*PodNetworkStatus, error)
}

// PodNetworkStatus contains status and addresses of that pod network
type PodNetworkStatus struct {
	// The IP of the pod network
	IPs []net.IP
}
