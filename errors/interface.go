package errors

import (
	"context"
)

type Translable interface {
	Translate(ctx context.Context) string
}

type Coder interface {
	Code() int
}

type HttpStatuser interface {
	HttpStatus() int
}

type Storage interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Exist(key string) bool
}

type I18nError interface {
	error
	Translable
	Wrap(cause error) I18nError
	Cause() error
	Unwrap() error
	// Format(s fmt.State, verb rune)
	Is(err error) bool
	Storage
	WithCode(code int)
	Coder
	WithHttpStatus(httpStatus int)
	HttpStatuser
}
