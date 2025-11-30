package internal

import (
	"context"
	"sync"

	"github.com/epkgs/i18n/errors"
	"github.com/epkgs/i18n/types"
	"golang.org/x/text/language"
)

// i18nBundle represents an internationalization bundle containing translations for different languages
type i18nBundle struct {
	Name  string
	trans map[language.Tag]map[string]string // language identifier -> default text -> translated text

	matcher  *Matcher
	load     Loader
	loadOnce *sync.Once
}

func NewBundle(name string, matcher *Matcher, loader Loader) types.Bundler {
	b := &i18nBundle{
		Name:     name,
		trans:    map[language.Tag]map[string]string{},
		matcher:  matcher,
		load:     loader,
		loadOnce: &sync.Once{},
	}

	return b
}

// Str creates and returns a new Stringer object for handling internationalized strings
//   - txt: the original text to be translated
//   - args: arguments used to replace placeholders in the text
//
// Returns a Stringer interface that can handle internationalization
func (b *i18nBundle) Str(txt string, args ...any) types.Stringer {
	return NewString(b, txt, args...)
}

// NStr selects singular or plural form of string based on quantity and formats it
//   - n: quantity value to determine singular/plural form. Accepts numeric types (int, float, etc.) and boolean.
//     For numeric types: singular form is used when value equals 1
//     For boolean: singular form is used when value is true
//     Other types default to plural form
//   - one: singular form string template with placeholders
//   - others: plural form string template with placeholders
//   - args: variable arguments for string formatting, replacing placeholders in templates
//
// Returns: internationalized Stringer interface based on quantity
func (b *i18nBundle) NStr(n any, one, others string, args ...any) types.Stringer {

	var isOne bool
	switch t := n.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		isOne = t == 1
	case bool:
		isOne = t
	default:
		isOne = false
	}

	if isOne {
		return b.Str(one, args...)
	}
	return b.Str(others, args...)
}

// Err creates and returns an internationalizable error object
//   - txt: the original error text to be translated
//   - args: arguments used to replace placeholders in the text
//
// Returns an types.Error interface that includes internationalization capabilities
func (b *i18nBundle) Err(txt string, args ...any) types.Error {
	return errors.New(b.Str(txt, args...))
}

// NErr creates an internationalized error based on quantity, selecting singular or plural form
//   - n: quantity value to determine singular/plural form. Accepts numeric types (int, float, etc.) and boolean.
//     For numeric types: singular form is used when value equals 1
//     For boolean: singular form is used when value is true
//     Other types default to plural form
//   - one: singular form error message template with placeholders
//   - others: plural form error message template with placeholders
//   - args: variable arguments for string formatting, replacing placeholders in templates
//
// Returns: internationalized Error interface based on quantity
func (b *i18nBundle) NErr(n any, one, others string, args ...any) types.Error {
	return errors.New(b.NStr(n, one, others, args...))
}

func (b *i18nBundle) SetDefaultLanguage(t language.Tag) bool {

	languages := b.matcher.Languages()
	idx := IndexOf(languages, t)

	if idx == 0 {
		return true
	}

	if idx == -1 {
		languages = append([]language.Tag{t}, languages...)
	} else {
		old := languages[0]
		languages[0] = t
		languages[idx] = old
	}

	b.matcher.SetLanguages(languages)

	return true
}

func (b *i18nBundle) transCtx(ctx context.Context, format string, args ...any) string {
	langs := GetAcceptLanguages(ctx)
	return b.transLangs(langs, format, args...)
}

// translate translates the given format string based on language preferences in the context
func (b *i18nBundle) transLangs(langs []string, format string, args ...any) string {

	// Initialize a slice to store parsed language tags
	tags := []language.Tag{}
	// Iterate through language codes and attempt to parse them into language tags
	for _, l := range langs {
		// Try to parse the current language code into a language tag. If successful, add it to the tags slice
		if t := ParseLanguageTag(l); t != language.Und {
			tags = append(tags, t)
		}
	}

	translated := b.getTranslation(tags, format)

	return Parse(translated, args...)
}

// getTranslation retrieves the translated text for the given original text based on language tags
func (b *i18nBundle) getTranslation(tags []language.Tag, key string) string {

	b.lazyLoad()

	lang := b.matcher.Match(tags...)

	trans, exist := b.trans[lang]
	if !exist {
		defaultLanguage := b.matcher.DefaultLanguage()
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
		b.trans = b.load(b.Name, b.matcher)
	})
}

func (b *i18nBundle) Reload() {
	b.loadOnce = &sync.Once{} // reset loadOnce
}
