package lxf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// As of LXD 5.1+ and 5.0.1+ volatile config entries are not allowed to be changed and thus have to be sent as-is when updating, introduced in https://github.com/lxc/lxd/commit/955c005042ab1baf77a4deb3c3d839da843b7529
// See also https://discuss.linuxcontainers.org/t/issue-creating-cloud-init-profiles-using-golang-client/13722/9
func Test_makeContainerConfig_KeepVolatile(t *testing.T) {
	t.Parallel()

	c := &Container{}
	c.Image = "foo"
	c.Config = map[string]string{
		"hello":                  "world",
		"volatile.idmap.current": `[{"Isuid":true,...}]`,
	}

	exp := map[string]string{
		"hello":                  "world",
		"volatile.base_image":    "foo",
		"volatile.idmap.current": `[{"Isuid":true,...}]`,
	}

	config := makeContainerConfig(c)

	act := map[string]string{}
	for k := range exp {
		act[k] = config[k]
	}

	assert.Exactly(t, exp, act)
}
