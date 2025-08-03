package i18n

import (
	"log"
	"os"
)

type Options struct {
	DefaultLang   string // default language, default is "en"
	ResourcesPath string // resources path, default is "locales"
}

type OptionsFunc func(opts *Options)

type I18n struct {
	opts   Options
	name   string
	items  map[string]*Item
	parser Parser
}

// NewCatalog 创建并返回一个新的I18n实例。
// 它接受一个名称参数和一个可变参数的OptionsFunc函数切片，
// 用于配置I18n实例的选项。
func NewCatalog(name string, fn ...OptionsFunc) *I18n {
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
	return &I18n{
		opts:   opts,
		name:   name,
		items:  make(map[string]*Item),
		parser: new(JsonParser),
	}
}

// New 创建并返回一个新的Item实例，该实例包含默认文本。
//
// 此方法接收一个参数defaultText，作为新创建Item的默认文本。
//
// 新的Item实例会被存储在I18n实例的items映射中，以默认文本为键。
//
// 返回值是新创建的Item实例的指针，允许直接对新实例进行操作。
func (i18n *I18n) New(defaultText string) *Item {
	// 创建一个新的Item实例，并将其存储在items映射中。
	i18n.items[defaultText] = newItem(i18n, defaultText)
	// 返回新创建的Item实例的指针。
	return i18n.items[defaultText]
}

// Name 返回I18n实例的名称。
//
// 该方法没有输入参数。
//
// 返回值:
//
//	string: I18n实例的名称。
func (i18n *I18n) Name() string {
	return i18n.name
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
//	*I18n: 返回I18n实例指针，支持链式调用。
func (i18n *I18n) AddTrans(lang string, defaultText, transText string) *I18n {
	// 检查默认文本是否已存在于翻译项中
	if _, exist := i18n.items[defaultText]; exist {
		// 如果存在，则为该翻译项添加新的语言翻译
		i18n.items[defaultText].AddTrans(lang, transText)
	}

	// 返回I18n实例指针，支持链式调用
	return i18n
}

// LoadTranslations 加载翻译资源
// 此函数读取指定资源路径下的所有目录，每个目录代表一种语言
// 对于每个语言目录，它会解析其中的翻译文件，并将翻译结果添加到I18n实例中
func (i18n *I18n) LoadTranslations() {
	// 读取资源路径下的所有目录
	rd, err := os.ReadDir(i18n.opts.ResourcesPath)
	if err != nil {
		// 目录不存在，直接返回
		return
	}

	// 遍历资源路径下的所有目录和文件
	for _, f := range rd {
		// 只处理目录，目录名为 lang ID
		if f.IsDir() {
			// 解析每个语言目录中的翻译文件
			trans, err := i18n.parser.Parse(i18n.opts.ResourcesPath, f.Name(), i18n.name)
			if err != nil {
				// 如果解析过程中发生错误，记录错误信息并继续处理下一个目录
				log.Println(err)
				continue
			}

			// 文件夹为语言名字
			lang := f.Name()
			// 将解析得到的翻译结果添加到I18n实例中
			for key, value := range trans {
				i18n.AddTrans(lang, key, value)
			}
		}
	}
}

// LoadTranslations 加载翻译资源。
// 该函数接受一个或多个I18n实例作为参数，并调用每个实例的LoadTranslations方法来加载翻译资源。
// 这使得在程序启动时可以预先加载多个语言环境的翻译资源，从而确保在需要时能够快速响应。
// 参数:
//
//	i18n ...*I18n - 一个或多个I18n实例，用于指定需要加载翻译资源的语言环境。
func LoadTranslations(i18n ...*I18n) {
	// 遍历每个传入的I18n实例。
	for _, i := range i18n {
		// 调用当前I18n实例的LoadTranslations方法来加载翻译资源。
		i.LoadTranslations()
	}
}
