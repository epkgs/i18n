package i18n

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
)

type Options struct {
	DefaultLang   string // default language, default is "en"
	ResourcesPath string // resources path, default is "locales"
}

type OptionsFunc func(opts *Options)

type Bundle struct {
	opts      Options
	name      string
	parser    Parser
	trans     map[string]map[string]string // 语言标识符 -> 默认文本 -> 翻译文本
	tranLangs []language.Tag
	matcher   language.Matcher
}

// NewBundle 创建并返回一个新的I18n实例。
// 它接受一个名称参数和一个可变参数的OptionsFunc函数切片，
// 用于配置I18n实例的选项。
func NewBundle(name string, fn ...OptionsFunc) *Bundle {
	// 初始化默认的选项配置。
	opts := Options{
		DefaultLang:   "en",
		ResourcesPath: "locales",
	}

	// 遍历可变参数中的函数，应用到选项配置上。
	for _, f := range fn {
		f(&opts)
	}

	// 创建并返回新的I18n实例。
	return &Bundle{
		opts:   opts,
		name:   name,
		parser: new(JsonParser),
		trans: map[string]map[string]string{
			formatLangID(opts.DefaultLang): make(map[string]string),
		},
		tranLangs: []language.Tag{
			language.MustParse(opts.DefaultLang),
		},
	}
}

func (b *Bundle) Sprintf(format string, args ...any) *Stringer {
	return newStringer(b, format, args...)
}

func (b *Bundle) getTransTxt(tags []language.Tag, orig string) string {
	txt := orig

	// 匹配语言
	var lang string
	if tag, exist := b.match(tags...); exist {
		lang = tag.String()
	} else {
		lang = b.opts.DefaultLang
	}

	// 格式化语言键
	lang = formatLangID(lang)

	trans, exist := b.trans[lang]
	if !exist {
		trans, exist = b.trans[b.opts.DefaultLang]
		if !exist {
			return orig
		}
	}

	if txt, exist = trans[txt]; !exist {
		return orig
	}

	return txt
}

func (b *Bundle) match(tags ...language.Tag) (tag language.Tag, exist bool) {
	if len(tags) == 0 {
		return language.Und, false
	}

	_, i, _ := b.getMatcher().Match(tags...)

	if i < 0 || i >= len(b.tranLangs) {
		return language.Und, false
	}

	return b.tranLangs[i], true
}

func (b *Bundle) getMatcher() language.Matcher {
	if b.matcher == nil {
		b.matcher = language.NewMatcher(b.tranLangs)
	}
	return b.matcher
}

// Name 返回I18n实例的名称。
//
// 该方法没有输入参数。
//
// 返回值:
//
//	string: I18n实例的名称。
func (b *Bundle) Name() string {
	return b.name
}

// AddTrans 添加或更新指定语言的翻译文本。
// 如果默认文本(defaultText)在翻译项中已存在，则为其添加新的语言翻译；
// 否则，将创建一个新的翻译项并添加该语言的翻译。
// 参数:
//
//	lang: 语言标识符，用于指定翻译所对应的语言。
//	defaultText: 默认文本，作为翻译项的唯一标识。
//	transText: 翻译文本，是defaultText在指定语言下的翻译。
//
// 返回值:
//
//	*Bundle: 返回I18n实例指针，支持链式调用。
func (b *Bundle) AddTrans(lang string, defaultText, transText string) *Bundle {

	// 格式化 lang
	lang = formatLangID(lang)
	// 检查该语言是否已初始化
	if _, exist := b.trans[lang]; !exist {

		if langTag, err := language.Parse(lang); err == nil { // 尝试解析语言代码，如果解析成功，则添加该语言标签并重置匹配器。
			b.tranLangs = append(b.tranLangs, langTag)
			b.matcher = nil
		}

		b.trans[lang] = make(map[string]string)
	}

	b.trans[lang][defaultText] = transText

	// 返回I18n实例指针，支持链式调用
	return b
}

// Load 加载翻译资源
// 此函数读取指定资源路径下的所有目录，每个目录代表一种语言
// 对于每个语言目录，它会解析其中的翻译文件，并将翻译结果添加到I18n实例中
func (b *Bundle) Load() {
	// 读取资源路径下的所有目录
	rd, err := os.ReadDir(b.opts.ResourcesPath)
	if err != nil {
		// 目录不存在，直接返回
		return
	}

	// 遍历资源路径下的所有目录和文件
	for _, f := range rd {
		// 只处理目录，目录名为 lang ID
		if f.IsDir() {

			folder := f.Name()

			// 解析每个语言目录中的翻译文件
			trans, err := b.parser.Parse(filepath.Join(b.opts.ResourcesPath, folder), b.name)
			if err != nil {
				// 如果解析过程中发生错误，记录错误信息并继续处理下一个目录
				log.Println(err)
				continue
			}

			for key, value := range trans {
				// 文件夹为语言名字
				b.AddTrans(folder, key, value)
			}
		}
	}
}

// Load 加载翻译资源。
// 该函数接受一个或多个I18n实例作为参数，并调用每个实例的LoadTranslations方法来加载翻译资源。
// 这使得在程序启动时可以预先加载多个语言环境的翻译资源，从而确保在需要时能够快速响应。
// 参数:
//
//	bundle ...*Bundle - 一个或多个I18n实例，用于指定需要加载翻译资源的语言环境。
func Load(bundle ...*Bundle) {
	// 遍历每个传入的I18n实例。
	for _, b := range bundle {
		// 调用当前I18n实例的LoadTranslations方法来加载翻译资源。
		b.Load()
	}
}
