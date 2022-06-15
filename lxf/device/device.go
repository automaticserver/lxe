package device

import (
	"fmt"
)

var (
	schema = map[string]Device{
		BlockType: &Block{},
		CharType:  &Char{},
		DiskType:  &Disk{},
		NicType:   &Nic{},
		NoneType:  &None{},
		ProxyType: &Proxy{},
	}
)

// Device must support mapping from lxd device map bidirectional
type Device interface {
	// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map. The key name must be passed through trimKeyName function to guarantee length compatibility.
	ToMap() (name string, options map[string]string)
	// FromMap loads assigned name (can be empty) and options
	FromMap(name string, options map[string]string) error
	// New creates a new empty device
	new() Device
}

// Detects and loads device by type
func Detect(name string, options map[string]string) (Device, error) {
	t, is := schema[options["type"]]
	if !is {
		return nil, fmt.Errorf("device type %w: %v", ErrNotSupported, options["type"])
	}

	d := t.new()

	err := d.FromMap(name, options)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Devices allows having a list of devices unique by name
type Devices []Device

// Upsert adds a device or overrides an entry if the key name exists
func (d *Devices) Upsert(a Device) {
	for k, e := range *d {
		eName, _ := e.ToMap()
		aName, _ := a.ToMap()

		if eName == aName {
			(*d)[k] = a

			return
		}
	}

	*d = append(*d, a)
}

// The maximum character length for device key names. Change to 64 once 5.0.1 is released (see below) as the effort to carry through the exact LXD version is currently not worth.
const (
	maxKeyNameLength             = 27
	middleSeparatorKeyNameLength = "--"
)

// Trims key name to allowed length by cutting in the middle. There is a bug in LXD 5.0.0 limiting device names to 27 characters (https://github.com/lxc/lxd/issues/10238). A PR now officially describes the allowed length to be 64 characters and the fix will be available in 5.0.1 and is already merged for 5.1 (https://github.com/lxc/lxd/pull/10251).
func trimKeyName(s string) string {
	if len(s) <= maxKeyNameLength {
		return s
	}

	partLen := maxKeyNameLength/2 - len(middleSeparatorKeyNameLength)/2 // nolint: gomnd

	// we can expect there are no multibyte characters in this string. By dividing two ints the result is still int and thus automatically floored
	left := s[:partLen]
	right := s[len(s)-partLen:]

	return fmt.Sprintf("%s%s%s", left, middleSeparatorKeyNameLength, right)
}
