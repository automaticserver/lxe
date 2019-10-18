package network

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/automaticserver/lxe/network/libcnifake"
	"github.com/stretchr/testify/assert"
)

var (
	testCNIDir       string
	testCNIbinPath   string = defaultCNIbinPath
	testCNIconfPath  string
	testCNINetNSPath string = defaultCNINetNSPath

	ctx = context.TODO()
)

func init() {
	var err error
	testCNIDir, err = ioutil.TempDir("", "cni")
	if err != nil {
		panic(err)
	}

	testCNIconfPath = filepath.Join(testCNIDir, defaultCNIconfPath)
	//testCNINetNSPath = filepath.Join(testCNIDir, defaultCNINetNSPath)

	err = os.MkdirAll(testCNIconfPath, 0700)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(testCNIconfPath, "99-lo.conf"), []byte(`
	{
		"cniVersion": "0.4.0",
		"name": "lo",
		"type": "loopback"
	}`), 0600)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(testCNINetNSPath, 0700)
	if err != nil {
		panic(err)
	}
}

func testNeedsRoot(t *testing.T) {
	// in order create, read, delete netns the user must have permission for it, for now check if user is root
	user, err := user.Current()
	assert.NoError(t, err)
	if user.Username != "root" {
		t.Skip("not user root")
	}
}

func TestInitPluginCNI(t *testing.T) {
	type args struct {
		conf ConfCNI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"defaults", args{ConfCNI{}}, false},
		{"testing", args{ConfCNI{testCNIbinPath, testCNIconfPath, testCNINetNSPath}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitPluginCNI(tt.args.conf)
			assert.False(t, (err != nil) != tt.wantErr)
			assert.NotNil(t, got.cni)
			assert.NotEmpty(t, got.conf)
		})
	}
}

func testCNIPlugin(t *testing.T) *cniPlugin {
	t.Log("testCNIDir", testCNIDir)
	cniPlugin, err := InitPluginCNI(ConfCNI{testCNIbinPath, testCNIconfPath, testCNINetNSPath})
	assert.NoError(t, err)
	return cniPlugin
}

func toFakeCniPlugin(t *testing.T, cniPlugin *cniPlugin) *libcnifake.FakeCNI {
	fake := &libcnifake.FakeCNI{}
	cniPlugin.cni = fake
	return fake
}

func Test_cniPlugin_PodNetwork(t *testing.T) {
	cniPlugin := testCNIPlugin(t)
	type args struct {
		namespace   string
		name        string
		id          string
		annotations map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{"default", "nginx", "foo", nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cniPlugin.PodNetwork(tt.args.namespace, tt.args.name, tt.args.id, tt.args.annotations)
			fmt.Print(err)
			assert.False(t, (err != nil) != tt.wantErr)
			assert.NotNil(t, got.netList)
			assert.NotNil(t, got.runtimeConf)
		})
	}
}

func Test_cniPodNetwork_Setup(t *testing.T) {
	testNeedsRoot(t)
	cniPlugin := testCNIPlugin(t)
	podNetwork, err := cniPlugin.PodNetwork("default", "nginx", "foo", nil)
	assert.NoError(t, err)
	fake := toFakeCniPlugin(t, cniPlugin)

	err = podNetwork.Setup(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, fake.AddNetworkListCallCount())
}

func Test_cniPodNetwork_Teardown_MissingNetwork(t *testing.T) {
	testNeedsRoot(t)
	cniPlugin := testCNIPlugin(t)
	podNetwork, err := cniPlugin.PodNetwork("default", "nginx", "foo", nil)
	assert.NoError(t, err)
	//fake := toFakeCniPlugin(t, cniPlugin)

	// CRI DelNetwork always tries to remove as good as possible without throwing error
	err = podNetwork.Teardown(ctx)
	assert.NoError(t, err)
	//assert.Equal(t, 1, fake.DelNetworkListCallCount())
}

func Test_cniPodNetwork_Status_MissingNetwork(t *testing.T) {
	testNeedsRoot(t)
	cniPlugin := testCNIPlugin(t)
	podNetwork, err := cniPlugin.PodNetwork("default", "nginx", "foo", nil)
	assert.NoError(t, err)
	//fake := toFakeCniPlugin(t, cniPlugin)

	got, err := podNetwork.Status(ctx)
	assert.Error(t, err)
	assert.Nil(t, got)
	//assert.Equal(t, 1, fake.CheckNetworkListCallCount())
}

func Test_cniPodNetwork_Status_WithNetwork(t *testing.T) {
	testNeedsRoot(t)
	cniPlugin := testCNIPlugin(t)
	podNetwork, err := cniPlugin.PodNetwork("default", "nginx", "foo", nil)
	assert.NoError(t, err)
	//fake := toFakeCniPlugin(t, cniPlugin)

	err = podNetwork.Setup(ctx)
	assert.NoError(t, err)
	//assert.Equal(t, 1, fake.AddNetworkListCallCount())

	got, err := podNetwork.Status(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	//assert.Equal(t, 1, fake.CheckNetworkListCallCount())
}

func Test_parseIPAddrShow(t *testing.T) {
	type args struct {
		output []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []net.IP
		wantErr bool
	}{
		{"", args{[]byte("eth0   UP     10.10.100.169/24 ")}, []net.IP{net.ParseIP("10.10.100.169")}, false},
		{"", args{[]byte("eth0   DOWN   10.10.100.169/24")}, []net.IP{net.ParseIP("10.10.100.169")}, false},
		{"", args{[]byte("10.10.100.169/24")}, nil, true},
		{"", args{[]byte("2: eth0    inet 10.10.100.169/24 brd 10.10.100.255 scope global dynamic noprefixroute eth0")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIPAddrShow(tt.args.output)
			assert.False(t, (err != nil) != tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
