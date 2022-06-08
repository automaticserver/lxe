package lxf

import (
	"errors"
	"net/http"

	"github.com/lxc/lxd/shared/api"
)

// ErrNotFound for CRI related checks
var ErrNotFound = errors.New("not found")

// Error compatibility layer for LXD and custom error for CRI specific checks. The LXD client returns errors as `api.StatusErrorf(resp.StatusCode, response.Error)` where StatusCode is from net/http and Error the string in the response. LXF returns errors defined in this package. Clients of this package should use these functions if they don't want to differentiate these error sources.

// Whether err is a not found error from lxd api response or this package
func IsNotFoundError(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}

	return api.StatusErrorCheck(err, http.StatusNotFound)
}
