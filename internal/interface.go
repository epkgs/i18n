package internal

import (
	"golang.org/x/text/language"
)

type Loader func(bundleName string, m *Matcher) map[language.Tag]map[string]string
