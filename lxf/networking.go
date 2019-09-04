package lxf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"github.com/lxc/lxd/shared/api"
)

// EnsureBridge ensures the bridge exists with the defined options
// cidr is an expected ipv4 cidr or can be empty to automatically assign a cidr
func (l *Client) EnsureBridge(name, cidr string, nat, createOnly bool) error {
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
		if err.Error() == ErrorLXDNotFound {
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
func (l *Client) FindFreeIPBridgeLXD(bridge string) (net.IP, error) {
	network, _, err := l.server.GetNetwork(bridge)
	if err != nil {
		return nil, err
	}
	if network.Config["ipv4.dhcp.ranges"] != "" {
		// actually we can now using FindFreeIP() below, but not good enough, as this field can yield multiple ranges
		return nil, fmt.Errorf("Not yet implemented to find an IP with explicitly set ip ranges `ipv4.dhcp.ranges` in bridge %v", bridge)
	}

	rawLeases, err := l.server.GetNetworkLeases(bridge)
	if err != nil {
		return nil, err
	}
	var leases []net.IP
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

// FindFreeIP tries to find an available IP address within given subnet, respecting reserved addresses in leases and
// must be between the start and end address. Network and broadcast IP are also reserved and automatically added to
// leases. If start or end is nil their closest available address from the subnet is selected.
func FindFreeIP(subnet *net.IPNet, leases []net.IP, start, end net.IP) net.IP {
	// put non-usable addresses also to leases, so they can't be selected
	networkIP := subnet.IP
	broadcastIP := make(net.IP, 4)
	for i := range broadcastIP {
		broadcastIP[i] = subnet.IP[i] | ^subnet.Mask[i]
	}
	leases = append(leases, networkIP, broadcastIP)

	// defaults for start and end to usable addresses if not explicitly defined
	if start == nil {
		start = networkIP
		start[3]++
	}
	if end == nil {
		end = broadcastIP
		end[3]--
	}

	// Until a usable IP is found...
	// TODO: detect if there's never a possible address and return nil?
	var ip net.IP
	for {
		// randomly select an ip address within the specified subnet
		trial := make(net.IP, 4)
		binary.LittleEndian.PutUint32(trial, rand.Uint32())
		for i, v := range trial {
			trial[i] = networkIP[i] + (v &^ subnet.Mask[i])
		}

		// not allowed if outside explicitly defined range
		if bytes.Compare(trial, start) < 0 || bytes.Compare(trial, end) > 0 {
			continue
		}

		// not allowed if already exists in current leases
		for _, lease := range leases {
			if trial.Equal(lease) {
				continue
			}
		}

		// IP is fine :)
		ip = trial
		break
	}

	return ip
}
