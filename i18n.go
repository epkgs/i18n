package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type I18n struct {
	cfg *Config

	defaultLanguage language.Tag
	limitLanguages  []language.Tag
	loader          Loader

	bundles map[string]Bundler
}

type Config struct {
	DefaultLanguage string
	Languages       []string
}

func newI18n(config ...func(c *Config)) *I18n {
	cfg := &Config{
		DefaultLanguage: "en",
		Languages:       []string{},
	}

	for _, f := range config {
		f(cfg)
	}

	return &I18n{
		cfg:             cfg,
		defaultLanguage: parseLanguageTag(cfg.DefaultLanguage),
		limitLanguages:  parseLanguageTags(cfg.Languages...),
		bundles:         map[string]Bundler{},
	}
}

func (n *I18n) SetDefault(langCode string) bool {
	for _, b := range n.bundles {
		if !b.SetDefault(langCode) {
			return false
		}
	}

	return true
}

func (n *I18n) Bundle(name string) Bundler {

	if b, ok := n.bundles[name]; ok {
		return b
	}

	b := newBundle(name, n.loader)

	n.bundles[name] = b
	return b
}

func (n *I18n) Reload() {
	for _, b := range n.bundles {
		b.Reload()
	}
}

func (n *I18n) generateLoader(filePaths []string, readFile func(file string) ([]byte, error)) Loader {

	return func(bundle string) (*Matcher, map[language.Tag]map[string]string) {

		limit := n.limitLanguages
		matcher := newMatcher(append([]language.Tag{n.defaultLanguage}, limit...))
		trans := map[language.Tag]map[string]string{}

		for _, fpath := range filePaths {

			if info, err := os.Stat(fpath); err != nil {
				continue
			} else {
				if info.IsDir() {
					continue
				}
			}

			dir, filename := filepath.Split(fpath)
			ext := filepath.Ext(filename)
			filebase := filename[:len(filename)-len(ext)]

			var lang string // language name
			var name string // bundle name
			if idx := strings.LastIndexByte(filebase, '.'); idx > 1 {
				lang = filebase[idx+1:]
				name = filebase[:idx]
			} else {
				lang = filepath.Base(dir)
				name = filebase
			}

			if bundle != name {
				continue // skip if bundle name does not match
			}

			tag, err := language.Parse(lang)
			if err != nil {
				continue
			}

			if len(limit) > 0 && !includes(limit, tag) {
				continue
			}

			_, i, _ := matcher.MatchOrAdd(tag)
			tag = matcher.Languages[i]

			var unmarshal func(data []byte, v any) error
			switch ext {
			case ".json":
				unmarshal = json.Unmarshal
			case ".yaml", ".yml":
				unmarshal = yaml.Unmarshal
			case ".toml", ".tml":
				unmarshal = toml.Unmarshal
			case ".ini":
				unmarshal = unmarshalINI
			default:
				continue
			}

			data, err := readFile(fpath)
			if err != nil {
				continue
			}

			keyValues := make(map[string]any)
			if err := unmarshal(data, &keyValues); err != nil {
				continue
			}

			if trans[tag] == nil {
				trans[tag] = make(map[string]string)
			}

			for key, value := range keyValues {
				if str, ok := value.(string); ok {
					trans[tag][key] = str
				}
			}
		}

		return matcher, trans
	}
}
