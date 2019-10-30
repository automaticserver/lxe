// nolint: dupl
package device

import "fmt"

const (
	BlockType = "unix-block"
)

// Block device representation https://lxd.readthedocs.io/en/latest/containers/#type-unix-block
type Block struct {
	KeyName string
	Path    string
	Source  string
}

func (d *Block) getName() string {
	var name string

	switch {
	case d.KeyName != "":
		name = d.KeyName
	case d.Path == "":
		name = fmt.Sprintf("%s-%s", BlockType, d.Source)
	default:
		name = fmt.Sprintf("%s-%s", BlockType, d.Path)
	}

	return name
}

// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map
func (d *Block) ToMap() (string, map[string]string) {
	return d.getName(), map[string]string{
		"type":   BlockType,
		"source": d.Source,
		"path":   d.Path,
	}
}

// FromMap creates a new device with assigned name (can be empty) and options
func (d *Block) FromMap(name string, options map[string]string) (Device, error) {
	return &Block{
		KeyName: name,
		Path:    options["path"],
		Source:  options["source"],
	}, nil
}
