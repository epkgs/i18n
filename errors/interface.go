package errors

import (
	"context"
	"fmt"
)

type Translatable interface {
	T(ctx context.Context) string
}

type Storager interface {
	Set(key string, value any)
	Get(key string, defaultValue any) any
	Has(key string) bool
}

type Error interface {
	error
	WithMsg(msg any) Error
	WithStack() Error
	Wrap(cause error) Error
	Cause() error
	Unwrap() error
	Format(s fmt.State, verb rune)
	Is(err error) bool

	fmt.Stringer
	Translatable
	Storager
}
