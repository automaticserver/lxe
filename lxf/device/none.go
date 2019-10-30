// nolint: dupl
package device

const (
	NoneType = "none"
)

// None device representation https://lxd.readthedocs.io/en/latest/containers/#type-none
type None struct {
	KeyName string
}

func (d *None) getName() string {
	return d.KeyName
}

// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map
func (d *None) ToMap() (string, map[string]string) {
	return d.getName(), map[string]string{
		"type": NoneType,
	}
}

// FromMap creates a new device with assigned name (can be empty) and options
func (d *None) FromMap(name string, options map[string]string) (Device, error) {
	return &None{
		KeyName: name,
	}, nil
}
