// nolint: dupl
package device

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestBlock_getName_KeyName(t *testing.T) {
	t.Parallel()

	d := &Block{KeyName: "foo"}
	assert.Equal(t, "foo", d.getName())
}

func TestBlock_getName_PathOnly(t *testing.T) {
	t.Parallel()

	d := &Block{Path: "/tmp/foo"}
	assert.Equal(t, BlockType+"-/tmp/foo", d.getName())
}

func TestBlock_getName_SourceOnly(t *testing.T) {
	t.Parallel()

	d := &Block{Source: "/tmp/bar"}
	assert.Equal(t, BlockType+"-/tmp/bar", d.getName())
}

func TestBlock_getName_PathAndSource(t *testing.T) {
	t.Parallel()

	d := &Block{Path: "/tmp/foo", Source: "/tmp/bar"}
	assert.Equal(t, BlockType+"-/tmp/foo", d.getName())
}

func TestBlock_getName_KeyNamePriority(t *testing.T) {
	t.Parallel()

	d := &Block{KeyName: "foo", Path: "/tmp/foo", Source: "/tmp/bar"}
	assert.Equal(t, "foo", d.getName())
}

func TestBlock_ToMap(t *testing.T) {
	t.Parallel()

	d := &Block{KeyName: "foo", Path: "bar", Source: "baz"}
	exp := map[string]string{"type": BlockType, "path": "bar", "source": "baz"}
	n, m := d.ToMap()
	assert.Equal(t, "foo", n)
	assert.Equal(t, exp, m)
}

func TestBlock_FromMap(t *testing.T) {
	t.Parallel()

	raw := map[string]string{"type": BlockType, "path": "bar", "source": "baz"}
	exp := &Block{KeyName: "foo", Path: "bar", Source: "baz"}
	d, err := schema[BlockType].FromMap("foo", raw)
	assert.NoError(t, err)
	assert.Exactly(t, exp, d)
}
