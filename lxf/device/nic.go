package device

const nicType = "nic"

// Nic device
type Nic struct {
	Name    string
	NicType string
	Parent  string
}

// ToMap serializes itself into a lxd device map entry
func (b Nic) ToMap() (map[string]string, error) {
	return map[string]string{
		"type":    nicType,
		"name":    b.Name,
		"nictype": b.NicType,
		"parent":  b.Parent,
	}, nil
}

// GetName will generate a uinique name for the device map
func (b Nic) GetName() string {
	return nicType + "-" + b.Name
}

// NicFromMap create a new nic from map entries
func NicFromMap(dev map[string]string) (Nic, error) {
	return Nic{
		Name:    dev["name"],
		NicType: dev["nictype"],
		Parent:  dev["parent"],
	}, nil
}
