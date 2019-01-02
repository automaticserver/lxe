package network

// CNI PoC for LXE
// TODO should support plugin architecture see docker-shim

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
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

// ReattachCNIInterface attaches a interface to a container using the result of the previous configuration
func ReattachCNIInterface(namespace string, sandboxname string, containerID string, processID int64, prevConf string) ([]byte, error) {
	prevResult := new(types020.Result)
	err := json.Unmarshal([]byte(prevConf), &prevResult)
	if err != nil {
		logger.Errorf("unable to unmarshal cniJSON %v: %v", prevConf, err)
	}

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

	net, err = libcni.InjectConf(net, map[string]interface{}{
		"prevResult": prevResult,
	})
	if err != nil {
		return nil, err
	}

	logger.Infof("DEBUG cniConf: %+v | %+v | %+v", *net.Network, string(net.Bytes), rt)

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

// ReattachCNIInterfaceOld attaches a interface to a container using the provided configuration
// namespace and sandboxname is not in use, but stays for now to match AttachCNIInterface
func ReattachCNIInterfaceOld(namespace string, sandboxname string, containerID string, processID int64, previousconf string) error {
	cniResult202 := new(types020.Result)
	err := json.Unmarshal([]byte(previousconf), &cniResult202)
	if err != nil {
		logger.Errorf("unable to unmarshal cniJSON %v: %v", previousconf, err)
	}

	logger.Debugf("reattaching CNI interface to container: %v", containerID)

	// see:
	// lxd source: container_lxc.go // func (c *containerLXC) createNetworkDevice(name string, m types.Device) (string, error)
	// https://github.com/p8952/bocker/blob/master/bocker
	// https://stackoverflow.com/questions/31265993/docker-networking-namespace-not-visible-in-ip-netns-list

	// create netDSDir if not present
	netNSDir := "/var/run/netns"
	if _, err := os.Stat(netNSDir); os.IsNotExist(err) {
		err = os.MkdirAll(netNSDir, 0750)
		if err != nil {
			return err
		}
	}

	// create symlink from container NS to default NS directory
	// ln -sfT /proc/$pid/ns/net /var/run/netns/$container_id
	err = cmdExec("ln", "-sfT",
		fmt.Sprintf("/proc/%d/ns/net", processID),
		fmt.Sprintf("/var/run/netns/%s", containerID))
	if err != nil {
		return err
	}
	// remove containers symlink from default NS directory
	defer os.Remove(fmt.Sprintf("/var/run/netns/%s", containerID)) //nolint: errcheck

	vethHost := "veth" + hex.EncodeToString([]byte(containerID))[0:6]
	vethCnt := DefaultInterface

	// ip link add dev veth0_$ID type veth peer name veth1_$ID
	err = cmdExec("ip", "link", "add", "dev", vethHost, "type", "veth", "peer", "name", vethCnt)
	if err != nil {
		return err
	}

	// ip link set dev veth0_$ID up
	err = cmdExec("ip", "link", "set", "dev", vethHost, "up")
	if err != nil {
		return err
	}

	// ip link set veth0_$ID master cni0
	err = cmdExec("ip", "link", "set", vethHost, "master", DefaultCNIInterface)
	if err != nil {
		return err
	}

	// ip link set veth1_$ID netns $ID
	err = cmdExec("ip", "link", "set", vethCnt, "netns", containerID)
	if err != nil {
		return err
	}

	ipAddr := cniResult202.IP4.IP.String()
	// ip netns exec netns_"$uuid" ip addr add 10.0.0."$ip"/24 dev veth1_"$uuid"
	err = cmdExec("ip", "netns", "exec", containerID, "ip", "addr", "add", ipAddr, "dev", vethCnt)
	if err != nil {
		return err
	}

	// ip netns exec netns_"$uuid" ip link set dev veth1_"$uuid" up
	err = cmdExec("ip", "netns", "exec", containerID, "ip", "link", "set", "dev", vethCnt, "up")
	if err != nil {
		return err
	}

	// ip netns exec netns_"$uuid" ip route add default via 10.0.0.1
	// default route is among those routes
	for _, r := range cniResult202.IP4.Routes {
		err = cmdExec("ip", "netns", "exec", containerID, "ip", "route", "add", r.Dst.String(), "via", r.GW.String())
		if err != nil {
			return err
		}
	}

	return nil
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
