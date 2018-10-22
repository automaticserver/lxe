package lxf

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
