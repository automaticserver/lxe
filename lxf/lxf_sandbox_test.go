package lxf

import (
	"strconv"
	"testing"
	"time"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func basicProfile(name string) *api.Profile {
	p := &api.Profile{Name: name}
	p.Config = map[string]string{}
	satisfyProfileSchema(satisfyProfileCri(p))

	return p
}

func TestClient_NewSandbox(t *testing.T) {
	t.Parallel()

	client, _ := testClient()

	exp := &Sandbox{}
	exp.client = client
	exp.Config = make(map[string]string)
	exp.NetworkConfig.Mode = NetworkNone
	exp.NetworkConfig.ModeData = make(map[string]string)

	s := client.NewSandbox()

	assert.Exactly(t, exp, s)
}

func TestClient_GetSandbox_Minimal(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetProfileReturns(basicProfile("foo"), "", nil)

	s, err := client.GetSandbox("foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo", s.ID)
	assert.Equal(t, "foo", fake.GetProfileArgsForCall(0))
	assert.Equal(t, 1, fake.GetProfileCallCount())
}

func TestClient_GetSandbox_Missing(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetProfileReturns(nil, "", shared.NewErrNotFound())

	s, err := client.GetSandbox("foo")

	var expected *Sandbox = nil

	assert.Error(t, err)
	assert.Exactly(t, expected, s)
	assert.Equal(t, 1, fake.GetProfileCallCount())
}

func TestClient_GetSandbox_NonCri(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetProfileReturns(&api.Profile{}, "", nil)

	s, err := client.GetSandbox("foo")

	var expected *Sandbox = nil

	assert.Error(t, err)
	assert.Exactly(t, expected, s)
	assert.Equal(t, 1, fake.GetProfileCallCount())
}

func TestClient_ListSandboxes_Minimal(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetProfilesReturns([]api.Profile{*basicProfile("foo"), *basicProfile("bar")}, nil)

	sl, err := client.ListSandboxes()
	assert.NoError(t, err)
	assert.Len(t, sl, 2)
	assert.Equal(t, 1, fake.GetProfilesCallCount())
}

func TestClient_ListSandboxes_Error(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetProfilesReturns([]api.Profile{*basicProfile("foo"), *basicProfile("bar")}, shared.NewErrNotFound())

	sl, err := client.ListSandboxes()
	assert.Error(t, err)
	assert.Len(t, sl, 0)
	assert.Equal(t, 1, fake.GetProfilesCallCount())
}

func TestClient_ListSandboxes_Missing(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetProfilesReturns([]api.Profile{}, nil)

	sl, err := client.ListSandboxes()
	assert.NoError(t, err)
	assert.Len(t, sl, 0)
	assert.Equal(t, 1, fake.GetProfilesCallCount())
}

func TestClient_ListSandboxes_NonCri(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetProfilesReturns([]api.Profile{{Name: "foo"}, {Name: "bar"}}, nil)

	sl, err := client.ListSandboxes()
	assert.NoError(t, err)
	assert.Len(t, sl, 0)
	assert.Equal(t, 1, fake.GetProfilesCallCount())
}

func TestClient_toSandbox_AllFieldsSuccessful(t *testing.T) {
	t.Parallel()

	client, _ := testClient()

	now := time.Unix(0, time.Now().UnixNano())

	p := &api.Profile{
		Name: "profileName",
		ProfilePut: api.ProfilePut{
			Config: map[string]string{
				cfgMetaName:                      "metaName",
				cfgMetaNamespace:                 "metaNamespace",
				cfgMetaUID:                       "metaUID",
				cfgMetaAttempt:                   "1",
				cfgCreatedAt:                     strconv.FormatInt(now.UnixNano(), 10),
				cfgHostname:                      "hostname",
				cfgLogDirectory:                  "logDirectory",
				cfgNetworkConfigNameservers:      "1.2.3.4,5.6.7.8",
				cfgNetworkConfigSearches:         "svc.local,local",
				cfgNetworkConfigMode:             "none",
				cfgLabels + ".alabel":            "aLabel",
				cfgAnnotations + ".anannotation": "anAnnotation",
				"something.else":                 "somethingElse",
				cfgState:                         "notready",
				cfgNetworkConfigModeData:         "mode: data",
			},
			Devices: map[string]map[string]string{
				"first": {
					"type": "none",
				},
			},
		},
		UsedBy: []string{"/1.0/containers/containerName"},
	}

	exp := &Sandbox{}
	exp.ID = "profileName"
	exp.ETag = "etag"
	exp.client = client
	exp.Devices = []device.Device{
		&device.None{KeyName: "first"},
	}
	exp.Config = map[string]string{"something.else": "somethingElse"}
	exp.UsedBy = []string{"containerName"}
	exp.Labels = map[string]string{"alabel": "aLabel"}
	exp.Annotations = map[string]string{"anannotation": "anAnnotation"}
	exp.CreatedAt = now
	exp.Metadata.Attempt = 1
	exp.Metadata.Name = "metaName"
	exp.Metadata.Namespace = "metaNamespace"
	exp.Metadata.UID = "metaUID"
	exp.Hostname = "hostname"
	exp.NetworkConfig.Nameservers = []string{"1.2.3.4", "5.6.7.8"}
	exp.NetworkConfig.Searches = []string{"svc.local", "local"}
	exp.NetworkConfig.Mode = NetworkNone
	exp.NetworkConfig.ModeData = map[string]string{"mode": "data"}
	exp.State = SandboxNotReady
	exp.LogDirectory = "logDirectory"

	s, err := client.toSandbox(p, "etag")
	assert.NoError(t, err)
	assert.Exactly(t, exp, s)
}
