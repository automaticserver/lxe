package cri // import "github.com/automaticserver/lxe/cri"

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/network"
	"github.com/automaticserver/lxe/shared"
	sharedLXD "github.com/lxc/lxd/shared"
	"golang.org/x/net/context"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func toCriStatusResponse(c *lxf.Container) *rtApi.ContainerStatusResponse {
	status := rtApi.ContainerStatus{
		Metadata: &rtApi.ContainerMetadata{
			Name:    c.Metadata.Name,
			Attempt: c.Metadata.Attempt,
		},
		State:       stateContainerAsCri(c.StateName),
		CreatedAt:   c.CreatedAt.UnixNano(),
		StartedAt:   c.StartedAt.UnixNano(),
		FinishedAt:  c.FinishedAt.UnixNano(),
		Id:          c.ID,
		Labels:      c.Labels,
		Annotations: c.Annotations,
		Image:       &rtApi.ImageSpec{Image: c.Image},
		ImageRef:    c.Image,
		Mounts:      []*rtApi.Mount{},
	}

	for _, dev := range c.Devices {
		switch d := dev.(type) {
		case *device.Block:
			status.Mounts = append(status.Mounts, &rtApi.Mount{
				ContainerPath:  d.Path,
				HostPath:       d.Source,
				Readonly:       false,                                                // probably always that?
				SelinuxRelabel: false,                                                // though don't know what this means
				Propagation:    rtApi.MountPropagation_PROPAGATION_HOST_TO_CONTAINER, // unsure
			})
		case *device.Disk:
			status.Mounts = append(status.Mounts, &rtApi.Mount{
				ContainerPath:  d.Path,
				HostPath:       d.Source,
				Readonly:       d.Readonly,
				SelinuxRelabel: false, // though don't know what this means
				Propagation:    rtApi.MountPropagation_PROPAGATION_PRIVATE,
			})
		}
	}

	return &rtApi.ContainerStatusResponse{
		Status: &status,
		Info:   map[string]string{},
	}
}

func toCriStats(c *lxf.Container) (*rtApi.ContainerStats, error) {
	st, err := c.State()
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixNano()

	cpu := rtApi.CpuUsage{
		Timestamp:            now,
		UsageCoreNanoSeconds: &rtApi.UInt64Value{Value: st.Stats.CPUUsage},
	}
	memory := rtApi.MemoryUsage{
		Timestamp:       now,
		WorkingSetBytes: &rtApi.UInt64Value{Value: st.Stats.MemoryUsage},
	}
	disk := rtApi.FilesystemUsage{
		Timestamp: now,
		FsId: &rtApi.FilesystemIdentifier{
			Mountpoint: path.Join(sharedLXD.VarPath("containers"), c.ID, "rootfs"),
		},
		UsedBytes:  &rtApi.UInt64Value{Value: st.Stats.FilesystemUsage}, // TODO: root seems not visible? or does it depend?
		InodesUsed: &rtApi.UInt64Value{Value: 0},                        // TODO: do we have to find out?
	}
	attribs := rtApi.ContainerAttributes{
		Id: c.ID,
		Metadata: &rtApi.ContainerMetadata{
			Name:    c.Metadata.Name,
			Attempt: c.Metadata.Attempt,
		},
		Labels:      c.Labels,
		Annotations: c.Annotations,
	}

	response := rtApi.ContainerStats{
		Cpu:           &cpu,
		Memory:        &memory,
		WritableLayer: &disk,
		Attributes:    &attribs,
	}

	return &response, nil
}

func toCriContainer(c *lxf.Container) *rtApi.Container {
	return &rtApi.Container{
		Id:           c.ID,
		PodSandboxId: c.SandboxID(),
		Image:        &rtApi.ImageSpec{Image: c.Image},
		ImageRef:     c.Image,
		CreatedAt:    c.CreatedAt.UnixNano(),
		State:        stateContainerAsCri(c.StateName),
		Metadata: &rtApi.ContainerMetadata{
			Name:    c.Metadata.Name,
			Attempt: c.Metadata.Attempt,
		},
		Labels:      c.Labels,
		Annotations: c.Annotations,
	}
}

func stateContainerAsCri(s lxf.ContainerStateName) rtApi.ContainerState {
	return rtApi.ContainerState(
		rtApi.ContainerState_value["CONTAINER_"+strings.ToUpper(s.String())])
}

func stateSandboxAsCri(s lxf.SandboxState) rtApi.PodSandboxState {
	return rtApi.PodSandboxState(
		rtApi.PodSandboxState_value["SANDBOX_"+strings.ToUpper(s.String())])
}

func nameSpaceOptionToString(no rtApi.NamespaceMode) string {
	return strings.ToLower(no.String())
}

func stringToNamespaceOption(s string) rtApi.NamespaceMode {
	return rtApi.NamespaceMode(rtApi.NamespaceMode_value[strings.ToUpper(s)])
}

// CompareFilterMap allows comparing two string maps
func CompareFilterMap(base map[string]string, filter map[string]string) bool {
	if filter == nil { // filter can be nil
		return true
	}

	for key := range filter {
		if base[key] != filter[key] {
			return false
		}
	}

	return true
}

