package cri

import (
	"fmt"
	"path"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"k8s.io/utils/exec"
)

// callTracing logs requests, responses and error returned by the handler. What gets logged is influenced by what error types the handler returns and the log level. This simplifies error logging in the CRI implementation.
func callTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := log.WithContext(ctx)
	method := path.Base(info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		// Depending on the error type the logging is influenced and the error return type modified
		switch e := err.(type) { // nolint: errorlint
		// The AnnotatedError uses the provided logger entry to set fields of the actual logger
		case AnnotatedError:
			log.WithError(e.Err).WithFields(e.Log.Data).Error(fmt.Sprintf("%s: %s", method, e.Msg))
			err = status.Error(e.Code, e.Error())
		// SilentErrors are useful for not implemented functions, still return the error to the caller!
		case SilentError:
			err = status.Error(e.Code, e.Error())
		// CodeExitError is a special wrapping of AnnotatedError and exec.CodeExitError
		// TODO: this can be made better
		case *exec.CodeExitError:
			a, is := e.Err.(AnnotatedError) // nolint: errorlint
			if is {
				log.WithError(a.Err).WithFields(a.Log.Data).Error(fmt.Sprintf("%s: %s", method, a.Msg))
			} else {
				log.Error(fmt.Sprintf("%s: %s", method, err.Error()))
			}
		// In any other case just log the error
		default:
			log.WithError(err).Error(fmt.Sprintf("%s: %s", method, "untyped error"))
		}
	}

	log.WithError(err).WithFields(logrus.Fields{
		"req":  req,
		"resp": resp,
	}).Trace(fmt.Sprintf("grpc %s", method))

	return resp, err
}
