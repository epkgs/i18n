package i18n

import "github.com/epkgs/i18n/internal"

// Bundler is an interface that defines the core functionality of the internationalization package.
type Bundler = internal.Bundler

// Translatable is an interface that provides translation capability
// Implementations of this interface can translate content based on context language preferences
type Translatable = internal.Translatable

// Stringer is an interface that combines fmt.Stringer and Translatable interfaces
// It provides both standard string representation and translation capabilities
type Stringer = internal.Stringer
