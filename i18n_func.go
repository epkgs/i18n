package i18n

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/epkgs/i18n/internal"
	"github.com/epkgs/i18n/types"
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
		languages[i] = internal.ParseLanguageTag(langCode)
		i++
	}

	n.loader = func(bundleName string, m *internal.Matcher) map[language.Tag]map[string]string {

		limit := n.limitLanguages
		trans := map[language.Tag]map[string]string{}
		for lang, kv := range langKeyValues {

			tag, err := language.Parse(lang)
			if err != nil {
				continue
			}
			if len(limit) > 0 && !internal.Includes(limit, tag) {
				continue
			}

			trans[tag] = kv
		}

		return trans
	}

	return n, nil
}

func SetDefaultLanguage(lang string) {
	defaultI18n.SetDefault(lang)
}

func Bundle(name string) types.Bundler {
	return defaultI18n.Bundle(name)
}

// Reload reloads translation resources for all bundles in the cache.
// It iterates through all bundle instances in the cache and calls their load method
// to reload translation files from the filesystem.
func Reload() {
	defaultI18n.Reload()
}
