package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

const (
	defaultCNIbinPath   = "/opt/cni/bin"
	defaultCNIconfPath  = "/etc/cni/net.d"
	defaultCNInetnsPath = "/run/netns"
)

// ConfCNI are configuration options for the cni plugin. All properties are optional and get a default value
type ConfCNI struct {
	BinPath   string
	ConfPath  string
	NetnsPath string
}

func (c *ConfCNI) setDefaults() {
	if c.BinPath == "" {
		c.BinPath = defaultCNIbinPath
	}

	if c.ConfPath == "" {
		c.ConfPath = defaultCNIconfPath
	}

	if c.NetnsPath == "" {
		c.NetnsPath = defaultCNInetnsPath
	}
}

// cniPlugin manages the pod networks using CNI
type cniPlugin struct {
	Plugin
	cni  libcni.CNI
	conf ConfCNI
}

// InitPluginCNI instantiates the cni plugin using the provided config
func InitPluginCNI(conf ConfCNI) (*cniPlugin, error) { // nolint: golint // intended to not export cniPlugin
	conf.setDefaults()

	return &cniPlugin{
		cni:  libcni.NewCNIConfig([]string{conf.BinPath}, nil),
		conf: conf,
	}, nil
}

// PodNetwork enters a pod network environment context
func (c *cniPlugin) PodNetwork(id string, annotations map[string]string) (PodNetwork, error) {
	netList, warnings, err := c.getCNINetworkConfig()
	if err != nil {
		return nil, errors.Append(err, warnings)
	}

	runtimeConf := c.getCNIRuntimeConf(id)

	return &cniPodNetwork{
		plugin:      c,
		netList:     netList,
		runtimeConf: runtimeConf,
		annotations: annotations,
	}, nil
}

// Status returns error if the plugin is in error state
func (c *cniPlugin) Status() error {
	return nil
}

// getCNINetworkConfig looks into the cni configuration dir for configs to load
func (c *cniPlugin) getCNINetworkConfig() (*libcni.NetworkConfigList, error, error) {
	confDir := c.conf.ConfPath

	files, err := libcni.ConfFiles(confDir, []string{".conf", ".conflist", ".json"})

	switch {
	case err != nil:
		return nil, nil, err
	case len(files) == 0:
		return nil, nil, errors.Errorf("No networks found in %s", confDir)
	}

	var warnings error

	sort.Strings(files)

	for _, confFile := range files {
		var confList *libcni.NetworkConfigList
		if strings.HasSuffix(confFile, ".conflist") {
			confList, err = libcni.ConfListFromFile(confFile)
			if err != nil {
				warnings = errors.Append(warnings, errors.Wrapf(err, "Error loading CNI config list file %s", confFile))
				continue
			}
		} else {
			conf, err := libcni.ConfFromFile(confFile)
			if err != nil {
				warnings = errors.Append(warnings, errors.Wrapf(err, "Error loading CNI config file %s", confFile))
				continue
			}
			// Ensure the config has a "type" so we know what plugin to run.
			// Also catches the case where somebody put a conflist into a conf file.
			if conf.Network.Type == "" {
				warnings = errors.Append(warnings, errors.Errorf("Error loading CNI config file %s: no 'type'; perhaps this is a .conflist?", confFile))
				continue
			}

			confList, err = libcni.ConfListFromConf(conf)
			if err != nil {
				warnings = errors.Append(warnings, errors.Wrapf(err, "Error converting CNI config file %s to list", confFile))
				continue
			}
		}

		if len(confList.Plugins) == 0 {
			warnings = errors.Append(warnings, errors.Errorf("CNI config list %s has no networks, skipping", confFile))
			continue
		}

		return confList, warnings, nil
	}

	return nil, warnings, errors.Errorf("No valid networks found in %s", confDir)
}

// getRuntimeConf returns common libcni runtime conf used to interact with the cni
func (c *cniPlugin) getCNIRuntimeConf(id string) *libcni.RuntimeConf {
	return &libcni.RuntimeConf{
		ContainerID: id,
		NetNS:       "",
		IfName:      DefaultInterface,
		Args:        [][2]string{
			// Removed, as they all seem to have no purpose
			// {"IgnoreUnknown", "1"},
			// {"K8S_POD_NAMESPACE", namespace},
			// {"K8S_POD_NAME", name},
			// {"K8S_POD_INFRA_CONTAINER_ID", id},
		},
	}
}

// cniPodNetwork is a pod network environment context. It handles the pod network by creating a netns and interfaces for
// the pod
type cniPodNetwork struct {
	DeviceHandler // noop
	plugin        *cniPlugin
	netList       *libcni.NetworkConfigList
	runtimeConf   *libcni.RuntimeConf
	annotations   map[string]string
}

// SetupPid creates the network for that pod. pid is the process id of the pod. The retured result bytes are provided
// for the other calls of this PodNetwork.
func (c *cniPodNetwork) SetupPid(_ context.Context, _ int64) ([]byte, error) {
	// TODO: As long as we haven't figured out to do 1:n podnetwork:container this method does nothing
	return nil, nil
}

// TeardownPid removes the network of that pod. Must tear down networking as good as possible, an error will only be
// logged and doesn't stop execution of further statements.
func (c *cniPodNetwork) TeardownPid(_ context.Context, _ []byte) error {
	// TODO: As long as we haven't figured out to do 1:n podnetwork:container this method does nothing
	return nil
}

// AttachPid attaches a container to the pod network. pid is the process id of the container. Can return arbitrary
// bytes or nic devices or both. The retured result bytes replace the existing one if provided.
func (c *cniPodNetwork) AttachPid(ctx context.Context, _ []byte, pid int64) ([]byte, error) {
	if pid == 0 {
		cID := c.runtimeConf.ContainerID

		out, err := exec.Command("ip", "netns", "add", cID).CombinedOutput()
		if err != nil {
			return nil, errors.Wrap(err, string(out))
		}

		defer exec.Command("ip", "netns", "delete", cID) // nolint: errcheck

		c.runtimeConf.NetNS = filepath.Join(c.plugin.conf.NetnsPath, cID)
	} else {
		c.runtimeConf.NetNS = fmt.Sprintf("/proc/%s/ns/net", strconv.FormatInt(pid, 10))
	}

	result, err := c.plugin.cni.AddNetworkList(ctx, c.netList, c.runtimeConf)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

// DetachPid detaches a container from the pod network. Must detach networking as good as possible, an error will only
// be logged and doesn't stop execution of further statements.
func (c *cniPodNetwork) DetachPid(ctx context.Context, _ []byte) error {
	c.runtimeConf.NetNS = ""
	return c.plugin.cni.DelNetworkList(ctx, c.netList, c.runtimeConf)
}

// StatusPid reports IP and any error with the network of that pod. Bytes can be nil if LXE thinks it never ran Setup
// and thus also pid is not set or weren't returned yet.
func (c *cniPodNetwork) StatusPid(ctx context.Context, previousresult []byte, _ int64) (*PodNetworkStatus, error) {
	if previousresult == nil {
		previousresult = []byte{}
	}

	resultInit, err := current.NewResult(previousresult)
	if err != nil {
		return nil, err
	}

	result, err := current.GetResult(resultInit)
	if err != nil {
		return nil, err
	}

	if len(result.IPs) == 0 {
		return nil, fmt.Errorf("no ip address found for %v", c.runtimeConf.ContainerID)
	}

	return &PodNetworkStatus{
		IPs: []net.IP{
			result.IPs[0].Address.IP,
		},
	}, nil
}
