// nolint: nestif
package lxf

import (
	"strconv"
	"time"

	"github.com/canonical/lxd/shared/api"
)

// Schema Version this package is currently expecting
const (
	cfgSchema             = "user.lxe.schema"
	SchemaVersionProfile  = zeroThree
	SchemaVersionInstance = zeroFive

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
		lxf: l.(*client), // nolint: forcetypeassert
	}
}

// IsSchemaCurrent checks if a object is in the current schema
func IsSchemaCurrent(i interface{}) bool {
	var (
		val string
		has bool
	)

	switch o := i.(type) {
	case api.Instance:
		if val, has = o.Config[cfgSchema]; !has {
			return false
		}

		return val == SchemaVersionInstance
	case *api.Instance:
		return IsSchemaCurrent(*o)
	case api.Profile:
		if val, has = o.Config[cfgSchema]; !has {
			return false
		}

		return val == SchemaVersionProfile
	case *api.Profile:
		return IsSchemaCurrent(*o)
	// Images are always schema conform for now
	case api.Image:
		return true
	case *api.Image:
		return IsSchemaCurrent(*o)

	default:
		return false
	}
}

// Ensure applies all migration steps from detected schema to current schema
func (m *MigrationWorkspace) Ensure() error { // nolint: gocognit, cyclop
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

	instances, err := m.lxf.server.GetInstances(api.InstanceTypeAny)
	if err != nil {
		return err
	}

	for k := range instances {
		// Since we want to work and modify the item directly, reference the entry
		i := &instances[k]

		// Ignore everything which is not created by lxe
		if i.Config[cfgIsCRI] == "" && i.Config[cfgOldIsContainer] == "" {
			continue
		}

		// TODO: or better compare to a copy of the entry?
		counter := 0

		if m.ensureInstanceZeroOne(i) {
			counter++
		}

		if m.ensureInstanceZeroTwo(i) {
			counter++
		}

		if m.ensureInstanceZeroThree(i) {
			counter++
		}

		if m.ensureInstanceZeroFour(i) {
			counter++
		}

		if m.ensureInstanceZeroFive(i) {
			counter++
		}

		// If something has changed, update it
		if counter > 0 {
			anyChanges = true

			err := m.lxf.opwait.UpdateInstance(i.Name, i.Writable(), etag)
			if err != nil {
				return err
			}
		}
	}

	if anyChanges {
		log.Warn("Migration changes applied successfully")
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

func (m *MigrationWorkspace) ensureInstanceZeroOne(i *api.Instance) bool {
	if i.Config[cfgSchema] == "" {
		i.Config[cfgSchema] = zeroOne

		return true
	}

	return false
}

// user.is_cri_container has moved to user.cri
// user.containerName has moved to user.metadata.Name
func (m *MigrationWorkspace) ensureInstanceZeroTwo(i *api.Instance) bool {
	if i.Config[cfgSchema] == zeroOne {
		i.Config[cfgIsCRI] = i.Config[cfgOldIsContainer]
		i.Config[cfgMetaName] = i.Config[cfgOldContainerName]
		i.Config[cfgSchema] = zeroTwo

		return true
	}

	return false
}

// createdDate can be missing
// autostart can be missing
// cleanup unused keys
func (m *MigrationWorkspace) ensureInstanceZeroThree(i *api.Instance) bool {
	if i.Config[cfgSchema] == zeroTwo {
		delete(i.Config, cfgOldIsContainer)
		delete(i.Config, cfgOldContainerName)

		if i.Config[cfgCreatedAt] == "" {
			if i.Config[cfgStartedAt] == "" {
				i.Config[cfgCreatedAt] = strconv.FormatInt(time.Now().UnixNano(), 10)
			} else {
				i.Config[cfgCreatedAt] = i.Config[cfgStartedAt]
			}
		}

		if i.Config[cfgStartedAt] == "" {
			i.Config[cfgStartedAt] = strconv.FormatInt(time.Time{}.UnixNano(), 10)
		}

		if i.Config[cfgFinishedAt] == "" {
			i.Config[cfgFinishedAt] = strconv.FormatInt(time.Time{}.UnixNano(), 10)
		}

		i.Config[cfgSchema] = zeroThree

		return true
	}

	return false
}

// boot.autostart is not managed by lxe anymore, keep field as-is
// WARNING: intentionally changed migration to 0.3 to not force-setting that field if
// someone is coming from 0.2 or below
func (m *MigrationWorkspace) ensureInstanceZeroFour(i *api.Instance) bool {
	if i.Config[cfgSchema] == zeroThree {
		i.Config[cfgSchema] = zeroFour

		return true
	}

	return false
}

// Implemented variable length of profiles. The order of profiles in schema <= 0.4 was wrong.
// Move the first profile, which was the sandbox, to the last position, otherwise preserve position
func (m *MigrationWorkspace) ensureInstanceZeroFive(i *api.Instance) bool {
	if i.Config[cfgSchema] == zeroFour {
		i.Profiles = append(i.Profiles[1:], i.Profiles[0])
		i.Config[cfgSchema] = zeroFive

		return true
	}

	return false
}
