package errorx

import (
	"github.com/epkgs/i18n"
)

type Factory struct {
	I18n *i18n.I18n
}

func NewFactory(name string, fn ...i18n.OptionsFunc) *Factory {
	return &Factory{
		I18n: i18n.New(name, fn...),
	}
}

func (f *Factory) newBuilder(code, httpStatus int, msg string) *Builder {
	item := f.I18n.NewItem(msg)

	err := &Error{
		Item:       item,
		code:       code,
		httpStatus: httpStatus,
	}

	return &Builder{def: err}
}

func (f *Factory) New(code, httpStatus int, txt string) BuilderA0 {
	return f.newBuilder(code, httpStatus, txt)
}

func (f *Factory) NewA1(code, httpStatus int, txt string) BuilderA1 {
	return f.newBuilder(code, httpStatus, txt)
}

func (f *Factory) NewA2(code, httpStatus int, txt string) BuilderA2 {
	return f.newBuilder(code, httpStatus, txt)
}

func (f *Factory) NewA3(code, httpStatus int, txt string) BuilderA3 {
	return f.newBuilder(code, httpStatus, txt)
}

func (f *Factory) NewA4(code, httpStatus int, txt string) BuilderA4 {
	return f.newBuilder(code, httpStatus, txt)
}

func (f *Factory) NewA5(code, httpStatus int, txt string) BuilderA5 {
	return f.newBuilder(code, httpStatus, txt)
}

func (f *Factory) NewAN(code, httpStatus int, txt string) BuilderAN {
	return f.newBuilder(code, httpStatus, txt)
}
