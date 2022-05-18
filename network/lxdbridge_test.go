package network

import (
	"testing"

	"github.com/automaticserver/lxe/lxf/lxdfakes"
	"github.com/automaticserver/lxe/shared"
	lxd "github.com/lxc/lxd/client"
	lxdApi "github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	testLXDBridge = "testbr0"
)

var (
	// verify interface satisfaction
	_ Plugin           = &lxdBridgePlugin{}
	_ PodNetwork       = &lxdBridgePodNetwork{}
	_ ContainerNetwork = &lxdBridgeContainerNetwork{}
)

func testLXDClient() (lxd.ContainerServer, *lxdfakes.FakeContainerServer) {
	fake := &lxdfakes.FakeContainerServer{}

	return fake, fake
}

func TestInitPluginLXDBridge_DefaultsAndCreate(t *testing.T) {
	t.Parallel()

	server, fake := testLXDClient()

	fake.GetNetworkReturns(nil, "", shared.NewErrNotFound())

	p, err := InitPluginLXDBridge(server, ConfLXDBridge{})
	assert.NoError(t, err)
	assert.Exactly(t, fake, p.server)
	assert.NotEmpty(t, p.conf.LXDBridge, "lxebr0 is the the default")

	assert.Equal(t, 1, fake.GetNetworkCallCount())
	assert.Equal(t, 1, fake.CreateNetworkCallCount())
	args := fake.CreateNetworkArgsForCall(0)
	assert.Equal(t, "auto", args.Config["ipv4.address"])
	assert.Equal(t, "true", args.Config["ipv4.dhcp"])
	assert.Equal(t, "false", args.Config["ipv4.nat"])
	assert.Equal(t, "port=0", args.Config["raw.dnsmasq"])
}

func TestInitPluginLXDBridge_DefinedAndUpdate(t *testing.T) {
	t.Parallel()

	server, fake := testLXDClient()
	cidr := "192.168.224.0/24"
	cidrExp := "192.168.224.1/24"

	fake.GetNetworkReturns(&lxdApi.Network{
		Type: "bridge",
		NetworkPut: lxdApi.NetworkPut{
			Config: make(map[string]string),
		},
	}, "", nil)

	p, err := InitPluginLXDBridge(server, ConfLXDBridge{LXDBridge: testLXDBridge, Cidr: cidr, Nat: true, CreateOnly: false})
	assert.NoError(t, err)
	assert.Exactly(t, fake, p.server)
	assert.NotEmpty(t, p.conf.Cidr)
	assert.NotEmpty(t, p.conf.Nat)

	assert.Equal(t, 1, fake.GetNetworkCallCount())
	assert.Equal(t, 0, fake.CreateNetworkCallCount())
	assert.Equal(t, 1, fake.UpdateNetworkCallCount())
	name, args, _ := fake.UpdateNetworkArgsForCall(0)
	assert.Equal(t, testLXDBridge, name)
	assert.Equal(t, cidrExp, args.Config["ipv4.address"])
	assert.Equal(t, "true", args.Config["ipv4.dhcp"])
	assert.Equal(t, "true", args.Config["ipv4.nat"])
	assert.Equal(t, "port=0", args.Config["raw.dnsmasq"])
}

func testLXDBridgePlugin() (*lxdBridgePlugin, *lxdfakes.FakeContainerServer) {
	client, fake := testLXDClient()

	return &lxdBridgePlugin{
		server: client,
		conf:   ConfLXDBridge{LXDBridge: testLXDBridge},
	}, fake
}

func Test_lxdBridgePlugin_PodNetwork(t *testing.T) {
	t.Parallel()

	plugin, _ := testLXDBridgePlugin()

	podNet, err := plugin.PodNetwork("foo", nil)
	assert.NoError(t, err)

	tPodNet := podNet.(*lxdBridgePodNetwork)
	assert.Equal(t, "foo", tPodNet.podID)
}

func Test_lxdBridgePlugin_UpdateRuntimeConfig(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()

	fake.GetNetworkReturns(nil, "", shared.NewErrNotFound())

	err := plugin.UpdateRuntimeConfig(&rtApi.RuntimeConfig{NetworkConfig: &rtApi.NetworkConfig{PodCidr: "192.168.224.0/24"}})
	assert.NoError(t, err)

	assert.Equal(t, 1, fake.CreateNetworkCallCount())
	args := fake.CreateNetworkArgsForCall(0)
	assert.Equal(t, testLXDBridge, args.Name)
	assert.Equal(t, "192.168.224.1/24", args.Config["ipv4.address"])
}

func Test_lxdBridgePlugin_ensureBridge_WrongNetworkTypeExists(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()

	fake.GetNetworkReturns(&lxdApi.Network{Type: "other"}, "", nil)

	err := plugin.ensureBridge()
	assert.Error(t, err)
	assert.Empty(t, fake.CreateNetworkCallCount())
	assert.Empty(t, fake.UpdateNetworkCallCount())
}

func Test_lxdBridgePlugin_ensureBridge_CreateOnly(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()
	plugin.conf.CreateOnly = true

	fake.GetNetworkReturns(&lxdApi.Network{Type: "bridge", Name: testLXDBridge}, "", nil)

	err := plugin.ensureBridge()
	assert.NoError(t, err)
	assert.Empty(t, fake.CreateNetworkCallCount())
	assert.Empty(t, fake.UpdateNetworkCallCount())
}

