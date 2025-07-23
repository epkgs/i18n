package i18n

import (
	"context"
)

// DefineError 定义一个新的错误。
// 该方法接受一个格式化字符串和一个或多个错误包装器作为参数，用于构建国际化错误消息。
// 参数:
//
//	format - 用于错误消息的格式化字符串。
//	wrappers - 可变参数，代表一系列的错误包装器，用于包装错误消息。
//
// 返回值:
//
//	返回一个新的错误定义对象，用于进一步的错误处理和国际化支持。
func (i18n *I18n) DefineError(format string, wrappers ...ErrorWrapFunc) *errorDefinition {
	return newErrorDefinition(i18n.New(format), wrappers...)
}

type ErrorWrapFunc func(error) error

type errorDefinition struct {
	t        *Item // i18n item
	base     error
	wrappers []ErrorWrapFunc
}

// newErrorDefinition 创建并初始化一个新的错误定义对象。
// 该函数接收一个Item类型的参数t，以及一个可变长参数wrappers，后者是由ErrorWrapper接口类型的对象组成。
// Item类型和ErrorWrapper接口的具体定义未在上下文中给出，因此这里不做具体说明。
// 函数返回一个指向errorDefinition类型的指针。
func newErrorDefinition(t *Item, wrappers ...ErrorWrapFunc) *errorDefinition {
	// 创建一个errorDefinition类型的对象并对其进行初始化。
	def := &errorDefinition{
		t:        t,
		wrappers: wrappers,
	}

	// 调用errorDefinition对象的New方法，传入一个空的上下文，以初始化其base字段。
	// 这里使用context.Background()作为参数，是因为newErrorDefinition函数没有上下文信息传入。
	def.base = def.New(context.Background())
	// 返回初始化后的errorDefinition对象。
	return def
}

// Base 返回错误定义的基础错误。
// 该方法允许访问错误定义内部的基础错误，以便在需要时进行进一步处理或检查。
func (d *errorDefinition) Base() error {
	return d.base
}

// New 创建并返回一个新的错误对象，该对象基于当前错误定义，并根据上下文和参数进行定制。
// 此函数允许在给定的上下文和参数下，对错误进行包装和处理，以便在不同的场景下提供更丰富的错误信息。
func (d *errorDefinition) New(ctx context.Context, args ...any) ErrorWrapper {
	// 创建初始错误对象，基于当前错误定义的类型、上下文和可变参数。
	var err error = newError(d.t, ctx, args...)
	// 遍历当前错误定义的所有错误包装器，对初始错误对象进行逐层包装。
	for _, wrapper := range d.wrappers {
		err = wrapper(err)
	}

	return err.(ErrorWrapper)
}

// With为errorDefinition添加错误包装器，以定制错误处理行为。
// 该方法接收一个或多个ErrorWrapper接口实现，将它们附加到errorDefinition的wrappers列表中，
// 并依次用这些包装器包装base错误，以实现错误的层次化管理和处理。
func (d *errorDefinition) With(wrappers ...ErrorWrapFunc) *errorDefinition {
	// 将新的包装器附加到已有的包装器列表中，以便后续处理。
	d.wrappers = append(d.wrappers, wrappers...)

	// 遍历所有包装器，逐个包装base错误。
	// 这一步是错误定制的核心，通过每个包装器的处理逻辑为base错误添加额外的上下文或修改其行为。
	for _, w := range wrappers {
		d.base = w(d.base)
	}

	// 返回修改后的errorDefinition，以支持链式调用和进一步的错误定制。
	return d
}

// Code 返回错误定义的代码标识。
// 如果基础错误实现了 Code() int 方法，则调用该方法获取错误代码。
// 否则，返回默认错误代码 1。
// 此方法用于在不同的错误处理场景中提供一致的错误代码识别。
func (d *errorDefinition) Code() int {
	// 检查基础错误是否具有 Code() int 方法。
	if err, ok := d.base.(interface{ Code() int }); ok {
		// 如果基础错误实现了 Code 方法，则调用并返回该方法的结果。
		return err.Code()
	}
	// 如果基础错误没有实现 Code 方法，则返回默认错误代码 1。
	return 1
}
