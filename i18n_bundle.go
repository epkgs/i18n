package i18n

import (
	"sync"

	"github.com/epkgs/i18n/errors"
	"golang.org/x/text/language"
)

// i18nBundle represents an internationalization bundle containing translations for different languages
type i18nBundle struct {
	Name  string
	trans map[language.Tag]map[string]string // language identifier -> default text -> translated text

	i18n     *I18n
	loadOnce *sync.Once
}

func newBundle(i18n *I18n, name string) Bundler {
	b := &i18nBundle{
		Name:     name,
		trans:    map[language.Tag]map[string]string{},
		i18n:     i18n,
		loadOnce: &sync.Once{},
	}

	return b
}

// Str creates and returns a new Stringer object for handling internationalized strings
//   - txt: the original text to be translated
//   - args: arguments used to replace placeholders in the text
//
// Returns a Stringer interface that can handle internationalization
func (b *i18nBundle) Str(txt string, args ...any) Stringer {
	return newString(b, txt, args...)
}

// NStr selects singular or plural form of string based on quantity and formats it
//   - isOne: boolean value to determine whether to use singular form
//   - one: singular form string template
//   - others: plural form string template
//   - args: variable arguments for string formatting
//
// Returns: translatable Stringer interface
func (b *i18nBundle) NStr(isOne bool, one, others string, args ...any) Stringer {
	if isOne {
		return b.Str(one, args...)
	}
	return b.Str(others, args...)
}

// Err creates and returns an internationalizable error object
//   - txt: the original error text to be translated
//   - args: arguments used to replace placeholders in the text
//
// Returns an errors.Error interface that includes internationalization capabilities
func (b *i18nBundle) Err(txt string, args ...any) errors.Error {
	return errors.New(b.Str(txt, args...))
}

// NErr selects singular or plural form of string based on quantity and formats it
//   - isOne: boolean value to determine whether to use singular form
//   - one: singular form string template
//   - others: plural form string template
//   - args: variable arguments for string formatting
//
// Returns: translatable errors.Error interface
func (b *i18nBundle) NErr(isOne bool, one, others string, args ...any) errors.Error {
	if isOne {
		return b.Err(one, args...)
	} else {
		return b.Err(others, args...)
	}
}

// translate translates the given format string based on language preferences in the context
func (b *i18nBundle) translate(langs []string, format string, args ...any) string {

	// Initialize a slice to store parsed language tags
	tags := []language.Tag{}
	// Iterate through language codes and attempt to parse them into language tags
	for _, l := range langs {
		// Try to parse the current language code into a language tag. If successful, add it to the tags slice
		if t := parseLanguageTag(l); t != language.Und {
			tags = append(tags, t)
		}
	}

	translated := b.getTranslation(tags, format)

	return parse(translated, args...)
}

// getTranslation retrieves the translated text for the given original text based on language tags
func (b *i18nBundle) getTranslation(tags []language.Tag, key string) string {

	b.lazyLoad()

	if len(b.i18n.matcher.Languages) == 0 {
		return key
	}

	_, i, _ := b.i18n.matcher.Match(tags...)
	lang := b.i18n.matcher.Languages[i]

	trans, exist := b.trans[lang]
	if !exist {
		defaultLanguage := b.i18n.matcher.Languages[0]
		trans, exist = b.trans[defaultLanguage]
		if !exist {
			return key
		}
	}

	if txt, exist := trans[key]; exist {
		return txt
	}

	return key
}

func (b *i18nBundle) lazyLoad() {
	b.loadOnce.Do(func() {
		b.trans = b.i18n.loader(b.Name)
	})
}

func (b *i18nBundle) Reload() {
	b.loadOnce = &sync.Once{} // reset loadOnce
}
