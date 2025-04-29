package lxo

import (
	"errors"
	"testing"

	lxdfakes "github.com/automaticserver/lxe/fakes/lxd/client"
	"github.com/canonical/lxd/shared/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLXO_StopInstance_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceStateReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.StopInstance("foo", 10, 0)
	require.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateInstanceStateCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_StopInstance_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceStateReturns(fakeOp, errors.New("something failed"))

	err := lxo.StopInstance("foo", 10, 0)
	require.Error(t, err)

	assert.Equal(t, 1, fake.UpdateInstanceStateCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_StopInstance_ForceSuccess(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceStateReturns(fakeOp, nil)
	fakeOp.WaitReturnsOnCall(0, errors.New("some error"))
	fakeOp.WaitReturnsOnCall(1, nil)

	err := lxo.StopInstance("foo", 5, 1)
	require.NoError(t, err)

	assert.Equal(t, 2, fake.UpdateInstanceStateCallCount())
	assert.Equal(t, 2, fakeOp.WaitCallCount())
}

func TestLXO_StopInstance_ForceFailed(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceStateReturns(fakeOp, nil)
	fakeOp.WaitReturnsOnCall(0, errors.New("some error"))
	fakeOp.WaitReturnsOnCall(1, errors.New("still error"))

	err := lxo.StopInstance("foo", 5, 1)
	require.Error(t, err)

	assert.Equal(t, 2, fake.UpdateInstanceStateCallCount())
	assert.Equal(t, 2, fakeOp.WaitCallCount())
}

func TestLXO_StopInstance_AlreadyStopped(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceStateReturns(fakeOp, nil)
	fakeOp.WaitReturnsOnCall(0, errors.New("The instance is already stopped"))

	err := lxo.StopInstance("foo", 5, 1)
	require.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateInstanceStateCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_StartInstance_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceStateReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.StartInstance("foo")
	require.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateInstanceStateCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_StartInstance_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceStateReturns(fakeOp, errors.New("something missing"))

	err := lxo.StartInstance("foo")
	require.Error(t, err)

	assert.Equal(t, 1, fake.UpdateInstanceStateCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_CreateInstance_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.CreateInstanceReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.CreateInstance(api.InstancesPost{})
	require.NoError(t, err)

	assert.Equal(t, 1, fake.CreateInstanceCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_CreateInstance_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.CreateInstanceReturns(fakeOp, errors.New("something failed"))

	err := lxo.CreateInstance(api.InstancesPost{})
	require.Error(t, err)

	assert.Equal(t, 1, fake.CreateInstanceCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_UpdateInstance_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.UpdateInstance("foo", api.InstancePut{}, "")
	require.NoError(t, err)

	assert.Equal(t, 1, fake.UpdateInstanceCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_UpdateInstance_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.UpdateInstanceReturns(fakeOp, errors.New("something failed"))

	err := lxo.UpdateInstance("foo", api.InstancePut{}, "")
	require.Error(t, err)

	assert.Equal(t, 1, fake.UpdateInstanceCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}

func TestLXO_DeleteInstance_Simple(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.DeleteInstanceReturns(fakeOp, nil)
	fakeOp.WaitReturns(nil)

	err := lxo.DeleteInstance("foo")
	require.NoError(t, err)

	assert.Equal(t, 1, fake.DeleteInstanceCallCount())
	assert.Equal(t, 1, fakeOp.WaitCallCount())
}

func TestLXO_DeleteInstance_Error(t *testing.T) {
	t.Parallel()

	lxo, fake := newFakeClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.DeleteInstanceReturns(fakeOp, errors.New("something failed"))

	err := lxo.DeleteInstance("foo")
	require.Error(t, err)

	assert.Equal(t, 1, fake.DeleteInstanceCallCount())
	assert.Equal(t, 0, fakeOp.WaitCallCount())
}
