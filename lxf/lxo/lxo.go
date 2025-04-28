package lxo

import (
	lxd "github.com/canonical/lxd/client"
)

// LXO abstracts some of the lxd calls with additional functionality like retrying, idempotency
// and some level of error recovery. Usage stays the same as lxd.InstanceServer
type LXO struct {
	server lxd.InstanceServer
}

// New creates LXO
func NewClient(server lxd.InstanceServer) *LXO {
	return &LXO{
		server: server,
	}
}
