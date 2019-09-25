package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericAddToMapErrorOnMultiple(t *testing.T) {
	disk := Disk{
		Path: "/",
	}
	disk2 := Disk{
		Path: "/foo",
	}

	s := map[string]map[string]string{}

	err := AddDisksToMap(s, disk)
	assert.NoError(t, err)

	// error when there exists one already conflicting
	err = AddDisksToMap(s, disk)
	assert.Error(t, err)

	// error when input is conflicting
	err = AddDisksToMap(s, disk2, disk2)
	assert.Error(t, err)
}

// nolint: dupl
func TestAddDiskToMap(t *testing.T) {
	disk := Disk{
		Path: "/",
	}
	s := map[string]map[string]string{}

	err := AddDisksToMap(s, disk)
	if err != nil {
		t.Errorf("could not serialize disk to map: %v", err)
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}

// nolint: dupl
func TestGetDisksFromMap(t *testing.T) {
	disk := Disk{
		Path: "/",
	}
	s := map[string]map[string]string{}

	err := AddDisksToMap(s, disk)
	if err != nil {
		t.Errorf("could not serialize disk to map: %v", err)
	}

	disks, err := GetDisksFromMap(s)

	if err != nil {
		t.Errorf("could not read disk from map, %v", err)
	}
	if len(disks) != 1 {
		t.Errorf("expected one disk but there are %v", len(disks))
	}
}

// nolint: dupl
func TestGetDisksWithOverrideAdd(t *testing.T) {
	var disks Disks
	disk := Disk{
		Path: "/",
	}

	disks.Add(disk)
	disks.Add(disk)

	assert.Len(t, disks, 1)
}

// nolint: dupl
func TestAddProxyToMap(t *testing.T) {
	proxy := Proxy{
		Destination: ProxyEndpoint{Protocol: ProtocolTCP, Port: 8000},
		Listen:      ProxyEndpoint{Protocol: ProtocolTCP, Port: 80},
	}
	s := map[string]map[string]string{}

	err := AddProxiesToMap(s, proxy)
	if err != nil {
		t.Errorf("could not serialize proxy to map: %v", err)
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}

// nolint: dupl
func TestGetProxiesFromMap(t *testing.T) {
	proxy := Proxy{
		Destination: ProxyEndpoint{Protocol: ProtocolTCP, Port: 8000},
		Listen:      ProxyEndpoint{Protocol: ProtocolTCP, Port: 80},
	}
	s := map[string]map[string]string{}

	err := AddProxiesToMap(s, proxy)
	if err != nil {
		t.Errorf("could not serialize proxy to map: %v", err)
	}

	proxies, err := GetProxiesFromMap(s)
	if err != nil {
		t.Errorf("could not read proxy from map, %v", err)
	}
	if len(proxies) != 1 {
		t.Errorf("expected one proxy but there are %v", len(proxies))
	}
}

// nolint: dupl
func TestGetProxiesWithOverrideAdd(t *testing.T) {
	var proxies Proxies
	proxy := Proxy{
		Destination: ProxyEndpoint{Protocol: ProtocolTCP, Port: 8000},
		Listen:      ProxyEndpoint{Protocol: ProtocolTCP, Port: 80},
	}

	proxies.Add(proxy)
	proxies.Add(proxy)

	assert.Len(t, proxies, 1)
}

// nolint: dupl
func TestAddBlocksToMap(t *testing.T) {
	block := Block{
		Path: "/",
	}
	s := map[string]map[string]string{}

	err := AddBlocksToMap(s, block)
	if err != nil {
		t.Errorf("could not serialize block to map: %v", err)
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}

// nolint: dupl
func TestGetBlocksFromMap(t *testing.T) {
	block := Block{
		Path: "/",
	}
	s := map[string]map[string]string{}

	err := AddBlocksToMap(s, block)
	if err != nil {
		t.Errorf("could not serialize block to map: %v", err)
	}

	blocks, err := GetBlocksFromMap(s)
	if err != nil {
		t.Errorf("could not read block from map, %v", err)
	}
	if len(blocks) != 1 {
		t.Errorf("expected one block but there are %v", len(blocks))
	}
}

// nolint: dupl
func TestGetBlocksWithOverrideAdd(t *testing.T) {
	var blocks Blocks
	block := Block{
		Path: "/",
	}

	blocks.Add(block)
	blocks.Add(block)

	assert.Len(t, blocks, 1)
}

// nolint: dupl
func TestAddNicsToMap(t *testing.T) {
	nic := Nic{
		Name: "eth0",
	}
	s := map[string]map[string]string{}

	err := AddNicsToMap(s, nic)
	if err != nil {
		t.Errorf("could not serialize nic to map: %v", err)
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}

// nolint: dupl
func TestGetNicsFromMap(t *testing.T) {
	nic := Nic{
		Name: "eth0",
	}
	s := map[string]map[string]string{}

	err := AddNicsToMap(s, nic)
	if err != nil {
		t.Errorf("could not serialize nic to map: %v", err)
	}

	nics, err := GetNicsFromMap(s)
	if err != nil {
		t.Errorf("could not read nic from map, %v", err)
	}
	if len(nics) != 1 {
		t.Errorf("expected one nic but there are %v", len(nics))
	}
}

// nolint: dupl
func TestGetNicsWithOverrideAdd(t *testing.T) {
	var nics Nics
	nic := Nic{
		Name: "eth0",
	}

	nics.Add(nic)
	nics.Add(nic)

	assert.Len(t, nics, 1)
}

// nolint: dupl
func TestAddNonesToMap(t *testing.T) {
	none := None{
		Name: "eth0",
	}
	s := map[string]map[string]string{}

	err := AddNonesToMap(s, none)
	if err != nil {
		t.Errorf("could not serialize none to map: %v", err)
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}

// nolint: dupl
func TestGetNonesFromMap(t *testing.T) {
	none := None{
		Name: "eth0",
	}
	s := map[string]map[string]string{}

	err := AddNonesToMap(s, none)
	if err != nil {
		t.Errorf("could not serialize none to map: %v", err)
	}

	nones, err := GetNonesFromMap(s)
	if err != nil {
		t.Errorf("could not read none from map, %v", err)
	}
	if len(nones) != 1 {
		t.Errorf("expected one none but there are %v", len(nones))
	}
}

// nolint: dupl
func TestGetNonesWithOverrideAdd(t *testing.T) {
	var nones Nones
	none := None{
		Name: "eth0",
	}

	nones.Add(none)
	nones.Add(none)

	assert.Len(t, nones, 1)
}
