package transport

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

type errorSet map[error]struct{}

func newErrorSet(errs ...error) errorSet {
	s := make(errorSet)
	for _, err := range errs {
		s[err] = struct{}{}
	}
	return s
}

func (s errorSet) Has(err error) bool {
	_, ok := s[err]
	return ok
}

var badRequestErrorCodes = newErrorSet()

var notFoundErrorCodes = newErrorSet()

var unauthorizedErrorCodes = newErrorSet()

var permissionDeniedErrorCodes = newErrorSet()

var internalErrorCodes = newErrorSet()

// getGRPCCode recursively unwraps joined errors and returns GRPC code by the first meaningful error
func getGRPCCode(err error) codes.Code {
	cause := errors.Cause(err)

	switch {
	case cause == nil:
		return codes.OK
	case isBadRequestError(cause):
		return codes.InvalidArgument
	case isNotFoundError(cause):
		return codes.NotFound
	case isUnauthorizedError(cause):
		return codes.Unauthenticated
	case isPermissionDeniedError(cause):
		return codes.PermissionDenied
	case isInternalError(cause):
		return codes.Internal
	}

	switch cause {
	case context.DeadlineExceeded:
		return codes.DeadlineExceeded
	case context.Canceled:
		return codes.Canceled
	default:
		return codes.Unknown
	}
}

func isWarnLevel(err error) bool {
	switch getGRPCCode(err) {
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.PermissionDenied,
		codes.InvalidArgument,
		codes.NotFound,
		codes.FailedPrecondition,
		codes.Unauthenticated:
		return true
	default:
		return false
	}
}

func isBadRequestError(cause error) bool {
	return badRequestErrorCodes.Has(cause)
}

func isNotFoundError(cause error) bool {
	return notFoundErrorCodes.Has(cause)
}

func isUnauthorizedError(cause error) bool {
	return unauthorizedErrorCodes.Has(cause)
}

func isPermissionDeniedError(cause error) bool {
	return permissionDeniedErrorCodes.Has(cause)
}

func isInternalError(cause error) bool {
	return internalErrorCodes.Has(cause)
}
