package lxf

import (
	"net/http"
	"testing"

	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

// Image not found error compatibility handling for lxd <= 5.0.0 (Issue #16). It was wrong to check for the error text "not found", as this has now changed during lxd 5.x development.
//
// https://github.com/lxc/lxd/blob/lxd-5.0.0/lxd/db/images.go#L768
// > return ErrNoSuchObject
// https://github.com/lxc/lxd/blob/lxd-5.0.0/lxd/db/errors.go#L18
// > ErrNoSuchObject = fmt.Errorf("No such object")
// https://github.com/lxc/lxd/blob/lxd-5.0.0/lxd/response/smart.go#L14
// > http.StatusNotFound:  {os.ErrNotExist, sql.ErrNoRows, db.ErrNoSuchObject},
// https://github.com/lxc/lxd/blob/lxd-5.0.0/lxd/response/smart.go#L40
// > return &errorResponse{httpStatusCode, http.StatusText(httpStatusCode)}
// https://github.com/lxc/lxd/blob/lxd-5.0.0/client/lxd.go#L228
// > return nil, "", api.StatusErrorf(resp.StatusCode, response.Error)
//
// errors can return status text "not found" but will have http.NotFound as code. This test succeeds if it executes the second call to the fake server, which is GetImage().
func Test_getRemoteImageFromAliasOrFingerprint_AliasNotFoundPreLXD5(t *testing.T) {
	t.Parallel()

	c, fake := testClient()

	fake.GetImageAliasReturns(nil, "", api.StatusErrorf(http.StatusNotFound, http.StatusText(http.StatusNotFound)))
	fake.GetImageReturns(&api.Image{Fingerprint: "abcdefg"}, "", nil)

	resp, err := getRemoteImageFromAliasOrFingerprint(c.server, "ubuntu:nextgen")

	assert.NoError(t, err)
	assert.Equal(t, "abcdefg", resp.Fingerprint)
	assert.Equal(t, 1, fake.GetImageAliasCallCount())
	assert.Equal(t, 1, fake.GetImageCallCount())
}

// Image not found error compatibility handling for lxd >= 5.1 (Issue #16). It was wrong to check for the error text "not found", as this has now changed during lxd 5.x development.
//
// https://github.com/lxc/lxd/blob/lxd-5.1/lxd/db/images.go#L769
// > return api.StatusErrorf(http.StatusNotFound, "Image alias not found")
// https://github.com/lxc/lxd/blob/lxd-5.1/lxd/response/smart.go#L27
// > return &errorResponse{statusCode, err.Error()}
// https://github.com/lxc/lxd/blob/lxd-5.1/client/lxd.go#L228
// > return nil, "", api.StatusErrorf(resp.StatusCode, response.Error)
//
// That error now definitely returns the error text with http.NotFound as code. This test succeeds if it executes the second call to the fake server, which is GetImage().
func Test_getRemoteImageFromAliasOrFingerprint_AliasNotFoundPostLXD5(t *testing.T) {
	t.Parallel()

	c, fake := testClient()

	fake.GetImageAliasReturns(nil, "", api.StatusErrorf(http.StatusNotFound, "Image alias not found"))
	fake.GetImageReturns(&api.Image{Fingerprint: "abcdefg"}, "", nil)

	resp, err := getRemoteImageFromAliasOrFingerprint(c.server, "ubuntu:nextgen")

	assert.NoError(t, err)
	assert.Equal(t, "abcdefg", resp.Fingerprint)
	assert.Equal(t, 1, fake.GetImageAliasCallCount())
	assert.Equal(t, 1, fake.GetImageCallCount())
}

// Image not found error compatibility handling for lxd <= 5.0.0 (Issue #16). It was wrong to check for the error text "not found", as this has now changed during lxd 5.x development.
//
// https://github.com/lxc/lxd/blob/lxd-5.0.0/lxd/db/images.go#L485
// > return api.StatusErrorf(http.StatusNotFound, "Image not found")
// https://github.com/lxc/lxd/blob/lxd-5.0.0/lxd/response/smart.go#L40
// > return &errorResponse{httpStatusCode, http.StatusText(httpStatusCode)}
// https://github.com/lxc/lxd/blob/lxd-5.0.0/client/lxd.go#L228
// > return nil, "", api.StatusErrorf(resp.StatusCode, response.Error)
func Test_getRemoteImageFromAliasOrFingerprint_ImageNotFoundPreLXD5(t *testing.T) {
	t.Parallel()

	c, fake := testClient()

	fake.GetImageAliasReturns(nil, "", api.StatusErrorf(http.StatusNotFound, http.StatusText(http.StatusNotFound)))
	fake.GetImageReturns(nil, "", api.StatusErrorf(http.StatusNotFound, "Image not found"))

	resp, err := getRemoteImageFromAliasOrFingerprint(c.server, "ubuntu:nextgen")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, IsNotFoundError(err))
	assert.Equal(t, 1, fake.GetImageAliasCallCount())
	assert.Equal(t, 1, fake.GetImageCallCount())
}

// Image not found error compatibility handling for lxd >= 5.1 (Issue #16). It was wrong to check for the error text "not found", as this has now changed during lxd 5.x development.
//
// https://github.com/lxc/lxd/blob/lxd-5.1/lxd/db/images.go#L486
// > return api.StatusErrorf(http.StatusNotFound, "Image not found")
// https://github.com/lxc/lxd/blob/lxd-5.1/lxd/response/smart.go#L27
// > return &errorResponse{statusCode, err.Error()}
// https://github.com/lxc/lxd/blob/lxd-5.1/client/lxd.go#L228
// > return nil, "", api.StatusErrorf(resp.StatusCode, response.Error)
func Test_getRemoteImageFromAliasOrFingerprint_ImageNotFoundPostLXD5(t *testing.T) {
	t.Parallel()

	c, fake := testClient()

	fake.GetImageAliasReturns(nil, "", api.StatusErrorf(http.StatusNotFound, "Image alias not found"))
	fake.GetImageReturns(nil, "", api.StatusErrorf(http.StatusNotFound, "Image not found"))

	resp, err := getRemoteImageFromAliasOrFingerprint(c.server, "ubuntu:nextgen")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, IsNotFoundError(err))
	assert.Equal(t, 1, fake.GetImageAliasCallCount())
	assert.Equal(t, 1, fake.GetImageCallCount())
}
