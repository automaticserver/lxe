package device

import (
	"fmt"
)

// Device defines needed methods for all devices for translation
// from and to maps
type Device interface {
	ToMap() (map[string]string, error)
	GetName() string
}

// AddToMap will serialize the provided devices into provided map
// Will complain with an error if a device can't be serialized due to invalid values
// or if the name of a device colides.
func AddToMap(m map[string]map[string]string, devs ...Device) error {
	additional := map[string]map[string]string{}

	for _, dev := range devs {
		name := dev.GetName()
		if _, has := m[name]; has {
			return fmt.Errorf("there are more than one devices with name %v", name)
		}
		if _, has := additional[name]; has {
			return fmt.Errorf("there are more than one devices with name %v", name)
		}
		dm, err := dev.ToMap()
		if err != nil {
			return err
		}
		additional[name] = dm
	}

	// now we are sure there are no errors so append all the entries
	for k, v := range additional {
		m[k] = v
	}
	return nil
}

// AddDisksToMap Will add the disks to the map
func AddDisksToMap(m map[string]map[string]string, disks ...Disk) error {
	devs := []Device{}
	for _, d := range disks {
		devs = append(devs, d)
	}
	return AddToMap(m, devs...)
}

// GetDisksFromMap will will add the proxies to the map
func GetDisksFromMap(maps map[string]map[string]string) (Disks, error) { // nolint: dupl
	disks := Disks{}
	for _, m := range maps {
		if m["type"] == diskType {
			p, err := DiskFromMap(m)
			if err != nil {
				return nil, err
			}
			disks.Add(p)
		}
	}
	return disks, nil
}

// AddProxiesToMap Will add the proxies to the map
func AddProxiesToMap(m map[string]map[string]string, proxies ...Proxy) error {
	devs := []Device{}
	for _, d := range proxies {
		devs = append(devs, d)
	}
	return AddToMap(m, devs...)
}

// GetProxiesFromMap will read all proxy devices from the map
func GetProxiesFromMap(maps map[string]map[string]string) (Proxies, error) { // nolint: dupl
	proxies := Proxies{}
	for _, m := range maps {
		if m["type"] == proxyType {
			p, err := ProxyFromMap(m)
			if err != nil {
				return nil, err
			}
			proxies.Add(p)
		}
	}
	return proxies, nil
}

// AddBlocksToMap Will add the block devices to the map
func AddBlocksToMap(m map[string]map[string]string, blocks ...Block) error {
	devs := []Device{}
	for _, d := range blocks {
		devs = append(devs, d)
	}
	return AddToMap(m, devs...)
}

// GetBlocksFromMap will read all proxy devices from the map
func GetBlocksFromMap(maps map[string]map[string]string) ([]Block, error) { // nolint: dupl
	blocks := Blocks{}
	for _, m := range maps {
		if m["type"] == blockType {
			p, err := BlockFromMap(m)
			if err != nil {
				return nil, err
			}
			blocks.Add(p)
		}
	}
	return blocks, nil
}

// AddNicsToMap Will add the nic devices to the map
func AddNicsToMap(m map[string]map[string]string, nics ...Nic) error {
	devs := []Device{}
	for _, d := range nics {
		devs = append(devs, d)
	}
	return AddToMap(m, devs...)
}

// GetNicsFromMap will read all proxy devices from the map
func GetNicsFromMap(maps map[string]map[string]string) ([]Nic, error) { // nolint: dupl
	nics := Nics{}
	for _, m := range maps {
		if m["type"] == nicType {
			p, err := NicFromMap(m)
			if err != nil {
				return nil, err
			}
			nics.Add(p)
		}
	}
	return nics, nil
}
