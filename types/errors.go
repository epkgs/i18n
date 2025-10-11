package types

import (
	"fmt"
)

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
	Translator
	Storager
}
