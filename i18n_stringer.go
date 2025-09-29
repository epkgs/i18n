package i18n

import (
	"context"
)

// i18nString represents an internationalizable string structure containing the context and parameters needed for translation
type i18nString struct {
	i18n Bundler // Pointer to the internationalization bundle instance for handling translations in different languages
	txt  string  // Raw text content
	args []any   // List of arguments passed to the text for formatting
}

// newString creates and returns a new i18nString instance
//   - i18n: Bundle instance used for internationalization
//   - txt: Text to be translated
//   - args: Arguments used to replace placeholders in the text
func newString(i18n Bundler, txt string, args ...any) Stringer {
	return &i18nString{
		i18n: i18n,
		txt:  txt,
		args: args,
	}
}

// i18nString implements the fmt.Stringer interface, returning a localized string after parameter replacement
// This method processes the s.txt template string with s.args parameters to generate the final string
// Returns the processed string
func (s *i18nString) String() string {
	return parse(s.txt, s.args...)
}

// T returns the translated version of the current string based on language preferences in the context
// ctx: Context containing language preferences to determine which language to translate into
// Returns the translated string
func (s *i18nString) T(ctx context.Context) string {
	langs := GetAcceptLanguages(ctx)
	return s.i18n.translate(langs, s.txt, s.args...)
}

func (s *i18nString) TL(langs ...string) string {
	return s.i18n.translate(langs, s.txt, s.args...)
}
