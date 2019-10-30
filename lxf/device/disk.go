// nolint: dupl
package device

import (
	"fmt"
	"strconv"
)

const (
	DiskType = "disk"
)

// Disk device representation https://lxd.readthedocs.io/en/latest/containers/#type-disk
type Disk struct {
	KeyName  string
	Path     string
	Source   string
	Pool     string
	Size     string
	Readonly bool
	Optional bool
}

func (d *Disk) getName() string {
	var name string

	switch {
	case d.KeyName != "":
		name = d.KeyName
	case d.Path == "":
		name = fmt.Sprintf("%s-%s", DiskType, d.Source)
	default:
		name = fmt.Sprintf("%s-%s", DiskType, d.Path)
	}

	return name
}

// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map
func (d *Disk) ToMap() (string, map[string]string) {
	return d.getName(), map[string]string{
		"type":     DiskType,
		"path":     d.Path,
		"source":   d.Source,
		"pool":     d.Pool,
		"size":     d.Size,
		"readonly": strconv.FormatBool(d.Readonly),
		"optional": strconv.FormatBool(d.Optional),
	}
}

// FromMap creates a new device with assigned name (can be empty) and options
func (d *Disk) FromMap(name string, options map[string]string) (Device, error) {
	return &Disk{
		KeyName:  name,
		Path:     options["path"],
		Source:   options["source"],
		Pool:     options["pool"],
		Size:     options["size"],
		Readonly: options["readonly"] == "true",
		Optional: options["optional"] == "true",
	}, nil
}
