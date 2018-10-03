package device

const blockType = "unix-block"

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
	return blockType + "-" + b.Source
}

// BlockFromMap create a new block from map entries
func BlockFromMap(dev map[string]string) (Block, error) {
	return Block{
		Path:   dev["path"],
		Source: dev["source"],
	}, nil
}
