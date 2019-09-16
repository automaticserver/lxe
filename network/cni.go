package network

// CNI PoC for LXE
// TODO should support plugin architecture see docker-shim

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/containernetworking/cni/libcni"
	types "github.com/containernetworking/cni/pkg/types/current"
	"github.com/lxc/lxd/shared/logger"
)

const (
	defaultCNIbinPath  = "/opt/cni/bin"
	defaultCNIconfPath = "/etc/cni/net.d"

	// DefaultInterface for all containers
	DefaultInterface = "eth0"
)

// getDefaultCNINetwork is borrowed from k8s' dockershim
func getDefaultCNINetwork(confDir string, binDirs []string) (libcni.CNI, *libcni.NetworkConfigList, error) {
	files, err := libcni.ConfFiles(confDir, []string{".conf", ".conflist", ".json"})
	switch {
	case err != nil:
		return nil, nil, err
	case len(files) == 0:
		return nil, nil, fmt.Errorf("No networks found in %s", confDir)
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

		return &libcni.CNIConfig{Path: binDirs}, confList, nil
	}
	return nil, nil, fmt.Errorf("No valid networks found in %s", confDir)
}

// AttachCNIInterface will setup the Pod Networking using CNI
func AttachCNIInterface(namespace string, sandboxname string, containerID string, processID int64) ([]byte, error) {
	cni, netList, err := getDefaultCNINetwork(defaultCNIconfPath, []string{defaultCNIbinPath})
	if err != nil {
		return nil, err
	}

	rt, err := getRuntimeConf(namespace, sandboxname, containerID, processID)
	if err != nil {
		return nil, err
	}

	addNetworkResult, err := cni.AddNetworkList(context.TODO(), netList, rt)
	if err != nil {
		return nil, err
	}

	result, err := types.NewResultFromResult(addNetworkResult)
	if err != nil {
		return nil, err
	}

	// TODO: since upgrading to cni v0.7.1 we could possibly skip saving the result for ourselves...

	resultstr, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return resultstr, nil
}

// DetachCNIInterface will teardown the Pod Networking
func DetachCNIInterface(namespace string, sandboxname string, containerID string, processID int64) error {
	cni, netList, err := getDefaultCNINetwork(defaultCNIconfPath, []string{defaultCNIbinPath})
	if err != nil {
		return err
	}

	rt, err := getRuntimeConf(namespace, sandboxname, containerID, processID)
	if err != nil {
		return err
	}

	return cni.DelNetworkList(context.TODO(), netList, rt)
}

// getRuntimeConf returns common libcni runtime conf used to interact with the cni binaries. If the processID is empty,
// NetNS will be empty, which is useful for deletion when the container doesn't exist anymore (see
// https://github.com/containernetworking/cni/pull/230)
func getRuntimeConf(namespace string, sandboxname string, containerID string, processID int64) (*libcni.RuntimeConf, error) {
	podNSPath := ""
	if processID > 0 {
		podNSPath = fmt.Sprintf("/proc/%s/ns/net", strconv.FormatInt(processID, 10))
	}

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

	return rt, nil
}

// Status of the Pod Networking
func Status() {
}

// Name of the Pod Networking
func Name() {
}
