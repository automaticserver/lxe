package cri

import (
	"context"
	"testing"

	crifakes "github.com/automaticserver/lxe/fakes/lxe/lxf"
	"github.com/canonical/lxd/lxc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// implements interface
var _ rtApi.ImageServiceServer = &ImageServer{}

var (
	ctx = context.TODO()
)

func testImageServer() (*ImageServer, *crifakes.FakeClient) {
	fake := &crifakes.FakeClient{}

	return &ImageServer{
		lxf:       fake,
		lxdConfig: config.DefaultConfig(),
	}, fake
}

func Test_ImageServer_PullImage(t *testing.T) {
	t.Parallel()

	s, fake := testImageServer()

	fake.PullImageReturns("something", nil)

	resp, err := s.PullImage(ctx, &rtApi.PullImageRequest{
		Image: &rtApi.ImageSpec{
			Image: "ubuntu/nextgen",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 1, fake.PullImageCallCount())
	assert.Equal(t, "something", resp.ImageRef)
	assert.Equal(t, "ubuntu:nextgen", fake.PullImageArgsForCall(0))
}
