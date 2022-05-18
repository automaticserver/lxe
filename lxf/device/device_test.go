package device

import (
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestDetect_UnknownType(t *testing.T) {
	t.Parallel()

	_, err := Detect("foo", map[string]string{"type": "foo"})
	assert.Error(t, err)
	assert.Equal(t, true, errors.Is(err, errors.NotSupported))
}

func TestDetect_KnownType(t *testing.T) {
	t.Parallel()

	n, err := Detect("foo", map[string]string{"type": "none"})
	assert.NoError(t, err)

	exp := &None{KeyName: "foo"}
	assert.Exactly(t, exp, n)
}

func TestDetect_SameTypeMultiple(t *testing.T) {
	t.Parallel()

	m, err := Detect("foo", map[string]string{"type": "none"})
	assert.NoError(t, err)

	n, err := Detect("bar", map[string]string{"type": "none"})
	assert.NoError(t, err)

	assert.NotEqual(t, m, n)
}

func TestDevices_Upsert_AddMultiple(t *testing.T) {
	t.Parallel()

	d := Devices{}
	d.Upsert(&None{KeyName: "foo"})
	d.Upsert(&None{KeyName: "bar"})

	assert.Len(t, d, 2)
}

func TestDevices_Upsert_Override(t *testing.T) {
	t.Parallel()

	d := Devices{}
	disk := &Disk{KeyName: "foo"}

	d.Upsert(&None{KeyName: "foo"})
	d.Upsert(disk)

	assert.Len(t, d, 1)
	assert.Exactly(t, disk, d[0])
}
