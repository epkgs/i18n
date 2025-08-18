package errors

import (
	"errors"
	"fmt"

	pkgErrors "github.com/pkg/errors"
)

var (
	As     = errors.As
	Is     = errors.Is
	Unwrap = errors.Unwrap
	Cause  = pkgErrors.Cause
)

func New(format fmt.Stringer) I18nError {
	return &Error{
		msg:   format,
		extra: map[string]any{},
		stack: callers(),
	}
}

func WithStack(err error) I18nError {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return New(String(err.Error()))
}
