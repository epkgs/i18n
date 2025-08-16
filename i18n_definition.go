package i18n

import (
	"context"
	"fmt"
	"reflect"

	"golang.org/x/text/language"
)

type Definition struct {
	bundle *Bundle
	format string
}

func newDefinition(bundle *Bundle, format string) *Definition {
	return &Definition{
		bundle: bundle,
		format: format,
	}
}

// T 根据上下文中的语言偏好翻译格式化字符串。
// 该函数首先从上下文中获取接受的语言列表，然后使用这些语言来查找最适合的翻译文本，
// 最后根据提供的参数格式化翻译后的文本。
//
// 参数:
//
//	ctx context.Context: 上下文对象，应包含语言偏好信息，可通过 WithAcceptLanguages 函数设置
//	args ...any: 可选的参数，用于格式化翻译后的文本
//
// 返回值:
//
//	string: 翻译并格式化后的文本
func (d *Definition) T(ctx context.Context, args ...any) string {
	return d.TLang(GetAcceptLanguages(ctx), args...)
}

// TLang 根据提供的语言列表翻译指定格式的字符串
// 它会依次尝试使用列表中的语言进行翻译，直到找到匹配的翻译内容
//
// 参数:
//   - langs []string: 语言代码列表，例如 []string{"zh-CN", "en-US"}
//   - args ...any: 可变参数，用于替换格式字符串中的占位符
//
// 返回值:
//   - string: 翻译后的字符串，如果找不到匹配的翻译则可能返回默认语言的翻译或原格式字符串
func (d *Definition) TLang(langs []string, args ...any) string {
	// 初始化一个语言标签切片，用于存储解析后的语言标签。
	tags := []language.Tag{}
	// 遍历输入的语言代码切片，尝试解析每个语言代码为语言标签。
	for _, l := range langs {
		// 尝试解析当前语言代码为语言标签。如果解析成功，则将标签添加到标签切片中。
		if t := parseLanguageTag(l); t != nil {
			tags = append(tags, *t)
		}
	}

	format := d.bundle.getTransTxt(tags, d.format)

	// 调用 translate 方法，使用解析后的语言标签和任何额外参数，以获取翻译内容。
	return d.translate(format, args...)
}

func (d *Definition) translate(format string, args ...any) string {

	// 无参数直接返回原始文本
	if len(args) == 0 {
		return format
	}

	if len(args) == 1 {

		arg1 := args[0]

		if arg1 == nil {
			return format
		}

		v := reflect.ValueOf(arg1)
		switch v.Kind() {
		case reflect.Struct:
			// 结构体为零值 或 空结构体（无字段），避免模板渲染失败
			if v.IsZero() || v.NumField() == 0 {
				return format
			}

			return parseTemplate(format, arg1)

		case reflect.Map:
			if v.Len() == 0 {
				return format
			}

			return parseTemplate(format, arg1)
		case reflect.Array, reflect.Slice:
			return fmt.Sprintf(format, (arg1.([]any))...)
		}
	}

	// 否则使用 fmt.Sprintf 处理顺序参数
	return fmt.Sprintf(format, args...)
}

type innerDefinition = Definition

type DefinitionF[Args any] struct {
	*innerDefinition
}

func newDefinitionF[Args any](bundle *Bundle, format string) *DefinitionF[Args] {
	return &DefinitionF[Args]{
		innerDefinition: newDefinition(bundle, format),
	}
}

func (d *DefinitionF[Args]) T(ctx context.Context, args Args) string {
	return d.innerDefinition.T(ctx, args)
}

func (d *DefinitionF[Args]) TLang(langs []string, args Args) string {
	return d.innerDefinition.TLang(langs, args)
}
