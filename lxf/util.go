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
func SetIfSet(s *map[string]string, key, value string) {
	if value != "" {
		(*s)[key] = value
	}
}

// AppendIfSet sets a key in a map[string]string with the value, if the value is not empty. And if there was
// already a value, append it after a newline
func AppendIfSet(s *map[string]string, key, value string) {
	if value != "" {
		if (*s)[key] == "" {
			SetIfSet(s, key, value)
		} else {
			SetIfSet(s, key, (*s)[key]+"\n"+value)
		}
	}
}
