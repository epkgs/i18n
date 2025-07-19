package i18n

import (
	"log"
	"os"
	"path"
)

type Options struct {
	DefaultLang string // default language, default is "en"
	LocalesPath string // locales path, default is "locales"
}

type OptionsFunc func(opts *Options)

type I18n struct {
	opts   Options
	name   string
	items  map[string]*Item
	parser Parser
}

func New(name string, fn ...OptionsFunc) *I18n {

	opts := Options{
		DefaultLang: "en",
		LocalesPath: "locales",
	}

	for _, f := range fn {
		f(&opts)
	}

	return &I18n{
		opts:   opts,
		name:   name,
		items:  make(map[string]*Item),
		parser: new(JsonParser),
	}
}

func (i18n *I18n) NewItem(defaultText string) *Item {
	i18n.items[defaultText] = newItem(i18n.opts.DefaultLang, defaultText)
	return i18n.items[defaultText]
}

func (i18n *I18n) AddTrans(lang string, defaultText, transText string) *I18n {
	if _, exist := i18n.items[defaultText]; exist {
		i18n.items[defaultText].AddTrans(lang, transText)
	}

	return i18n
}

func (i18n *I18n) LoadLocales() {
	rd, err := os.ReadDir(i18n.opts.LocalesPath)
	if err != nil {
		log.Printf("read locales path error: %v", err)
		// 读取出错则直接返回
		return
	}

	for _, f := range rd {
		if f.IsDir() {

			trans, err := i18n.parser.Parse(path.Join(i18n.opts.LocalesPath, f.Name()), i18n.name)
			if err != nil {
				log.Println(err)
				continue
			}

			// 文件夹为语言名字
			lang := f.Name()
			for key, value := range trans {
				i18n.AddTrans(lang, key, value)
			}
		}
	}
}

func LoadLocales(i18n ...*I18n) {
	for _, i := range i18n {
		i.LoadLocales()
	}
}
