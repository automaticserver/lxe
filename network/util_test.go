package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFreeIP_CanFind(t *testing.T) {
	t.Parallel()

	_, ipNet, err := net.ParseCIDR("192.168.224.0/30")
	assert.NoError(t, err)

	for i := 0; i < 20; i++ {
		ip := FindFreeIP(ipNet, nil, nil, nil)
		assert.NotNil(t, ip)
	}
}

func TestFindFreeIP_ExcludesLeases(t *testing.T) {
	t.Parallel()

	_, ipNet, err := net.ParseCIDR("192.168.224.0/30")
	assert.NoError(t, err)

	leases := []net.IP{net.ParseIP("192.168.224.1")}

	for i := 0; i < 10; i++ {
		ip := FindFreeIP(ipNet, leases, nil, nil)
		assert.Equal(t, "192.168.224.2", ip.String())
	}
}

func TestFindFreeIP_RespectRange(t *testing.T) {
	t.Parallel()

	_, ipNet, err := net.ParseCIDR("192.168.224.0/30")
	assert.NoError(t, err)

	start := net.ParseIP("192.168.224.2")
	end := start

	for i := 0; i < 10; i++ {
		ip := FindFreeIP(ipNet, nil, start, end)
		assert.Equal(t, start.String(), ip.String())
	}
}

// TODO: Timeout or inability to find a valid ip to return an error
