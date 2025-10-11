package internal

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"github.com/epkgs/i18n/errors"
	"golang.org/x/text/language"
	"gopkg.in/ini.v1"
)

// templateParser is a template parser for internationalization
var templateParser = template.New("i18n")

// parseTemplate parses a template message with the given argument
func parseTemplate(msg string, arg1 any) string {

	// Parse struct or map using text/template
	tmpl, err := templateParser.Parse(msg)
	if err != nil {
		return msg // Fallback on parse failure
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, arg1); err != nil {
		return msg // Fallback on execution failure
	}
	return buf.String()
}

// languageTagCache caches parsed language tags
var languageTagCache = make(map[string]language.Tag)

// parseLanguageTag parses a language string into a language.Tag and caches the result
func ParseLanguageTag(lang string) language.Tag {
	if _, exist := languageTagCache[lang]; !exist {
		t, e := language.Parse(lang)
		if e != nil {
			languageTagCache[lang] = language.Und
		} else {
			languageTagCache[lang] = t
		}
	}

	return languageTagCache[lang]
}

func ParseLanguageTags(langs ...string) []language.Tag {
	tags := make([]language.Tag, len(langs))

	for i, l := range langs {
		tags[i] = ParseLanguageTag(l)
	}

	return tags
}

// parse processes a translated string with the given arguments
// It handles different argument types appropriately:
//   - Single struct or map: uses template parsing
//   - Single slice or array: expands elements as separate arguments
//   - Multiple arguments or other types: uses standard fmt.Sprintf
func Parse(transleted string, args ...any) string {

	if len(args) == 0 {
		return transleted
	}

	if len(args) == 1 {

		arg1 := args[0]

		v := reflect.ValueOf(arg1)

		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return transleted
			}
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			// Zero value struct or empty struct (no fields), avoid template rendering failure
			if v.IsZero() || v.NumField() == 0 {
				return transleted
			}

			return parseTemplate(transleted, arg1)

		case reflect.Map:
			if v.Len() == 0 {
				return transleted
			}

			return parseTemplate(transleted, arg1)
		case reflect.Array, reflect.Slice:
			if v.Len() == 0 {
				return transleted
			}

			// Convert array/slice to []any
			slices := make([]any, v.Len())
			for i := 0; i < v.Len(); i++ {
				slices[i] = v.Index(i).Interface()
			}
			return Parse(transleted, slices...)

		default:
			return fmt.Sprintf(transleted, arg1)
		}
	}

	return fmt.Sprintf(transleted, args...)
}

func UnmarshalINI(data []byte, val any) error {
	f, err := ini.Load(data)
	if err != nil {
		return err
	}

	m, ok := (val.(*map[string]any))
	if !ok {
		return errors.New("val is not a pointer to map[string]any")
	}

	section, err := f.GetSection(ini.DefaultSection)
	if err != nil {
		return err
	}

	for _, key := range section.Keys() {
		(*m)[key.Name()] = key.Value()
	}

	return nil
}

func IndexOf[T comparable](slice []T, val T) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func Includes[T comparable](slice []T, val T) bool {
	return IndexOf(slice, val) >= 0
}
