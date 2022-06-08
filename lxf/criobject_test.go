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

func getCRIImage(cri string) api.Image {
	i := api.Image{}
	i.Properties = make(map[string]string)

	if cri != "" {
		i.Properties[cfgIsCRI] = cri
	}

	return i
}

func TestIsCRIEmpty(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c := l.IsCRI(getCRIContainer(""))
	assert.Equal(t, false, c)

	p := l.IsCRI(getCRIProfile(""))
	assert.Equal(t, false, p)

	i := l.IsCRI(getCRIImage(""))
	assert.Equal(t, false, i)

	e := l.IsCRI(fmt.Errorf("some wrong object"))
	assert.Equal(t, false, e)
}

func TestIsCRIFalse(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c := l.IsCRI(getCRIContainer("false"))
	assert.Equal(t, false, c)

	p := l.IsCRI(getCRIProfile("False"))
	assert.Equal(t, false, p)

	i := l.IsCRI(getCRIImage("FALSE"))
	assert.Equal(t, false, i)
}

func TestIsCRIWrong(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c := l.IsCRI(getCRIContainer("no"))
	assert.Equal(t, false, c)

	p := l.IsCRI(getCRIProfile("yes"))
	assert.Equal(t, false, p)

	i := l.IsCRI(getCRIImage("maybe"))
	assert.Equal(t, false, i)
}

func TestIsCRITrue(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c := l.IsCRI(getCRIContainer("true"))
	assert.Equal(t, true, c)

	p := l.IsCRI(getCRIProfile("True"))
	assert.Equal(t, true, p)

	i := l.IsCRI(getCRIImage("TRUE"))
	assert.Equal(t, true, i)
}

func TestIsCRIPointer(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c1 := getCRIContainer("true")
	c := l.IsCRI(&c1)
	assert.Equal(t, true, c)

	p1 := getCRIProfile("True")
	p := l.IsCRI(&p1)
	assert.Equal(t, true, p)

	i1 := getCRIImage("TRUE")
	i := l.IsCRI(&i1)
	assert.Equal(t, true, i)
}

func satisfyContainerCri(ct *api.Container) *api.Container {
	ct.Config[cfgIsCRI] = "true"

	return ct
}

func satisfyProfileCri(p *api.Profile) *api.Profile {
	p.Config[cfgIsCRI] = "true"

	return p
}
