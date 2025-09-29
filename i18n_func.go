package i18n

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
)

var defaultI18n *I18n

func init() {
	defaultI18n, _ = NewDir("locales")
}

func NewDir(dir string, config ...func(c *Config)) (*I18n, error) {
	return NewGlob(filepath.Join(dir, "*/*"), config...)
}

func NewGlob(pattern string, config ...func(c *Config)) (*I18n, error) {
	return NewFS(os.DirFS("."), pattern, config...)
}

func NewFS(fileSystem fs.FS, pattern string, config ...func(c *Config)) (*I18n, error) {

	n := newI18n(config...)

	assets, err := fs.Glob(fileSystem, pattern)
	if err != nil {
		return nil, err
	}

	n.loader = n.generateLoader(assets, func(file string) ([]byte, error) {
		return fs.ReadFile(fileSystem, file)
	})

	return n, nil
}

func NewKV(langKeyValues map[string]map[string]string, config ...func(c *Config)) (*I18n, error) {
	n := newI18n(config...)

	languages := make([]language.Tag, len(langKeyValues))
	i := 0
	for langCode := range langKeyValues {
		languages[i] = parseLanguageTag(langCode)
		i++
	}

	n.loader = func(bundle string) (*Matcher, map[language.Tag]map[string]string) {

		limit := n.limitLanguages
		matcher := newMatcher(append([]language.Tag{n.defaultLanguage}, languages...))
		trans := map[language.Tag]map[string]string{}
		for lang, kv := range langKeyValues {

			tag, err := language.Parse(lang)
			if err != nil {
				continue
			}
			if len(limit) > 0 && !includes(limit, tag) {
				continue
			}

			trans[tag] = kv
		}

		return matcher, trans
	}

	return n, nil
}

func SetDefaultLanguage(lang string) {
	defaultI18n.SetDefault(lang)
}

func Bundle(name string) Bundler {
	return defaultI18n.Bundle(name)
}

// Reload reloads translation resources for all bundles in the cache.
// It iterates through all bundle instances in the cache and calls their load method
// to reload translation files from the filesystem.
func Reload() {
	defaultI18n.Reload()
}
