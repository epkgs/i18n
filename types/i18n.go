package types

import (
	"context"
	"fmt"

	"golang.org/x/text/language"
)

// Bundler is an interface that defines the core functionality of the internationalization package.
type Bundler interface {
	// Str Returns a translatable string instance containing the given text.
	// text: The text to translate.
	// args: Arguments passed to the formatted string.
	Str(text string, args ...any) Stringer

	// Err Returns a translatable error instance containing the given text.
	// text: Error message text
	// args: Arguments passed to the error message
	Err(text string, args ...any) Error

	// SetDefaultLanguage Sets the default language
	// lang: Language tag to set as default
	// Returns whether the setting was successful
	SetDefaultLanguage(lang language.Tag) bool

	// Reload Reloads the language pack data
	Reload()
}

// Translator is an interface that provides translation capability
// Implementations of this interface can translate content based on context language preferences
type Translator interface {
	// T returns the translated string based on the language preferences in the context
	T(ctx context.Context) string

	// TL returns the translated string based on the specified language preferences
	TL(langs ...string) string
}

// Stringer is an interface that combines fmt.Stringer and Translator interfaces
// It provides both standard string representation and translation capabilities
type Stringer interface {
	fmt.Stringer
	Translator
}
