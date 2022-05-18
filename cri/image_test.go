package cri

import (
	"context"
	"testing"

	"github.com/automaticserver/lxe/cri/crifakes"
	"github.com/stretchr/testify/assert"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

var (
	ctx = context.TODO()
)

func testImageServer() (*ImageServer, *crifakes.FakeClient) {
	fake := &crifakes.FakeClient{}

	return &ImageServer{
		lxf: fake,
	}, fake
}

func TestImageServer_PullImage(t *testing.T) {
	t.Parallel()

	s, fake := testImageServer()

	fake.PullImageReturns("something", nil)

	resp, err := s.PullImage(ctx, &rtApi.PullImageRequest{
		Image: &rtApi.ImageSpec{
			Image: "an/image",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, fake.PullImageCallCount())
	assert.Equal(t, "an/image", fake.PullImageArgsForCall(0))
	assert.Equal(t, "something", resp.ImageRef)
}
