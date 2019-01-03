package network

// CNI PoC for LXE
// TODO should support plugin architecture see docker-shim

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/020"
	"github.com/lxc/lxd/shared/logger"
)

const (
	defaultCNIbinPath  = "/opt/cni/bin"
	defaultCNIconfPath = "/etc/cni/net.d"

	// DefaultInterface for all containers
	DefaultInterface = "eth0"

	// DefaultCNIInterface the default CNI Interface (usually a bridge)
	DefaultCNIInterface = "cni0"
)

type cniNetwork struct {
	name          string
	NetworkConfig *libcni.NetworkConfigList
	CNIConfig     libcni.CNI
}

// getDefaultCNINetwork is borrowed from k8s' dockershim
func getDefaultCNINetwork(confDir string, binDirs []string) (*cniNetwork, error) {
	files, err := libcni.ConfFiles(confDir, []string{".conf", ".conflist", ".json"})
	switch {
	case err != nil:
		return nil, err
	case len(files) == 0:
		return nil, fmt.Errorf("No networks found in %s", confDir)
	}

	sort.Strings(files)
	for _, confFile := range files {
		var confList *libcni.NetworkConfigList
		if strings.HasSuffix(confFile, ".conflist") {
			confList, err = libcni.ConfListFromFile(confFile)
			if err != nil {
				logger.Errorf("Error loading CNI config list file %s: %v", confFile, err)
				continue
			}
		} else {
			conf, err := libcni.ConfFromFile(confFile)
			if err != nil {
				logger.Errorf("Error loading CNI config file %s: %v", confFile, err)
				continue
			}
			// Ensure the config has a "type" so we know what plugin to run.
			// Also catches the case where somebody put a conflist into a conf file.
			if conf.Network.Type == "" {
				logger.Errorf("Error loading CNI config file %s: no 'type'; perhaps this is a .conflist?", confFile)
				continue
			}

			confList, err = libcni.ConfListFromConf(conf)
			if err != nil {
				logger.Errorf("Error converting CNI config file %s to list: %v", confFile, err)
				continue
			}

		}
		if len(confList.Plugins) == 0 {
			logger.Errorf("CNI config list %s has no networks, skipping", confFile)
			continue
		}

		network := &cniNetwork{
			name:          confList.Name,
			NetworkConfig: confList,
			CNIConfig:     &libcni.CNIConfig{Path: binDirs},
		}
		return network, nil
	}
	return nil, fmt.Errorf("No valid networks found in %s", confDir)
}

// AttachCNIInterface will setup the Pod Networking using CNI
func AttachCNIInterface(namespace string, sandboxname string, containerID string, processID int64) ([]byte, error) {

	cniNetwork, err := getDefaultCNINetwork(defaultCNIconfPath, []string{defaultCNIbinPath})
	if err != nil {
		return nil, err
	}

	cni := cniNetwork.CNIConfig
	net := cniNetwork.NetworkConfig.Plugins[0]

	podNSPath := fmt.Sprintf("/proc/%s/ns/net", strconv.FormatInt(processID, 10))

	rt := &libcni.RuntimeConf{
		ContainerID: containerID,
		NetNS:       podNSPath,
		IfName:      DefaultInterface,
		Args: [][2]string{
			{"IgnoreUnknown", "1"},
			{"K8S_POD_NAMESPACE", namespace},
			{"K8S_POD_NAME", sandboxname},
			{"K8S_POD_INFRA_CONTAINER_ID", containerID},
		},
	}

	addNetworkResult, err := cni.AddNetwork(net, rt)
	if err != nil {
		return nil, err
	}

	result, err := types020.GetResult(addNetworkResult)
	if err != nil {
		return nil, err
	}

	resultstr, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return resultstr, nil
}

// DetachCNIInterface will teardown the Pod Networking
func DetachCNIInterface() {
}

// Status of the Pod Networking
func Status() {
}

// Name of the Pod Networking
func Name() {
}

func cmdExec(cmd string, args ...string) error {
	c := exec.Command(cmd, args...) // nolint: gosec
	logger.Debugf("%v", c.Args)
	err := c.Run()
	if err != nil {
		return fmt.Errorf("%s: %v", err.Error(), c)
	}
	return nil
}
