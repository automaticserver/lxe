package lxf

import (
	"fmt"
	"testing"

	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func getCRIContainer(cri string) api.Container {
	c := getSchemaContainer(SchemaVersionContainer)
	if cri != "" {
		c.Config[cfgIsCRI] = cri
	}

	return c
}

func getCRIProfile(cri string) api.Profile {
	p := getSchemaProfile(SchemaVersionProfile)
	if cri != "" {
		p.Config[cfgIsCRI] = cri
	}

	return p
}

func TestIsCRIEmpty(t *testing.T) {
	t.Parallel()

	c := IsCRI(getCRIContainer(""))
	assert.Equal(t, false, c)

	p := IsCRI(getCRIProfile(""))
	assert.Equal(t, false, p)

	e := IsCRI(fmt.Errorf("some wrong object"))
	assert.Equal(t, false, e)
}

func TestIsCRIFalse(t *testing.T) {
	t.Parallel()

	c := IsCRI(getCRIContainer("false"))
	assert.Equal(t, false, c)

	p := IsCRI(getCRIProfile("false"))
	assert.Equal(t, false, p)
}

func TestIsCRIWrong(t *testing.T) {
	t.Parallel()

	c := IsCRI(getCRIContainer("no"))
	assert.Equal(t, false, c)

	p := IsCRI(getCRIProfile("yes"))
	assert.Equal(t, false, p)
}

func TestIsCRITrue(t *testing.T) {
	t.Parallel()

	c := IsCRI(getCRIContainer("true"))
	assert.Equal(t, true, c)

	p := IsCRI(getCRIProfile("True"))
	assert.Equal(t, true, p)
}

func TestIsCRIPointer(t *testing.T) {
	t.Parallel()

	c1 := getCRIContainer("true")
	c := IsCRI(&c1)
	assert.Equal(t, true, c)

	p1 := getCRIProfile("True")
	p := IsCRI(&p1)
	assert.Equal(t, true, p)
}
