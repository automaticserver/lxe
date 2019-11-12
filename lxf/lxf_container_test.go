package lxf

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/lxc/lxd/shared/api"
	opencontainers "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"
)

func basicContainer(name, sandbox string) *api.Container {
	c := &api.Container{Name: name}
	c.Profiles = []string{sandbox}
	c.Config = map[string]string{}
	satisfyContainerSchema(satisfyContainerCri(c))

	return c
}

func TestClient_NewContainer(t *testing.T) {
	t.Parallel()

	client, _ := testClient()

	exp := &Container{}
	exp.client = client
	exp.Profiles = []string{
		"default",
		"sandboxID",
	}

	s := client.NewContainer("sandboxID", "default")

	assert.Exactly(t, exp, s)
}

func TestClient_GetContainer_Minimal(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetContainerReturns(basicContainer("foo", "bar"), "", nil)

	s, err := client.GetContainer("foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo", s.ID)
	assert.Equal(t, "foo", fake.GetContainerArgsForCall(0))
	assert.Equal(t, 1, fake.GetContainerCallCount())
}

func TestClient_GetContainer_Missing(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetContainerReturns(nil, "", fmt.Errorf(ErrorLXDNotFound))

	s, err := client.GetContainer("foo")

	var expected *Container = nil

	assert.Error(t, err)
	assert.Exactly(t, expected, s)
	assert.Equal(t, 1, fake.GetContainerCallCount())
}

func TestClient_GetContainer_NonCri(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetContainerReturns(&api.Container{}, "", nil)

	s, err := client.GetContainer("foo")

	var expected *Container = nil

	assert.Error(t, err)
	assert.Exactly(t, expected, s)
	assert.Equal(t, 1, fake.GetContainerCallCount())
}

func TestClient_ListContainers_Minimal(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetContainersReturns([]api.Container{*basicContainer("foo", "default"), *basicContainer("bar", "default")}, nil)

	sl, err := client.ListContainers()
	assert.NoError(t, err)
	assert.Len(t, sl, 2)
	assert.Equal(t, 1, fake.GetContainersCallCount())
}

func TestClient_ListContainers_Error(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetContainersReturns([]api.Container{*basicContainer("foo", "default"), *basicContainer("bar", "default")}, fmt.Errorf(ErrorLXDNotFound))

	sl, err := client.ListContainers()
	assert.Error(t, err)
	assert.Len(t, sl, 0)
	assert.Equal(t, 1, fake.GetContainersCallCount())
}

func TestClient_ListContainers_Missing(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetContainersReturns([]api.Container{}, nil)

	sl, err := client.ListContainers()
	assert.NoError(t, err)
	assert.Len(t, sl, 0)
	assert.Equal(t, 1, fake.GetContainersCallCount())
}

func TestClient_ListContainers_NonCri(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	fake.GetContainersReturns([]api.Container{{Name: "foo"}, {Name: "bar"}}, nil)

	sl, err := client.ListContainers()
	assert.NoError(t, err)
	assert.Len(t, sl, 0)
	assert.Equal(t, 1, fake.GetContainersCallCount())
}

func TestClient_toContainer_AllFieldsSuccessful(t *testing.T) {
	t.Parallel()

	client, _ := testClient()

	now := time.Unix(0, time.Now().UnixNano())
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	ct := &api.Container{
		Name: "containerName",
		ContainerPut: api.ContainerPut{
			Config: map[string]string{
				cfgVolatileBaseImage:             "image",
				cfgMetaName:                      "metaName",
				cfgMetaAttempt:                   "1",
				cfgLabels + ".alabel":            "aLabel",
				cfgAnnotations + ".anannotation": "anAnnotation",
				"something.else":                 "somethingElse",
				cfgLogPath:                       "logPath",
				cfgCreatedAt:                     strconv.FormatInt(now.UnixNano(), 10),
				cfgStartedAt:                     strconv.FormatInt(past.UnixNano(), 10),
				cfgFinishedAt:                    strconv.FormatInt(future.UnixNano(), 10),
				cfgEnvironmentPrefix + ".data":   "content",
				cfgSecurityPrivileged:            "true",
				cfgCloudInitUserData:             "userData",
				cfgCloudInitMetaData:             "metaData",
				cfgCloudInitNetworkConfig:        "networkConfig",
				cfgResourcesCPUShares:            "600",
				cfgResourcesCPUQuota:             "300",
				cfgResourcesCPUPeriod:            "100",
				cfgResourcesMemoryLimit:          "1234567",
			},
			Devices: map[string]map[string]string{
				"first": {
					"type": "none",
				},
			},
			Profiles: []string{"profile"},
		},
		StatusCode: api.Aborting,
	}

	exp := &Container{}
	exp.ID = "containerName"
	exp.ETag = "etag"
	exp.client = client
	exp.Devices = []device.Device{
		&device.None{KeyName: "first"},
	}
	exp.Config = map[string]string{"something.else": "somethingElse"}
	exp.Profiles = []string{"profile"}
	exp.Image = "image"
	exp.Privileged = true
	exp.Environment = map[string]string{"data": "content"}
	exp.Labels = map[string]string{"alabel": "aLabel"}
	exp.Annotations = map[string]string{"anannotation": "anAnnotation"}
	exp.Metadata.Attempt = 1
	exp.Metadata.Name = "metaName"
	exp.CreatedAt = now
	exp.StartedAt = past
	exp.FinishedAt = future
	exp.StateName = ContainerStateExited
	exp.LogPath = "logPath"
	exp.CloudInitUserData = "userData"
	exp.CloudInitMetaData = "metaData"
	exp.CloudInitNetworkConfig = "networkConfig"

	var shares uint64 = 600
	var quota int64 = 300
	var period uint64 = 100
	var memory int64 = 1234567

	exp.Resources = &opencontainers.LinuxResources{
		CPU: &opencontainers.LinuxCPU{
			Shares: &shares,
			Quota:  &quota,
			Period: &period,
		},
		Memory: &opencontainers.LinuxMemory{
			Limit: &memory,
		},
	}

	c, err := client.toContainer(ct, "etag")
	assert.NoError(t, err)
	assert.Exactly(t, exp, c)
}

// TODO lifecycle event handler, but first network modes need an interface
