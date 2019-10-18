package network

import (
	"context"
	"encoding/json"
	"net"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"emperror.dev/errors"
	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

const (
	defaultCNIbinPath   = "/opt/cni/bin"
	defaultCNIconfPath  = "/etc/cni/net.d"
	defaultCNINetNSPath = "/run/netns"
)

// ConfCNI are configuration options for the cni plugin. All properties are optional and get a default value
type ConfCNI struct {
	BinPath   string
	ConfPath  string
	NetNSPath string
}

func (c *ConfCNI) setDefaults() {
	if c.BinPath == "" {
		c.BinPath = defaultCNIbinPath
	}
	if c.ConfPath == "" {
		c.ConfPath = defaultCNIconfPath
	}
	if c.NetNSPath == "" {
		c.NetNSPath = defaultCNINetNSPath
	}
}

// cniPlugin manages the pod networks using CNI
type cniPlugin struct {
	NetworkPlugin
	cni  libcni.CNI
	conf ConfCNI
}

// InitPluginCNI instantiates the cni plugin using the provided config
func InitPluginCNI(conf ConfCNI) (*cniPlugin, error) {
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
	runtimeConf, err := c.getCNIRuntimeConf(namespace, name, id)
	if err != nil {
		return nil, err
	}
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
func (c *cniPlugin) getCNIRuntimeConf(namespace string, name string, id string) (*libcni.RuntimeConf, error) {
	netNSPath := c.conf.NetNSPath
	return &libcni.RuntimeConf{
		ContainerID: id,
		NetNS:       filepath.Join(netNSPath, id),
		IfName:      DefaultInterface,
		Args:        [][2]string{
			// Removed, as they all seem to have no purpose
			// {"IgnoreUnknown", "1"},
			// {"K8S_POD_NAMESPACE", namespace},
			// {"K8S_POD_NAME", name},
			// {"K8S_POD_INFRA_CONTAINER_ID", id},
		},
	}, nil
}

// cniPodNetwork is a pod network environment context. It handles the pod network by creating a netns and interfaces for
// the pod
type cniPodNetwork struct {
	plugin      *cniPlugin
	netList     *libcni.NetworkConfigList
	runtimeConf *libcni.RuntimeConf
	annotations map[string]string
}

// Setup creates the network for that pod and the result is saved
func (c *cniPodNetwork) Setup(ctx context.Context) ([]byte, error) {
	out, err := exec.Command("ip", "netns", "add", c.runtimeConf.ContainerID).CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(out))
	}

	addResult, err := c.plugin.cni.AddNetworkList(ctx, c.netList, c.runtimeConf)
	if err != nil {
		return nil, err
	}
	result, err := current.NewResultFromResult(addResult)
	if err != nil {
		return nil, err
	}
	// TODO: since upgrading to cni v0.7.1 we could possibly skip saving the result for ourselves...
	return json.Marshal(result)
}

// Teardown removes the network for that pod
func (c *cniPodNetwork) Teardown(ctx context.Context) error {
	err := c.plugin.cni.DelNetworkList(ctx, c.netList, c.runtimeConf)
	// if the netns is not exists retry without netns path
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		c.runtimeConf.NetNS = ""
		err = c.plugin.cni.DelNetworkList(ctx, c.netList, c.runtimeConf)
	}
	return err
}

// Status reports IP and any error with the network of that pod
func (c *cniPodNetwork) Status(ctx context.Context) (*PodNetworkStatus, error) {
	err := c.plugin.cni.CheckNetworkList(ctx, c.netList, c.runtimeConf)
	if err != nil {
		return nil, err
	}

	out, err := exec.Command("ip", "netns", "exec", c.runtimeConf.ContainerID, "-br", "-o", "-4", "addr", "show", "dev", DefaultInterface, "scope", "global").CombinedOutput()
	//out, err := exec.Command("nsenter", fmt.Sprintf("--net=%s", filepath.Join(c.runtimeConf.NetNS)), "-F", "--", "ip", "-br", "-o", "-4", "addr", "show", "dev", DefaultInterface, "scope", "global").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, string(out))
	}

	ips, err := parseIPAddrShow(out)
	if err != nil {
		return nil, err
	}

	return &PodNetworkStatus{
		IPs: ips,
	}, nil
}

func parseIPAddrShow(output []byte) ([]net.IP, error) {
	var ips []net.IP
	lines := strings.Split(string(output), "\n")
	for _, l := range lines {
		f := strings.Fields(l)
		if len(f) < 3 {
			return nil, errors.Errorf("Unexpected address output %s ", l)
		}
		ip, _, err := net.ParseCIDR(f[2])
		if err != nil {
			return nil, errors.Errorf("CNI failed to parse ip from output %s due to %v", output, err)
		}
		ips = append(ips, ip)
	}
	return ips, nil
}
