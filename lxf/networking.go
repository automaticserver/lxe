package lxf

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"github.com/lxc/lxd/shared/api"
)

// EnsureBridge ensures the bridge exists with the defined options
// cidr is an expected ipv4 cidr or can be empty to automatically assign a cidr
func (l *LXF) EnsureBridge(name, cidr string, nat, createOnly bool) error {
	var address string
	if cidr == "" {
		address = "auto"
	} else {
		// Always use first address in range for the bridge
		_, net, err := net.ParseCIDR(cidr)
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
			"ipv4.nat":     strconv.FormatBool(true),
			"ipv6.address": "none",
			// We don't need to recieve a DNS in DHCP, Kubernetes' DNS is always set
			// disables dns (option -p: https://linux.die.net/man/8/dnsmasq)
			// > Listen on <port> instead of the standard DNS port (53). Setting this to
			// > zero completely disables DNS function, leaving only DHCP and/or TFTP.
			"raw.dnsmasq": `port=0`,
		},
	}

	network, ETag, err := l.server.GetNetwork(name)
	if err != nil {
		if IsErrorNotFound(err) {
			return l.server.CreateNetwork(api.NetworksPost{
				Name:       name,
				Type:       "bridge",
				Managed:    true,
				NetworkPut: put,
			})
		}

		return err
	}
	if network.Type != "bridge" {
		return fmt.Errorf("Expected %v to be a bridge, but is %v", name, network.Type)
	}

	// don't update when only creation is requested
	if createOnly {
		return nil
	}

	for k, v := range put.Config {
		network.Config[k] = v
	}
	return l.server.UpdateNetwork(name, network.Writable(), ETag)
}

// FindFreeIP generates a IP within the range of the provided lxd managed bridge which does
// not exist in the current leases
func (l *LXF) FindFreeIP(bridge string) (net.IP, error) {
	network, _, err := l.server.GetNetwork(bridge)
	if err != nil {
		return nil, err
	}
	if network.Config["ipv4.dhcp.ranges"] != "" {
		return nil, fmt.Errorf("Not yet implemented to find an IP with explicitly set ip ranges `ipv4.dhcp.ranges` in bridge %v", bridge)
	}

	leases, err := l.server.GetNetworkLeases(bridge)
	if err != nil {
		return nil, err
	}

	bridgeIP, bridgeNet, err := net.ParseCIDR(network.Config["ipv4.address"])
	if err != nil {
		return nil, err
	}

	broadcastIP := make(net.IP, 4)
	for i := range broadcastIP {
		broadcastIP[i] = bridgeNet.IP[i] | ^bridgeNet.Mask[i]
	}

	var ip net.IP
	// Until a usable IP is found...
	for {
		// select randomly an ip address within the specified range
		randIP := make(net.IP, 4)
		binary.LittleEndian.PutUint32(randIP, rand.Uint32())
		for i, v := range randIP {
			randIP[i] = bridgeNet.IP[i] + (v &^ bridgeNet.Mask[i])
		}

		// not allowed to be the bridge ip, network address or broadcast address
		if randIP.String() == bridgeIP.String() || randIP.String() == bridgeNet.IP.String() || randIP.String() == broadcastIP.String() {
			continue
		}
		// not allowed to exist in current leases
		for _, lease := range leases {
			if randIP.String() == lease.Address {
				continue
			}
		}
		// if reached here, ip is fine
		ip = randIP
		break
	}
	return ip, nil
}
