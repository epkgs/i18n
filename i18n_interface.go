package i18n

import (
	"context"
	"fmt"
)

// Translatable is an interface that provides translation capability
// Implementations of this interface can translate content based on context language preferences
type Translatable interface {
	// T returns the translated string based on the language preferences in the context
	T(ctx context.Context) string
}

// Stringer is an interface that combines fmt.Stringer and Translatable interfaces
// It provides both standard string representation and translation capabilities
type Stringer interface {
	fmt.Stringer
	Translatable
}
