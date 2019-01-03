package lxf

import "github.com/lxc/lxe/lxf/device"

// nolint: gosec #nosec (no sensitive data)

const (
	// ErrorMultipleFound is the error string in cases, where it's unexpected to find multiple
	ErrorMultipleFound = "multiple found"
	// ErrorNotFound is the error strin gin cases, where it's unexpected to find nothing
	ErrorNotFound = "not found"
)

// LXDObject contains common properties of containers and sandboxes without CRI influence
type LXDObject struct {
	// ID is a unique generated ID and is read-only
	ID string
	// ETag uniquely identifies user modifiable content of this resource, prevents race conditions when saving
	// see: https://lxd.readthedocs.io/en/latest/api-extensions/#etag
	ETag string
	// ETag uniquely identifies user modifiable content of the state of this resource, prevents race conditions when
	// trying to modify the state
	//StateETag string
	// Devices
	Proxies []device.Proxy
	Disks   []device.Disk
	Blocks  []device.Block
	Nics    []device.Nic
	// Config contains options not provided by a own property
	Config map[string]string
}
