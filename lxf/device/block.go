package device

const (
	blockType = "unix-block"
)

// Blocks holds slice of Block
// Use it if you want to Add() a entry non-conflicting (see Add())
type Blocks []Block

// Add a entry to the slice, if the name is the same, will overwrite the existing entry
func (bs Blocks) Add(b Block) {
	for k, e := range bs {
		if e.GetName() == b.GetName() {
			bs[k] = b
			return
		}
	}
	bs = append(bs, b)
}

// Block device
type Block struct {
	Path   string
	Source string
}

// ToMap serializes itself into a lxd device map entry
func (b Block) ToMap() (map[string]string, error) {
	return map[string]string{
		"type":   blockType,
		"source": b.Source,
		"path":   b.Path,
	}, nil
}

// GetName will generate a uinique name for the device map
func (b Block) GetName() string {
	return blockType + "-" + b.Path
}

// BlockFromMap create a new block from map entries
func BlockFromMap(dev map[string]string) (Block, error) {
	return Block{
		Path:   dev["path"],
		Source: dev["source"],
	}, nil
}
