package lxf

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
