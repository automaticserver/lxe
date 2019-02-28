package device

const (
	nicType = "nic"
)

// Nics holds slice of Nic
// Use it if you want to Add() a entry non-conflicting (see Add())
type Nics []Nic

// Add a entry to the slice, if the name is the same, will overwrite the existing entry
func (ns *Nics) Add(n Nic) {
	for k, e := range *ns {
		if e.GetName() == n.GetName() {
			(*ns)[k] = n
			return
		}
	}
	*ns = append(*ns, n)
}

// Nic device
type Nic struct {
	Name        string
	NicType     string
	Parent      string
	IPv4Address string
}

// ToMap serializes itself into a lxd device map entry
func (b Nic) ToMap() (map[string]string, error) {
	return map[string]string{
		"type":         nicType,
		"name":         b.Name,
		"nictype":      b.NicType,
		"parent":       b.Parent,
		"ipv4.address": b.IPv4Address,
	}, nil
}

// GetName will generate a uinique name for the device map
func (b Nic) GetName() string {
	return nicType + "-" + b.Name
}

// NicFromMap create a new nic from map entries
func NicFromMap(dev map[string]string) (Nic, error) {
	return Nic{
		Name:        dev["name"],
		NicType:     dev["nictype"],
		Parent:      dev["parent"],
		IPv4Address: dev["ipv4.address"],
	}, nil
}
