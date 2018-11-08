package lxf

import (
	"strconv"
	"time"

	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxe/lxf/lxo"
)

// Schema Version this package is currently expecting
const (
	cfgSchema              = "user.lxe.schema"
	SchemaVersionProfile   = "0.3"
	SchemaVersionContainer = "0.3"

	cfgOldIsSandbox     = "user.is_cri_sandbox"
	cfgOldIsContainer   = "user.is_cri_container"
	cfgOldContainerName = "user.containerName"
)

// MigrationWorkspace manages schema of lxd objects
type MigrationWorkspace struct {
	lxf *LXF
}

// Migration initializes the migration workspace
func (l *LXF) Migration() *MigrationWorkspace {
	return &MigrationWorkspace{
		lxf: l,
	}
}

// IsSchemaCurrent checks if a object is in the current schema
func IsSchemaCurrent(i interface{}) bool {
	var val string
	var has bool

	switch o := i.(type) {
	case api.Container:
		if val, has = o.Config[cfgSchema]; !has {
			return false
		}
		return val == SchemaVersionContainer
	case api.Profile:
		if val, has = o.Config[cfgSchema]; !has {
			return false
		}
		return val == SchemaVersionProfile
	default:
		return false
	}
}

// Ensure applies all migration steps from detected schema to current schema
func (m *MigrationWorkspace) Ensure() error {
	profiles, err := m.lxf.server.GetProfiles()
	if err != nil {
		return err
	}

	for k := range profiles {
		// Since we want to work and modify the item directly, reference the entry
		p := &profiles[k]

		// Ignore everything which is not created by lxe
		if p.Config[cfgIsCRI] == "" && p.Config[cfgOldIsSandbox] == "" {
			continue
		}

		// TODO: or better compare to a copy of the entry?
		counter := 0
		if m.ensureProfileZeroOne(p) {
			counter++
		}
		if m.ensureProfileZeroTwo(p) {
			counter++
		}
		if m.ensureProfileZeroThree(p) {
			counter++
		}

		// If something has changed, update it
		if counter > 0 {
			err = m.lxf.server.UpdateProfile(p.Name, p.Writable(), "")
			if err != nil {
				return err
			}
		}
	}

	containers, err := m.lxf.server.GetContainers()
	if err != nil {
		return err
	}

	for k := range containers {
		// Since we want to work and modify the item directly, reference the entry
		c := &containers[k]

		// Ignore everything which is not created by lxe
		if c.Config[cfgIsCRI] == "" && c.Config[cfgOldIsContainer] == "" {
			continue
		}

		// TODO: or better compare to a copy of the entry?
		counter := 0
		if m.ensureContainerZeroOne(c) {
			counter++
		}
		if m.ensureContainerZeroTwo(c) {
			counter++
		}
		if m.ensureContainerZeroThree(c) {
			counter++
		}

		// If something has changed, update it
		if counter > 0 {
			err := lxo.UpdateContainer(m.lxf.server, c.Name, c.Writable())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// All the following functions return true, if they have changed something, otherwise false

func (m *MigrationWorkspace) ensureProfileZeroOne(p *api.Profile) bool {
	if p.Config[cfgSchema] == "" {
		p.Config[cfgMetaUID] = p.Name
		p.Config[cfgSchema] = "0.1"
		return true
	}
	return false
}

// user.is_cri_sandbox has moved to user.cri
func (m *MigrationWorkspace) ensureProfileZeroTwo(p *api.Profile) bool {
	if p.Config[cfgSchema] == "0.1" {
		p.Config[cfgIsCRI] = p.Config[cfgOldIsSandbox]
		p.Config[cfgSchema] = "0.2"
		return true
	}
	return false
}

// cleanup unused keys
func (m *MigrationWorkspace) ensureProfileZeroThree(p *api.Profile) bool {
	if p.Config[cfgSchema] == "0.2" {
		delete(p.Config, cfgOldIsSandbox)
		p.Config[cfgSchema] = "0.3"
		return true
	}
	return false
}

func (m *MigrationWorkspace) ensureContainerZeroOne(c *api.Container) bool {
	if c.Config[cfgSchema] == "" {
		c.Config[cfgSchema] = "0.1"
		return true
	}
	return false
}

// user.is_cri_container has moved to user.cri
// user.containerName has moved to user.metadata.Name
func (m *MigrationWorkspace) ensureContainerZeroTwo(c *api.Container) bool {
	if c.Config[cfgSchema] == "0.1" {
		c.Config[cfgIsCRI] = c.Config[cfgOldIsContainer]
		c.Config[cfgMetaName] = c.Config[cfgOldContainerName]
		c.Config[cfgSchema] = "0.2"
		return true
	}
	return false
}

// createdDate can be missing
// autostart can be missing
// cleanup unused keys
func (m *MigrationWorkspace) ensureContainerZeroThree(c *api.Container) bool {
	if c.Config[cfgSchema] == "0.2" {
		delete(c.Config, cfgOldIsContainer)
		delete(c.Config, cfgOldContainerName)
		c.Config[cfgAutoStartOnBoot] = strconv.FormatBool(true)
		if c.Config[cfgCreatedAt] == "" {
			if c.Config[cfgStartedAt] == "" {
				c.Config[cfgCreatedAt] = strconv.FormatInt(time.Now().UnixNano(), 10)
			} else {
				c.Config[cfgCreatedAt] = c.Config[cfgStartedAt]
			}
		}
		c.Config[cfgAutoStartOnBoot] = strconv.FormatBool(true)
		c.Config[cfgAutoStartOnBoot] = strconv.FormatBool(true)
		c.Config[cfgAutoStartOnBoot] = strconv.FormatBool(true)
		c.Config[cfgSchema] = "0.3"
		return true
	}
	return false
}
