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
func Definef[Args any, E error](bundle *i18n.Bundle, format string, wrapper Wrapper[E]) *DefinitionF[E, Args] {
	// 创建一个errorDefinition类型的对象并对其进行初始化。
	d := &DefinitionF[E, Args]{
		i18n:    bundle.Define(format),
		wrapper: wrapper,
	}

	d.base = d.New(context.Background(), *new(Args))

	// 返回初始化后的errorDefinition对象。
	return d
}

// Define 定义一个简单的错误定义。 New 函数无须填入args
//
// 参数：
//   - bundle是国际化对象，用于处理多语言错误消息。
//   - msg是错误消息的格式字符串，用于动态生成错误消息。
//   - wrapper是一个错误包装函数，用于将原始错误包装成新的错误类型。
func Define[E error](bundle *i18n.Bundle, msg string, wrapper Wrapper[E]) *Definition[E] {
	// 调用 Definef 函数创建一个通用的错误定义。
	def := Definef[None](bundle, msg, wrapper)
	// 将通用错误定义封装到 Definition 结构中并返回。
	return &Definition[E]{def}
}

var WrapSelf = func(err error) error { return err }

type None = struct{}

type Wrapper[E error] func(error) E

// Definition Formatter
type DefinitionF[E error, Args any] struct {
	// bundle *i18n.Bundle // i18n bundle
	// format string

	i18n *i18n.Definition

	base    error
	wrapper Wrapper[E]
}

func (d *DefinitionF[E, Args]) New(ctx context.Context, args Args) E {
	err := errors.New(d.i18n.T(ctx, args))
	return d.wrapper(err)
}

func (d *DefinitionF[E, Args]) Is(err error) bool {
	return errors.Is(err, d.base)
}

func (d *DefinitionF[E, Args]) As(target any) bool {
	return errors.As(d.base, target)
}

func (d *DefinitionF[E, Args]) Code() int {
	if code, ok := Code(d.base); ok {
		return code
	}
	return 1
}

func (d *DefinitionF[E, Args]) HttpStatus() int {
	if httpStatus, ok := HttpStatus(d.base); ok {
		return httpStatus
	}
	return 200
}

type innerDefinition[E error] = DefinitionF[E, None]

// Definition
type Definition[E error] struct {
	*innerDefinition[E] // lowercase to avoid external access
}

func (d *Definition[E]) New(ctx context.Context) E {
	return d.innerDefinition.New(ctx, None{})
}
