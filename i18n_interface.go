package i18n

import (
	"context"
	"fmt"

	"github.com/epkgs/i18n/errors"
	"golang.org/x/text/language"
)

type Bundler interface {
	translate(langs []string, format string, args ...any) string

	// Str returns a stringer instance with the given text
	Str(text string, args ...any) Stringer

	// Err returns an error instance with the given text
	Err(text string, args ...any) errors.Error

	Reload()
}

type Loader func(bundle string) map[language.Tag]map[string]string

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

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}
