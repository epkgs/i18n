package errors

import (
	"github.com/epkgs/i18n"
)

type ErrorConstructor[E error] func(code int, t *i18n.Item, httpStatus int) E

type Builder[E error] struct {
	I18n *i18n.I18n

	constructor ErrorConstructor[E]
}

func NewBuilderCustom[E error](name string, constructor ErrorConstructor[E], fn ...i18n.OptionsFunc) *Builder[E] {
	return &Builder[E]{
		I18n:        i18n.New(name, fn...),
		constructor: constructor,
	}
}

func NewBuilder(name string, fn ...i18n.OptionsFunc) *Builder[*Error] {
	return NewBuilderCustom(name, New, fn...)
}

func (b *Builder[E]) New(code int, txt string, httpStatus int) E {
	return b.constructor(code, b.I18n.New(txt), httpStatus)
}
