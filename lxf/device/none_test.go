// nolint: dupl
package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNone_getName_KeyName(t *testing.T) {
	t.Parallel()

	d := &None{KeyName: "foo"}
	assert.Equal(t, "foo", d.getName())
}

func TestNone_ToMap(t *testing.T) {
	t.Parallel()

	d := &None{KeyName: "foo"}
	exp := map[string]string{"type": NoneType}
	n, m := d.ToMap()
	assert.Equal(t, "foo", n)
	assert.Equal(t, exp, m)
}

func TestNone_FromMap(t *testing.T) {
	t.Parallel()

	raw := map[string]string{"type": NoneType}
	exp := &None{KeyName: "foo"}
	d := &None{}
	err := d.FromMap("foo", raw)
	assert.NoError(t, err)
	assert.Exactly(t, exp, d)
}
