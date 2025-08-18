package i18n

import (
	"context"

	"golang.org/x/text/language"
)

type Stringer struct {
	bundle *Bundle
	text   string
	args   []any
}

func newStringer(bundle *Bundle, text string, args ...any) *Stringer {
	return &Stringer{
		bundle: bundle,
		text:   text,
		args:   args,
	}
}

func (s *Stringer) String() string {
	return format(s.text, s.args...)
}

func (s *Stringer) Translate(ctx context.Context) string {
	langs := GetAcceptLanguages(ctx)

	// 初始化一个语言标签切片，用于存储解析后的语言标签。
	tags := []language.Tag{}
	// 遍历输入的语言代码切片，尝试解析每个语言代码为语言标签。
	for _, l := range langs {
		// 尝试解析当前语言代码为语言标签。如果解析成功，则将标签添加到标签切片中。
		if t := parseLanguageTag(l); t != nil {
			tags = append(tags, *t)
		}
	}

	translated := s.bundle.getTransTxt(tags, s.text)

	return format(translated, s.args...)
}
