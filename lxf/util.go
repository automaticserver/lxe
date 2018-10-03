package lxf

import (
	"bytes"
	"encoding/base32"
	"fmt"
)

var (
	b32lowerEncoder = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")
)

// WriteCloserBuffer decorates a byte buffer with the Closer interface.
type WriteCloserBuffer struct {
	*bytes.Buffer
}

// EmptyAnnotationWarning is returned when the annotation given on call for the function
// is actually empty.
type EmptyAnnotationWarning struct {
	Where string
}

func (e *EmptyAnnotationWarning) Error() string {
	return fmt.Sprintf("Empty Annotation found for %s", e.Where)
}

// Close does nothing
func (m WriteCloserBuffer) Close() error {
	fmt.Println("closed closer")
	return nil
}

// NewWriteCloserBuffer creates a write closer buffer
func NewWriteCloserBuffer() *WriteCloserBuffer {
	return &WriteCloserBuffer{&bytes.Buffer{}}
}

// SetIfSet sets a key in a map[string]string with the value, if the value is not empty
func SetIfSet(s *map[string]string, where string, what string) error {
	if what != "" {
		(*s)[where] = what
		return nil
	}
	return &EmptyAnnotationWarning{where}
}
