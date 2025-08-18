package errors

import (
	"context"
	"fmt"
	"io"
)

type String string

func (s String) String() string {
	return string(s)
}

type Error struct {
	msg   fmt.Stringer
	args  []any
	extra map[string]any

	cause error
	stack *stack
}

func New(format fmt.Stringer) *Error {
	return Errorf(format)
}

func Errorf(format fmt.Stringer, args ...any) *Error {

	err := &Error{
		msg:   format,
		args:  args,
		extra: map[string]any{},
		stack: callers(),
	}

	return err
}

func (e *Error) String() string {
	return e.msg.String()
}

func (e *Error) Error() string {
	if e.cause != nil {
		return e.String() + ": " + e.cause.Error()
	}

	return e.String()
}

func (e *Error) Translate(ctx context.Context) string {
	if tran, ok := e.msg.(Translable); ok {
		return tran.Translate(ctx)
	}
	return e.String()
}

func (e *Error) Wrap(cause error) I18nError {
	e.cause = cause
	return e
}

func (e *Error) Cause() error { return e.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", e.Cause())
			io.WriteString(s, e.String())
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.String())
	case 'q':
		fmt.Fprintf(s, "%q", e.String())
	}
}

func (e *Error) Is(err error) bool {
	er, ok := err.(*Error)
	if !ok {
		return false
	}

	if er.msg != e.msg {
		return false
	}

	// compareMaps checks if two maps[string]any are equal in keys and values.
	compareMaps := func(a, b map[string]any) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		if len(a) != len(b) {
			return false
		}
		for k, v := range a {
			if bv, ok := b[k]; !ok || bv != v {
				return false
			}
		}
		return true
	}

	return compareMaps(er.extra, e.extra)
}

func (e *Error) Set(key string, value any) {
	e.extra[key] = value
}

func (e *Error) Get(key string) (any, bool) {
	v, ok := e.extra[key]
	return v, ok
}

func (e *Error) Exist(key string) bool {
	_, ok := e.extra[key]
	return ok
}
