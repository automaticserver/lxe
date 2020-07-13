package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	DefaultCNIbinPath   = "/opt/cni/bin"
	DefaultCNIconfPath  = "/etc/cni/net.d"
	defaultCNInetnsPath = "/run/netns"
)

var (
	ErrNoUpdateRuntimeConfig = errors.New("cniPlugin can't update runtime config")
	ErrNoNetworksFound       = errors.New("No valid networks found")
)

// ConfCNI are configuration options for the cni plugin. All properties are optional and get a default value
type ConfCNI struct {
	BinPath   string
	ConfPath  string
	NetnsPath string
}

func (c *ConfCNI) setDefaults() {
	if c.BinPath == "" {
		c.BinPath = DefaultCNIbinPath
	}

	if c.ConfPath == "" {
		c.ConfPath = DefaultCNIconfPath
	}

	if c.NetnsPath == "" {
		c.NetnsPath = defaultCNInetnsPath
	}
}

// cniPlugin manages the pod networks using CNI
type cniPlugin struct {
	noopPlugin // every method not implemented is noop
	cni        libcni.CNI
	conf       ConfCNI
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
func (p *cniPlugin) PodNetwork(id string, annotations map[string]string) (PodNetwork, error) {
	netList, warnings, err := p.getCNINetworkConfig()
	if err != nil {
		return nil, fmt.Errorf("%w, %v", err, warnings)
	}

	runtimeConf := p.getCNIRuntimeConf(id)

	return &cniPodNetwork{
		plugin:      p,
		netList:     netList,
		runtimeConf: runtimeConf,
		annotations: annotations,
	}, nil
}

// UpdateRuntimeConfig is called when there are updates to the configuration which the plugin might need to apply
func (p *cniPlugin) UpdateRuntimeConfig(_ *rtApi.RuntimeConfig) error {
	return ErrNoUpdateRuntimeConfig
}

// getCNINetworkConfig looks into the cni configuration dir for configs to load
func (p *cniPlugin) getCNINetworkConfig() (*libcni.NetworkConfigList, error, error) {
	confDir := p.conf.ConfPath

	files, err := libcni.ConfFiles(confDir, []string{".conf", ".conflist", ".json"})

	switch {
	case err != nil:
		return nil, nil, err
	case len(files) == 0:
		return nil, nil, fmt.Errorf("%w in %s", ErrNoNetworksFound, confDir)
	}

	var warnings error

	sort.Strings(files)

	for _, confFile := range files {
		var confList *libcni.NetworkConfigList
		if strings.HasSuffix(confFile, ".conflist") { // nolint: nestif
			confList, err = libcni.ConfListFromFile(confFile)
			if err != nil {
				warnings = fmt.Errorf("%v: %w, error loading CNI config list file %s", warnings, err, confFile)
				continue
			}
		} else {
			conf, err := libcni.ConfFromFile(confFile)
			if err != nil {
				warnings = fmt.Errorf("%v: %w, error loading CNI config file %s", warnings, err, confFile)
				continue
			}
			// Ensure the config has a "type" so we know what plugin to run.
			// Also catches the case where somebody put a conflist into a conf file.
			if conf.Network.Type == "" {
				warnings = fmt.Errorf("%w: error loading CNI config file %s: no 'type'; perhaps this is a .conflist?", warnings, confFile)
				continue
			}

			confList, err = libcni.ConfListFromConf(conf)
			if err != nil {
				warnings = fmt.Errorf("%v: %w, error converting CNI config file %s to list", warnings, err, confFile)
				continue
			}
		}

		if len(confList.Plugins) == 0 {
			warnings = fmt.Errorf("%w: CNI config list %s has no networks, skipping", warnings, confFile)
			continue
		}

		return confList, warnings, nil
	}

	return nil, warnings, fmt.Errorf("%w in %s", ErrNoNetworksFound, confDir)
}

// getRuntimeConf returns common libcni runtime conf used to interact with the cni
func (p *cniPlugin) getCNIRuntimeConf(id string) *libcni.RuntimeConf {
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

// cniPodNetwork is a pod network environment context
type cniPodNetwork struct {
	noopPodNetwork // every method not implemented is noop
	plugin         *cniPlugin
	netList        *libcni.NetworkConfigList
	runtimeConf    *libcni.RuntimeConf
	annotations    map[string]string
}

// ContainerNetwork enters a container network environment context
func (s *cniPodNetwork) ContainerNetwork(id string, annotations map[string]string) (ContainerNetwork, error) {
	return &cniContainerNetwork{
		pod:         s,
		cid:         id,
		annotations: annotations,
	}, nil
}

// Status reports IP and any error with the network of that pod
func (s *cniPodNetwork) Status(ctx context.Context, prop *PropertiesRunning) (*Status, error) {
	ips, err := s.ips([]byte(prop.Data["result"]))
	if err != nil {
		return nil, err
	}

	return &Status{IPs: ips}, nil
}

// Setup creates the network interface for the provided netfile
func (s *cniPodNetwork) setup(ctx context.Context, netfile string) (types.Result, error) {
	s.runtimeConf.NetNS = netfile

	prevResult, err := s.plugin.cni.AddNetworkList(ctx, s.netList, s.runtimeConf)
	if err != nil {
		return nil, err
	}

	// convert the result to the current cni version
	return current.NewResultFromResult(prevResult)
}

// Teardown removes the network compeletely as good as possible
func (s *cniPodNetwork) teardown(ctx context.Context) error {
	s.runtimeConf.NetNS = ""
	return s.plugin.cni.DelNetworkList(ctx, s.netList, s.runtimeConf)
}

// Get ips of that result
func (s *cniPodNetwork) ips(previousresult []byte) ([]net.IP, error) {
	if previousresult == nil {
		previousresult = []byte{}
	}

	prevResult, err := current.NewResult(previousresult)
	if err != nil {
		return nil, err
	}

	// convert the result to the current cni version
	result, err := current.NewResultFromResult(prevResult)
	if err != nil {
		return nil, err
	}

	if len(result.IPs) == 0 {
		return nil, fmt.Errorf("%w: for %v", &net.AddrError{Err: "missing address"}, s.runtimeConf.ContainerID)
	}

	if result.IPs[0].Address.IP == nil {
		return nil, fmt.Errorf("%w: for %v", &net.AddrError{Err: "invalid address"}, s.runtimeConf.ContainerID)
	}

	return []net.IP{result.IPs[0].Address.IP}, nil
}

// cniContainerNetwork is a container network environment context
type cniContainerNetwork struct {
	noopContainerNetwork // every method not implemented is noop
	pod                  *cniPodNetwork
	cid                  string
	annotations          map[string]string
}

// WhenStarted is called when the container is started.
func (c *cniContainerNetwork) WhenStarted(ctx context.Context, prop *PropertiesRunning) (*Result, error) {
	// TODO: As long as we haven't figured out to do 1:n podnetwork:container this method goes up to pod
	result, err := c.pod.setup(ctx, fmt.Sprintf("/proc/%s/ns/net", strconv.FormatInt(prop.Pid, 10)))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return &Result{Data: map[string]string{"result": string(b)}}, nil
}

// WhenDeleted is called when the container is deleted. If tearing down here, must tear down as good as possible. Must
// tear down here if not implemented for WhenStopped. If an error is returned it will only be logged
func (c *cniContainerNetwork) WhenDeleted(ctx context.Context, prop *Properties) error {
	// TODO: As long as we haven't figured out to do 1:n podnetwork:container this method goes up to pod
	return c.pod.teardown(ctx)
}
