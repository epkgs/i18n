package internal

import (
	"context"
	"fmt"

	"github.com/epkgs/i18n/errors"
	"golang.org/x/text/language"
)

type Loader func(bundleName string, m *Matcher) map[language.Tag]map[string]string

// Bundler is an interface that defines the core functionality of the internationalization package.
type Bundler interface {
	// transCtx Translates a formatted string according to the language preference in the context.
	// ctx: Context containing the language preference.
	// format: The formatted string to translate.
	// args: Arguments passed to the formatted string.
	transCtx(ctx context.Context, format string, args ...any) string

	// transLangs Translates a formatted string according to the specified list of languages.
	// langs: List of language tags, sorted by priority.
	// format: The formatted string to translate.
	// args: Arguments passed to the formatted string.
	transLangs(langs []string, format string, args ...any) string

	// Str Returns a translatable string instance containing the given text.
	// text: The text to translate.
	// args: Arguments passed to the formatted string.
	Str(text string, args ...any) Stringer

	// Err Returns a translatable error instance containing the given text.
	// text: Error message text
	// args: Arguments passed to the error message
	Err(text string, args ...any) errors.Error

	// SetDefaultLanguage Sets the default language
	// lang: Language tag to set as default
	// Returns whether the setting was successful
	SetDefaultLanguage(lang language.Tag) bool

	// Reload Reloads the language pack data
	Reload()
}

// Translatable is an interface that provides translation capability
// Implementations of this interface can translate content based on context language preferences
type Translatable interface {
	// T returns the translated string based on the language preferences in the context
	T(ctx context.Context) string

	// TL returns the translated string based on the specified language preferences
	TL(langs ...string) string
}

// Stringer is an interface that combines fmt.Stringer and Translatable interfaces
// It provides both standard string representation and translation capabilities
type Stringer interface {
	fmt.Stringer
	Translatable
}
