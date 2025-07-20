package errors

import (
	"context"

	"github.com/epkgs/i18n"
)

type Error struct {
	t *i18n.Item // i18n item

	code       int // 0 success, other error
	httpStatus int // http status code

	ctx  context.Context // context
	args []any           // args for i18n

	extra any // extra data for custom error
}

func newError(t *i18n.Item, code, httpStatus int) *Error {
	return &Error{
		t:          t,
		code:       code,
		httpStatus: httpStatus,
	}
}

func (e *Error) Error() string {

	if e.ctx != nil {
		return e.t.T(e.ctx, e.args...)
	}
	return e.t.String()
}

func (e *Error) Clone() *Error {
	return &Error{
		t:          e.t,
		code:       e.code,
		httpStatus: e.httpStatus,
		ctx:        e.ctx,
		args:       e.args,
		extra:      e.extra,
	}
}

func (e *Error) AddTrans(lang, text string) {
	e.t.AddTrans(lang, text)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) HttpStatus() int {
	return e.httpStatus
}

// return new error with context
func (e *Error) WithContext(ctx context.Context) *Error {
	err := e.Clone()
	err.ctx = ctx
	return err
}

// return new error with args
func (e *Error) WithArgs(args ...any) *Error {
	err := e.Clone()
	err.args = args
	return err
}

// return new error with extra
func (e *Error) WithExtra(extra any) *Error {
	err := e.Clone()
	err.extra = extra

	return err
}
