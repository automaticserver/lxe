// nolint: nestif
package lxf // import "github.com/automaticserver/lxe/lxf"

import (
	"strconv"
	"time"

	"github.com/lxc/lxd/shared/api"
)

// Schema Version this package is currently expecting
const (
	cfgSchema              = "user.lxe.schema"
	SchemaVersionProfile   = zeroThree
	SchemaVersionContainer = zeroFive

	cfgOldIsSandbox     = "user.is_cri_sandbox"
	cfgOldIsContainer   = "user.is_cri_container"
	cfgOldContainerName = "user.containerName"

	// make linter goconst happy, didn't want to disable it
	zeroOne   = "0.1"
	zeroTwo   = "0.2"
	zeroThree = "0.3"
	zeroFour  = "0.4"
	zeroFive  = "0.5"
)

// MigrationWorkspace manages schema of lxd objects
type MigrationWorkspace struct {
	lxf *client
}

// Migration initializes the migration workspace
func NewMigrationWorkspace(l Client) *MigrationWorkspace {
	return &MigrationWorkspace{
		lxf: l.(*client),
	}
}

// IsSchemaCurrent checks if a object is in the current schema
func IsSchemaCurrent(i interface{}) bool {
	var (
		val string
		has bool
	)

	switch o := i.(type) {
	case api.Container:
		if val, has = o.Config[cfgSchema]; !has {
			return false
		}

		return val == SchemaVersionContainer
	case *api.Container:
		return IsSchemaCurrent(*o)
	case api.Profile:
		if val, has = o.Config[cfgSchema]; !has {
			return false
		}

		return val == SchemaVersionProfile
	case *api.Profile:
		return IsSchemaCurrent(*o)
	default:
		return false
	}
}

// Ensure applies all migration steps from detected schema to current schema
func (m *MigrationWorkspace) Ensure() error { // nolint: gocognit
	profiles, err := m.lxf.server.GetProfiles()
	if err != nil {
		return err
	}

	anyChanges := false

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
			anyChanges = true

			err = m.lxf.server.UpdateProfile(p.Name, p.Writable(), "")
			if err != nil {
				return err
			}
		}
	}

	var etag string

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

		if m.ensureContainerZeroFour(c) {
			counter++
		}

		if m.ensureContainerZeroFive(c) {
			counter++
		}

		// If something has changed, update it
		if counter > 0 {
			anyChanges = true

			err := m.lxf.opwait.UpdateContainer(c.Name, c.Writable(), etag)
			if err != nil {
				return err
			}
		}
	}

	if anyChanges {
		log.Warnf("Migration changes applied successfully")
	}

	return nil
}

// All the following functions return true, if they have changed something, otherwise false

func (m *MigrationWorkspace) ensureProfileZeroOne(p *api.Profile) bool {
	if p.Config[cfgSchema] == "" {
		p.Config[cfgMetaUID] = p.Name
		p.Config[cfgSchema] = zeroOne

		return true
	}

	return false
}

// user.is_cri_sandbox has moved to user.cri
func (m *MigrationWorkspace) ensureProfileZeroTwo(p *api.Profile) bool {
	if p.Config[cfgSchema] == zeroOne {
		p.Config[cfgIsCRI] = p.Config[cfgOldIsSandbox]
		p.Config[cfgSchema] = zeroTwo

		return true
	}

	return false
}

// cleanup unused keys
func (m *MigrationWorkspace) ensureProfileZeroThree(p *api.Profile) bool {
	if p.Config[cfgSchema] == zeroTwo {
		delete(p.Config, cfgOldIsSandbox)
		p.Config[cfgSchema] = zeroThree

		return true
	}

	return false
}

func (m *MigrationWorkspace) ensureContainerZeroOne(c *api.Container) bool {
	if c.Config[cfgSchema] == "" {
		c.Config[cfgSchema] = zeroOne
		return true
	}

	return false
}

// user.is_cri_container has moved to user.cri
// user.containerName has moved to user.metadata.Name
func (m *MigrationWorkspace) ensureContainerZeroTwo(c *api.Container) bool {
	if c.Config[cfgSchema] == zeroOne {
		c.Config[cfgIsCRI] = c.Config[cfgOldIsContainer]
		c.Config[cfgMetaName] = c.Config[cfgOldContainerName]
		c.Config[cfgSchema] = zeroTwo

		return true
	}

	return false
}

// createdDate can be missing
// autostart can be missing
// cleanup unused keys
func (m *MigrationWorkspace) ensureContainerZeroThree(c *api.Container) bool {
	if c.Config[cfgSchema] == zeroTwo {
		delete(c.Config, cfgOldIsContainer)
		delete(c.Config, cfgOldContainerName)

		if c.Config[cfgCreatedAt] == "" {
			if c.Config[cfgStartedAt] == "" {
				c.Config[cfgCreatedAt] = strconv.FormatInt(time.Now().UnixNano(), 10)
			} else {
				c.Config[cfgCreatedAt] = c.Config[cfgStartedAt]
			}
		}

		if c.Config[cfgStartedAt] == "" {
			c.Config[cfgStartedAt] = strconv.FormatInt(time.Time{}.UnixNano(), 10)
		}

		if c.Config[cfgFinishedAt] == "" {
			c.Config[cfgFinishedAt] = strconv.FormatInt(time.Time{}.UnixNano(), 10)
		}

		c.Config[cfgSchema] = zeroThree

		return true
	}

	return false
}

// boot.autostart is not managed by lxe anymore, keep field as-is
// WARNING: intentionally changed migration to 0.3 to not force-setting that field if
// someone is coming from 0.2 or below
func (m *MigrationWorkspace) ensureContainerZeroFour(c *api.Container) bool {
	if c.Config[cfgSchema] == zeroThree {
		c.Config[cfgSchema] = zeroFour
		return true
	}

	return false
}

// Implemented variable length of profiles. The order of profiles in schema <= 0.4 was wrong.
// Move the first profile, which was the sandbox, to the last position, otherwise preserve position
func (m *MigrationWorkspace) ensureContainerZeroFive(c *api.Container) bool {
	if c.Config[cfgSchema] == zeroFour {
		c.Profiles = append(c.Profiles[1:], c.Profiles[0])
		c.Config[cfgSchema] = zeroFive

		return true
	}

	return false
}
