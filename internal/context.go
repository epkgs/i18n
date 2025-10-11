package internal

import (
	"context"
)

type (
	acceptLanguagesCtx struct{}
)

// WithAcceptLanguages returns a context with accepted languages.
// This function is mainly used to add one or more accepted language codes to the context,
// so that subsequent operations can access these language preferences.
// Parameters:
//
//   - ctx: The input context, usually a request context.
//   - acceptLanguages: One or more strings representing accepted language codes.
//
// Return value:
//   - Returns a context with accepted languages.
func WithAcceptLanguages(ctx context.Context, acceptLanguages ...string) context.Context {
	// Use context.WithValue to add acceptLanguages to ctx.
	// acceptLanguagesCtx{} is used as the key, which is a type-safe approach to avoid key conflicts.
	return context.WithValue(ctx, acceptLanguagesCtx{}, acceptLanguages)
}

// GetAcceptLanguages retrieves the list of accepted languages from the context.
// This function is mainly used to extract the list of accepted languages from the given context object.
// If the list of accepted languages exists in the context, it is converted to a string slice and returned;
// otherwise, nil is returned, indicating that no accepted languages list was found in the context.
//
// Parameters:
//
//   - ctx context.Context: The context object used to pass request-scoped data.
//
// Returns:
//   - []string: The list of accepted languages, or nil if not found.
func GetAcceptLanguages(ctx context.Context) []string {
	// Retrieve the list of accepted languages from the context.
	v := ctx.Value(acceptLanguagesCtx{})
	// Check if the retrieved result is nil.
	if v != nil {
		// If not nil, assert it as a string slice and return.
		return v.([]string)
	}
	// If nil, return nil.
	return nil
}
