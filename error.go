package i18n

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Error struct {
	t *Item

	ctx  context.Context
	args []any

	cause error // 原始错误
}

func newError(t *Item, ctx context.Context, args ...any) *Error {
	return &Error{
		t:    t,
		ctx:  ctx,
		args: args,
	}
}

func (e *Error) Error() string {
	if e.ctx != nil {
		return e.t.T(e.ctx, e.args...)
	}

	return e.t.T(context.Background(), e.args...)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) WithCause(err error) {
	e.cause = err
}

func (e *Error) Is(err error) bool {
	err2, ok := err.(*Error)
	if !ok {
		return false
	}

	if e.t.I18n.Name() != err2.t.I18n.Name() {
		return false
	}

	if e.t.String() != err2.t.String() {
		return false
	}

	return true
}

// As 实现 errors.As 接口
func (e *Error) As(target any) bool {
	if target == nil {
		return false
	}

	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return false
	}

	if reflect.TypeOf(e).AssignableTo(targetType.Elem()) {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(e))
		return true
	}

	// 如果有 cause，递归调用 errors.As
	if e.cause != nil {
		return errors.As(e.cause, target)
	}

	return false
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s\n", e.Error())
			if e.cause != nil {
				fmt.Fprintf(s, "Cause: %+v\n", e.cause)
			}
			return
		}
		fallthrough
	default:
		fmt.Fprintf(s, "%s", e.Error())
	}
}
