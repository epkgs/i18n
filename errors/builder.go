package errors

import (
	"github.com/epkgs/i18n"
)

type Builder struct {
	I18n *i18n.I18n
}

func NewBuilder(name string, fn ...i18n.OptionsFunc) *Builder {
	return &Builder{
		I18n: i18n.New(name, fn...),
	}
}

func (b *Builder) New(code int, txt string, httpStatus int) *Error {
	return New(b.I18n.New(txt), code, httpStatus)
}
