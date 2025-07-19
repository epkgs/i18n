package errorx

import (
	"context"

	"github.com/epkgs/i18n"
)

type Error struct {
	*i18n.Item

	code       int
	httpStatus int

	ctx  context.Context
	args []any
}

func (e *Error) Error() string {

	if e.ctx != nil {
		return e.T(e.ctx, e.args...)
	}
	return e.String()
}

func (e *Error) AddTrans(lang, text string) *Error {
	e.Item.AddTrans(lang, text)
	return e
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) HttpStatus() int {
	return e.httpStatus
}

func (e *Error) WithContext(ctx context.Context) *Error {
	e.ctx = ctx
	return e
}

func (e *Error) WithArgs(args ...any) *Error {
	e.args = args
	return e
}
