package device

import (
	"fmt"
)

const (
	NicType = "nic"
)

// Nic device representation https://lxd.readthedocs.io/en/latest/containers/#type-nic
type Nic struct {
	KeyName     string
	Name        string
	NicType     string
	Parent      string
	IPv4Address string
}

func (d *Nic) getName() string {
	var name string

	switch {
	case d.KeyName != "":
		name = d.KeyName
	default:
		name = fmt.Sprintf("%s-%s", NicType, d.Name)
	}

	return name
}

// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map
func (d *Nic) ToMap() (string, map[string]string) {
	return d.getName(), map[string]string{
		"type":         NicType,
		"name":         d.Name,
		"nictype":      d.NicType,
		"parent":       d.Parent,
		"ipv4.address": d.IPv4Address,
	}
}

// FromMap loads assigned name (can be empty) and options
func (d *Nic) FromMap(name string, options map[string]string) error {
	d.KeyName = name
	d.Name = options["name"]
	d.NicType = options["nictype"]
	d.Parent = options["parent"]
	d.IPv4Address = options["ipv4.address"]

	return nil
}

// New creates a new empty device
func (d *Nic) new() Device {
	return &Nic{}
}
