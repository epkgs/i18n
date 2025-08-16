package i18n

import (
	"context"
	"errors"
)

// Errorf 根据上下文中的语言偏好将格式化消息字符串翻译后包装为error返回
// 该函数结合了格式化功能和翻译功能，首先使用提供的参数格式化消息字符串，
// 然后将格式化后的字符串根据语言偏好进行翻译，最后包装为error类型返回
//
// 参数:
//   - ctx context.Context: 上下文对象，应包含语言偏好信息，可通过 WithAcceptLanguages 函数设置
//   - format string: 需要翻译的原始格式化模板字符串
//   - args ...any: 可变参数，用于替换格式化模板中的占位符
//
// 返回值:
//   - error: 包含翻译后格式化消息的error实例
func (b *Bundle) Errorf(ctx context.Context, format string, args ...any) error {
	return errors.New(b.Define(format).T(ctx, args...))
}
