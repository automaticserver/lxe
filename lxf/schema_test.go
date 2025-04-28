package lxf

import (
	"fmt"
	"testing"

	"github.com/canonical/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func getSchemaInstance(schema string) api.Instance {
	c := api.Instance{}
	c.Config = make(map[string]string)

	if schema != "" {
		c.Config[cfgSchema] = schema
	}

	return c
}

func getSchemaProfile(schema string) api.Profile {
	p := api.Profile{}
	p.Config = make(map[string]string)

	if schema != "" {
		p.Config[cfgSchema] = schema
	}

	return p
}

func TestIsSchemaEmpty(t *testing.T) {
	t.Parallel()

	c := IsSchemaCurrent(getSchemaInstance(""))
	assert.Equal(t, false, c)

	p := IsSchemaCurrent(getSchemaProfile(""))
	assert.Equal(t, false, p)

	e := IsSchemaCurrent(fmt.Errorf("some wrong object"))
	assert.Equal(t, false, e)
}

func TestIsSchemaWrong(t *testing.T) {
	t.Parallel()

	c := IsSchemaCurrent(getSchemaInstance("0.0"))
	assert.Equal(t, false, c)

	p := IsSchemaCurrent(getSchemaProfile("0.0"))
	assert.Equal(t, false, p)
}

func TestIsSchemaCurrent(t *testing.T) {
	t.Parallel()

	c := IsSchemaCurrent(getSchemaInstance(SchemaVersionInstance))
	assert.Equal(t, true, c)

	p := IsSchemaCurrent(getSchemaProfile(SchemaVersionProfile))
	assert.Equal(t, true, p)
}

func TestIsSchemaPointer(t *testing.T) {
	t.Parallel()

	c1 := getSchemaInstance(SchemaVersionInstance)
	c := IsSchemaCurrent(&c1)
	assert.Equal(t, true, c)

	p1 := getSchemaProfile(SchemaVersionProfile)
	p := IsSchemaCurrent(&p1)
	assert.Equal(t, true, p)
}

func satisfyInstanceSchema(ct *api.Instance) *api.Instance {
	ct.Config[cfgSchema] = SchemaVersionInstance

	return ct
}

func satisfyProfileSchema(p *api.Profile) *api.Profile {
	p.Config[cfgSchema] = SchemaVersionProfile

	return p
}
