package lxf

import (
	"errors"
	"testing"

	"github.com/automaticserver/lxe/lxf/lxdfakes"
	"github.com/automaticserver/lxe/lxf/lxo"
	"github.com/lxc/lxd/lxc/config"
	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func testClient() (*Client, *lxdfakes.FakeContainerServer) {
	fake := &lxdfakes.FakeContainerServer{}

	return &Client{
		server: fake,
		config: &config.Config{},
		opwait: lxo.NewClient(fake),
	}, fake
}

func TestClient_GetRuntimeInfo_Ok(t *testing.T) {
	client, fake := testClient()
	fake.GetServerReturns(&api.Server{
		ServerUntrusted: api.ServerUntrusted{
			APIVersion: "a.b",
		},
	}, "", nil)

	info, err := client.GetRuntimeInfo()
	assert.NoError(t, err)
	assert.Equal(t, 1, fake.GetServerCallCount())
	assert.Exactly(t, "a.b.0", info.Version)
}

func TestClient_GetRuntimeInfo_Error(t *testing.T) {
	client, fake := testClient()
	fake.GetServerReturns(nil, "", errors.New("some connection error"))

	_, err := client.GetRuntimeInfo()
	assert.Error(t, err)
	assert.Equal(t, 1, fake.GetServerCallCount())
}

// func TestConnection(t *testing.T) {
// 	_, err := lxf.NewClient("", os.Getenv("HOME")+"/.config/lxc/config.yml")
// 	if err != nil {
// 		t.Errorf("failed to set up connection %v", err)
// 	}
// }

// func TestConnectionWithInvalidSocket(t *testing.T) {
// 	_, err := lxf.NewClient("/var/lib/lxd/unix.invalidsocket",
// 		os.Getenv("HOME")+"/.config/lxc/config.yml")
// 	if err == nil {
// 		t.Errorf("invalid socket should return an error")
// 	}
// }

// func NewTestClient(t *testing.T) *Client {
// 	client, err := NewClient("", os.Getenv("HOME")+"/.config/lxc/config.yml")
// 	assert.NoError(t, err)
// 	return client
// }

// // clear tries to remove all resources to provide a clean test base
// func (lt *lxfTest) clear() {
// 	for _, ct := range lt.listContainers() {
// 		lt.stopContainer(ct.ID)
// 		lt.deleteContainer(ct.ID)
// 	}
// 	for _, sb := range lt.listSandboxes() {
// 		lt.stopSandbox(sb.ID)
// 		lt.deleteSandbox(sb.ID)
// 	}

// 	// remove all image aliases except the busybox one
// 	for _, im := range lt.listImages("") {
// 		for _, al := range im.Aliases {
// 			keep := false
// 			for _, k := range keepImages {
// 				if k == al {
// 					keep = true
// 					continue
// 				}
// 			}
// 			if !keep {
// 				fmt.Println("remove image", al)
// 				lt.removeImage(al)
// 			}
// 		}
// 	}
// }
