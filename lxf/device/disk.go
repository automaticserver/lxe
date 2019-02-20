package device

import (
	"strconv"
)

const diskType = "disk"

// Disk mounts a host path into the container
type Disk struct {
	Path     string
	Readonly bool
	Pool     string
	Source   string
	Optional bool
}

// ToMap serializes itself into a map. Will return an error if the data
// is inconsistent/invalid in some way
func (d Disk) ToMap() (map[string]string, error) {
	def := map[string]string{
		"type":     diskType,
		"path":     d.Path,
		"source":   d.Source,
		"readonly": strconv.FormatBool(d.Readonly),
		"optional": strconv.FormatBool(d.Optional),
	}
	if d.Pool != "" {
		def["pool"] = d.Pool
	}
	return def, nil
}

// GetName will return the path with prefix
func (d Disk) GetName() string {
	return diskType + "-" + d.Path
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
