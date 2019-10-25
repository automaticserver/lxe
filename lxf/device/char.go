package device

const (
	CharType = "char"
)

// Chars holds slice of Char
// Use it if you want to Add() a entry non-conflicting (see Add())
type Chars []Char

// Add a entry to the slice, if the name is the same, will overwrite the existing entry
func (ds *Chars) Add(d Char) {
	for k, e := range *ds {
		if e.GetName() == d.GetName() {
			(*ds)[k] = d
			return
		}
	}

	*ds = append(*ds, d)
}

// Char mounts a host path into the container
type Char struct {
	Path   string
	Source string
}

// ToMap serializes itself into a map. Will return an error if the data
// is inconsistent/invalid in some way
func (d Char) ToMap() (map[string]string, error) {
	return map[string]string{
		"type":   CharType,
		"path":   d.Path,
		"source": d.Source,
	}, nil
}

// GetName will return the path with prefix
func (d Char) GetName() string {
	suffix := d.Path
	if d.Path == "" {
		suffix = d.Source
	}

	return CharType + "-" + suffix
}

// CharFromMap crrate a new char from map entries
func CharFromMap(dev map[string]string) (Char, error) {
	return Char{
		Path:   dev["path"],
		Source: dev["source"],
	}, nil
}
