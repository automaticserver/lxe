package lxf

import (
	"fmt"
)

type lxfError struct {
	ID     string
	Reason error
}

func (e *lxfError) error(name string) string {
	return fmt.Sprintf("%s %s: %s", name, e.Reason.Error(), e.ID)
}

func newLxfError(id string, reason error) *lxfError {
	return &lxfError{
		ID:     id,
		Reason: reason,
	}
}

// ContainerError is an error type for errors related to containers
type ContainerError struct {
	*lxfError
}

func (e *ContainerError) Error() string {
	return e.error("container")
}

// NewContainerError creates a new SandboxError
func NewContainerError(id string, reason error) *ContainerError {
	return &ContainerError{
		lxfError: newLxfError(id, reason),
	}
}

// SandboxError is an error type for errors related to sandboxes
type SandboxError struct {
	*lxfError
}

func (e *SandboxError) Error() string {
	return e.error("sandbox")
}

// NewSandboxError creates a new SandboxError
func NewSandboxError(id string, reason error) *SandboxError {
	return &SandboxError{
		lxfError: newLxfError(id, reason),
	}
}

// ImageError is an error type for errors related to images
type ImageError struct {
	*lxfError
}

func (e *ImageError) Error() string {
	return e.error("image")
}

// NewImageError creates a new ImageError
func NewImageError(id string, reason error) *ImageError {
	return &ImageError{
		lxfError: newLxfError(id, reason),
	}
}
