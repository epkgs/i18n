package errorx

import (
	"context"
	"errors"

	"github.com/epkgs/i18n"
)

// Definef 函数用于创建一个新的错误定义。
// 它允许开发者通过指定国际化信息和错误包装函数来定义特定的错误类型。
//
// 参数：
//   - i18n是国际化对象，用于处理多语言错误消息。
//   - format是错误消息的格式字符串，用于动态生成错误消息。
//   - wrapper是一个错误包装函数，用于将原始错误包装成新的错误类型。
//
// 返回值：
// 指向Definition结构的指针，该结构包含了国际化错误信息和错误包装函数。
func Definef[Args any, E error](i18n *i18n.I18n, format string, wrapper Wrapper[E]) *Definition[E, Args] {
	// 调用newDefinition函数创建一个新的错误定义。
	// i18n.New(format)生成一个新的国际化错误消息对象。
	// wrapper作为错误定义的一部分，用于后续对错误进行包装。
	return newDefinition[E, Args](i18n.New(format), wrapper)
}

// Define 定义一个简单的错误定义。 New 函数无须填入args
//
// 参数：
//   - i18n是国际化对象，用于处理多语言错误消息。
//   - format是错误消息的格式字符串，用于动态生成错误消息。
//   - wrapper是一个错误包装函数，用于将原始错误包装成新的错误类型。
func Define[E error](i18n *i18n.I18n, format string, wrapper Wrapper[E]) *DefinitionSimple[E] {
	// 调用 Definef 函数创建一个通用的错误定义。
	def := Definef[None](i18n, format, wrapper)
	// 将通用错误定义封装到 DefinitionSimple 结构中并返回。
	return &DefinitionSimple[E]{def}
}

var WrapSelf = func(err *Error) *Error { return err }

type None = struct{}

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

func (d *Definition[E, Args]) New(ctx context.Context, args Args) E {
	err := newError(d.t, ctx, args)
	return d.wrapper(err)
}

func (d *Definition[E, Args]) Is(err error) bool {
	return errors.Is(err, d.base)
}

func (d *Definition[E, Args]) As(target any) bool {
	return errors.As(d.base, target)
}

func (d *Definition[E, Args]) Code() int {
	if code, ok := ErrorCode(d.base); ok {
		return code
	}
	return 1
}

type DefinitionSimple[E error] struct {
	*Definition[E, None]
}

func (d *DefinitionSimple[E]) New(ctx context.Context) E {
	return d.Definition.New(ctx, None{})
}
