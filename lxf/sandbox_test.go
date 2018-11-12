package lxf_test

import (
	"reflect"
	"testing"

	"github.com/lxc/lxe/lxf"
	"github.com/lxc/lxe/lxf/device"
)

func TestCreateSandbox(t *testing.T) {
	lt := newLXFTest(t)

	lt.createSandbox(&lxf.Sandbox{
		Hostname: "affective",
	})
}

func TestGetSandbox(t *testing.T) {
	lt := newLXFTest(t)

	r := &lxf.Sandbox{
		Hostname: "affective",
	}
	lt.createSandbox(r)

	s := lt.getSandbox(r.ID)

	if s.Metadata.Name != r.ID {
		t.Errorf("sandbox has name %v instead of roosevelt", s.Metadata.Name)
	}
}

func TestGetInexistentSandbox(t *testing.T) {
	lt := newLXFTest(t)

	s, err := lt.lxf.GetSandbox("roosevelt")

	if err != nil {
		t.Errorf("missing sandbox should not produce an error but got %v", err)
	}

	if s != nil {
		t.Errorf("missing sandbox should be nil but is not")
	}
}

func TestDeleteSandbox(t *testing.T) {
	lt := newLXFTest(t)

	r := &lxf.Sandbox{
		Hostname: "affective",
	}
	lt.createSandbox(r)

	lt.deleteSandbox(r.ID)

	sbs := lt.listSandboxes()

	if len(sbs) != 0 {
		lt.t.Errorf("expected sandbox to be deleted but there are still %v", len(sbs))
	}
}

func TestListSandboxes(t *testing.T) {
	lt := newLXFTest(t)

	r1 := &lxf.Sandbox{
		Hostname: "affective1",
	}
	lt.createSandbox(r1)
	r2 := &lxf.Sandbox{
		Hostname: "affective2",
	}
	lt.createSandbox(r2)

	sbs := lt.listSandboxes()

	if len(sbs) != 2 {
		lt.t.Errorf("expected number of sandboxes to be 2 but they are %v", len(sbs))
	}

}

