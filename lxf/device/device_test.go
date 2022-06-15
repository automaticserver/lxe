package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetect_UnknownType(t *testing.T) {
	t.Parallel()

	_, err := Detect("foo", map[string]string{"type": "foo"})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotSupported)
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

func Test_trimKeyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		arg  string
		want string
	}{
		{"disk-/var/lib/short", "disk-/var/lib/short"},
		{"disk-/var/lib/something/long", "disk-/var/li--mething/long"}, // remove when maxKeyNameLength=64
		// {"disk-/var/lib/something/so/very/long/it/exceeds/sixtyfour/characters", "disk-/var/lib/something/so/very--it/exceeds/sixtyfour/characters"}, // add when maxKeyNameLength=64
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.arg, func(t *testing.T) {
			t.Parallel()

			if got := trimKeyName(tt.arg); got != tt.want {
				t.Errorf("trimKeyName() = %v, want %v", got, tt.want)
			}
		})
	}
}
