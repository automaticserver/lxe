package lxf

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/automaticserver/lxe/lxf/device"
	yaml "gopkg.in/yaml.v2"
)

// NewSandbox creates a local representation of a sandbox
func (l *Client) NewSandbox() *Sandbox {
	s := &Sandbox{}
	s.client = l
	return s
}

// GetSandbox will find a sandbox by id and return it.
func (l *Client) GetSandbox(id string) (*Sandbox, error) {
	p, ETag, err := l.server.GetProfile(id)
	if err != nil {
		return nil, NewSandboxError(id, err)
	}

	if !IsCRI(p) {
		return nil, NewSandboxError(id, fmt.Errorf(ErrorLXDNotFound))
	}
	return l.toSandbox(p, ETag)
}

// ListSandboxes will return a list with all the available sandboxes
func (l *Client) ListSandboxes() ([]*Sandbox, error) {
	ETag := ""
	ps, err := l.server.GetProfiles()
	if err != nil {
		return nil, NewSandboxError("lxdApi", err)
	}

	sl := []*Sandbox{}
	for _, p := range ps {
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
func (l *Client) toSandbox(p *api.Profile, ETag string) (*Sandbox, error) {
	attempts, err := strconv.ParseUint(p.Config[cfgMetaAttempt], 10, 32)
	if err != nil {
		return nil, err
	}
	createdAt, err := strconv.ParseInt(p.Config[cfgCreatedAt], 10, 64)
	if err != nil {
		return nil, err
	}

	s := &Sandbox{}
	s.client = l

	s.ID = p.Name
	s.ETag = ETag
	s.Hostname = p.Config[cfgHostname]
	s.LogDirectory = p.Config[cfgLogDirectory]
	s.Metadata = SandboxMetadata{
		Attempt:   uint32(attempts),
		Name:      p.Config[cfgMetaName],
		Namespace: p.Config[cfgMetaNamespace],
		UID:       p.Config[cfgMetaUID],
	}
	s.NetworkConfig = NetworkConfig{
		Nameservers: strings.Split(p.Config[cfgNetworkConfigNameservers], ","),
		Searches:    strings.Split(p.Config[cfgNetworkConfigSearches], ","),
		Mode:        getNetworkMode(p.Config[cfgNetworkConfigMode]),
		// ModeData:    make(map[string]string),
	}
	s.Labels = sandboxConfigStore.StripedPrefixMap(p.Config, cfgLabels)
	s.Annotations = sandboxConfigStore.StripedPrefixMap(p.Config, cfgAnnotations)
	s.Config = sandboxConfigStore.UnreservedMap(p.Config)
	s.State = getSandboxState(p.Config[cfgState])
	s.CreatedAt = time.Unix(0, createdAt)

	err = yaml.Unmarshal([]byte(p.Config[cfgNetworkConfigModeData]), &s.NetworkConfig.ModeData)
	if err != nil {
		return nil, err
	}
	if len(s.NetworkConfig.ModeData) == 0 {
		s.NetworkConfig.ModeData = make(map[string]string)
	}

	// cloud-init network config & vendor-data are write-only so not read

	// get devices
	s.Proxies, err = device.GetProxiesFromMap(p.Devices)
	if err != nil {
		return nil, err
	}
	s.Disks, err = device.GetDisksFromMap(p.Devices)
	if err != nil {
		return nil, err
	}
	s.Blocks, err = device.GetBlocksFromMap(p.Devices)
	if err != nil {
		return nil, err
	}
	s.Nics, err = device.GetNicsFromMap(p.Devices)
	if err != nil {
		return nil, err
	}
	s.Nones, err = device.GetNonesFromMap(p.Devices)
	if err != nil {
		return nil, err
	}

	// get containers using this sandbox
	for _, name := range p.UsedBy {
		name = strings.TrimPrefix(name, "/1.0/containers/")
		name = strings.TrimSuffix(name, "?project=default")
		if strings.Contains(name, shared.SnapshotDelimiter) {
			// this is a snapshot so dont parse this entry
			continue
		}
		s.UsedBy = append(s.UsedBy, name)
	}

	return s, nil
}
