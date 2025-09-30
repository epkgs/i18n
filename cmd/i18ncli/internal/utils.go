package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

// getCallArgString 从方法调用中提取字符串入参
func getCallArgString(callExpr *ast.CallExpr, pos int) string {
	if len(callExpr.Args) > 0 {
		if lit, isLit := callExpr.Args[pos].(*ast.BasicLit); isLit && lit.Kind == token.STRING {
			return unquote(lit.Value)
		}
	}
	return ""
}

// findModule looks for go.mod file starting from searchPath and going up
// Returns the module path and the directory containing go.mod
func findModule(searchPath string) (module, moduleDir string, err error) {
	dir, err := filepath.Abs(searchPath)
	if err != nil {
		return "", "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// Found go.mod, read it to get module path
			content, err := os.ReadFile(goModPath)
			if err != nil {
				return "", "", fmt.Errorf("failed to read go.mod: %w", err)
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "module ") {
					module = strings.TrimSpace(strings.TrimPrefix(line, "module "))
					return module, dir, nil
				}
			}
			return "", "", fmt.Errorf("module directive not found in go.mod")
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root
			break
		}
		dir = parent
	}

	return "", "", fmt.Errorf("go.mod file not found")
}

func findI18nImportAliases(f *ast.File) map[string]bool {

	i18nAliases := make(map[string]bool)

	i18nPkgPath := "github.com/epkgs/i18n"

	for _, imp := range f.Imports {
		if unquote(imp.Path.Value) == i18nPkgPath {
			// Check if it has an alias
			if imp.Name != nil {
				i18nAliases[imp.Name.Name] = true
			}
			i18nAliases[filepath.Base(i18nPkgPath)] = true
		}
	}

	return i18nAliases
}

func findPkgByID(f *ast.File, pkgID string) string {
	for _, imp := range f.Imports {
		if imp.Name != nil {
			if imp.Name.Name == pkgID {
				return unquote(imp.Path.Value)
			}
		}

		pkgPath := unquote(imp.Path.Value)
		id := filepath.Base(pkgPath)
		if id == pkgID {
			return pkgPath
		}
	}
	return ""
}

func unquote(str string) string {

	if len(str) >= 2 &&
		((str[0] == '"' && str[len(str)-1] == '"') || (str[0] == '`' && str[len(str)-1] == '`')) {
		// Unquote the string
		unquoted, err := strconv.Unquote(str)
		if err != nil {
			// If unquoting fails, just remove the quotes
			str = str[1 : len(str)-1]
		} else {
			str = unquoted
		}
	}

	return str
}

// extractBundleName 从函数调用中提取bundle名称
func extractBundleName(callExpr *ast.CallExpr, i18nAliases map[string]bool) string {
	if selector, isSelector := callExpr.Fun.(*ast.SelectorExpr); isSelector {
		if selIdent, isIdent := selector.X.(*ast.Ident); isIdent {
			// 检查是否是i18n包或其别名的Bundle方法
			// 支持以下情况:
			// 1. i18n.Bundle("bundleName")
			// 2. 别名形式: i18nAlias.Bundle("bundleName")
			// 3. 包完整路径的别名形式
			if i18nAliases[selIdent.Name] && selector.Sel.Name == "Bundle" {
				if len(callExpr.Args) > 0 {
					if lit, isLit := callExpr.Args[0].(*ast.BasicLit); isLit && lit.Kind == token.STRING {
						// 更健壮地处理字符串字面量
						bundleName := strings.Trim(lit.Value, "\"")
						// 处理原始字符串字面量
						if strings.HasPrefix(lit.Value, "`") && strings.HasSuffix(lit.Value, "`") {
							bundleName = strings.Trim(lit.Value, "`")
						}
						return bundleName
					}
				}
			}
		}
	}
	return ""
}

func marshalINI(v any) ([]byte, error) {
	// 创建一个新的INI文件对象
	cfg := ini.Empty()

	// 获取默认section
	section, err := cfg.GetSection("")
	if err != nil {
		return nil, err
	}

	// 检查输入是否为map[string]any类型
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("input is not a map[string]any")
	}

	// 将map中的键值对写入INI文件
	for key, value := range m {
		_, err := section.NewKey(key, fmt.Sprintf("%v", value))
		if err != nil {
			return nil, err
		}
	}

	// 将INI内容写入字节缓冲区
	var buf strings.Builder
	_, err = cfg.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
