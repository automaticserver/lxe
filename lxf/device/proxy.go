package device

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	proxyType = "proxy"
	// ProtocolUndefined is not a valid protocol
	ProtocolUndefined = Protocol(0)
	// ProtocolTCP makes the endpoint use TCP
	ProtocolTCP = Protocol(1)
	// ProtocolUDP makes the endpoint use UDP
	ProtocolUDP = Protocol(2)
)

// Proxies holds slice of Proxy
// Use it if you want to Add() a entry non-conflicting (see Add())
type Proxies []Proxy

// Add a entry to the slice, if the name is the same, will overwrite the existing entry
func (ps *Proxies) Add(p Proxy) {
	for k, e := range *ps {
		if e.GetName() == p.GetName() {
			(*ps)[k] = p
			return
		}
	}
	*ps = append(*ps, p)
}

// Proxy defines a lxd proxy device
type Proxy struct {
	Listen      ProxyEndpoint
	Destination ProxyEndpoint
}

// Protocol defines the type of a proxy endpoint
type Protocol int

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

// ToMap will serialize itself into a lxd device map
func (p Proxy) ToMap() (map[string]string, error) {
	return map[string]string{
		"type":    proxyType,
		"listen":  p.Listen.String(),
		"connect": p.Destination.String(),
	}, nil
}

// GetName will generate a uinique name for the device map
func (p Proxy) GetName() string {
	return proxyType + "-" + p.Destination.String()
}

// ProxyFromMap constructs a Proxy from provided map values
func ProxyFromMap(dev map[string]string) (Proxy, error) {
	l, err := NewProxyEndpoint(dev["listen"])
	if err != nil {
		return Proxy{}, err
	}
	d, err := NewProxyEndpoint(dev["connect"])
	if err != nil {
		return Proxy{}, err
	}

	return Proxy{
		Listen:      l,
		Destination: d,
	}, nil
}

func newProtocol(str string) (Protocol, error) {
	if i, has := protMapNameVal[str]; has && str != "undefined" {
		return i, nil
	}
	return ProtocolUndefined, fmt.Errorf("unknown protocol, %v", str)
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
func NewProxyEndpoint(str string) (ProxyEndpoint, error) {
	parts := strings.Split(str, ":")
	if len(parts) != 3 {
		return ProxyEndpoint{}, fmt.Errorf("Proxy endpoint must be delimited by two colons (::), we were given: `%v`", str)
	}

	prot, err := newProtocol(parts[0])
	if err != nil {
		return ProxyEndpoint{}, err
	}

	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return ProxyEndpoint{}, fmt.Errorf("Port must be an int not %v", parts[2])
	}

	return ProxyEndpoint{
		Protocol: prot,
		Address:  parts[1],
		Port:     port,
	}, nil
}

func (p *ProxyEndpoint) String() string {
	return p.Protocol.String() + ":" + p.Address + ":" + strconv.Itoa(p.Port)
}
