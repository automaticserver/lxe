// nolint: dupl
package device

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

const (
	ProxyType = "proxy"
)

// Proxy device representation https://lxd.readthedocs.io/en/latest/containers/#type-proxy
type Proxy struct {
	KeyName     string
	Listen      *ProxyEndpoint
	Destination *ProxyEndpoint
}

func (d *Proxy) getName() string {
	var name string

	switch {
	case d.KeyName != "":
		name = d.KeyName
	default:
		name = fmt.Sprintf("%v-%v", ProxyType, d.Listen.String())
	}

	return name
}

// ToMap returns assigned name or if unset the type specific unique name and serializes the options into a lxd device map
func (d *Proxy) ToMap() (string, map[string]string) {
	return d.getName(), map[string]string{
		"type":    ProxyType,
		"listen":  d.Listen.String(),
		"connect": d.Destination.String(),
	}
}

// New creates a new empty device
func (d *Proxy) new() Device {
	return &Proxy{}
}

// FromMap loads assigned name (can be empty) and options
func (d *Proxy) FromMap(name string, options map[string]string) error {
	var err error

	d.KeyName = name

	d.Listen, err = NewProxyEndpoint(options["listen"])
	if err != nil {
		return err
	}

	d.Destination, err = NewProxyEndpoint(options["connect"])
	if err != nil {
		return err
	}

	return nil
}

// Protocol defines the type of a proxy endpoint
type Protocol int

const (
	// ProtocolUndefined is not a valid protocol
	ProtocolUndefined = Protocol(0)
	// ProtocolTCP makes the endpoint use TCP
	ProtocolTCP = Protocol(1)
	// ProtocolUDP makes the endpoint use UDP
	ProtocolUDP = Protocol(2)
)

var (
	protMapNameVal = map[string]Protocol{
		"undefined": ProtocolUndefined,
		"tcp":       ProtocolTCP,
		"udp":       ProtocolUDP,
	}
	protMapValName = map[Protocol]string{
		ProtocolUndefined: "undefined",
		ProtocolTCP:       "tcp",
		ProtocolUDP:       "udp",
	}
)

func newProtocol(str string) (Protocol, error) {
	if i, has := protMapNameVal[str]; has && str != protMapValName[ProtocolUndefined] {
		return i, nil
	}

	return ProtocolUndefined, errors.NotValidf("unknown protocol: %v", str)
}

func (p Protocol) String() string {
	return protMapValName[p]
}

// ProxyEndpoint defines an enpoint with protocol, address and port
type ProxyEndpoint struct {
	Protocol Protocol
	Address  string
	Port     int
}

// NewProxyEndpoint parses a string of the form protocol:address:port
// protocol: tcp|udp
// address: ip or empty
// port: uiint16
// TODO verify and document allowed format values
func NewProxyEndpoint(str string) (*ProxyEndpoint, error) {
	parts := strings.Split(str, ":")
	if len(parts) != 3 {
		return nil, errors.NotValidf("proxy endpoint must be delimited by two colons (::), we were given: `%v`", str)
	}

	prot, err := newProtocol(parts[0])
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errors.NotValidf("port must be an int not %v", parts[2])
	}

	return &ProxyEndpoint{
		Protocol: prot,
		Address:  parts[1],
		Port:     port,
	}, nil
}

func (p *ProxyEndpoint) String() string {
	return p.Protocol.String() + ":" + p.Address + ":" + strconv.Itoa(p.Port)
}
