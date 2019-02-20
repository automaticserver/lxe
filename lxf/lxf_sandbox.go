package lxf

import "fmt"

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
		return nil, err
	}

	sandboxes := []*Sandbox{}
	for _, p := range ps {
		if !IsCRI(p) {
			continue
		}
		sb, err2 := l.toSandbox(&p, ETag)
		if err2 != nil {
			return nil, err2
		}
		sandboxes = append(sandboxes, sb)
	}

	return sandboxes, nil
}
