package lxo

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/automaticserver/lxe/lxf/lxdfakes"
)

func newFakeClient() (*LXO, *lxdfakes.FakeContainerServer) {
	fake := &lxdfakes.FakeContainerServer{}

	return &LXO{
		server: fake,
	}, fake
}

func TestNewClient(t *testing.T) {
	fake := &lxdfakes.FakeContainerServer{}

	lxo := NewClient(fake)
	assert.NotNil(t, lxo)

	assert.Exactly(t, fake, lxo.server)
}
