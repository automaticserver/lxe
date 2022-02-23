package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxy_getName_KeyName(t *testing.T) {
	t.Parallel()

	d := &Proxy{KeyName: "foo"}
	assert.Equal(t, "foo", d.getName())
}

func TestProxy_getName_ListenOnly(t *testing.T) {
	t.Parallel()

	d := &Proxy{Listen: &ProxyEndpoint{Protocol: ProtocolTCP, Address: "baz", Port: 22}}
	assert.Equal(t, ProxyType+"-tcp:baz:22", d.getName())
}

func TestProxy_getName_KeyNamePriority(t *testing.T) {
	t.Parallel()

	d := &Proxy{KeyName: "foo", Listen: &ProxyEndpoint{Protocol: ProtocolTCP, Address: "baz", Port: 22}}
	assert.Equal(t, "foo", d.getName())
}

func TestProxy_ToMap(t *testing.T) {
	t.Parallel()

	d := &Proxy{KeyName: "foo", Listen: &ProxyEndpoint{Protocol: ProtocolTCP, Address: "baz", Port: 22}, Destination: &ProxyEndpoint{Protocol: ProtocolUDP, Address: "cba", Port: 33}}
	exp := map[string]string{"type": ProxyType, "listen": "tcp:baz:22", "connect": "udp:cba:33"}
	n, m := d.ToMap()
	assert.Equal(t, "foo", n)
	assert.Equal(t, exp, m)
}

func TestProxy_FromMap(t *testing.T) {
	t.Parallel()

	raw := map[string]string{"type": ProxyType, "listen": "tcp:baz:22", "connect": "udp:cba:33"}
	exp := &Proxy{KeyName: "foo", Listen: &ProxyEndpoint{Protocol: ProtocolTCP, Address: "baz", Port: 22}, Destination: &ProxyEndpoint{Protocol: ProtocolUDP, Address: "cba", Port: 33}}
	d := &Proxy{}
	err := d.FromMap("foo", raw)
	assert.NoError(t, err)
	assert.Exactly(t, exp, d)
}

func Test_newProtocol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    Protocol
		wantErr bool
	}{
		{"tcp", ProtocolTCP, false},
		{"udp", ProtocolUDP, false},
		{"", ProtocolUndefined, true},
		{"undefined", ProtocolUndefined, true},
		{"foo", ProtocolUndefined, true},
	}
	for _, tt := range tests {
		tt := tt // pin!

		t.Run("", func(t *testing.T) {
			t.Parallel()

			got, err := newProtocol(tt.input)
			assert.False(t, (err != nil) != tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewProxyEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    *ProxyEndpoint
		wantErr bool
	}{
		{"", nil, true},
		{"foo:bar:baz", nil, true},
		{"foo:bar", nil, true},
		{"foo:25", nil, true},
		{"foo:bar:25", nil, true},
		{"tcp:bar:25", &ProxyEndpoint{Protocol: ProtocolTCP, Address: "bar", Port: 25}, false},
		{"udp:baz:35", &ProxyEndpoint{Protocol: ProtocolUDP, Address: "baz", Port: 35}, false},
		{":baz:35", nil, true},
		{"udp:baz:foo", nil, true},
	}
	for _, tt := range tests {
		tt := tt // pin!

		t.Run("", func(t *testing.T) {
			t.Parallel()

			got, err := NewProxyEndpoint(tt.input)
			assert.False(t, (err != nil) != tt.wantErr)
			assert.Exactly(t, tt.want, got)
		})
	}
}
