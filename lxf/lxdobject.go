package lxf

import "github.com/lxc/lxe/lxf/device"

// nolint: gosec #nosec (no sensitive data)

const (
	// ErrorMultipleFound is the error string in cases, where it's unexpected to find multiple
	// ErrorMultipleFound = "multiple found"

	// ErrorLXDNotFound is the error string a LXD request returns, when nothing is found
	// Unfortunately there is no constant in the lxd source we could've used
	ErrorLXDNotFound = "not found"
)

// LXDObject contains common properties of containers and sandboxes without CRI influence
type LXDObject struct {
	// client holds the lxf.Client representing as a lxd client
	client *Client
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
