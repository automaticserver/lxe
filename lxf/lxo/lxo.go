package lxo

import (
	lxd "github.com/lxc/lxd/client"
)

// LXO abstracts some of the lxd calls with additional functionality like retrying, idempotency
// and some level of error recovery. Usage stays the same as lxd.ContainerServer
type LXO struct {
	server lxd.ContainerServer
}

// New creates LXO
func NewClient(server lxd.ContainerServer) *LXO {
	return &LXO{
		server: server,
	}
}