func Test_lxdBridgePlugin_ensureBridge_CorrectIPRangeBridgeIP(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()
	plugin.conf.Cidr = "192.168.224.0/24"
	cidrExp := "192.168.224.1/24"

	fake.GetNetworkReturns(nil, "", shared.NewErrNotFound())

	err := plugin.ensureBridge()
	assert.NoError(t, err)
	assert.Equal(t, 1, fake.CreateNetworkCallCount())

	args := fake.CreateNetworkArgsForCall(0)
	assert.Equal(t, testLXDBridge, args.Name)
	assert.Equal(t, cidrExp, args.Config["ipv4.address"])
}

func Test_lxdBridgePlugin_ensureBridge_CorrectIPRangeAuto(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()
	plugin.conf.Cidr = ""

	fake.GetNetworkReturns(nil, "", shared.NewErrNotFound())

	err := plugin.ensureBridge()
	assert.NoError(t, err)
	assert.Equal(t, 1, fake.CreateNetworkCallCount())

	args := fake.CreateNetworkArgsForCall(0)
	assert.Equal(t, testLXDBridge, args.Name)
	assert.Equal(t, "auto", args.Config["ipv4.address"])
}

func Test_lxdBridgePlugin_findFreeIP_Simple(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()

	fake.GetNetworkReturns(&lxdApi.Network{
		Type: "bridge",
		Name: testLXDBridge,
		NetworkPut: lxdApi.NetworkPut{
			Config: map[string]string{
				"ipv4.address":     "192.168.224.1/30",
				"ipv4.dhcp.ranges": "",
			},
		},
	}, "", nil)
	fake.GetNetworkLeasesReturns([]lxdApi.NetworkLease{}, nil)

	ip, err := plugin.findFreeIP()
	assert.NoError(t, err)
	assert.Equal(t, "192.168.224.2", ip.String())
}

func Test_lxdBridgePlugin_findFreeIP_WithLeases(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()

	fake.GetNetworkReturns(&lxdApi.Network{
		Type: "bridge",
		Name: testLXDBridge,
		NetworkPut: lxdApi.NetworkPut{
			Config: map[string]string{
				"ipv4.address":     "192.168.224.1/29",
				"ipv4.dhcp.ranges": "",
			},
		},
	}, "", nil)
	fake.GetNetworkLeasesReturns([]lxdApi.NetworkLease{
		{Address: "192.168.224.2"},
		{Address: "192.168.224.3"},
		{Address: "192.168.224.4"},
		{Address: "192.168.224.5"},
	}, nil)

	ip, err := plugin.findFreeIP()
	assert.NoError(t, err)
	assert.Equal(t, "192.168.224.6", ip.String())
}

func Test_lxdBridgePlugin_findFreeIP_NoRangeSupportYet(t *testing.T) {
	t.Parallel()

	plugin, fake := testLXDBridgePlugin()

	fake.GetNetworkReturns(&lxdApi.Network{
		Type: "bridge",
		Name: testLXDBridge,
		NetworkPut: lxdApi.NetworkPut{
			Config: map[string]string{
				"ipv4.address":     "192.168.224.1/29",
				"ipv4.dhcp.ranges": "192.168.224.2-192.168.224.3,192.168.224.4-192.168.224.6",
			},
		},
	}, "", nil)

	_, err := plugin.findFreeIP()
	assert.Error(t, err)
}

func testLXDBridgePodNetwork() (*lxdBridgePodNetwork, *lxdfakes.FakeContainerServer) {
	plugin, fake := testLXDBridgePlugin()

	return &lxdBridgePodNetwork{
		plugin: plugin,
		podID:  "hello",
	}, fake
}

func Test_lxdBridgePodNetwork_ContainerNetwork(t *testing.T) {
	t.Parallel()

	podNet, _ := testLXDBridgePodNetwork()

	contNet, err := podNet.ContainerNetwork("foo", nil)
	assert.NoError(t, err)

	tContNet := contNet.(*lxdBridgeContainerNetwork)
	assert.Equal(t, "foo", tContNet.cid)
}

func Test_lxdBridgePodNetwork_Status_NoData(t *testing.T) {
	t.Parallel()

	podNet, _ := testLXDBridgePodNetwork()

	status, err := podNet.Status(ctx, &PropertiesRunning{})
	assert.Error(t, err)
	assert.Nil(t, status)
}

func Test_lxdBridgePodNetwork_Status_InvalidData(t *testing.T) {
	t.Parallel()

	podNet, _ := testLXDBridgePodNetwork()

	status, err := podNet.Status(ctx, &PropertiesRunning{Properties: Properties{Data: map[string]string{"interface-address": "bar"}}})
	assert.Error(t, err)
	assert.Nil(t, status)
}

func Test_lxdBridgePodNetwork_Status_Simple(t *testing.T) {
	t.Parallel()

	podNet, _ := testLXDBridgePodNetwork()

	status, err := podNet.Status(ctx, &PropertiesRunning{Properties: Properties{Data: map[string]string{"interface-address": "192.168.224.2"}}})
	assert.NoError(t, err)
	assert.Equal(t, "192.168.224.2", status.IPs[0].String())
}

func Test_lxdBridgePodNetwork_WhenCreated_Simple(t *testing.T) {
	t.Parallel()

	podNet, fake := testLXDBridgePodNetwork()

	fake.GetNetworkReturns(&lxdApi.Network{
		Type: "bridge",
		Name: testLXDBridge,
		NetworkPut: lxdApi.NetworkPut{
			Config: map[string]string{
				"ipv4.address":     "192.168.224.1/30",
				"ipv4.dhcp.ranges": "",
			},
		},
	}, "", nil)
	fake.GetNetworkLeasesReturns([]lxdApi.NetworkLease{}, nil)

	res, err := podNet.WhenCreated(ctx, &Properties{})
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Data["interface-address"])
	assert.NotEmpty(t, res.Nics[0].IPv4Address)
}
