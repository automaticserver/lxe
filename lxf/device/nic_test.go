// nolint: dupl
package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNic_getName_KeyName(t *testing.T) {
	t.Parallel()

	d := &Nic{KeyName: "foo"}
	assert.Equal(t, "foo", d.getName())
}

func TestNic_getName_NameOnly(t *testing.T) {
	t.Parallel()

	d := &Nic{Name: "ethX"}
	assert.Equal(t, NicType+"-ethX", d.getName())
}

func TestNic_getName_KeyNamePriority(t *testing.T) {
	t.Parallel()

	d := &Nic{KeyName: "foo", Name: "ethX"}
	assert.Equal(t, "foo", d.getName())
}

func TestNic_ToMap(t *testing.T) {
	t.Parallel()

	d := &Nic{KeyName: "foo", Name: "ethX", NicType: "bridge", Parent: "brX", IPv4Address: "1.2.3.4"}
	exp := map[string]string{"type": NicType, "name": "ethX", "nictype": "bridge", "parent": "brX", "ipv4.address": "1.2.3.4"}
	n, m := d.ToMap()
	assert.Equal(t, "foo", n)
	assert.Equal(t, exp, m)
}

func TestNic_FromMap(t *testing.T) {
	t.Parallel()

	raw := map[string]string{"type": NicType, "name": "ethX", "nictype": "bridge", "parent": "brX", "ipv4.address": "1.2.3.4"}
	exp := &Nic{KeyName: "foo", Name: "ethX", NicType: "bridge", Parent: "brX", IPv4Address: "1.2.3.4"}
	d := &Nic{}
	err := d.FromMap("foo", raw)
	assert.NoError(t, err)
	assert.Exactly(t, exp, d)
}
