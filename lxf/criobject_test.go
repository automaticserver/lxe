package lxf

import (
	"fmt"
	"testing"

	"github.com/canonical/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func getCRIInstance(cri string) api.Instance {
	c := getSchemaInstance(SchemaVersionInstance)
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

	c := l.IsCRI(getCRIInstance(""))
	assert.False(t, c)

	p := l.IsCRI(getCRIProfile(""))
	assert.False(t, p)

	i := l.IsCRI(getCRIImage(""))
	assert.False(t, i)

	e := l.IsCRI(fmt.Errorf("some wrong object"))
	assert.False(t, e)
}

func TestIsCRIFalse(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c := l.IsCRI(getCRIInstance("false"))
	assert.False(t, c)

	p := l.IsCRI(getCRIProfile("False"))
	assert.False(t, p)

	i := l.IsCRI(getCRIImage("FALSE"))
	assert.False(t, i)
}

func TestIsCRIWrong(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c := l.IsCRI(getCRIInstance("no"))
	assert.False(t, c)

	p := l.IsCRI(getCRIProfile("yes"))
	assert.False(t, p)

	i := l.IsCRI(getCRIImage("maybe"))
	assert.False(t, i)
}

func TestIsCRITrue(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c := l.IsCRI(getCRIInstance("true"))
	assert.True(t, c)

	p := l.IsCRI(getCRIProfile("True"))
	assert.True(t, p)

	i := l.IsCRI(getCRIImage("TRUE"))
	assert.True(t, i)
}

func TestIsCRIPointer(t *testing.T) {
	t.Parallel()

	l, _ := testClient()

	c1 := getCRIInstance("true")
	c := l.IsCRI(&c1)
	assert.True(t, c)

	p1 := getCRIProfile("True")
	p := l.IsCRI(&p1)
	assert.True(t, p)

	i1 := getCRIImage("TRUE")
	i := l.IsCRI(&i1)
	assert.True(t, i)
}

func satisfyContainerCri(ct *api.Instance) *api.Instance {
	ct.Config[cfgIsCRI] = "true"

	return ct
}

func satisfyProfileCri(p *api.Profile) *api.Profile {
	p.Config[cfgIsCRI] = "true"

	return p
}
