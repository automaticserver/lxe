package network // import "github.com/automaticserver/lxe/network"

import (
	"context"
	"errors"
	"fmt"

	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

var (
	ErrNoop = errors.New("noop")
	// verify interface satisfaction
	_ Plugin           = &noopPlugin{}
	_ PodNetwork       = &noopPodNetwork{}
	_ ContainerNetwork = &noopContainerNetwork{}
)

// noopPlugin implements Plugin without doing anything
type noopPlugin struct{}

// InitPluginNoop instantiates the noop plugin
func InitPluginNoop() (*noopPlugin, error) { // nolint: golint // intended to not export
	return &noopPlugin{}, nil
}

// PodNetwork enters a pod network environment context
func (p *noopPlugin) PodNetwork(_ string, _ map[string]string) (PodNetwork, error) {
	return &noopPodNetwork{}, nil
}

func (p *noopPlugin) Status() error {
	return fmt.Errorf("%w plugin is never running", ErrNoop)
}

// UpdateRuntimeConfig is called when there are updates to the configuration which the plugin might need to apply
func (p *noopPlugin) UpdateRuntimeConfig(_ *rtApi.RuntimeConfig) error {
	return fmt.Errorf("%w plugin can't update runtime config", ErrNoop)
}

// cniPodNetwork is a pod network environment context
type noopPodNetwork struct{}

// ContainerNetwork enters a container network environment context
func (s *noopPodNetwork) ContainerNetwork(_ string, _ map[string]string) (ContainerNetwork, error) {
	return &noopContainerNetwork{}, nil
}

// Status reports IP and any error with the network of that pod
func (s *noopPodNetwork) Status(_ context.Context, _ *PropertiesRunning) (*Status, error) {
	return nil, nil
}

// WhenCreated is called when the pod is created
func (s *noopPodNetwork) WhenCreated(_ context.Context, _ *Properties) (*Result, error) {
	return nil, nil
}

// WhenStarted is called when the pod is started.
func (s *noopPodNetwork) WhenStarted(_ context.Context, _ *PropertiesRunning) (*Result, error) {
	return nil, nil
}

// WhenStopped is called when the pod is stopped. If tearing down here, must tear down as good as possible. Must tear
// down here if not implemented for WhenDeleted. If an error is returned it will only be logged
func (s *noopPodNetwork) WhenStopped(_ context.Context, _ *Properties) error {
	return nil
}

// WhenDeleted is called when the pod is deleted. If tearing down here, must tear down as good as possible. Must tear
// down here if not implemented for WhenStopped. If an error is returned it will only be logged
func (s *noopPodNetwork) WhenDeleted(_ context.Context, _ *Properties) error {
	return nil
}

// noopContainerNetwork is a container network environment context
type noopContainerNetwork struct{}

// WhenCreated is called when the container is created
func (c *noopContainerNetwork) WhenCreated(_ context.Context, _ *Properties) (*Result, error) {
	return nil, nil
}

// WhenStarted is called when the container is started
func (c *noopContainerNetwork) WhenStarted(_ context.Context, _ *PropertiesRunning) (*Result, error) {
	return nil, nil
}

// WhenStopped is called when the container is stopped. If tearing down here, must tear down as good as possible. Must tear
// down here if not implemented for WhenDeleted. If an error is returned it will only be logged
func (c *noopContainerNetwork) WhenStopped(_ context.Context, _ *Properties) error {
	return nil
}

// WhenDeleted is called when the container is deleted. If tearing down here, must tear down as good as possible. Must tear
// down here if not implemented for WhenStopped. If an error is returned it will only be logged
func (c *noopContainerNetwork) WhenDeleted(_ context.Context, _ *Properties) error {
	return nil
}
