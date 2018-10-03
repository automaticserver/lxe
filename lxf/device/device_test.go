package device

import "testing"

// nolint: dupl
func TestAddDiskToMap(t *testing.T) {
	disk := &Disk{}
	s := map[string]map[string]string{}

	err := AddToMap(s, disk)
	if err != nil {
		t.Errorf("could not serialize disk to map")
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}

// nolint: dupl
func TestGetDisksFromMap(t *testing.T) {
	disk := &Disk{}
	s := map[string]map[string]string{}

	err := AddToMap(s, disk)
	if err != nil {
		t.Errorf("could not serialize disk to map")
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
func TestAddProxyToMap(t *testing.T) {
	proxy := &Proxy{}
	s := map[string]map[string]string{}

	err := AddToMap(s, proxy)
	if err != nil {
		t.Errorf("could not serialize proxy to map")
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}

// nolint: dupl
func TestGetProxiesFromMap(t *testing.T) {
	proxy := &Proxy{
		Destination: ProxyEndpoint{Protocol: ProtocolTCP, Port: 8000},
		Listen:      ProxyEndpoint{Protocol: ProtocolTCP, Port: 80},
	}
	s := map[string]map[string]string{}

	err := AddToMap(s, proxy)
	if err != nil {
		t.Errorf("could not serialize proxy to map")
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
func TestAddBlockToMap(t *testing.T) {
	block := &Block{}
	s := map[string]map[string]string{}

	err := AddToMap(s, block)
	if err != nil {
		t.Errorf("could not serialize block to map")
	}

	if len(s) != 1 {
		t.Errorf("device map should have one entry")
	}
}
