package errorx

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/epkgs/i18n"
)

// Error 结构体用于封装一个错误的详细信息。
// 它不仅包含原始错误（cause），还可以包括错误发生的上下文（ctx）和相关参数（args），
// 以便于在日志或错误处理中提供更多的错误细节。
type Error struct {
	t *i18n.Item

	ctx  context.Context
	args []any

	cause error // 原始错误，即引发Error结构体封装的最初错误。
}

// newError 创建并返回一个新的Error对象。
// 该函数接收一个Item指针t、一个context.Context类型的ctx，以及一个可变参数args。
// 参数t用于指定错误相关的项，ctx用于传递请求范围的上下文信息，args用于传递额外的错误信息。
// 返回值是一个Error对象，包含了传入的t、ctx和args信息。
func newError(t *i18n.Item, ctx context.Context, args ...any) *Error {
	return &Error{
		t:    t,
		ctx:  ctx,
		args: args,
	}
}

// Error 实现了 error 接口，返回错误的字符串表示。
// 该方法首先检查错误实例 e 是否包含上下文信息(ctx)。
// 如果包含上下文信息，则使用该上下文信息和任何额外的参数(args)来获取本地化错误信息。
// 如果没有上下文信息，则使用一个空白的上下文背景和额外的参数来获取默认的本地化错误信息。
func (e *Error) Error() string {

	// 检查是否存在上下文信息
	if e.ctx != nil {
		// 使用给定的上下文和参数获取本地化错误信息
		return e.t.T(e.ctx, e.args...)
	}

	// 如果没有上下文信息，使用空白上下文背景获取默认本地化错误信息
	return e.t.T(context.Background(), e.args...)
}

// Unwrap 返回错误的底层原因（cause）。
// 此方法允许错误处理机制能够访问Error类型内部封装的实际错误。
// 参数: 无
// 返回值: error，代表错误的底层原因。
func (e *Error) Unwrap() error {
	return e.cause
}

// Wrap 设置Error类型的cause字段
// 该方法用于将一个错误标记为另一个错误的直接原因，便于错误追踪和处理
// 参数:
//
//	err error: 导致当前错误的原始错误，不能为空
func (e *Error) Wrap(err error) error {
	e.cause = err
	return e
}

// Is 检查两个错误是否相等。
// 该方法首先验证传入的错误是否为 *Error 类型，然后比较错误的国际化名称和字符串表示是否完全相同。
// 如果所有比较都相等，则认为两个错误相等。
func (e *Error) Is(err error) bool {
	// 将传入的错误尝试转换为 *Error 类型，并检查转换是否成功
	err2, ok := err.(*Error)
	if !ok {
		// 如果转换不成功，说明类型不匹配，直接返回 false
		return false
	}

	// 比较两个错误的国际化名称是否相同
	if e.t.I18n.Name() != err2.t.I18n.Name() {
		// 如果国际化名称不同，返回 false
		return false
	}

	// 比较两个错误的字符串表示是否相同
	if e.t.String() != err2.t.String() {
		// 如果字符串表示不同，返回 false
		return false
	}

	// 所有比较都相等，返回 true，表示两个错误相等
	return true
}

// As 尝试将错误转换为指定的目标类型。
// 如果错误类型可以转换为由 target 引用的类型，则返回 true。
// 此方法允许在错误链中检查特定类型的错误。
func (e *Error) As(target any) bool {
	// 检查 target 是否为空，为空则无法进行转换。
	if target == nil {
		return false
	}

	// 获取 target 的类型信息，确保它是一个指向结构体的指针。
	targetType := reflect.TypeOf(target)
	// 如果 target 不是一个结构体指针，则返回 false。
	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return false
	}

	// 检查当前错误类型是否可以直接赋值给 target 所指向的类型。
	if reflect.TypeOf(e).AssignableTo(targetType.Elem()) {
		// 如果可以，将当前错误值设置到 target 中。
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(e))
		return true
	}

	// 如果有 cause，递归调用 errors.As
	// 尝试将 cause 转换为 target 类型。
	if e.cause != nil {
		return errors.As(e.cause, target)
	}

	// 如果以上所有检查都未通过，则返回 false。
	return false
}

// Error 类型的 Format 方法用于自定义错误信息的格式化输出。
// 此方法是通过实现 fmt.Formatter 接口来达到格式化错误信息的目的。
// 参数 s 是格式化状态，用于控制输出格式；verb 是格式化动词，决定如何格式化。
func (e *Error) Format(s fmt.State, verb rune) {
	// 根据格式化动词处理不同的格式化情况。
	switch verb {
	case 'v':
		// 当动词为 'v' 且格式化状态 s 的 '+' 标志被设置时，
		// 输出错误的详细信息，包括错误原因（如果有的话）。
		if s.Flag('+') {
			fmt.Fprintf(s, "%s\n", e.Error())
			if e.cause != nil {
				fmt.Fprintf(s, "Cause: %+v\n", e.cause)
			}
			return
		}
		// 如果 '+' 标志未被设置，则继续执行默认的格式化行为。
		fallthrough
	default:
		// 对于其他所有情况，包括动词不是 'v' 或者 '+' 标志未被设置，
		// 只输出错误的基本信息，不包括错误原因。
		fmt.Fprintf(s, "%s", e.Error())
	}
}
