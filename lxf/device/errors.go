package device

import (
	"errors"
)

var (
	ErrNotSupported = errors.New("not supported")
	ErrNotValid     = errors.New("not valid")
)
