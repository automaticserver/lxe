package lxf

import (
	"strconv"

	"github.com/lxc/lxd/shared/api"
)

// nolint: gosec #nosec (no sensitive data)

const (
	cfgIsCRI       = "user.cri"
	cfgLabels      = "user.labels"
	cfgAnnotations = "user.annotations"

	cfgMetaAttempt   = "user.metadata.attempt"
	cfgMetaName      = "user.metadata.name"
	cfgMetaNamespace = "user.metadata.namespace"
	cfgMetaUID       = "user.metadata.uid"
)

// CRIObject contains common properties of containers and sandboxes
type CRIObject struct {
	LXDObject
	// Labels and Annotations to be saved provided by CRI
	Labels      map[string]string
	Annotations map[string]string
}

// IsCRI checks if a object is a cri object
func IsCRI(i interface{}) bool {
	if !IsSchemaCurrent(i) {
		return false
	}

	var val string
	var has bool

	switch o := i.(type) {
	case api.Container:
		if val, has = o.Config[cfgIsCRI]; !has {
			return false
		}
	case api.Profile:
		if val, has = o.Config[cfgIsCRI]; !has {
			return false
		}
	default:
		return false
	}

	is, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return is
}
