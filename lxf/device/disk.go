package device

import (
	"strconv"
)

const (
	DiskType = "disk"
)

// Disks holds slice of Disk
// Use it if you want to Add() a entry non-conflicting (see Add())
type Disks []Disk

// Add a entry to the slice, if the name is the same, will overwrite the existing entry
func (ds *Disks) Add(d Disk) {
	for k, e := range *ds {
		if e.GetName() == d.GetName() {
			(*ds)[k] = d
			return
		}
	}

	*ds = append(*ds, d)
}

// Disk mounts a host path into the container
type Disk struct {
	Path     string
	Pool     string
	Source   string
	Size     string
	Readonly bool
	Optional bool
}

// ToMap serializes itself into a map. Will return an error if the data
// is inconsistent/invalid in some way
func (d Disk) ToMap() (map[string]string, error) {
	def := map[string]string{
		"type":     DiskType,
		"path":     d.Path,
		"source":   d.Source,
		"readonly": strconv.FormatBool(d.Readonly),
		"optional": strconv.FormatBool(d.Optional),
	}

	if d.Pool != "" {
		def["pool"] = d.Pool
	}

	if d.Size != "" {
		def["size"] = d.Size
	}

	return def, nil
}

// GetName will return the path with prefix
func (d Disk) GetName() string {
	return DiskType + "-" + d.Path
}

// DiskFromMap crrate a new disk from map entries
func DiskFromMap(dev map[string]string) (Disk, error) {
	return Disk{
		Path:     dev["path"],
		Source:   dev["source"],
		Pool:     dev["pool"],
		Readonly: dev["readonly"] == "true",
		Optional: dev["optional"] == "true",
	}, nil
}
