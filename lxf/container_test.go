package lxf_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/lxc/lxe/lxf"
)

const (
	imgBusybox = "critest.asag.io/busybox/1.28"
)

func TestCreateContainer(t *testing.T) {
	lt := newLXFTest(t)

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})
}

func TestContainerAttributes(t *testing.T) {
	lt := newLXFTest(t)

	ld := "/var/log/veryloggylog"

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		LogPath: ld,
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	ct := lt.getContainer("roosevelt")

	if ct.LogPath != ld {
		t.Errorf("expected log directory to be %v but is %v", ld, ct.LogPath)
	}

	if !strings.HasPrefix(ct.Image, "9445061f9fad") {
		t.Errorf("expected image to have prefix %v but is %v", imgBusybox, ct.Image)
	}
}

func TestContainerLabels(t *testing.T) { // nolint:dupl
	lt := newLXFTest(t)

	lb := map[string]string{
		"foo":       "bar",
		"delicious": "coffee",
	}
	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
		Labels:  lb,
	})

	ct := lt.getContainer("roosevelt")

	if !reflect.DeepEqual(ct.Labels, lb) {
		t.Errorf("expected labels to be %v but is %v", lb, ct.Labels)
	}
}

func TestContainerAnnotations(t *testing.T) { // nolint:dupl
	lt := newLXFTest(t)

	an := map[string]string{
		"foo":       "bar-annotation",
		"delicious": "coffee",
	}
	lt.createContainer(&lxf.Container{
		Name:        "roosevelt",
		Sandbox:     setUpSandbox(lt, "roosevelt"),
		Image:       imgBusybox,
		Annotations: an,
	})

	ct := lt.getContainer("roosevelt")

	if !reflect.DeepEqual(ct.Annotations, an) {
		t.Errorf("expected labels to be %v but is %v", an, ct.Annotations)
	}
}

func TestContainerCloudInit(t *testing.T) { // nolint:dupl
	lt := newLXFTest(t)

	cloudinit := `
#cloud-config
write_files:
  - path: /tmp/cloud-init-test-file
    owner: root:root
    permissions: '0644'
    content: |
      blueberry
`

	lt.createContainer(&lxf.Container{
		Name:              "roosevelt",
		Sandbox:           setUpSandbox(lt, "roosevelt"),
		Image:             imgBusybox,
		CloudInitUserData: cloudinit,
	})

	ct := lt.getContainer("roosevelt")

	if !reflect.DeepEqual(ct.CloudInitUserData, cloudinit) {
		t.Errorf("expected cloudinit to be %v but is %v", cloudinit, ct.CloudInitUserData)
	}
}

func TestUpdateContainer(t *testing.T) {
	lt := newLXFTest(t)

	an1 := map[string]string{
		"foo":       "bar-annotation",
		"delicious": "coffee",
	}
	an2 := map[string]string{
		"delicious": "iced-coffee",
	}
	lt.createContainer(&lxf.Container{
		Name:        "roosevelt",
		Sandbox:     setUpSandbox(lt, "roosevelt"),
		Image:       imgBusybox,
		Annotations: an1,
	})

	ct := lt.getContainer("roosevelt")
	if !reflect.DeepEqual(ct.Annotations, an1) {
		t.Errorf("expected labels to be %v but is %v", an1, ct.Annotations)
	}

	ct.Annotations = an2
	lt.updateContainer(ct)

	ct = lt.getContainer("roosevelt")
	if !reflect.DeepEqual(ct.Annotations, an2) {
		t.Errorf("expected labels to be %v but is %v", an2, ct.Annotations)
	}
}

func TestStartContainer(t *testing.T) {
	lt := newLXFTest(t)

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")
}

func TestStopContainer(t *testing.T) {
	lt := newLXFTest(t)

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")
	lt.stopContainer("roosevelt")
}

func TestContainerWithReallyLongName(t *testing.T) {
	lt := newLXFTest(t)
	name := strings.Repeat("roosevelt", 10) + "_reallylong"

	lt.createContainer(&lxf.Container{
		Name:    name,
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer(name)
	ct := lt.getContainer(name)

	if ct.Name != name {
		t.Errorf("expected container name to be %v but is %v", name, ct.Name)
	}
}

func TestContainerStateTransitions(t *testing.T) {
	lt := newLXFTest(t)
	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	ct := lt.getContainer("roosevelt")

	if ct.State != lxf.ContainerStateCreated {
		t.Errorf("state of created container must be 'created' but is '%v'", ct.State)
	}

	lt.startContainer("roosevelt")
	ct = lt.getContainer("roosevelt")
	if ct.State != lxf.ContainerStateRunning {
		t.Errorf("state of running container must be 'running' but is '%v'", ct.State)
	}

	lt.stopContainer("roosevelt")
	ct = lt.getContainer("roosevelt")
	if ct.State != lxf.ContainerStateExited {
		t.Errorf("state of running container must be 'exited' but is '%v'", ct.State)
	}

}

func TestContainerPid(t *testing.T) {
	lt := newLXFTest(t)

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")
	ct := lt.getContainer("roosevelt")

	if ct.Pid == 0 {
		t.Errorf("expected pid to not be 0 but is %v", ct.Pid)
	}

}

func TestContainerAttempts(t *testing.T) {
	lt := newLXFTest(t)

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	ct := lt.getContainer("roosevelt")

	if ct.Metadata.Attempt != 0 {
		t.Errorf("attempts of unstarted container should be 0 but is %v", ct.Metadata.Attempt)
	}

	lt.startContainer("roosevelt")
	ct = lt.getContainer("roosevelt")

	if ct.Metadata.Attempt != 1 {
		t.Errorf("attempts of started container should be 1 but is %v", ct.Metadata.Attempt)
	}
}

func TestContainerStartedAt(t *testing.T) {
	lt := newLXFTest(t)

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	ct := lt.getContainer("roosevelt")

	if !ct.StartedAt.IsZero() {
		t.Errorf("started at of unstarted container should be 0 but is %v", ct.StartedAt)
	}

	lt.startContainer("roosevelt")
	ct = lt.getContainer("roosevelt")

	diff := time.Since(ct.StartedAt)
	if diff > time.Second {
		t.Errorf("started date should be less than 1 second old but is %v", diff)
	}
}

func TestContainerCreatedAt(t *testing.T) {
	lt := newLXFTest(t)

	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	ct := lt.getContainer("roosevelt")
	diff := time.Since(ct.CreatedAt)
	if diff > time.Second*3 {
		t.Errorf("created at should be less than 3 seconds old but is %v", diff)
	}
}

func setUpSandbox(lt *lxfTest, name string) *lxf.Sandbox {
	r := &lxf.Sandbox{
		Hostname: name + "-sb-hn",
		Metadata: lxf.SandboxMetadata{
			UID: name + "-sb",
		},
	}
	lt.createSandbox(r)

	return lt.getSandbox(r.ID)
}
