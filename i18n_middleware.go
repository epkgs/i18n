package i18n

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// Language identifier sources
const (
	headerAcceptLanguage = "Accept-Language"
	queryLang            = "lang"
	cookieLang           = "lang"
)

// GinMiddleware is a Gin framework middleware for handling internationalization
// defaultLangs: default languages to use when no language is specified in query, cookie, or accept-language header
func GinMiddleware(defaultLangs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer c.Next()

		langs := []string{}

		// Get language identifier
		// Search order: 1. URL parameter 2. Cookie 3. Accept-Language header 4. Default language
		if lang := c.Query(queryLang); lang != "" {
			langs = append(langs, lang)

			// If language is set via URL parameter, save it to cookie
			defer func() {
				c.SetCookie(cookieLang, lang, 0, "/", "", false, true)
			}()
		}

		if lang, _ := c.Cookie(cookieLang); lang != "" {
			langs = append(langs, lang)
		}

		if acceptedLangs := parseAcceptLanguages(c.GetHeader(headerAcceptLanguage)); len(acceptedLangs) > 0 {
			langs = append(langs, acceptedLangs...)
		}

		// If no language was obtained or the language is not supported, use the default language
		if len(langs) == 0 {
			langs = defaultLangs
		}

		// Set language to context
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(WithAcceptLanguages(ctx, langs...))
	}
}

// parseAcceptLanguages parses the Accept-Language header and returns the best matching languages
func parseAcceptLanguages(acceptLanguage string) []string {
	langs := []string{}

	if langTags, _, err := language.ParseAcceptLanguage(acceptLanguage); err == nil {
		for _, tag := range langTags {
			langs = append(langs, tag.String())
		}
	}

	return langs
}
