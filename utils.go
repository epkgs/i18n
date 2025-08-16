package i18n

import (
	"bytes"
	"os"
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
