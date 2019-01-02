package lxf

import (
	"net"
	"strconv"

	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
)

// EnsureBridge ensures the bridge exists with the defined options
func (l *LXF) EnsureBridge(name, cidr string, nat bool) error {
	logger.Infof("DEBUG EnsureBridge: %v | %v | %v", name, cidr, nat)

	// Always use first address in range for the bridge
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	net.IP[3]++

	put := api.NetworkPut{
		Description: "managed by LXE, default bridge",
		Config: map[string]string{
			"ipv4.address": net.String(),
			"ipv4.dhcp":    strconv.FormatBool(true),
			"ipv4.nat":     strconv.FormatBool(true),
			"ipv6.address": "none",
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

	for k, v := range put.Config {
		network.Config[k] = v
	}
	return l.server.UpdateNetwork(name, network.Writable(), ETag)
}
