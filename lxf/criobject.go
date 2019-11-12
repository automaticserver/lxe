package lxf

import (
	"strconv"
	"time"

	"github.com/lxc/lxd/shared/api"
)

const (
	cfgIsCRI         = "user.cri"
	cfgLabels        = "user.labels"
	cfgAnnotations   = "user.annotations"
	cfgState         = "user.state"
	cfgMetadata      = "user.metadata"
	cfgMetaAttempt   = cfgMetadata + ".attempt"
	cfgMetaName      = cfgMetadata + ".name"
	cfgMetaNamespace = cfgMetadata + ".namespace"
	cfgMetaUID       = cfgMetadata + ".uid"
	cfgVolatile      = "volatile"
)

var (
	reservedConfigCRI = []string{
		cfgSchema,
		cfgIsCRI,
		cfgCreatedAt,
	}
	reservedConfigPrefixesCRI = []string{
		cfgVolatile,
		cfgLabels,
		cfgAnnotations,
		cfgMetadata,
	}
)

// CRIObject contains common properties of containers and sandboxes
type CRIObject struct {
	// Labels and Annotations to be saved provided by CRI
	Labels      map[string]string
	Annotations map[string]string
	// CreatedAt is when the resource was created
	CreatedAt time.Time
}

// IsCRI checks if a object is a cri object
func IsCRI(i interface{}) bool {
	if !IsSchemaCurrent(i) {
		return false
	}

	var (
		val string
		has bool
	)

	switch o := i.(type) {
	case api.Container:
		if val, has = o.Config[cfgIsCRI]; !has {
			return false
		}
	case *api.Container:
		return IsCRI(*o)
	case api.Profile:
		if val, has = o.Config[cfgIsCRI]; !has {
			return false
		}
	case *api.Profile:
		return IsCRI(*o)
	default:
		return false
	}

	is, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return is
}
