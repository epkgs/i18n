package errorx

import (
	"context"

	"github.com/epkgs/i18n"
)

func Define[Args any, E error](i18n *i18n.I18n, format string, wrapper Wrapper[E]) *Definition[E, Args] {
	return newDefinition[E, Args](i18n.New(format), wrapper)
}

func DefineSimple(i18n *i18n.I18n, format string) *Definition[*Error, struct{}] {

	wrapper := func(err *Error) *Error {
		return err
	}

	return Define[struct{}](i18n, format, wrapper)
}

type Wrapper[E error] func(*Error) E

type Definition[E error, Args any] struct {
	t       *i18n.Item // i18n item
	base    E
	wrapper Wrapper[E]
}

// newDefinition 创建并初始化一个新的错误定义对象。
// 该函数接收一个Item类型的参数t，以及一个可变长参数wrappers，后者是由ErrorWrapper接口类型的对象组成。
// Item类型和ErrorWrapper接口的具体定义未在上下文中给出，因此这里不做具体说明。
// 函数返回一个指向errorDefinition类型的指针。
func newDefinition[E error, Args any](t *i18n.Item, wrapper Wrapper[E]) *Definition[E, Args] {
	// 创建一个errorDefinition类型的对象并对其进行初始化。
	d := &Definition[E, Args]{
		t:       t,
		wrapper: wrapper,
	}

	d.base = d.New(context.Background(), *new(Args))

	// 返回初始化后的errorDefinition对象。
	return d
}

// Definition 返回错误定义的基础错误。
// 该方法允许访问错误定义内部的基础错误，以便在需要时进行进一步处理或检查。
func (d *Definition[E, Args]) Definition() E {
	return d.base
}

func (d *Definition[E, Args]) New(ctx context.Context, args Args) E {
	err := newError(d.t, ctx, args)
	return d.wrapper(err)
}

func (d *Definition[E, Args]) Code() int {
	if code, ok := ErrorCode(d.Definition()); ok {
		return code
	}
	return 1
}