// getLXDConfigPath tries to find the remote configuration file path
func getLXDConfigPath(cfg *Config) (string, error) {
	configPath := cfg.LXDRemoteConfig

	if cfg.LXDRemoteConfig == "" {
		// Equality to lxd expected, github.com/lxc/lxd/lxc/main.go:56
		var configDir string

		switch {
		case os.Getenv("LXD_CONF") != "":
			configDir = os.Getenv("LXD_CONF")
		case os.Getenv("HOME") != "":
			configDir = path.Join(os.Getenv("HOME"), ".config", "lxc")
		default:
			user, err := user.Current()
			if err != nil {
				return "", err
			}

			configDir = path.Join(user.HomeDir, ".config", "lxc")
		}

		configPath = os.ExpandEnv(path.Join(configDir, "config.yml"))
	}

	return configPath, nil
}

func (s RuntimeServer) stopContainers(sb *lxf.Sandbox) error {
	cl, err := sb.Containers()
	if err != nil {
		return err
	}

	for _, c := range cl {
		err := s.stopContainer(c, 30)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s RuntimeServer) stopContainer(c *lxf.Container, timeout int) error {
	// if container is not running, no stopping needed
	if c.StateName != lxf.ContainerStateRunning {
		return nil
	}

	err := c.Stop(timeout)
	if err != nil {
		if shared.IsErrNotFound(err) {
			return nil
		}

		return err
	}

	return nil
}

func (s RuntimeServer) deleteContainers(ctx context.Context, sb *lxf.Sandbox) error {
	cl, err := sb.Containers()
	if err != nil {
		return err
	}

	for _, c := range cl {
		err = s.deleteContainer(ctx, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s RuntimeServer) deleteContainer(ctx context.Context, c *lxf.Container) error {
	err := c.Delete()
	if err != nil {
		if shared.IsErrNotFound(err) {
			return nil
		}

		return err
	}

	sb, err := c.Sandbox()
	if err != nil {
		return err
	}

	// remove network
	if sb.NetworkConfig.Mode != lxf.NetworkHost {
		podNet, err := s.network.PodNetwork(sb.ID, sb.Annotations)
		if err == nil { // force cleanup, we don't care about error, but only enter if there's no error
			contNet, err := podNet.ContainerNetwork(c.ID, c.Annotations)
			if err == nil { // dito
				_ = contNet.WhenDeleted(ctx, &network.Properties{Data: sb.NetworkConfig.ModeData})
			}
		}
	}

	return nil
}

var NetworkSetupTimeout = 30 * time.Second

// ContainerStarted implements lxf.EventHandler interface
func (s RuntimeServer) ContainerStarted(c *lxf.Container) error {
	sb, err := c.Sandbox()
	if err != nil {
		return err
	}

	if sb.NetworkConfig.Mode != lxf.NetworkHost { // nolint: nestif
		st, err := c.State()
		if err != nil {
			return err
		}

		podNet, err := s.network.PodNetwork(sb.ID, sb.Annotations)
		if err != nil {
			return fmt.Errorf("can't enter pod network context: %w", err)
		}

		contNet, err := podNet.ContainerNetwork(c.ID, c.Annotations)
		if err != nil {
			return fmt.Errorf("can't enter container network context: %w", err)
		}

		ctx, _ := context.WithTimeout(context.Background(), NetworkSetupTimeout)

		res, err := contNet.WhenStarted(ctx, &network.PropertiesRunning{
			Properties: network.Properties{
				Data: sb.NetworkConfig.ModeData,
			},
			Pid: st.Pid,
		})
		if err != nil {
			return fmt.Errorf("can't start container network: %w", err)
		}

		err = s.handleNetworkResult(sb, res)
		if err != nil {
			return fmt.Errorf("unable to save create container network result: %w", err)
		}
	}

	return nil
}

// ContainerStopped implements lxf.EventHandler interface
func (s *RuntimeServer) ContainerStopped(c *lxf.Container) error {
	sb, err := c.Sandbox()
	if err != nil {
		return err
	}

	// stop network
	if sb.NetworkConfig.Mode != lxf.NetworkHost {
		podNet, err := s.network.PodNetwork(sb.ID, sb.Annotations)
		if err == nil { // force cleanup, we don't care about error, but only enter if there's no error
			contNet, err := podNet.ContainerNetwork(c.ID, c.Annotations)
			if err == nil { // dito
				ctx, _ := context.WithTimeout(context.Background(), NetworkSetupTimeout)
				_ = contNet.WhenStopped(ctx, &network.Properties{Data: sb.NetworkConfig.ModeData})
			}
		}
	}

	return nil
}

func (s *RuntimeServer) handleNetworkResult(sb *lxf.Sandbox, res *network.Result) error {
	if res != nil {
		if len(res.Data) > 0 {
			sb.NetworkConfig.ModeData = res.Data
		}

		for _, n := range res.Nics {
			n := n
			sb.Devices.Upsert(&n)
		}

		sb.CloudInitNetworkConfigEntries = append(sb.CloudInitNetworkConfigEntries, res.NetworkConfigEntries...)

		return sb.Apply()
	}

	return nil
}
