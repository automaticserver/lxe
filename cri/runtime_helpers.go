package cri

import (
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxe/lxf"
	rtApi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

func toCriStatusResponse(c *lxf.Container) *rtApi.ContainerStatusResponse {
	status := rtApi.ContainerStatus{
		Metadata: &rtApi.ContainerMetadata{
			Name:    c.Metadata.Name,
			Attempt: uint32(c.Metadata.Attempt),
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

	for _, d := range c.Disks {
		status.Mounts = append(status.Mounts, &rtApi.Mount{
			ContainerPath:  d.Path,
			HostPath:       d.Source,
			Readonly:       d.Readonly,
			SelinuxRelabel: false, // though don't know what this means
			Propagation:    rtApi.MountPropagation_PROPAGATION_PRIVATE,
		})
	}

	for _, b := range c.Blocks {
		status.Mounts = append(status.Mounts, &rtApi.Mount{
			ContainerPath:  b.Path,
			HostPath:       b.Source,
			Readonly:       false,                                                // probably always that?
			SelinuxRelabel: false,                                                // though don't know what this means
			Propagation:    rtApi.MountPropagation_PROPAGATION_HOST_TO_CONTAINER, // unsure
		})
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
			Mountpoint: path.Join(shared.VarPath("container"), c.ID, "rootfs"),
		},
		UsedBytes:  &rtApi.UInt64Value{Value: st.Stats.FilesystemUsage}, // TODO: root seems not visible? or does it depend?
		InodesUsed: &rtApi.UInt64Value{Value: 0},                        // TODO: do we have to find out?
	}
	attribs := rtApi.ContainerAttributes{
		Id: c.ID,
		Metadata: &rtApi.ContainerMetadata{
			Name:    c.Metadata.Name,
			Attempt: uint32(c.Metadata.Attempt),
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
		PodSandboxId: c.Profiles[0],
		Image:        &rtApi.ImageSpec{Image: c.Image},
		ImageRef:     c.Image,
		CreatedAt:    c.CreatedAt.UnixNano(),
		State:        stateContainerAsCri(c.StateName),
		Metadata: &rtApi.ContainerMetadata{
			Name:    c.Metadata.Name,
			Attempt: uint32(c.Metadata.Attempt),
		},
		Labels:      c.Labels,
		Annotations: c.Annotations,
	}
	// TODO: more fields?
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
func getLXDConfigPath(cfg *LXEConfig) (string, error) {
	configPath := cfg.LXDRemoteConfig
	if cfg.LXDRemoteConfig == "" {
		// Copied from github.com/lxc/lxd/lxc/main.go:56, since there it is unexported
		var configDir string
		if os.Getenv("LXD_CONF") != "" {
			configDir = os.Getenv("LXD_CONF")
		} else if os.Getenv("HOME") != "" {
			configDir = path.Join(os.Getenv("HOME"), ".config", "lxc")
		} else {
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
		if lxf.IsContainerNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

func (s RuntimeServer) deleteContainers(sb *lxf.Sandbox) error {
	cl, err := sb.Containers()
	if err != nil {
		return err
	}
	for _, c := range cl {
		err = s.deleteContainer(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s RuntimeServer) deleteContainer(c *lxf.Container) error {
	err := c.Delete()
	if err != nil {
		if lxf.IsContainerNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}
