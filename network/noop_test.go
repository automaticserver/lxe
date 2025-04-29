package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitPluginNoop(t *testing.T) {
	t.Parallel()

	plugin, err := InitPluginNoop()
	require.NoError(t, err)
	assert.NotNil(t, plugin)
}

func Test_noopPlugin_PodNetwork(t *testing.T) {
	t.Parallel()

	plugin := &noopPlugin{}
	podNet, err := plugin.PodNetwork("", nil)
	require.NoError(t, err)
	assert.NotNil(t, podNet)
}

func Test_noopPlugin_Status(t *testing.T) {
	t.Parallel()

	plugin := &noopPlugin{}
	err := plugin.Status()
	require.Error(t, err)
}

func Test_noopPlugin_UpdateRuntimeConfig(t *testing.T) {
	t.Parallel()

	plugin := &noopPlugin{}
	err := plugin.UpdateRuntimeConfig(nil)
	require.Error(t, err)
}

func Test_noopPodNetwork_ContainerNetwork(t *testing.T) {
	t.Parallel()

	podNet := &noopPodNetwork{}
	contNet, err := podNet.ContainerNetwork("", nil)
	require.NoError(t, err)
	assert.NotNil(t, contNet)
}

func Test_noopPodNetwork_Status(t *testing.T) {
	t.Parallel()

	podNet := &noopPodNetwork{}
	status, err := podNet.Status(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, status)
}

func Test_noopPodNetwork_WhenCreated(t *testing.T) {
	t.Parallel()

	podNet := &noopPodNetwork{}
	res, err := podNet.WhenCreated(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, res)
}

func Test_noopPodNetwork_WhenStarted(t *testing.T) {
	t.Parallel()

	podNet := &noopPodNetwork{}
	res, err := podNet.WhenStarted(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, res)
}

func Test_noopPodNetwork_WhenStopped(t *testing.T) {
	t.Parallel()

	podNet := &noopPodNetwork{}
	err := podNet.WhenStopped(ctx, nil)
	require.NoError(t, err)
}

func Test_noopPodNetwork_WhenDeleted(t *testing.T) {
	t.Parallel()

	podNet := &noopPodNetwork{}
	err := podNet.WhenDeleted(ctx, nil)
	require.NoError(t, err)
}

func Test_noopContainerNetwork_WhenCreated(t *testing.T) {
	t.Parallel()

	contNet := &noopContainerNetwork{}
	res, err := contNet.WhenCreated(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, res)
}

func Test_noopContainerNetwork_WhenStarted(t *testing.T) {
	t.Parallel()

	contNet := &noopContainerNetwork{}
	res, err := contNet.WhenStarted(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, res)
}

func Test_noopContainerNetwork_WhenStopped(t *testing.T) {
	t.Parallel()

	contNet := &noopContainerNetwork{}
	err := contNet.WhenStopped(ctx, nil)
	require.NoError(t, err)
}

func Test_noopContainerNetwork_WhenDeleted(t *testing.T) {
	t.Parallel()

	contNet := &noopContainerNetwork{}
	err := contNet.WhenDeleted(ctx, nil)
	require.NoError(t, err)
}
