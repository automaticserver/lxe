package cri

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

// Provide possibility to annotate errors for logging. The grpc CallTracer will try to match the returned error and log accordingly.
type AnnotatedError struct {
	Log  *logrus.Entry
	Code codes.Code
	Err  error
	Msg  string
}

func (e AnnotatedError) Error() string {
	return e.String()
}

func (e AnnotatedError) String() string {
	if e.Msg != "" {
		return fmt.Sprintf("%s: %s: %v", e.Msg, e.Err, e.Log.Data)
	}

	return fmt.Sprintf("%s: %v", e.Err, e.Log.Data)
}

func AnnErr(log *logrus.Entry, code codes.Code, err error, msg string) error {
	return AnnotatedError{log, code, err, msg}
}

// Some errors should not be logged, so we can differentiate that by type
type SilentError struct {
	AnnotatedError
}

func SilErr(log *logrus.Entry, code codes.Code, err error, msg string) error {
	return SilentError{AnnotatedError{log, code, err, msg}}
}
