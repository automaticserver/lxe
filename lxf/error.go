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

type ContainerError struct {
	*lxfError
}

func (e *ContainerError) Error() string {
	return e.error("container")
}

func NewContainerError(id string, reason error) *ContainerError {
	return &ContainerError{
		lxfError: newLxfError(id, reason),
	}
}

type SandboxError struct {
	*lxfError
}

func (e *SandboxError) Error() string {
	return e.error("sandbox")
}

func NewSandboxError(id string, reason error) *SandboxError {
	return &SandboxError{
		lxfError: newLxfError(id, reason),
	}
}

type ImageError struct {
	*lxfError
}

func (e *ImageError) Error() string {
	return e.error("image")
}

func NewImageError(id string, reason error) *ImageError {
	return &ImageError{
		lxfError: newLxfError(id, reason),
	}
}
