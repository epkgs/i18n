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

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

var langIDCache = map[string]string{}

func formatLangID(lang string) string {
	if id, ok := langIDCache[lang]; ok {
		return id
	}
	id := strings.Replace(lang, "-", "_", -1)
	langIDCache[lang] = id
	return id
}

var paser = template.New("i18n")

func parseTemplate(msg string, arg1 any) string {

	// 使用 text/template 解析结构体或 map
	tmpl, err := paser.Parse(msg)
	if err != nil {
		return msg // 解析失败回退
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, arg1); err != nil {
		return msg // 执行失败回退
	}
	return buf.String()
}

var languageTagCache = make(map[string]*language.Tag)

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
			// 结构体为零值 或 空结构体（无字段），避免模板渲染失败
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

			// 将数组/切片转换为 []any
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