func TestSandboxAttributes(t *testing.T) {
	lt := newLXFTest(t)

	ld := "/var/log/veryloggylog"
	hn := "affectivepeanut"

	r := &lxf.Sandbox{
		Hostname:     hn,
		LogDirectory: ld,
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if sb.LogDirectory != ld {
		t.Errorf("expected log directory to be %v but is %v", ld, sb.LogDirectory)
	}

	if sb.Hostname != hn {
		t.Errorf("expected hostname to be %v but is %v", hn, sb.Hostname)
	}
}

func TestSandboxMetadata(t *testing.T) {
	lt := newLXFTest(t)

	r := &lxf.Sandbox{
		Hostname: "affective",
		Metadata: lxf.SandboxMetadata{
			Attempt:   5,
			Name:      "roosevelt",
			Namespace: "bar",
		},
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if sb.Metadata.Attempt != 5 {
		t.Errorf("expected attempt to be 5 but is %v", sb.Metadata.Attempt)
	}

	if sb.Metadata.Name != "foo" {
		t.Errorf("expected name to be foo but is %v", sb.Metadata.Name)
	}

	if sb.Metadata.Namespace != "bar" {
		t.Errorf("expected namespace to be bar but is %v", sb.Metadata.Namespace)
	}
}

func TestSandboxNetworkConfig(t *testing.T) {
	lt := newLXFTest(t)

	ns := []string{"a.foo.b"}
	sr := []string{"foo", "bar"}

	r := &lxf.Sandbox{
		Hostname: "affective",
		NetworkConfig: lxf.NetworkConfig{
			Nameservers: ns,
			Searches:    sr,
		},
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if !reflect.DeepEqual(sb.NetworkConfig.Nameservers, ns) {
		t.Errorf("expected nameservers to be %v but is %v", ns, sb.NetworkConfig.Nameservers)
	}

	if !reflect.DeepEqual(sb.NetworkConfig.Searches, sr) {
		t.Errorf("expected searches to be %v but is %v", sr, sb.NetworkConfig.Searches)
	}
}

func TestSandboxLabels(t *testing.T) { // nolint:dupl
	lt := newLXFTest(t)

	lb := map[string]string{
		"foo":       "bar",
		"delicious": "coffee",
	}
	r := &lxf.Sandbox{
		CRIObject: lxf.CRIObject{
			Labels: lb,
		},
		Hostname: "affective",
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if !reflect.DeepEqual(sb.Labels, lb) {
		t.Errorf("expected labels to be %v but is %v", lb, sb.Labels)
	}
}

func TestSandboxAnnotations(t *testing.T) { // nolint:dupl
	lt := newLXFTest(t)

	an := map[string]string{
		"foo":       "bar-annotation",
		"delicious": "coffee",
	}
	r := &lxf.Sandbox{
		CRIObject: lxf.CRIObject{
			Annotations: an,
		},
		Hostname: "affective",
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if !reflect.DeepEqual(sb.Annotations, an) {
		t.Errorf("expected labels to be %v but is %v", an, sb.Annotations)
	}
}

func TestSandboxProxies(t *testing.T) {
	lt := newLXFTest(t)

	proxies := []device.Proxy{
		device.Proxy{
			Destination: device.ProxyEndpoint{
				Protocol: device.ProtocolTCP,
				Address:  "",
				Port:     80,
			},
			Listen: device.ProxyEndpoint{
				Protocol: device.ProtocolTCP,
				Address:  "127.0.0.1",
				Port:     8080,
			},
		},
	}

	r := &lxf.Sandbox{
		CRIObject: lxf.CRIObject{
			LXDObject: lxf.LXDObject{
				Proxies: proxies,
			},
		},
		Hostname: "affective",
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if !reflect.DeepEqual(sb.Proxies, proxies) {
		t.Errorf("expected proxies to be %v but is %v", proxies, sb.Proxies)
	}
}

func TestSandboxConfigs(t *testing.T) {
	lt := newLXFTest(t)

	cfg := map[string]string{
		"user.linux.cgroup_parent": "foobar",
	}

	r := &lxf.Sandbox{
		CRIObject: lxf.CRIObject{
			LXDObject: lxf.LXDObject{
				Config: cfg,
			},
		},
		Hostname: "affective",
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if !reflect.DeepEqual(sb.Config, cfg) {
		t.Errorf("expected config to be %v but is %v", cfg, sb.Config)
	}

}

func TestSandboxDisk(t *testing.T) {
	lt := newLXFTest(t)

	disks := []device.Disk{
		device.Disk{
			Path:     "/",
			Readonly: true,
			Pool:     "default",
		},
	}

	r := &lxf.Sandbox{
		CRIObject: lxf.CRIObject{
			LXDObject: lxf.LXDObject{
				Disks: disks,
			},
		},
		Hostname: "affective",
	}
	lt.createSandbox(r)

	sb := lt.getSandbox(r.ID)

	if !reflect.DeepEqual(sb.Disks, disks) {
		t.Errorf("expected disks to be %v but is %v", disks, sb.Disks)
	}
}

func TestUsedBy(t *testing.T) {
	lt := newLXFTest(t)

	c := &lxf.Container{
		CRIObject: lxf.CRIObject{
			LXDObject: lxf.LXDObject{
				ID: "roosevelt",
			},
		},
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	}
	lt.createContainer(c)

	lt.startContainer(c.ID)

	sb := lt.getSandbox(c.Sandbox.ID)
	if len(sb.UsedBy) != 1 {
		t.Fatalf("expected used-by to be of length 1 but is %v", len(sb.UsedBy))
	}

	if sb.UsedBy[0] != c.ID {
		t.Errorf("expected used-by to be 'roosevelt' but is '%v'", sb.UsedBy[0])
	}
}

func TestSandboxStateCycle(t *testing.T) {
	lt := newLXFTest(t)

	r := setUpSandbox(lt, "roosevelt")

	sb := lt.getSandbox(r.ID)

	if sb.State != lxf.SandboxReady {
		t.Errorf("new sandbox expected to be in state 'ready' but it is '%v'", sb.State)
	}

	lt.stopSandbox("roosevelt-sb")

	sb = lt.getSandbox("roosevelt-sb")
	if sb.State != lxf.SandboxNotReady {
		t.Errorf("stoped sandbox expected to be in state 'not ready' but it is '%v'", sb.State)
	}

}
