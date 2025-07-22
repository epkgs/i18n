package i18n

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

// Item 结构体用于表示具有国际化功能的项目。
// 它包含用于文本本地化的必要信息和工具。
type Item struct {
	I18n  *I18n             // I18n 字段用于存储指向 I18n 实例的指针，可能包含全局的国际化设置或状态。
	trans map[string]string // language -> text，一个映射表，用于将语言标签映射到对应的本地化文本。

	matcher  language.Matcher // Matcher 实例，用于根据用户请求的语言标签和可用的语言资源，找到最佳匹配的语言。
	langTags []language.Tag   // langTags 是一个数组，包含该项目支持的所有语言标签。
}

// newItem 创建并初始化一个新的Item实例，用于处理国际化文本。
//
// 该函数接收一个I18n对象指针和一个默认文本字符串作为参数。
//
// I18n对象用于提供国际化选项和默认语言设置。
//
// 默认文本参数用于设置Item的默认翻译文本。
//
// 函数返回一个指向新创建并初始化的Item对象的指针。
func newItem(i18n *I18n, defaultText string) *Item {
	// 创建一个新的Item实例，并将其I18n字段设置为传入的I18n对象。
	// 初始化一个空的trans映射，用于存储不同语言的翻译文本。
	item := &Item{
		I18n:  i18n,
		trans: make(map[string]string),
	}

	// 使用I18n对象中的默认语言设置和传入的默认文本，
	// 通过AddTrans方法为Item实例添加默认翻译文本，并返回该实例。
	return item.AddTrans(i18n.opts.DefaultLang, defaultText)
}

// AddTrans 为 Item 添加或更新指定语言的翻译文本。
//
// 该方法接受语言代码和对应的文本作为参数，语言代码以连字符或下划线分隔。
//
// 如果是新添加的语言翻译，它还会更新 Item 的语言标签和匹配器。
func (item *Item) AddTrans(lang string, text string) *Item {
	// 将语言代码中的连字符统一替换为下划线，以保持语言代码格式的一致性。
	formattedLang := strings.Replace(lang, "-", "_", -1)

	// 检查 item.trans 中是否已存在该语言的翻译。
	_, exist := item.trans[formattedLang]

	// 无论 exist与否，都添加或更新该语言的翻译文本。
	item.trans[formattedLang] = text

	// 如果是新添加的语言翻译，则进一步处理。
	if !exist {
		// 尝试解析语言代码，如果解析成功，则添加该语言标签并重置匹配器。
		langTag, err := language.Parse(formattedLang)
		if err == nil {
			// 将解析成功的新语言标签添加到 langTags 列表中。
			item.langTags = append(item.langTags, langTag)
			// 重置 matcher 为 nil，以便在下次使用时重新创建，确保语言匹配的准确性。
			item.matcher = nil
		}
	}

	// 返回更新后的 Item 实例。
	return item
}

// Item的String方法返回项的字符串表示，主要用于国际化文本展示。
// 该方法首先尝试根据默认语言获取文本，如果找不到，则返回空字符串。
// 这允许在不牺牲可读性的情况下，为不同的语言环境提供灵活的文本展示。
//
// 参数: 无
// 返回值:
//   - string: 如果找到了默认语言对应的文本，则返回该文本；否则返回空字符串。
func (item *Item) String() string {
	// 尝试从trans映射中获取默认语言对应的文本。
	// trans映射存储了不同语言环境下的文本。
	// opts.DefaultLang指定了默认的语言环境。
	if txt, exist := item.trans[item.I18n.opts.DefaultLang]; exist {
		return txt
	}
	// 如果默认语言的文本不存在，则返回空字符串。
	// 这避免了返回nil或错误处理，简化了调用方的逻辑。
	return ""
}

func (item *Item) match(tags ...language.Tag) (tag language.Tag, exist bool) {
	if len(tags) == 0 {
		return language.Und, false
	}

	_, i, _ := item.getMatcher().Match(tags...)

	if i < 0 || i >= len(item.langTags) {
		return language.Und, false
	}

	return item.langTags[i], true
}

// Item.T 根据请求的上下文和参数，使用适当的语言进行翻译。
// 该方法主要目的是根据给定的上下文获取接受的语言列表，并调用 TLang 方法进行翻译。
// 参数:
//   - ctx context.Context: 请求的上下文，用于提取接受的语言信息。
//   - args ...any: 可变参数列表，传递给 TLang 方法，用于指定翻译的具体内容和其他相关信息。
//
// 返回值:
//   - string: 返回翻译后的字符串。
func (item *Item) T(ctx context.Context, args ...any) string {
	return item.TLang(GetAcceptLanguages(ctx), args...)
}

// TTag 根据提供的标签数组和可选参数，返回相应的语言字符串。
// 该方法主要用于根据用户请求的语言标签来选择合适的语言字符串。
// 参数:
//
//	tags []language.Tag: 一个语言标签数组，表示用户请求的语言偏好。
//	args ...any: 可变参数，用于格式化字符串中的占位符。
//
// 返回值:
//
//	string: 根据语言标签和可选参数格式化后的字符串。
func (item *Item) TTag(tags []language.Tag, args ...any) string {
	var lang string

	// 尝试匹配提供的语言标签，如果找到匹配的标签，则使用它。
	if tag, exist := item.match(tags...); exist {
		lang = tag.String()
	} else {
		// 如果没有找到匹配的语言标签，则使用默认语言。
		lang = item.I18n.opts.DefaultLang
	}

	// 将语言标签中的连字符替换为下划线，以适应某些特定的语言代码格式。
	key := strings.Replace(lang, "-", "_", -1)

	// 根据语言标签获取对应的翻译消息。
	msg := item.trans[key]

	// 如果提供了可选参数，则使用这些参数格式化消息。
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}

	// 如果没有提供可选参数，则直接返回翻译消息。
	return msg
}

// TLang 根据给定的语言切片获取与之匹配的语言标签。
// 该方法首先将输入的语言字符串转换为语言标签，然后调用 TTag 方法
// 以找到与这些标签匹配的翻译内容。
// 参数:
//
//	langs - 一个字符串切片，代表优先级排序的语言代码。
//	args  - 可变参数，可包含额外的信息，如翻译的默认值或其他相关数据。
//
// 返回值:
//
//	一个字符串，代表根据提供的语言偏好和额外参数找到的翻译内容。
func (item *Item) TLang(langs []string, args ...any) string {
	// 初始化一个语言标签切片，用于存储解析后的语言标签。
	tags := []language.Tag{}
	// 遍历输入的语言代码切片，尝试解析每个语言代码为语言标签。
	for _, l := range langs {
		// 尝试解析当前语言代码为语言标签。如果解析成功，则将标签添加到标签切片中。
		if t := parseTag(l); t != nil {
			tags = append(tags, *t)
		}
	}

	// 调用 TTag 方法，使用解析后的语言标签和任何额外参数，以获取翻译内容。
	return item.TTag(tags, args...)
}

func (item *Item) getMatcher() language.Matcher {
	if item.matcher == nil {
		item.matcher = language.NewMatcher(item.langTags)
	}
	return item.matcher
}

var languageTagCache = make(map[string]*language.Tag)

func parseTag(lang string) *language.Tag {
	if _, exist := languageTagCache[lang]; !exist {
		t, e := language.Parse(lang)
		if e != nil {
			languageTagCache[lang] = nil
		} else {
			languageTagCache[lang] = &t
		}
	}

	return languageTagCache[lang]
}
