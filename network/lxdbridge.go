package network

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/network/cloudinit"
	"github.com/automaticserver/lxe/shared"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	DefaultLXDBridge = "lxdbr0"
)

var (
	ErrNotBridge = errors.New("not a bridge")
)

// ConfLXDBridge are configuration options for the LXDBridge plugin. All properties are optional and get a default value
type ConfLXDBridge struct {
	LXDBridge  string
	Cidr       string
	Nat        bool
	CreateOnly bool
}

func (c *ConfLXDBridge) setDefaults() {
	if c.LXDBridge == "" {
		c.LXDBridge = DefaultLXDBridge
	}
}

// lxdBridgePlugin manages the pod networks using LXDBridge
type lxdBridgePlugin struct {
	noopPlugin // every method not implemented is noop
	server     lxd.ContainerServer
	conf       ConfLXDBridge
}

// InitPluginLXDBridge instantiates the LXDBridge plugin using the provided config
func InitPluginLXDBridge(server lxd.ContainerServer, conf ConfLXDBridge) (*lxdBridgePlugin, error) { // nolint: golint, revive // intended to not export lxdBridgePlugin
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
			// We don't need to receive a DNS in DHCP, Kubernetes' DNS is always set by requesting a mount for resolv.conf.
			// This disables dns in dnsmasq (option -p: https://linux.die.net/man/8/dnsmasq)
			"raw.dnsmasq": `port=0`,
		},
	}

	network, ETag, err := p.server.GetNetwork(p.conf.LXDBridge)
	if err != nil {
		if shared.IsErrNotFound(err) {
			return p.server.CreateNetwork(api.NetworksPost{
				Name:       p.conf.LXDBridge,
				Type:       "bridge",
				NetworkPut: put,
			})
		}

		return err
	} else if network.Type != "bridge" {
		return fmt.Errorf("%w: %v, but is %v", ErrNotBridge, p.conf.LXDBridge, network.Type)
	}

	// don't update when only creation is requested
	// TODO: Should we return an error if the bridge settings e.g. cidr would change?
	if p.conf.CreateOnly {
		return nil
	}

	for k, v := range put.Config {
		network.Config[k] = v
	}

	return p.server.UpdateNetwork(p.conf.LXDBridge, network.Writable(), ETag)
}

var ErrNotImplemented = errors.New("not implemented")

// findFreeIP generates a IP within the range of the provided lxd managed bridge which does
// not exist in the current leases
func (p *lxdBridgePlugin) findFreeIP() (net.IP, error) {
	network, _, err := p.server.GetNetwork(p.conf.LXDBridge)
	if err != nil {
		return nil, err
	} else if network.Config["ipv4.dhcp.ranges"] != "" {
		// actually we can now using FindFreeIP(), but not good enough, as this field can yield multiple ranges
		return nil, fmt.Errorf("%w to find an IP with explicitly set ip ranges `ipv4.dhcp.ranges` in bridge %v", ErrNotImplemented, p.conf.LXDBridge)
	}

	rawLeases, err := p.server.GetNetworkLeases(p.conf.LXDBridge)
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
		return nil, &net.AddrError{Addr: prop.Data["interface-address"], Err: "missing"}
	}

	ip := net.ParseIP(prop.Data["interface-address"])
	if ip == nil {
		return nil, &net.ParseError{Type: "IP address", Text: prop.Data["interface-address"]}
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
		// 	"bridge":            s.plugin.conf.LXDBridge,
		"interface-address": randIP.String(), // except this for IP return shortcut in Status
		// 	"physical-type":     "dhcp",
	}
	r.Nics = []device.Nic{
		{
			Name:        DefaultInterface,
			NicType:     "bridged",
			Parent:      s.plugin.conf.LXDBridge,
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
