package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/orderedmap"
)

var (
	g_DefaultLanguage = "en"      // Default value
	g_ResourcesDir    = "locales" // Default value
)

type Var struct {
	Name        string // The variable name
	FilePath    string // The file where this variable is defined
	PackageName string // The package identifier name
	PackagePath string // The package path
}

// Bundle represents a variable that holds a Bundle
type Bundle struct {
	Name string

	vars []Var

	trans map[string]struct{}
}

func newBundle(name string) *Bundle {
	return &Bundle{
		Name:  name,
		vars:  []Var{},
		trans: make(map[string]struct{}),
	}
}

func (b *Bundle) AddVarDefine(varName string, filePath string, packageName string, packagePath string) {
	exist := false
	for _, v := range b.vars {
		if v.Name == varName && v.FilePath == filePath {
			exist = true
			break
		}
	}

	if !exist {
		b.vars = append(b.vars, Var{
			Name:        varName,
			FilePath:    filePath,
			PackageName: packageName,
			PackagePath: packagePath,
		})
	}
}

func (b *Bundle) AddTrans(key string) {
	b.trans[key] = struct{}{}
}

func (b *Bundle) GenerateTranslationFile(baseDir string, langs ...string) error {

	var resPath string
	if filepath.IsAbs(g_ResourcesDir) {
		resPath = g_ResourcesDir
	} else {
		absBaseDir, _ := filepath.Abs(baseDir)
		resPath = filepath.Join(absBaseDir, g_ResourcesDir)
	}

	uniqueLangs := map[string]string{} // lang ID -> folder name
	for _, lang := range langs {
		id := formatLangID(lang)
		uniqueLangs[id] = id
	}

	// Traverse the subdirectories of resPath and use directory names as language IDs
	if rd, err := os.ReadDir(resPath); err == nil {
		for _, f := range rd {
			if f.IsDir() {
				folder := f.Name()
				id := formatLangID(folder)
				uniqueLangs[id] = folder
			}
		}
	}

	if len(uniqueLangs) == 0 {
		uniqueLangs[g_DefaultLanguage] = g_DefaultLanguage
	}

	for _, folder := range uniqueLangs {
		langDir := filepath.Join(resPath, folder)
		if err := os.MkdirAll(langDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", langDir, err)
		}

		filePath := filepath.Join(langDir, b.Name+".json")

		// Check if file exists
		translations := orderedmap.New()
		translations.SetEscapeHTML(false)
		if content, err := os.ReadFile(filePath); err == nil {
			// File exists, parse existing content
			if err := json.Unmarshal(content, translations); err != nil {
				return fmt.Errorf("failed to parse existing file %s: %w", filePath, err)
			}
		}

		// changed mark
		changed := false

		// Add format strings as both keys and values
		for txt := range b.trans {
			// Only add if not already present
			if _, exists := translations.Get(txt); !exists {
				translations.Set(txt, txt)
				changed = true // Mark as changed
			}
		}

		if !changed {
			continue // No changes, skip
		}

		// Write to file
		data, err := json.MarshalIndent(translations, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal translations: %w", err)
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

	}

	return nil
}
