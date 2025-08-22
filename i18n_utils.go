package i18n

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/text/language"
)

// isFileExist checks if a file exists at the given path
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// langIDCache caches formatted language identifiers
var langIDCache = map[string]string{}

// formatLangID formats a language identifier by replacing hyphens with underscores
// and caches the result for performance
func formatLangID(lang string) string {
	if id, ok := langIDCache[lang]; ok {
		return id
	}
	id := strings.Replace(lang, "-", "_", -1)
	langIDCache[lang] = id
	return id
}

// paser is a template parser for internationalization
var paser = template.New("i18n")

// parseTemplate parses a template message with the given argument
func parseTemplate(msg string, arg1 any) string {

	// Parse struct or map using text/template
	tmpl, err := paser.Parse(msg)
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
var languageTagCache = make(map[string]*language.Tag)

// parseLanguageTag parses a language string into a language.Tag and caches the result
func parseLanguageTag(lang string) *language.Tag {
	if _, exist := languageTagCache[lang]; !exist {
		t, e := language.Parse(lang)
		if e != nil {
			languageTagCache[lang] = nil
		} else {
			languageTagCache[lang] = &t
		}
	}

	return languageTagCache[lang]
}

// parse processes a translated string with the given arguments
// It handles different argument types appropriately:
//   - Single struct or map: uses template parsing
//   - Single slice or array: expands elements as separate arguments
//   - Multiple arguments or other types: uses standard fmt.Sprintf
func parse(transleted string, args ...any) string {

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
			return parse(transleted, slices...)

		default:
			return fmt.Sprintf(transleted, arg1)
		}
	}

	return fmt.Sprintf(transleted, args...)
}
