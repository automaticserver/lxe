package lxf

import (
	"fmt"
	"testing"

	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func getSchemaContainer(schema string) api.Container {
	c := api.Container{}
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

	c := IsSchemaCurrent(getSchemaContainer(""))
	assert.Equal(t, false, c)

	p := IsSchemaCurrent(getSchemaProfile(""))
	assert.Equal(t, false, p)

	e := IsSchemaCurrent(fmt.Errorf("some wrong object"))
	assert.Equal(t, false, e)
}

func TestIsSchemaWrong(t *testing.T) {
	t.Parallel()

	c := IsSchemaCurrent(getSchemaContainer("0.0"))
	assert.Equal(t, false, c)

	p := IsSchemaCurrent(getSchemaProfile("0.0"))
	assert.Equal(t, false, p)
}

func TestIsSchemaCurrent(t *testing.T) {
	t.Parallel()

	c := IsSchemaCurrent(getSchemaContainer(SchemaVersionContainer))
	assert.Equal(t, true, c)

	p := IsSchemaCurrent(getSchemaProfile(SchemaVersionProfile))
	assert.Equal(t, true, p)
}

func TestIsSchemaPointer(t *testing.T) {
	t.Parallel()

	c1 := getSchemaContainer(SchemaVersionContainer)
	c := IsSchemaCurrent(&c1)
	assert.Equal(t, true, c)

	p1 := getSchemaProfile(SchemaVersionProfile)
	p := IsSchemaCurrent(&p1)
	assert.Equal(t, true, p)
}
