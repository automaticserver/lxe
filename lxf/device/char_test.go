// nolint: dupl
package device

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestChar_getName_KeyName(t *testing.T) {
	t.Parallel()

	d := &Char{KeyName: "foo"}
	assert.Equal(t, "foo", d.getName())
}

func TestChar_getName_PathOnly(t *testing.T) {
	t.Parallel()

	d := &Char{Path: "/tmp/foo"}
	assert.Equal(t, CharType+"-/tmp/foo", d.getName())
}

func TestChar_getName_SourceOnly(t *testing.T) {
	t.Parallel()

	d := &Char{Source: "/tmp/bar"}
	assert.Equal(t, CharType+"-/tmp/bar", d.getName())
}

func TestChar_getName_PathAndSource(t *testing.T) {
	t.Parallel()

	d := &Char{Path: "/tmp/foo", Source: "/tmp/bar"}
	assert.Equal(t, CharType+"-/tmp/foo", d.getName())
}

func TestChar_getName_KeyNamePriority(t *testing.T) {
	t.Parallel()

	d := &Char{KeyName: "foo", Path: "/tmp/foo", Source: "/tmp/bar"}
	assert.Equal(t, "foo", d.getName())
}

func TestChar_ToMap(t *testing.T) {
	t.Parallel()

	d := &Char{KeyName: "foo", Path: "bar", Source: "baz"}
	exp := map[string]string{"type": CharType, "path": "bar", "source": "baz"}
	n, m := d.ToMap()
	assert.Equal(t, "foo", n)
	assert.Equal(t, exp, m)
}

func TestChar_FromMap(t *testing.T) {
	t.Parallel()

	raw := map[string]string{"type": CharType, "path": "bar", "source": "baz"}
	exp := &Char{KeyName: "foo", Path: "bar", Source: "baz"}
	d, err := schema[CharType].FromMap("foo", raw)
	assert.NoError(t, err)
	assert.Exactly(t, exp, d)
}
