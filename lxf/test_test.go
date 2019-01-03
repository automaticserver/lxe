package lxf_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/lxc/lxe/lxf"
)

var keepImages = []string{
	"critest.asag.io/busybox/1.28:latest",
}

// lxfTest is a facade for testing lxf which provides the same methods
// but swallows returned errors and lets the test fail if not nil
type lxfTest struct {
	t   *testing.T
	lxf *lxf.LXF
}

// newLXFTest returns a lxf test facade using the provided test context
// to fail on error
func newLXFTest(t *testing.T) *lxfTest {
	lxs, err := lxf.New("", os.Getenv("HOME")+"/.config/lxc/config.yml")
	if err != nil {
		t.Fatalf("could not create lx facade, %v", err)
	}

	lxft := &lxfTest{
		t:   t,
		lxf: lxs,
	}

	lxft.clear()
	lxft.pullImage("critest.asag.io/busybox/1.28")
	return lxft
}

// clear tries to remove all resources to provide a clean test base
func (lt *lxfTest) clear() {
	for _, ct := range lt.listContainers() {
		lt.stopContainer(ct.ID)
		lt.deleteContainer(ct.ID)
	}
	for _, sb := range lt.listSandboxes() {
		lt.stopSandbox(sb.ID)
		lt.deleteSandbox(sb.ID)
	}

	// remove all image aliases except the busybox one
	for _, im := range lt.listImages("") {
		for _, al := range im.Aliases {
			keep := false
			for _, k := range keepImages {
				if k == al {
					keep = true
					continue
				}
			}
			if !keep {
				fmt.Println("remove image", al)
				lt.removeImage(al)
			}
		}
	}
}

func (lt *lxfTest) createSandbox(s *lxf.Sandbox) {
	err := lt.lxf.CreateSandbox(s)
	if err != nil {
		lt.t.Errorf("failed to create sandbox %v, %v", s.ID, err)
	}
}

func (lt *lxfTest) getSandbox(id string) *lxf.Sandbox {
	s, err := lt.lxf.GetSandbox(id)
	if err != nil {
		lt.t.Errorf("failed to get sandbox '%v;, %v", id, err)
	}
	return s
}

func (lt *lxfTest) stopSandbox(id string) {
	err := lt.lxf.StopSandbox(id)
	if err != nil {
		lt.t.Errorf("failed to get sandbox '%v;, %v", id, err)
	}
}

func (lt *lxfTest) getContainer(id string) *lxf.Container {
	s, err := lt.lxf.GetContainer(id)
	if err != nil {
		lt.t.Errorf("failed to get container '%v;, %v", id, err)
	}
	return s
}

func (lt *lxfTest) listSandboxes() []*lxf.Sandbox {
	s, err := lt.lxf.ListSandboxes()
	if err != nil {
		lt.t.Errorf("failed to list sandboxes, %v", err)
	}
	return s
}

func (lt *lxfTest) listContainers() []*lxf.Container {
	s, err := lt.lxf.ListContainers()
	if err != nil {
		lt.t.Errorf("failed to list containers, %v", err)
	}
	return s
}

func (lt *lxfTest) deleteSandbox(id string) {
	err := lt.lxf.DeleteSandbox(id)
	if err != nil {
		lt.t.Errorf("failed to delete sandbox '%v', %v", id, err)
	}
}

func (lt *lxfTest) deleteContainer(id string) {
	err := lt.lxf.DeleteContainer(id)
	if err != nil {
		lt.t.Errorf("failed to delete container '%v', %v", id, err)
	}
}

func (lt *lxfTest) createContainer(s *lxf.Container) {
	err := lt.lxf.CreateContainer(s)
	if err != nil {
		lt.t.Errorf("failed to create container, %v", err)
	}
}

func (lt *lxfTest) updateContainer(s *lxf.Container) {
	err := lt.lxf.UpdateContainer(s)
	if err != nil {
		lt.t.Errorf("failed to update container, %v", err)
	}
}

func (lt *lxfTest) startContainer(id string) {
	err := lt.lxf.StartContainer(id)
	if err != nil {
		lt.t.Errorf("failed to start container '%v', %v", id, err)
	}
}

func (lt *lxfTest) stopContainer(id string) {
	err := lt.lxf.StopContainer(id, 30)
	if err != nil {
		lt.t.Errorf("failed to stop container '%v', %v", id, err)
	}
}

func (lt *lxfTest) execSync(id string, cmd []string) *lxf.ExecResponse {
	r, err := lt.lxf.ExecSync(id, cmd)
	if err != nil {
		lt.t.Fatalf("failed to exec sync on container '%v', %v", id, err)
	}
	return r
}

func (lt *lxfTest) exec(cid string, cmd []string,
	stdin io.Reader, stdout, stderr io.WriteCloser) int {

	code, err := lt.lxf.Exec(cid, cmd, stdin, stdout, stderr)
	if err != nil {
		lt.t.Fatalf("failed to exec on container '%v', %v", cid, err)
	}
	return code
}

func (lt *lxfTest) listImages(filter string) []lxf.Image {
	imgs, err := lt.lxf.ListImages(filter)
	if err != nil {
		lt.t.Errorf("failed to list images, %v", err)
	}
	return imgs
}

func (lt *lxfTest) pullImage(name string) string {
	hash, err := lt.lxf.PullImage(name)
	if err != nil {
		lt.t.Errorf("failed to list images, %v", err)
	}
	return hash
}

func (lt *lxfTest) removeImage(name string) {
	err := lt.lxf.RemoveImage(name)
	if err != nil {
		lt.t.Errorf("failed to remove images, %v", err)
	}
}

func (lt *lxfTest) getImage(name string) *lxf.Image {
	i, err := lt.lxf.GetImage(name)
	if err != nil {
		lt.t.Errorf("failed to get image, %v", err)
	}
	return i
}
