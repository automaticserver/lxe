package device

const (
	NoneType = "none"
)

// None holds slice of None
// Use it if you want to Add() a entry non-conflicting (see Add())
type Nones []None

// Add a entry to the slice, if the name is the same, will overwrite the existing entry
func (ns *Nones) Add(n None) {
	for k, e := range *ns {
		if e.GetName() == n.GetName() {
			(*ns)[k] = n
			return
		}
	}

	*ns = append(*ns, n)
}

// None device allows disabling inherited devices
type None struct {
	// Name of the device to disable
	Name string
}

// ToMap serializes itself into a lxd device map entry
func (n None) ToMap() (map[string]string, error) {
	return map[string]string{
		"type": NoneType,
	}, nil
}

// GetName will generate a uinique name for the device map
func (n None) GetName() string {
	return n.Name
}

// NoneFromMap create a new none from map entries
func NoneFromMap(dev map[string]string, name string) (None, error) {
	return None{
		Name: name,
	}, nil
}
