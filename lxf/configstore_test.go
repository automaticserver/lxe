package lxf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigStore_WithReserved(t *testing.T) {
	t.Parallel()

	cs := NewConfigStore().WithReserved("foo", "foo.bar")

	reserved := []string{"foo", "foo.bar"}
	notreserved := []string{"foo.baz", "bar"}

	for _, s := range reserved {
		assert.True(t, cs.IsReserved(s), "should be reserved")
	}

	for _, s := range notreserved {
		assert.False(t, cs.IsReserved(s), "should not be reserved")
	}
}

func TestConfigStore_WithReservedPrefixes(t *testing.T) {
	t.Parallel()

	cs := NewConfigStore().WithReservedPrefixes("foo", "zoo.bar")

	reserved := []string{"foo", "foo.bar", "zoo.bar", "zoo.bar.tree"}
	notreserved := []string{"fool", "zoo"}

	for _, s := range reserved {
		assert.True(t, cs.IsReserved(s), "should be reserved")
	}

	for _, s := range notreserved {
		assert.False(t, cs.IsReserved(s), "should not be reserved")
	}
}

func TestConfigStore_IsReserved(t *testing.T) {
	t.Parallel()

	cs := NewConfigStore().WithReserved("tree", "kanu.manu").WithReservedPrefixes("foo", "zoo.bar")

	assert.True(t, cs.IsReserved("tree"))
	assert.True(t, cs.IsReserved("kanu.manu"))
	assert.True(t, cs.IsReserved("foo"))
	assert.True(t, cs.IsReserved("foo.baz"))
	assert.True(t, cs.IsReserved("zoo.bar"))
	assert.True(t, cs.IsReserved("zoo.bar.baz"))

	assert.False(t, cs.IsReserved("tree.leaves"))
	assert.False(t, cs.IsReserved("green.tree"))
	assert.False(t, cs.IsReserved("kanu.manu.velo"))
	assert.False(t, cs.IsReserved("velo.kanu.manu"))
	assert.False(t, cs.IsReserved("random"))
	assert.False(t, cs.IsReserved("random.string"))
}

func TestConfigStore_UnreservedMap(t *testing.T) {
	t.Parallel()

	cs := NewConfigStore().WithReserved("tree", "kanu.manu").WithReservedPrefixes("foo", "zoo.bar")

	unres := cs.UnreservedMap(map[string]string{
		"tree":             "no",
		"kanu.manu":        "no",
		"foo":              "no",
		"foo.lala":         "no",
		"foo.zoo.bar.yeah": "no",
		"foola":            "yes",
		"jouse":            "yes",
	})

	for _, v := range unres {
		assert.Equal(t, "yes", v, "should be in the unreserved result")
	}
}

func TestConfigStore_StripedPrefixMap(t *testing.T) {
	t.Parallel()

	cs := NewConfigStore()

	combined := map[string]string{
		"io.example.foo":     "yes",
		"io.example.bar.baz": "yes",
		"io.example":         "no",
		"other":              "no",
	}

	expected := map[string]string{
		"foo":     "yes",
		"bar.baz": "yes",
	}

	assert.Equal(t, expected, cs.StripedPrefixMap(combined, "io.example"))
}
