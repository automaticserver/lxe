// nolint: dupl
package device

import "fmt"

const (
	CharType = "unix-char"
)

// Char device representation https://lxd.readthedocs.io/en/latest/containers/#type-unix-char
type Char struct {
	KeyName string
	Path    string
	Source  string
}

func (d *Char) getName() string {
	var name string

	switch {
	case d.KeyName != "":
		name = d.KeyName
	case d.Path == "":
		name = fmt.Sprintf("%s-%s", CharType, d.Source)
	default:
		name = fmt.Sprintf("%s-%s", CharType, d.Path)
	}

	return name
}

// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map
func (d *Char) ToMap() (string, map[string]string) {
	return d.getName(), map[string]string{
		"type":   CharType,
		"source": d.Source,
		"path":   d.Path,
	}
}

// FromMap creates a new device with assigned name (can be empty) and options
func (d *Char) FromMap(name string, options map[string]string) (Device, error) {
	return &Char{
		KeyName: name,
		Path:    options["path"],
		Source:  options["source"],
	}, nil
}
