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
	// fieldLXEBridge is the key name to specify the bridge to be used as parent
	// TODO: to be removed once specifyable with CNI
	fieldLXEBridge = "x-lxe-bridge"
	// fieldLXEAdditionalLXDConfig is the name of the field which contains various additional lxd config options
	fieldLXEAdditionalLXDConfig = "x-lxe-additional-lxd-config"
)

// LXDAdditionalConfig contains additional config options not present in PodSpec
// Key names and values must match the key names specified by LXD
type AdditionalLXDConfig map[string]string

func toCriStatusResponse(ct *lxf.Container) *rtApi.ContainerStatusResponse {
	status := rtApi.ContainerStatus{
		Metadata: &rtApi.ContainerMetadata{
			Name:    ct.Metadata.Name,
			Attempt: uint32(ct.Metadata.Attempt),
		},
		State: rtApi.ContainerState(
			rtApi.ContainerState_value["CONTAINER_"+strings.ToUpper(ct.State.String())]),
		CreatedAt:   ct.CreatedAt.UnixNano(),
		StartedAt:   ct.StartedAt.UnixNano(),
		FinishedAt:  ct.FinishedAt.UnixNano(),
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
	return &rtApi.Container{
		Id:           ct.ID,
		PodSandboxId: ct.Sandbox.ID,
		CreatedAt:    ct.CreatedAt.UnixNano(),
		State: rtApi.ContainerState(
			rtApi.ContainerState_value["CONTAINER_"+strings.ToUpper(ct.State.String())]),
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
