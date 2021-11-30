package cloudinit

// NetworkConfig is used as root element to serialize to cloud config
type NetworkConfig struct {
	Version int           `json:"version"`
	Config  []interface{} `json:"config"`
}

// NetworkConfigEntry is an entry in the v1 network config
type NetworkConfigEntry struct {
	Type string `json:"type"`
}

// NetworkConfigEntryNameserver is a nameserver entry
type NetworkConfigEntryNameserver struct {
	NetworkConfigEntry
	Address []string `json:"address"`
	Search  []string `json:"search"`
}

// NetworkConfigEntryPhysical is a nameserver entry
type NetworkConfigEntryPhysical struct {
	NetworkConfigEntry
	Name    string                             `json:"name"`
	Subnets []NetworkConfigEntryPhysicalSubnet `json:"subnets"`
}

// NetworkConfigEntryPhysicalSubnet is a subnet entry in the v1 network config of a physical device
type NetworkConfigEntryPhysicalSubnet struct {
	Type string `json:"type"`
}
