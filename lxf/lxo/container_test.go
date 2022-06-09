package lxo

import (
	"errors"
	"testing"

	lxdfakes "github.com/automaticserver/lxe/fakes/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func TestLXO_StopContainer_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerStateReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.StopContainer("foo", 10, 0)
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateContainerStateCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_StopContainer_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerStateReturns(fakeOp, errors.New("something failed"))

	err := lxo.StopContainer("foo", 10, 0)
	assert.Error(t, err)

	assert.Equal(t, 1, fake.UpdateContainerStateCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_StopContainer_ForceSuccess(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerStateReturns(fakeOp, nil)
	fakeOp.WaitReturnsOnCall(0, errors.New("some error"))
	fakeOp.WaitReturnsOnCall(1, nil)

	err := lxo.StopContainer("foo", 5, 1)
	assert.NoError(t, err)

	assert.Equal(t, 2, fake.UpdateContainerStateCallCount())
	assert.Equal(t, 2, fakeOp.WaitCallCount())
}

func TestLXO_StopContainer_ForceFailed(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerStateReturns(fakeOp, nil)
	fakeOp.WaitReturnsOnCall(0, errors.New("some error"))
	fakeOp.WaitReturnsOnCall(1, errors.New("still error"))

	err := lxo.StopContainer("foo", 5, 1)
	assert.Error(t, err)

	assert.Equal(t, 2, fake.UpdateContainerStateCallCount())
	assert.Equal(t, 2, fakeOp.WaitCallCount())
}

func TestLXO_StopContainer_AlreadyStopped(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerStateReturns(fakeOp, nil)
	fakeOp.WaitReturnsOnCall(0, errors.New("The container is already stopped"))

	err := lxo.StopContainer("foo", 5, 1)
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateContainerStateCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_StartContainer_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerStateReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.StartContainer("foo")
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateContainerStateCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_StartContainer_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerStateReturns(fakeOp, errors.New("something missing"))

	err := lxo.StartContainer("foo")
	assert.Error(t, err)

	assert.Equal(t, 1, fake.UpdateContainerStateCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_CreateContainer_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.CreateContainerReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.CreateContainer(api.ContainersPost{})
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.CreateContainerCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_CreateContainer_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.CreateContainerReturns(fakeOp, errors.New("something failed"))

	err := lxo.CreateContainer(api.ContainersPost{})
	assert.Error(t, err)

	assert.Equal(t, 1, fake.CreateContainerCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_UpdateContainer_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.UpdateContainer("foo", api.ContainerPut{}, "")
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateContainerCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_UpdateContainer_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateContainerReturns(fakeOp, errors.New("something failed"))

	err := lxo.UpdateContainer("foo", api.ContainerPut{}, "")
	assert.Error(t, err)

	assert.Equal(t, 1, fake.UpdateContainerCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_DeleteContainer_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.DeleteContainerReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.DeleteContainer("foo")
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.DeleteContainerCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_DeleteContainer_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.DeleteContainerReturns(fakeOp, errors.New("something failed"))

	err := lxo.DeleteContainer("foo")
	assert.Error(t, err)

	assert.Equal(t, 1, fake.DeleteContainerCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}
