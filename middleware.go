package i18n

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func GinMiddleware(defaultLangs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		langs := defaultLangs

		langTags, _, err := language.ParseAcceptLanguage(c.Request.Header.Get("Accept-Language"))
		if err == nil {
			strs := []string{}
			for _, tag := range langTags {
				strs = append(strs, tag.String())
			}
			langs = strs
		}
		ctx := WithAcceptLanguages(c.Request.Context(), langs...)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
