package device // import "github.com/automaticserver/lxe/lxf/device"

import (
	"github.com/juju/errors"
)

var (
	schema = map[string]Device{
		BlockType: &Block{},
		CharType:  &Char{},
		DiskType:  &Disk{},
		NicType:   &Nic{},
		NoneType:  &None{},
		ProxyType: &Proxy{},
	}
)

// Device must support mapping from lxd device map bidirectional
type Device interface {
	// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map
	ToMap() (name string, options map[string]string)
	// FromMap loads assigned name (can be empty) and options
	FromMap(name string, options map[string]string) error
	// New creates a new empty device
	new() Device
}

// Detects and loads device by type
func Detect(name string, options map[string]string) (Device, error) {
	t, is := schema[options["type"]]
	if !is {
		return nil, errors.NotSupportedf("unknown device type: %v", options["type"])
	}

	d := t.new()

	err := d.FromMap(name, options)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Devices allows having a list of devices unique by name
type Devices []Device

// Upsert adds a device or overrides an entry if the key name exists
func (d *Devices) Upsert(a Device) {
	for k, e := range *d {
		eName, _ := e.ToMap()
		aName, _ := a.ToMap()

		if eName == aName {
			(*d)[k] = a
			return
		}
	}

	*d = append(*d, a)
}
