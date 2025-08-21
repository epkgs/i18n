package i18n

import (
	"context"
	"reflect"
)

type Definitionf[Arg any] struct {
	i18n  *Bundle
	one   string
	other string
}

// Definef 创建并返回一个新的 Definitionf 实例，用于支持参数化复数形式的本地化文本定义
//
// 参数：
//   - i18n 是 Bundle 类型的指针，提供本地化功能支持
//   - other 是默认的复数形式文本模板
//   - one 是可选的单数形式文本模板，通过可变参数传入，最多使用第一个值
//
// 返回值：
//   - 指定 Arg 泛型参数类型的新 Definitionf 实例指针
func Definef[Arg any](i18n *Bundle, other string, one ...string) *Definitionf[Arg] {
	def := &Definitionf[Arg]{
		i18n:  i18n,
		other: other,
	}

	if len(one) > 0 {
		def.one = one[0]
	}

	return def
}

// T 根据上下文和参数执行本地化文本翻译
//
// 参数：
//   - ctx 上下文，用于获取语言环境等信息
//   - args 泛型参数，用于替换翻译文本中的占位符
//   - num 可选的数字参数，用于决定使用单数还是复数形式。1 或 false ：单数形式，其他情况：复数形式
//
// 返回值：
//   - 翻译后的文本字符串
func (d *Definitionf[Arg]) T(ctx context.Context, args Arg, num ...any) string {
	txt := d.other
	if len(num) > 0 && d.one != "" {
		txt = pluralize(num[0], d.other, d.one)
	}
	return d.i18n.translate(ctx, txt, args)
}

type Definition struct {
	i18n *Bundle
	txt  string
}

func Define(i18n *Bundle, txt string) *Definition {
	return &Definition{
		i18n: i18n,
		txt:  txt,
	}
}

func (d *Definition) T(ctx context.Context) string {
	return d.i18n.translate(ctx, d.txt)
}

func (d *Definition) String() string {
	return d.i18n.translate(context.Background(), d.txt)
}

func pluralize(n any, plural, singular string) string {

	v := reflect.ValueOf(n)
	switch v.Kind() {
	case reflect.Bool:
		if !v.Bool() {
			return singular
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 1 {
			return singular
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v.Uint() == 1 {
			return singular
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 1 {
			return singular
		}
	}

	return plural
}
