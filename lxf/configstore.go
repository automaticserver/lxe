package lxf

import "strings"

// ConfigStore contains the rules to split a config into different
// maps.
type ConfigStore struct {
	reserved         []string
	reservedPrefixes []string
}

// NewConfigStore initialises a new ConfigStore
func NewConfigStore() *ConfigStore {
	return &ConfigStore{}
}

// WithReserved creates a new configstore with given reserved keys added
func (c *ConfigStore) WithReserved(keys ...string) *ConfigStore {
	return &ConfigStore{
		reserved:         append(append([]string{}, c.reserved...), keys...),
		reservedPrefixes: c.reservedPrefixes,
	}
}

// WithReservedPrefixes creates a new configstore with given reserved key prefixes added.
// A trailing dot will automatically be used to match only full namespaces.
func (c *ConfigStore) WithReservedPrefixes(keys ...string) *ConfigStore {
	return &ConfigStore{
		reserved:         c.reserved,
		reservedPrefixes: append(append([]string{}, c.reservedPrefixes...), keys...),
	}
}

// IsReserved checks if the key is either reserved or has a reserved prefix
func (c *ConfigStore) IsReserved(key string) bool {
	for _, res := range c.reserved {
		if key == res {
			return true
		}
	}

	for _, res := range c.reservedPrefixes {
		if key == res || strings.HasPrefix(key, res+".") {
			return true
		}
	}

	return false
}

// UnreservedMap returns a map with all unresrved entries
func (c *ConfigStore) UnreservedMap(m map[string]string) map[string]string {
	r := make(map[string]string)

	for k, v := range m {
		if !c.IsReserved(k) {
			r[k] = v
		}
	}

	return r
}

// StripedPrefixMap filters out all the keys with given prefix and returns them
// whereas the keys are striped from the prefix. The . is added implicitly.
func (c *ConfigStore) StripedPrefixMap(m map[string]string, prefix string) map[string]string {
	r := make(map[string]string)

	for k, v := range m {
		if strings.HasPrefix(k, prefix+".") {
			r[strings.TrimPrefix(k, prefix+".")] = v
		}
	}

	return r
}
