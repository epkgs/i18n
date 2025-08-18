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
	return Errorf(format)
}

func Errorf(format fmt.Stringer, args ...any) I18nError {

	err := &Error{
		msg:   format,
		args:  args,
		extra: map[string]any{},
		stack: callers(),
	}

	return err
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

func Wrapf(err error, format fmt.Stringer, args ...any) I18nError {
	return Errorf(format, args...).Wrap(err)
}
