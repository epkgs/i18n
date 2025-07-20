package i18n

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// 语言标识符的来源
const (
	headerAcceptLanguage = "Accept-Language"
	queryLang            = "lang"
	cookieLang           = "lang"
)

// Gin middleware
// defaultLangs 默认语言，可多个。当从 query，cookie，accept-language 都没有时，使用默认语言列表
func GinMiddleware(defaultLangs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer c.Next()

		langs := []string{}

		// 获取语言标识
		// 查找顺序：1. URL参数 2. Cookie 3. Accept-Language头 4. 默认语言
		if lang := c.Query(queryLang); lang != "" {
			langs = append(langs, lang)

			// 如果请求是通过URL参数设置语言的，将其保存到cookie
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

		// 如果没有获取到语言或语言不受支持，使用默认语言
		if len(langs) == 0 {
			langs = defaultLangs
		}

		// 将语言设置到上下文
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(WithAcceptLanguages(ctx, langs...))
	}
}

// 解析Accept-Language头并返回最佳匹配的语言
func parseAcceptLanguages(acceptLanguage string) []string {
	langs := []string{}

	if langTags, _, err := language.ParseAcceptLanguage(acceptLanguage); err == nil {
		for _, tag := range langTags {
			langs = append(langs, tag.String())
		}
	}

	return langs
}
