package errors

import (
	"context"
)

type Translable interface {
	Translate(ctx context.Context) string
}

type Storage interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Exist(key string) bool
}

type I18nError interface {
	error
	Translable
	WithStack() I18nError
	Wrap(cause error) I18nError
	Cause() error
	Unwrap() error
	// Format(s fmt.State, verb rune)
	Is(err error) bool
	Storage
}

type Definition[Args any] interface {
	New(args Args) I18nError
}
