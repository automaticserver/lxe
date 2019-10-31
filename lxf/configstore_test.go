package lxf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigStoreKeyReserved(t *testing.T) {
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

func TestNewConfigStorePrefixReserved(t *testing.T) {
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

func TestConfigStoreUnreserved(t *testing.T) {
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
