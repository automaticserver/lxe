package network

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

const (
	defaultCNIbinPath  = "/opt/cni/bin"
	defaultCNIconfPath = "/etc/cni/net.d"
)

// ConfCNI are configuration options for the cni plugin. All properties are optional and get a default value
type ConfCNI struct {
	BinPath  string
	ConfPath string
}

func (c *ConfCNI) setDefaults() {
	if c.BinPath == "" {
		c.BinPath = defaultCNIbinPath
	}

	if c.ConfPath == "" {
		c.ConfPath = defaultCNIconfPath
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
func (c *cniPlugin) PodNetwork(namespace, name, id string, annotations map[string]string) (PodNetwork, error) {
	netList, warnings, err := c.getCNINetworkConfig()
	if err != nil {
		return nil, errors.Append(err, warnings)
	}

	runtimeConf := c.getCNIRuntimeConf(namespace, name, id)

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
func (c *cniPlugin) getCNIRuntimeConf(_ string, _ string, id string) *libcni.RuntimeConf {
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
	plugin      *cniPlugin
	netList     *libcni.NetworkConfigList
	runtimeConf *libcni.RuntimeConf
	annotations map[string]string
}

// Setup creates the network for that pod and the result is saved. pid is the process id of the pod.
func (c *cniPodNetwork) Setup(ctx context.Context, _ int64) ([]byte, error) {
	// TODO: As long as we haven't figured out to do 1:n podnetwork:container this method does nothing
	return nil, nil
}

// Teardown removes the network of that pod. pid is the process id of the pod, but might be missing. Must tear down
// networking as good as possible, an error will only be logged and doesn't stop execution of further statements
func (c *cniPodNetwork) Teardown(ctx context.Context, _ []byte, _ int64) error {
	// TODO: As long as we haven't figured out to do 1:n podnetwork:container this method does nothing
	return nil
}

// Attach a container to the pod network. pid is the process id of the container.
func (c *cniPodNetwork) Attach(ctx context.Context, _ []byte, pid int64) error {
	c.runtimeConf.NetNS = fmt.Sprintf("/proc/%s/ns/net", strconv.FormatInt(pid, 10))
	_, err := c.plugin.cni.AddNetworkList(ctx, c.netList, c.runtimeConf)

	return err
}

// Detach a container from the pod network. pid is the process id of the container, but might be missing. Must detach
// networking as good as possible, an error will only be logged and doesn't stop execution of further statements
func (c *cniPodNetwork) Detach(ctx context.Context, _ []byte, pid int64) error {
	c.runtimeConf.NetNS = ""
	if pid != 0 {
		c.runtimeConf.NetNS = fmt.Sprintf("/proc/%s/ns/net", strconv.FormatInt(pid, 10))
	}

	return c.plugin.cni.DelNetworkList(ctx, c.netList, c.runtimeConf)
}

// Status reports IP and any error with the network of that pod. Bytes can be nil if LXE thinks it never ran Setup and
// thus also pid is not set
func (c *cniPodNetwork) Status(ctx context.Context, _ []byte, _ int64) (*PodNetworkStatus, error) {
	resultCached, err := c.plugin.cni.GetNetworkListCachedResult(c.netList, c.runtimeConf)
	if err != nil {
		return nil, err
	}

	result, err := current.GetResult(resultCached)
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
