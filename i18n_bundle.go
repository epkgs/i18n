package i18n

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/epkgs/i18n/errors"
	"golang.org/x/text/language"
)

// Options represents configuration options for the internationalization bundle
type Options struct {
	DefaultLang   string // default language, default is "en"
	ResourcesPath string // resources path, default is "locales"
}

// OptionsFunc is a function type used to configure Options
type OptionsFunc func(opts *Options)

// Bundle represents an internationalization bundle containing translations for different languages
type Bundle struct {
	opts      Options
	name      string
	trans     map[string]map[string]string // language identifier -> default text -> translated text
	tranLangs []language.Tag
	matcher   language.Matcher
}

var bundleCache = sync.Map{}

// NewBundle creates and returns a new internationalization bundle instance
// It accepts a name parameter and variadic OptionsFunc functions to configure the bundle options
func NewBundle(name string, fn ...OptionsFunc) *Bundle {
	// Initialize default options
	opts := Options{
		DefaultLang:   "en",
		ResourcesPath: "locales",
	}

	// Apply option functions to the options
	for _, f := range fn {
		f(&opts)
	}

	key := opts.ResourcesPath + "." + name

	if v, ok := bundleCache.Load(key); ok {
		return v.(*Bundle)
	}

	// Create and return new bundle instance
	b := &Bundle{
		opts: opts,
		name: name,
		trans: map[string]map[string]string{
			formatLangID(opts.DefaultLang): make(map[string]string),
		},
		tranLangs: []language.Tag{
			language.MustParse(opts.DefaultLang),
		},
	}

	// auto load trans resources
	b.load()

	bundleCache.Store(b.opts.ResourcesPath+"."+b.name, b)

	return b
}

// Str creates and returns a new Stringer object for handling internationalized strings
//   - txt: the original text to be translated
//   - args: arguments used to replace placeholders in the text
//
// Returns a Stringer interface that can handle internationalization
func (b *Bundle) Str(txt string, args ...any) Stringer {
	return newString(b, txt, args...)
}

// Err creates and returns an internationalizable error object
//   - txt: the original error text to be translated
//   - args: arguments used to replace placeholders in the text
//
// Returns an errors.Error interface that includes internationalization capabilities
func (b *Bundle) Err(txt string, args ...any) errors.Error {
	return errors.New(b.Str(txt, args...))
}

// translate translates the given format string based on language preferences in the context
func (b *Bundle) translate(ctx context.Context, format string, args ...any) string {
	langs := GetAcceptLanguages(ctx)

	// Initialize a slice to store parsed language tags
	tags := []language.Tag{}
	// Iterate through language codes and attempt to parse them into language tags
	for _, l := range langs {
		// Try to parse the current language code into a language tag. If successful, add it to the tags slice
		if t := parseLanguageTag(l); t != nil {
			tags = append(tags, *t)
		}
	}

	translated := b.getTransTxt(tags, format)

	return parse(translated, args...)
}

// getTransTxt retrieves the translated text for the given original text based on language tags
func (b *Bundle) getTransTxt(tags []language.Tag, key string) string {

	// Match language
	var lang string
	if tag, exist := b.match(tags...); exist {
		lang = tag.String()
	} else {
		lang = b.opts.DefaultLang
	}

	// Format language key
	lang = formatLangID(lang)

	trans, exist := b.trans[lang]
	if !exist {
		trans, exist = b.trans[b.opts.DefaultLang]
		if !exist {
			return key
		}
	}

	if txt, exist := trans[key]; exist {
		return txt
	}

	return key
}

// match finds the best matching language tag from the provided tags
func (b *Bundle) match(tags ...language.Tag) (tag language.Tag, exist bool) {
	if len(tags) == 0 {
		return language.Und, false
	}

	_, i, _ := b.getMatcher().Match(tags...)

	if i < 0 || i >= len(b.tranLangs) {
		return language.Und, false
	}

	return b.tranLangs[i], true
}

// getMatcher creates and returns a language matcher if not already created
func (b *Bundle) getMatcher() language.Matcher {
	if b.matcher == nil {
		b.matcher = language.NewMatcher(b.tranLangs)
	}
	return b.matcher
}

// Name returns the name of the bundle instance
func (b *Bundle) Name() string {
	return b.name
}

// WithTranslations manually adds translation texts for a specified language
//   - lang: language identifier for the translation
//   - defaultText: default text as a unique identifier for the translation entry
//   - transText: translated text for the defaultText in the specified language
//
// Returns the bundle instance pointer to support method chaining
func (b *Bundle) WithTranslations(lang string, trans map[string]string) *Bundle {

	if len(trans) == 0 {
		return b
	}

	// Format language
	lang = formatLangID(lang)
	// Check if the language has been initialized
	if _, exist := b.trans[lang]; !exist {

		if langTag, err := language.Parse(lang); err == nil { // Try to parse the language code, if successful, add the language tag and reset the matcher
			b.tranLangs = append(b.tranLangs, langTag)
			b.matcher = nil
		}

		b.trans[lang] = make(map[string]string)
	}

	for k, v := range trans {
		b.trans[lang][k] = v
	}

	// Return bundle instance pointer to support method chaining
	return b
}

// Reload reloads all translation resources from the filesystem.
// It clears the existing translations and language tags, then loads
// all translations again by calling the load method.
func (b *Bundle) Reload() {
	b.trans = make(map[string]map[string]string)
	b.tranLangs = []language.Tag{}

	b.load()
}

// load loads translation resources
// This function reads all directories under the specified resource path, where each directory represents a language
// For each language directory, it parses the translation files and adds the translation results to the bundle instance
func (b *Bundle) load() {
	// Read all directories under the resource path
	rd, err := os.ReadDir(b.opts.ResourcesPath)
	if err != nil {
		// Directory does not exist, return directly
		return
	}

	// Iterate through all directories and files under the resource path
	for _, f := range rd {
		// Only process directories, where directory names are lang IDs
		if f.IsDir() {

			folder := f.Name()

			file := filepath.Join(b.opts.ResourcesPath, folder, b.name+".json")

			// Check if file exists, if not return nil map and nil error
			if !isFileExist(file) {
				continue
			}

			// Read file content
			byts, err := os.ReadFile(file)
			if err != nil {
				// If an error occurs while reading the file, return the error
				log.Printf("read file %s: %v", file, err)
				continue
			}

			var trans map[string]string

			// Parse JSON data
			if err := json.Unmarshal(byts, &trans); err != nil {
				// If an error occurs while parsing JSON data, return the error
				log.Printf("unmarshal file %s: %v", file, err)
				continue
			}

			// Folder name is the language name
			b.WithTranslations(folder, trans)
		}
	}
}

// Reload reloads translation resources for all bundles in the cache.
// It iterates through all bundle instances in the cache and calls their load method
// to reload translation files from the filesystem.
func Reload() {
	bundleCache.Range(func(key, value any) bool {
		if b, ok := value.(*Bundle); ok {
			b.Reload()
		}
		return true
	})
}
