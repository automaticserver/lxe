package lxo

import (
	"errors"
	"testing"

	"github.com/automaticserver/lxe/lxf/lxdfakes"
	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func TestLXO_CopyImage_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeRemoteOperation{}
	sourceFake := &lxdfakes.FakeImageServer{}

	fake.CopyImageReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.CopyImage(sourceFake, api.Image{}, nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.CopyImageCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_CopyImage_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeRemoteOperation{}
	sourceFake := &lxdfakes.FakeImageServer{}

	fake.CopyImageReturns(fakeOp, errors.New("something failed"))

	err := lxo.CopyImage(sourceFake, api.Image{}, nil)
	assert.Error(t, err)

	assert.Equal(t, 1, fake.CopyImageCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_DeleteImage_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.DeleteImageReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.DeleteImage("foo")
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.DeleteImageCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_DeleteImage_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.DeleteImageReturns(fakeOp, errors.New("something failed"))

	err := lxo.DeleteImage("foo")
	assert.Error(t, err)

	assert.Equal(t, 1, fake.DeleteImageCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_CreateImage_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.CreateImageReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)
	fakeOp.GetReturns(api.Operation{Metadata: map[string]any{"fingerprint": "abcdefg"}})

	fingerprint, err := lxo.CreateImage(api.ImagesPost{}, nil)
	assert.NoError(t, err)

	assert.Equal(t, "abcdefg", fingerprint)
	assert.Equal(t, 1, fake.CreateImageCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
	assert.Equal(t, 1, fakeOp.GetCallCount())
}

func TestLXO_CreateImage_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.CreateImageReturns(fakeOp, errors.New("something failed"))

	fingerprint, err := lxo.CreateImage(api.ImagesPost{}, nil)
	assert.Error(t, err)

	assert.Equal(t, "", fingerprint)
	assert.Equal(t, 1, fake.CreateImageCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
	assert.Equal(t, 0, fakeOp.GetCallCount())
}
