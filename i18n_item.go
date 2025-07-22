package i18n

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

type Item struct {
	I18n  *I18n
	trans map[string]string // language -> text

	matcher  language.Matcher
	langTags []language.Tag
}

func newItem(i18n *I18n, defaultText string) *Item {

	item := &Item{
		I18n:  i18n,
		trans: make(map[string]string),
	}

	return item.AddTrans(i18n.opts.DefaultLang, defaultText)
}

func (item *Item) AddTrans(lang string, text string) *Item {

	formattedLang := strings.Replace(lang, "-", "_", -1)

	_, exist := item.trans[formattedLang]

	item.trans[formattedLang] = text

	if !exist {
		// 如果是新添加的语言，则添加 langTag 并重置 matcher
		langTag, err := language.Parse(formattedLang)
		if err == nil {
			item.langTags = append(item.langTags, langTag)
			item.matcher = nil // Reset matcher to nil so it will be recreated on next use
		}
	}

	return item
}

func (item *Item) String() string {
	if txt, exist := item.trans[item.I18n.opts.DefaultLang]; exist {
		return txt
	}
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

func (item *Item) T(ctx context.Context, args ...any) string {
	return item.TLang(GetAcceptLanguages(ctx), args...)
}

func (item *Item) TTag(tags []language.Tag, args ...any) string {
	var lang string

	if tag, exist := item.match(tags...); exist {
		lang = tag.String()
	} else {
		lang = item.I18n.opts.DefaultLang
	}

	key := strings.Replace(lang, "-", "_", -1)

	msg := item.trans[key]

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}

	return msg
}

func (item *Item) TLang(langs []string, args ...any) string {

	tags := []language.Tag{}
	for _, l := range langs {
		if t := parseTag(l); t != nil {
			tags = append(tags, *t)
		}
	}

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
