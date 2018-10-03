package lxf

import (
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxe/lxf/lxo"
)

// Schema Version this package is currently expecting
const (
	cfgSchema              = "user.lxe.schema"
	SchemaVersionProfile   = "0.1"
	SchemaVersionContainer = "0.1"
)

type migrationWorkspace struct {
	lxf *LXF
}

// Migration initializes the migration workspace
func migration(l *LXF) *migrationWorkspace {
	return &migrationWorkspace{
		lxf: l,
	}
}

func (m *migrationWorkspace) ensure() error {
	profiles, err := m.lxf.server.GetProfiles()
	if err != nil {
		return err
	}

	for k := range profiles {
		// Since we want to work and modify the item directly, reference the entry
		p := &profiles[k]

		// Ignore everything which is not created by lxe
		if p.Config[cfgIsSandbox] == "" {
			continue
		}

		// TODO: or better compare to a copy of the entry?
		counter := 0
		if m.ensureProfileZeroOne(p) {
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
		if c.Config[cfgIsContainer] == "" {
			continue
		}

		// TODO: or better compare to a copy of the entry?
		counter := 0
		if m.ensureContainerZeroOne(c) {
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

func (m *migrationWorkspace) ensureProfileZeroOne(p *api.Profile) bool {
	if p.Config[cfgSchema] == "" {
		p.Config[cfgMetaUID] = p.Name
		p.Config[cfgSchema] = "0.1"
		return true
	}
	return false
}

func (m *migrationWorkspace) ensureContainerZeroOne(p *api.Container) bool {
	if p.Config[cfgSchema] == "" {
		p.Config[cfgSchema] = "0.1"
		return true
	}
	return false
}
