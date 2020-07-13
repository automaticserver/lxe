package device // import "github.com/automaticserver/lxe/lxf/device"

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

// FromMap loads assigned name (can be empty) and options
func (d *None) FromMap(name string, options map[string]string) error {
	d.KeyName = name
	return nil
}

// New creates a new empty device
func (d *None) new() Device {
	return &None{}
}
