package cri

import (
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/lxc/lxe/lxf"
	rtApi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

const (
	// The following fields are not covered by PodSpec but sometimes required to be defined

	// fieldLXEBridge is the key name to specify the bridge to be used as parent
	fieldLXEBridge = "x-lxe-bridge"
	// fieldLXENamespaces is the key name to specify namespaces
	fieldLXENamespaces = "x-lxe-namespaces"
	// fieldLXEKernelModules is the key name to specify kernel modules
	fieldLXEKernelModules = "x-lxe-kernel-modules"
	// fieldLXENesting is the key name to specify to allow nesting
	fieldLXENesting = "x-lxe-nesting"
	// fieldLXERawMounts is the key name to add raw mounts to the lxc config
	fieldLXERawMounts = "x-lxe-raw-mounts"
)

func toCriStatusResponse(ct *lxf.Container) *rtApi.ContainerStatusResponse {
	key := "CONTAINER_" + strings.ToUpper(string(ct.State))
	status := rtApi.ContainerStatus{
		Metadata: &rtApi.ContainerMetadata{
			Name:    ct.Metadata.Name,
			Attempt: uint32(ct.Metadata.Attempt),
		},
		State:       rtApi.ContainerState(rtApi.ContainerState_value[key]),
		CreatedAt:   ct.CreatedAt.UnixNano(),
		StartedAt:   ct.StartedAt.UnixNano(),
		Id:          ct.ID,
		Labels:      ct.Labels,
		Annotations: ct.Annotations,
		Image:       &rtApi.ImageSpec{Image: ct.Image},
		ImageRef:    ct.Image,
	}

	return &rtApi.ContainerStatusResponse{
		Status: &status,
		Info:   map[string]string{},
	}
}

func toCriStats(lxdStats *lxf.Container) *rtApi.ContainerStats {
	now := time.Now().UnixNano()

	cpu := rtApi.CpuUsage{
		Timestamp:            now,
		UsageCoreNanoSeconds: &rtApi.UInt64Value{Value: lxdStats.Stats.CPUUsage},
	}
	memory := rtApi.MemoryUsage{
		Timestamp:       now,
		WorkingSetBytes: &rtApi.UInt64Value{Value: lxdStats.Stats.MemoryUsage},
	}
	disk := rtApi.FilesystemUsage{
		Timestamp: now,
		UsedBytes: &rtApi.UInt64Value{Value: lxdStats.Stats.FilesystemUsage},
	}
	attribs := rtApi.ContainerAttributes{
		Id: lxdStats.ID,
		Metadata: &rtApi.ContainerMetadata{
			Name:    lxdStats.Metadata.Name,
			Attempt: uint32(lxdStats.Metadata.Attempt),
		},
		Labels:      lxdStats.Labels,
		Annotations: lxdStats.Annotations,
	}

	response := rtApi.ContainerStats{
		Cpu:           &cpu,
		Memory:        &memory,
		WritableLayer: &disk,
		Attributes:    &attribs,
	}
	return &response
}

func toCriContainer(ct *lxf.Container) *rtApi.Container {
	stateKey := "CONTAINER_" + strings.ToUpper(string(ct.State))
	state := rtApi.ContainerState(rtApi.ContainerState_value[stateKey])

	return &rtApi.Container{
		Id:           ct.ID,
		PodSandboxId: ct.Sandbox.ID,
		CreatedAt:    ct.CreatedAt.UnixNano(),
		State:        state,
		Metadata: &rtApi.ContainerMetadata{
			Name:    ct.Metadata.Name,
			Attempt: uint32(ct.Metadata.Attempt),
		},
		Labels:      ct.Labels,
		Annotations: ct.Annotations,
	}
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
