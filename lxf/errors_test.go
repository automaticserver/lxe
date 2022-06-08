package lxf

import (
	"net/http"
	"testing"

	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func Test_IsNotFoundError_API(t *testing.T) {
	t.Parallel()

	err := api.StatusErrorf(http.StatusNotFound, "Image not found")
	is := IsNotFoundError(err)
	assert.True(t, is)
}

func Test_IsNotFoundError_APIOther(t *testing.T) {
	t.Parallel()

	err := api.StatusErrorf(http.StatusForbidden, "Image not found")
	is := IsNotFoundError(err)
	assert.False(t, is)
}

func Test_IsNotFoundError_LXF(t *testing.T) {
	t.Parallel()

	err := ErrNotFound
	is := IsNotFoundError(err)
	assert.True(t, is)
}

func Test_IsNotFoundError_Other(t *testing.T) {
	t.Parallel()

	err := ErrUsage
	is := IsNotFoundError(err)
	assert.False(t, is)
}
