package lxo

import (
	"testing"

	"github.com/automaticserver/lxe/lxf/lxdfakes"
	"github.com/stretchr/testify/assert"
)

func newFakeClient() (*LXO, *lxdfakes.FakeContainerServer) {
	fake := &lxdfakes.FakeContainerServer{}

	return &LXO{
		server: fake,
	}, fake
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	fake := &lxdfakes.FakeContainerServer{}

	lxo := NewClient(fake)
	assert.NotNil(t, lxo)

	assert.Exactly(t, fake, lxo.server)
}
