package i18n

import (
	"context"
)

type (
	acceptLanguagesCtx struct{}
)

// WithAcceptLanguages 返回带有接受语言环境的上下文。
// 该函数主要用于将一个或多个接受的语言代码添加到上下文当中，以便于后续的操作可以访问这些语言偏好。
// 参数:
//
//	ctx: 输入的上下文，通常是一个请求的上下文。
//	acceptLanguages: 一个或多个表示接受的语言代码的字符串。
//
// 返回值:
//
//	返回一个带有接受语言环境的上下文。
func WithAcceptLanguages(ctx context.Context, acceptLanguages ...string) context.Context {
	// 使用context.WithValue将acceptLanguages添加到ctx中。
	// acceptLanguagesCtx{}用作键，这是一种类型安全的做法，避免了键的冲突。
	return context.WithValue(ctx, acceptLanguagesCtx{}, acceptLanguages)
}

// GetAcceptLanguages 从上下文中获取接受的语言列表。
// 该函数主要用于从给定的上下文对象中提取出接受的语言列表。
// 如果上下文中存在接受的语言列表，则将其转换为字符串切片并返回；
// 否则，返回nil，表示没有在上下文中找到接受的语言列表。
//
// 参数:
//
//	ctx context.Context: 上下文对象，用于传递请求范围的数据。
//
// 返回值:
//
//	[]string: 接受的语言列表，如果没有找到，则为nil。
func GetAcceptLanguages(ctx context.Context) []string {
	// 从上下文中获取接受的语言列表。
	v := ctx.Value(acceptLanguagesCtx{})
	// 检查获取的结果是否为空。
	if v != nil {
		// 如果不为空，则断言其为字符串切片并返回。
		return v.([]string)
	}
	// 如果为空，则返回nil。
	return nil
}
