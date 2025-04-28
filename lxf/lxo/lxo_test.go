package lxo

import (
	"testing"

	lxdfakes "github.com/automaticserver/lxe/fakes/lxd/client"
	"github.com/stretchr/testify/assert"
)

func newFakeClient() (*LXO, *lxdfakes.FakeInstanceServer) {
	fake := &lxdfakes.FakeInstanceServer{}

	return &LXO{
		server: fake,
	}, fake
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	fake := &lxdfakes.FakeInstanceServer{}

	lxo := NewClient(fake)
	assert.NotNil(t, lxo)

	assert.Exactly(t, fake, lxo.server)
}
