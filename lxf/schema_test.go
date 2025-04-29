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
	assert.False(t, c)

	p := IsSchemaCurrent(getSchemaProfile(""))
	assert.False(t, p)

	e := IsSchemaCurrent(fmt.Errorf("some wrong object"))
	assert.False(t, e)
}

func TestIsSchemaWrong(t *testing.T) {
	t.Parallel()

	c := IsSchemaCurrent(getSchemaInstance("0.0"))
	assert.False(t, c)

	p := IsSchemaCurrent(getSchemaProfile("0.0"))
	assert.False(t, p)
}

func TestIsSchemaCurrent(t *testing.T) {
	t.Parallel()

	c := IsSchemaCurrent(getSchemaInstance(SchemaVersionInstance))
	assert.True(t, c)

	p := IsSchemaCurrent(getSchemaProfile(SchemaVersionProfile))
	assert.True(t, p)
}

func TestIsSchemaPointer(t *testing.T) {
	t.Parallel()

	c1 := getSchemaInstance(SchemaVersionInstance)
	c := IsSchemaCurrent(&c1)
	assert.True(t, c)

	p1 := getSchemaProfile(SchemaVersionProfile)
	p := IsSchemaCurrent(&p1)
	assert.True(t, p)
}

func satisfyInstanceSchema(ct *api.Instance) *api.Instance {
	ct.Config[cfgSchema] = SchemaVersionInstance

	return ct
}

func satisfyProfileSchema(p *api.Profile) *api.Profile {
	p.Config[cfgSchema] = SchemaVersionProfile

	return p
}
