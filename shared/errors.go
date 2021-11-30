package shared

import (
	"errors"
)

// ExitCodeUnspecified is used for unspecified and unrecoverable errors
// 2-10 are reserved for future use
const ExitCodeUnspecified = 1

// ExitCodeSchemaMigrationFailure is returned when an error during schema migration happened
const ExitCodeSchemaMigrationFailure = 11

// In some situations in the code, we give the function caller a simple "not found" error, while technically it was not the LXD api returning that. This is the case, when we filter the response and if no result is left we imitate the "not found" (e.g. IsCRI()). Since LXD doesn't offer any error type, we also have to be able to have a similar function to errors.Is for the "not found" without allowing it to be used directly so we cleanly can differentiate between application error and api error.

// ErrLXDNotFound is the error string a LXD request returns, when nothing is found, so far for many/all resource types
const LXDNotFound = "not found"

// errLXDNotFound represents an LXD sourced error with string "not found", only allowed to be used in cases where we extend or imitate the behaviour of the LXD API!
var errLXDNotFound = errors.New(LXDNotFound)

// To minimize repetition and to also comply to err113 linter and the golang 1.13 error usage: These functions offer an Is check to exactly the LXD "not found" error and a way to imitate the error without allowing it to be use directly, since the error type is intentionally not exported from this package

func IsErrNotFound(err error) bool {
	if errors.Is(err, errLXDNotFound) {
		return true
	}

	return err.Error() == LXDNotFound
}

func NewErrNotFound() error {
	return errLXDNotFound
}
