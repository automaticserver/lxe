package lxf

import (
	"fmt"
)

const (
	// ErrorLXDNotFound is the error string a LXD request returns, when nothing is found
	// Unfortunately there is no constant in the lxd source we could've used
	ErrorLXDNotFound = "not found"
)

type lxfError struct {
	ID     string
	Reason error
}

func (e lxfError) error(name string) string {
	return fmt.Sprintf("%s %s: %s", name, e.Reason.Error(), e.ID)
}

func newLxfError(id string, reason error) lxfError {
	return lxfError{
		ID:     id,
		Reason: reason,
	}
}

// ContainerError is an error type for errors related to containers
type ContainerError struct {
	lxfError
}

func (e ContainerError) Error() string {
	return e.error("container")
}

// NewContainerError creates a new SandboxError
func NewContainerError(id string, reason error) ContainerError {
	return ContainerError{
		lxfError: newLxfError(id, reason),
	}
}

// IsContainerNotFound checks if error is of type ContainerError where its previous error was not found in LXD
func IsContainerNotFound(err error) bool {
	if serr, ok := err.(ContainerError); ok {
		return serr.Reason.Error() == ErrorLXDNotFound
	}

	return false
}

// SandboxError is an error type for errors related to sandboxes
type SandboxError struct {
	lxfError
}

func (e SandboxError) Error() string {
	return e.error("sandbox")
}

// NewSandboxError creates a new SandboxError
func NewSandboxError(id string, reason error) SandboxError {
	return SandboxError{
		lxfError: newLxfError(id, reason),
	}
}

// IsSandboxNotFound checks if error is of type SandboxError where its previous error was not found in LXD
func IsSandboxNotFound(err error) bool {
	if serr, ok := err.(SandboxError); ok {
		return serr.Reason.Error() == ErrorLXDNotFound
	}

	return false
}

// ImageError is an error type for errors related to images
type ImageError struct {
	lxfError
}

func (e ImageError) Error() string {
	return e.error("image")
}

// NewImageError creates a new ImageError
func NewImageError(id string, reason error) ImageError {
	return ImageError{
		lxfError: newLxfError(id, reason),
	}
}

// IsImageNotFound checks if error is of type ImageError where its previous error was not found in LXD
func IsImageNotFound(err error) bool {
	if serr, ok := err.(ImageError); ok {
		return serr.Reason.Error() == ErrorLXDNotFound
	}

	return false
}
