package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

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
		start = net.IPv4(networkIP[0], networkIP[1], networkIP[2], networkIP[3]+1)
	}

	if end == nil {
		end = net.IPv4(broadcastIP[0], broadcastIP[1], broadcastIP[2], broadcastIP[3]-1)
	}

	// Until a usable IP is found...
	// TODO: detect if there's never a possible address and return nil?
	var ip net.IP
OUTER:
	for {
		// randomly select an ip address within the specified subnet
		trialB := make([]byte, 4)
		binary.LittleEndian.PutUint32(trialB, rand.Uint32())
		for i, v := range trialB {
			trialB[i] = subnet.IP[i] + (v &^ subnet.Mask[i])
		}
		trial := net.IPv4(trialB[0], trialB[1], trialB[2], trialB[3])

		// not allowed if outside explicitly defined range
		if bytes.Compare(trial, start) < 0 || bytes.Compare(trial, end) > 0 {
			fmt.Printf("compare: trial %v, start %v, end %v\n", trial, start, end)
			continue
		}

		// not allowed if already exists in current leases
		for _, lease := range leases {
			if trial.Equal(lease) {
				continue OUTER
			}
		}

		// IP is fine :)
		ip = trial
		break
	}

	return ip
}
