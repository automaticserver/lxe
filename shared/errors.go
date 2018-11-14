package shared

const (
	// ExitCodeUnspecified is used for unspecified and unrecoverable errors
	// 2-10 are reserved for future use
	ExitCodeUnspecified = 1

	// ExitCodeSchemaMigrationFailure is returned when an error during schema migration happened
	ExitCodeSchemaMigrationFailure = 11
)
