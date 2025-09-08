package internal

import (
	"strings"

	"golang.org/x/text/language"
)

// langIDCache caches formatted language identifiers
var langIDCache = map[string]string{}

// formatLangID formats a language identifier by replacing hyphens with underscores
// and caches the result for performance
func formatLangID(lang string) string {
	id, ok := langIDCache[lang]
	if ok {
		return id
	}

	tag, err := language.Parse(lang)
	if err == nil {
		id = tag.String()
	} else {
		id = strings.Replace(lang, "_", "-", -1)
	}

	langIDCache[lang] = id
	return id
}
