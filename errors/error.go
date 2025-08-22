package errors

import (
	"context"
	"errors"
	"fmt"
	"io"

	pkgErrors "github.com/pkg/errors"
)

var (
	As     = errors.As
	Is     = errors.Is
	Unwrap = errors.Unwrap
	Cause  = pkgErrors.Cause
)

// String is a string type that implements the Stringer interface
type String string

// String returns the string representation of the String type
func (s String) String() string {
	return string(s)
}

func toStringer(str any) fmt.Stringer {
	switch s := str.(type) {
	case fmt.Stringer:
		return s
	case string:
		return String(s)
	default:
		return String(fmt.Sprintf("%v", s))
	}
}

// i18nError represents an internationalized error with additional metadata
type i18nError struct {
	msg   fmt.Stringer
	extra map[string]any

	cause error
	stack *stack
}

// New creates and returns a new Error with the given message
// The message can be a string, fmt.Stringer, or any other type
func New(msg any) Error {
	return &i18nError{
		msg:   toStringer(msg),
		extra: map[string]any{},
		stack: callers(),
	}
}

// Errorf creates and returns a new Error with formatted message
// It uses fmt.Sprintf to format the message with the given arguments
func Errorf(format string, args ...any) Error {
	return New(fmt.Sprintf(format, args...))
}

// WithStack wraps an error with stack trace information
// If the error is already an Error, it adds stack trace to it
// Otherwise, it creates a new Error with the given error as cause
func WithStack(err error) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return e.WithStack()
	}
	return Wrap(err, "")
}

// Wrap creates a new Error that wraps another error with an additional message
// If the error is nil, it returns nil
func Wrap(err error, message string) Error {
	if err == nil {
		return nil
	}
	return New(message).Wrap(err)
}

// String returns the string representation of the error message
func (e *i18nError) String() string {
	return e.msg.String()
}

// Error returns the full error message including cause if present
func (e *i18nError) Error() string {
	msg := e.msg.String()
	if e.cause != nil {
		if msg != "" {
			msg += ": "
		}
		return msg + e.cause.Error()
	}

	return e.msg.String()
}

// T returns the translated error message based on context language preferences
// If the message does not support translation, it returns the default string representation
func (e *i18nError) T(ctx context.Context) string {
	if tran, ok := e.msg.(interface {
		T(ctx context.Context) string
	}); ok {
		return tran.T(ctx)
	}
	return e.msg.String()
}

// clone creates and returns a copy of the error
func (e *i18nError) clone() *i18nError {
	err := &i18nError{
		msg:   e.msg,
		extra: make(map[string]any, len(e.extra)),
		cause: e.cause,
		stack: e.stack,
	}

	for i, v := range e.extra {
		err.extra[i] = v
	}

	return err
}

// WithMessage creates a new error instance with the same properties as the current error but with a different message content.
// It takes a message parameter of any type, converts it to a string representation,
// and creates a new error object that retains all properties of the original error
// (such as stack trace, extra data, etc.) but uses the new message content.
func (e *i18nError) WithMessage(msg any) Error {
	err := e.clone()
	err.msg = toStringer(msg)
	return err
}

// WithStack returns a copy of the error with a new stack trace
func (e *i18nError) WithStack() Error {
	err := e.clone()
	err.stack = callers()
	return err
}

// Wrap returns a copy of the error with the given cause error
// It also adds a new stack trace
func (e *i18nError) Wrap(cause error) Error {
	err := e.clone()
	err.cause = cause
	err.stack = callers()
	return err
}

// Cause returns the underlying cause error
func (e *i18nError) Cause() error { return e.cause }

// Unwrap provides compatibility for Go 1.13 error chains
func (e *i18nError) Unwrap() error { return e.cause }

// Format implements the fmt.Formatter interface to provide custom error formatting
// With %+v it prints the error with stack trace
func (e *i18nError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if e.cause != nil {
				fmt.Fprintf(s, "%+v\n", e.cause)
			}
			io.WriteString(s, e.msg.String())
			if e.stack != nil {
				e.stack.Format(s, verb)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg.String())
	case 'q':
		fmt.Fprintf(s, "%q", e.msg.String())
	}
}

// Is checks if the error is equivalent to another error
// Two errors are equivalent if they have the same message and extra data
func (e *i18nError) Is(err error) bool {
	er, ok := err.(*i18nError)
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

// Set stores a key-value pair in the error's extra data
func (e *i18nError) Set(key string, value any) {
	e.extra[key] = value
}

// Get retrieves a value from the error's extra data by key
// If the key does not exist, it returns the provided default value
func (e *i18nError) Get(key string, defaultValue any) any {
	v, ok := e.extra[key]
	if !ok {
		return defaultValue
	}
	return v
}

// Has checks if a key exists in the error's extra data
func (e *i18nError) Has(key string) bool {
	_, ok := e.extra[key]
	return ok
}
