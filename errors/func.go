package errors

import (
	"errors"

	pkgErrors "github.com/pkg/errors"
)

var (
	As     = errors.As
	Is     = errors.Is
	Unwrap = errors.Unwrap
	Cause  = pkgErrors.Cause
)

func WithStack(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return New(String(err.Error()))
}
