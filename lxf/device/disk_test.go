// nolint: dupl
package device

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestDisk_getName_KeyName(t *testing.T) {
	t.Parallel()

	d := &Disk{KeyName: "foo"}
	assert.Equal(t, "foo", d.getName())
}

func TestDisk_getName_PathOnly(t *testing.T) {
	t.Parallel()

	d := &Disk{Path: "/tmp/foo"}
	assert.Equal(t, DiskType+"-/tmp/foo", d.getName())
}

func TestDisk_getName_SourceOnly(t *testing.T) {
	t.Parallel()

	d := &Disk{Source: "/tmp/bar"}
	assert.Equal(t, DiskType+"-/tmp/bar", d.getName())
}

func TestDisk_getName_PathAndSource(t *testing.T) {
	t.Parallel()

	d := &Disk{Path: "/tmp/foo", Source: "/tmp/bar"}
	assert.Equal(t, DiskType+"-/tmp/foo", d.getName())
}

func TestDisk_getName_KeyNamePriority(t *testing.T) {
	t.Parallel()

	d := &Disk{KeyName: "foo", Path: "/tmp/foo", Source: "/tmp/bar"}
	assert.Equal(t, "foo", d.getName())
}

func TestDisk_ToMap(t *testing.T) {
	t.Parallel()

	d := &Disk{KeyName: "foo", Path: "bar", Source: "baz", Pool: "pool", Size: "size", Readonly: true, Optional: true}
	exp := map[string]string{"type": DiskType, "path": "bar", "source": "baz", "pool": "pool", "size": "size", "readonly": "true", "optional": "true"}
	n, m := d.ToMap()
	assert.Equal(t, "foo", n)
	assert.Equal(t, exp, m)
}

func TestDisk_FromMap(t *testing.T) {
	t.Parallel()

	raw := map[string]string{"type": DiskType, "path": "bar", "source": "baz", "pool": "pool", "size": "size", "readonly": "true", "optional": "true"}
	exp := &Disk{KeyName: "foo", Path: "bar", Source: "baz", Pool: "pool", Size: "size", Readonly: true, Optional: true}
	d := &Disk{}
	err := d.FromMap("foo", raw)
	assert.NoError(t, err)
	assert.Exactly(t, exp, d)
}
