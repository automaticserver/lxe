package network

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/network/cloudinit"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	rtApi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

const (
	// ErrorLXDNotFound is the error string a LXD request returns, when nothing is found. Unfortunately there is no
	// constant in the lxd source we could've used
	ErrorLXDNotFound = "not found"
	defaultLxdBridge = "lxebr0"
)

// ConfLxdBridge are configuration options for the LxdBridge plugin. All properties are optional and get a default value
type ConfLxdBridge struct {
	LxdBridge  string
	Cidr       string
	Nat        bool
	CreateOnly bool
}

func (c *ConfLxdBridge) setDefaults() {
	if c.LxdBridge == "" {
		c.LxdBridge = defaultLxdBridge
	}
}

// lxdBridgePlugin manages the pod networks using LxdBridge
type lxdBridgePlugin struct {
	noopPlugin // every method not implemented is noop
	server     lxd.ContainerServer
	conf       ConfLxdBridge
}

// InitPluginLxdBridge instantiates the LxdBridge plugin using the provided config
func InitPluginLxdBridge(server lxd.ContainerServer, conf ConfLxdBridge) (*lxdBridgePlugin, error) { // nolint: golint // intended to not export lxdBridgePlugin
	conf.setDefaults()

	p := &lxdBridgePlugin{
		server: server,
		conf:   conf,
	}

	err := p.ensureBridge()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// PodNetwork enters a pod network environment context
func (p *lxdBridgePlugin) PodNetwork(id string, annotations map[string]string) (PodNetwork, error) {
	return &lxdBridgePodNetwork{
		plugin:      p,
		podID:       id,
		annotations: annotations,
	}, nil
}

// UpdateRuntimeConfig is called when there are updates to the configuration which the plugin might need to apply
func (p *lxdBridgePlugin) UpdateRuntimeConfig(conf *rtApi.RuntimeConfig) error {
	if cidr := conf.GetNetworkConfig().GetPodCidr(); cidr != "" {
		p.conf.Cidr = cidr
		return p.ensureBridge()
	}

	return nil
}

// EnsureBridge ensures the bridge exists with the defined options. Cidr is an expected ipv4 cidr or can be empty to
// automatically assign a cidr
func (p *lxdBridgePlugin) ensureBridge() error {
	var address string
	if p.conf.Cidr == "" {
		address = "auto"
	} else {
		// Always use first address in range for the bridge
		_, net, err := net.ParseCIDR(p.conf.Cidr)
		if err != nil {
			return err
		}
		net.IP[3]++
		address = net.String()
	}

	put := api.NetworkPut{
		Description: "managed by LXE, default bridge",
		Config: map[string]string{
			"ipv4.address": address,
			"ipv4.dhcp":    strconv.FormatBool(true),
			"ipv4.nat":     strconv.FormatBool(p.conf.Nat),
			"ipv6.address": "none",
			// We don't need to receive a DNS in DHCP, Kubernetes' DNS is always set my requesting a mount for resolv.conf.
			// This disables dns in dnsmasq (option -p: https://linux.die.net/man/8/dnsmasq)
			"raw.dnsmasq": `port=0`,
		},
	}

	network, ETag, err := p.server.GetNetwork(p.conf.LxdBridge)
	if err != nil {
		if err.Error() == ErrorLXDNotFound {
			return p.server.CreateNetwork(api.NetworksPost{
				Name:       p.conf.LxdBridge,
				Type:       "bridge",
				Managed:    true,
				NetworkPut: put,
			})
		}

		return err
	} else if network.Type != "bridge" {
		return fmt.Errorf("expected %v to be a bridge, but is %v", p.conf.LxdBridge, network.Type)
	}

	// don't update when only creation is requested
	// TODO: Should we return an error if the bridge settings e.g. cidr would change?
	if p.conf.CreateOnly {
		return nil
	}

	for k, v := range put.Config {
		network.Config[k] = v
	}

	return p.server.UpdateNetwork(p.conf.LxdBridge, network.Writable(), ETag)
}

// findFreeIP generates a IP within the range of the provided lxd managed bridge which does
// not exist in the current leases
func (p *lxdBridgePlugin) findFreeIP() (net.IP, error) {
	network, _, err := p.server.GetNetwork(p.conf.LxdBridge)
	if err != nil {
		return nil, err
	} else if network.Config["ipv4.dhcp.ranges"] != "" {
		// actually we can now using findFreeIP() below, but not good enough, as this field can yield multiple ranges
		return nil, fmt.Errorf("not yet implemented to find an IP with explicitly set ip ranges `ipv4.dhcp.ranges` in bridge %v", p.conf.LxdBridge)
	}

	rawLeases, err := p.server.GetNetworkLeases(p.conf.LxdBridge)
	if err != nil {
		return nil, err
	}

	leases := []net.IP{}

	for _, rawIP := range rawLeases {
		leases = append(leases, net.ParseIP(rawIP.Address))
	}

	bridgeIP, bridgeNet, err := net.ParseCIDR(network.Config["ipv4.address"])
	if err != nil {
		return nil, err
	}

	leases = append(leases, bridgeIP) // also exclude bridge ip

	return FindFreeIP(bridgeNet, leases, nil, nil), nil
}

// lxdBridgePodNetwork is a pod network environment context
type lxdBridgePodNetwork struct {
	noopPodNetwork // every method not implemented is noop
	plugin         *lxdBridgePlugin
	podID          string
	annotations    map[string]string
}

// ContainerNetwork enters a container network environment context
func (s *lxdBridgePodNetwork) ContainerNetwork(id string, annotations map[string]string) (ContainerNetwork, error) {
	return &lxdBridgeContainerNetwork{
		pod:         s,
		cid:         id,
		annotations: annotations,
	}, nil
}

// Status reports IP and any error with the network of that pod
func (s *lxdBridgePodNetwork) Status(ctx context.Context, prop *PropertiesRunning) (*Status, error) {
	if prop.Data["interface-address"] == "" {
		return nil, fmt.Errorf("no ip address found")
	}

	ip := net.ParseIP(prop.Data["interface-address"])
	if ip == nil {
		return nil, fmt.Errorf("invalid ip address format")
	}

	return &Status{
		IPs: []net.IP{ip},
	}, nil
}

// WhenCreated is called when the pod is created.
func (s *lxdBridgePodNetwork) WhenCreated(ctx context.Context, prop *Properties) (*Result, error) {
	// default is to use the predefined lxd bridge managed by lxe
	randIP, err := s.plugin.findFreeIP()
	if err != nil {
		return nil, err
	}

	r := &Result{}
	// TODO: Remove, I think we don't/shouldn't need that anymore
	r.Data = map[string]string{
		// 	"bridge":            s.plugin.conf.LxdBridge,
		"interface-address": randIP.String(), // except this for IP return shortcut in Status
		// 	"physical-type":     "dhcp",
	}
	r.Nics = []device.Nic{
		{
			Name:        DefaultInterface,
			NicType:     "bridged",
			Parent:      s.plugin.conf.LxdBridge,
			IPv4Address: randIP.String(),
		},
	}
	r.NetworkConfigEntries = []cloudinit.NetworkConfigEntryPhysical{
		{
			NetworkConfigEntry: cloudinit.NetworkConfigEntry{
				Type: "physical",
			},
			Name: DefaultInterface,
			Subnets: []cloudinit.NetworkConfigEntryPhysicalSubnet{
				{
					Type: "dhcp",
				},
			},
		},
	}

	return r, nil
}

// lxdBridgeContainerNetwork is a container network environment context
type lxdBridgeContainerNetwork struct {
	noopContainerNetwork // every method not implemented is noop
	pod                  *lxdBridgePodNetwork
	cid                  string
	annotations          map[string]string
}
