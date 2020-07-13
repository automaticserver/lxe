package lxf // import "github.com/automaticserver/lxe/lxf"

import (
	"encoding/base32"
)

var (
	b32lowerEncoder = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")
)

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
