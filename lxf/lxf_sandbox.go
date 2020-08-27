package lxf // import "github.com/automaticserver/lxe/lxf"

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/shared"
	"github.com/lxc/lxd/shared/api"
	yaml "gopkg.in/yaml.v2"
)

// NewSandbox creates a local representation of a sandbox
func (l *client) NewSandbox() *Sandbox {
	s := &Sandbox{}
	s.client = l
	s.Config = make(map[string]string)
	s.NetworkConfig.Mode = NetworkNone
	s.NetworkConfig.ModeData = make(map[string]string)

	return s
}

// GetSandbox will find a sandbox by id and return it.
func (l *client) GetSandbox(id string) (*Sandbox, error) {
	p, ETag, err := l.server.GetProfile(id)
	if err != nil {
		return nil, err
	}

	if !IsCRI(p) {
		return nil, fmt.Errorf("sandbox %w: %s", shared.NewErrNotFound(), id)
	}

	return l.toSandbox(p, ETag)
}

// ListSandboxes will return a list with all the available sandboxes
func (l *client) ListSandboxes() ([]*Sandbox, error) {
	var ETag string

	ps, err := l.server.GetProfiles()
	if err != nil {
		return nil, err
	}

	var sl = []*Sandbox{}

	for _, p := range ps {
		p := p // pin!
		if !IsCRI(p) {
			continue
		}

		s, err := l.toSandbox(&p, ETag)
		if err != nil {
			return nil, err
		}

		sl = append(sl, s)
	}

	return sl, nil
}

// toSandbox will take a profile and convert it to a sandbox.
func (l *client) toSandbox(p *api.Profile, etag string) (*Sandbox, error) {
	var err error

	var attempt uint64
	if attemptS, is := p.Config[cfgMetaAttempt]; is {
		attempt, err = strconv.ParseUint(attemptS, 10, 32)
		if err != nil {
			return nil, err
		}
	}

	createdAt := time.Time{}.UnixNano()
	if createdAtS, is := p.Config[cfgCreatedAt]; is {
		createdAt, err = strconv.ParseInt(createdAtS, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	s := &Sandbox{}
	s.client = l

	s.ID = p.Name
	s.ETag = etag
	s.Hostname = p.Config[cfgHostname]
	s.LogDirectory = p.Config[cfgLogDirectory]
	s.Metadata = SandboxMetadata{
		Attempt:   uint32(attempt),
		Name:      p.Config[cfgMetaName],
		Namespace: p.Config[cfgMetaNamespace],
		UID:       p.Config[cfgMetaUID],
	}
	s.NetworkConfig = NetworkConfig{
		Nameservers: strings.Split(p.Config[cfgNetworkConfigNameservers], ","),
		Searches:    strings.Split(p.Config[cfgNetworkConfigSearches], ","),
		Mode:        getNetworkMode(p.Config[cfgNetworkConfigMode]),
		ModeData:    make(map[string]string),
	}
	s.Labels = sandboxConfigStore.StrippedPrefixMap(p.Config, cfgLabels)
	s.Annotations = sandboxConfigStore.StrippedPrefixMap(p.Config, cfgAnnotations)
	s.Config = sandboxConfigStore.UnreservedMap(p.Config)
	s.State = getSandboxState(p.Config[cfgState])
	s.CreatedAt = time.Unix(0, createdAt)

	err = yaml.Unmarshal([]byte(p.Config[cfgNetworkConfigModeData]), &s.NetworkConfig.ModeData)
	if err != nil {
		return nil, err
	}

	// cloud-init network config & vendor-data are write-only so not read

	// get devices
	for name, options := range p.Devices {
		d, err := device.Detect(name, options)
		if err != nil {
			return nil, err
		}

		s.Devices.Upsert(d)
	}

	// get containers using this sandbox
	for _, selflink := range p.UsedBy {
		name := GetContainerIDFromSelflink(selflink)

		if name != "" {
			s.UsedBy = append(s.UsedBy, name)
		}
	}

	return s, nil
}
