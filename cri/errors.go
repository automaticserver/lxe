package cri

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Provide possibility to annotate errors for logging. The grpc CallTracer will try to match the returned error and log accordingly.
type AnnotatedError struct {
	Log *logrus.Entry
	Err error
	Msg string
}

func (e AnnotatedError) Error() string {
	return e.String()
}

func (e AnnotatedError) String() string {
	return fmt.Sprintf("%s: %v", e.Err, e.Log.Data)
}

func AnnErr(log *logrus.Entry, err error, msg string) error {
	return AnnotatedError{log, err, msg}
}

// Some errors should not be logged, so we can differentiate that by type
type SilentError struct {
	AnnotatedError
}

func SilErr(log *logrus.Entry, err error, msg string) error {
	return SilentError{AnnotatedError{log, err, msg}}
}
