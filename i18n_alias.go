package i18n

import "github.com/epkgs/i18n/internal"

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
var WithAcceptLanguages = internal.WithAcceptLanguages

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
var GetAcceptLanguages = internal.GetAcceptLanguages
