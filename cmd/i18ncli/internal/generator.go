package internal

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/iancoleman/orderedmap"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Generator struct {
	BaseDir   string             // Generator 工作的根目录。（可能是 module 的子目录）
	Module    string             // module name
	ModuleDir string             // module directory (go.mod path)
	Bundles   map[string]*Bundle // name => bundle
}

// ParsedFile 存储已解析的文件信息
type ParsedFile struct {
	FilePath  string
	Pkg       string // 当前文件的包名
	Ast       *ast.File
	Fset      *token.FileSet
	I18nAlias map[string]bool
}

func NewGenerator(baseDir string) *Generator {
	module, moduleDir, err := findModule(baseDir)
	if err != nil {
		panic(err)
	}

	return &Generator{
		BaseDir:   baseDir,
		Module:    module,
		ModuleDir: moduleDir,
		Bundles:   make(map[string]*Bundle),
	}
}

func (g *Generator) Walk() error {
	// 收集所有需要处理的Go文件并预解析
	parsedFiles, err := g.parseFiles()
	if err != nil {
		return err
	}

	// 找出所有的 i18n.Bundle("name") 赋值和其他相关赋值
	for _, f := range parsedFiles {
		g.collectBundles(f)
	}

	// 扫描：查找 Str 和 Err 方法调用
	for _, f := range parsedFiles {
		g.collectBundleTranslations(f)
	}

	return nil
}

func (g *Generator) GenerateTranslationFiles(fileType, resDir string, langs ...string) error {

	var resourceDir string
	if filepath.IsAbs(resDir) {
		resourceDir = resDir
	} else {
		resourceDir = filepath.Join(g.BaseDir, resDir)
	}

	langMap := map[language.Tag]struct{}{}
	for _, lang := range langs {
		tag := language.Make(lang)
		langMap[tag] = struct{}{}
	}

	if rd, err := os.ReadDir(resourceDir); err == nil {
		for _, f := range rd {
			if f.IsDir() {
				folder := f.Name()
				tag := language.Make(folder)
				langMap[tag] = struct{}{}
			}
		}
	}

	if len(langMap) == 0 {
		return fmt.Errorf("no language specified")
	}

	var marshal func(v any) ([]byte, error)
	switch fileType {
	case "yaml", "yml":
		marshal = yaml.Marshal
	case "toml", "tml":
		marshal = toml.Marshal
	case "ini":
		marshal = marshalINI
	case "json":
		fallthrough
	default:
		marshal = func(v any) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		}
	}

	for lang := range langMap {
		langDir := filepath.Join(resourceDir, lang.String())
		if err := os.MkdirAll(langDir, 0755); err != nil {
			log.Printf("[ERROR] create dir %s: %v", langDir, err)
			continue // 忽略错误
		}

		for _, bundle := range g.Bundles {

			filePath := filepath.Join(langDir, bundle.Name+"."+fileType)

			// Check if file exists
			translations := orderedmap.New()
			translations.SetEscapeHTML(false)
			if content, err := os.ReadFile(filePath); err == nil {
				// File exists, parse existing content
				if err := json.Unmarshal(content, translations); err != nil {
					log.Printf("[ERROR] parse file %s: %v", filePath, err)
				}
			}

			// changed mark
			changed := false
			// Add format strings as both keys and values
			for txt := range bundle.Trans {
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
			data, err := marshal(translations)
			if err != nil {
				log.Printf("[ERROR] marshal translations: %v", err)
				continue
			}

			if err := os.WriteFile(filePath, data, 0644); err != nil {
				log.Printf("[ERROR] write file %s: %v", filePath, err)
				continue
			}

		}

	}

	return nil
}

func (g *Generator) collectBundles(f *ParsedFile) {
	ast.Inspect(f.Ast, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			// 处理赋值语句:
			// 1. user := i18n.Bundle("user") (短变量声明)
			// 2. user = i18n.Bundle("user") (赋值语句)
			// 3. user := locales.User (从其他包导入的bundle变量)
			for i, lhs := range stmt.Lhs {
				if len(stmt.Rhs) > i {
					if ident, isIdent := lhs.(*ast.Ident); isIdent {
						if callExpr, isCall := stmt.Rhs[i].(*ast.CallExpr); isCall {
							// 处理 i18n.Bundle("user") 这样的调用
							if bundleName := extractBundleName(callExpr, f.I18nAlias); bundleName != "" {
								// 将变量信息添加到bundle中
								bundle := g.getBundleOrNew(bundleName)
								pkg, _ := g.getFullPkgPath(filepath.Dir(f.FilePath))
								bundle.AddVarDefine(ident.Name, pkg, f.FilePath)
							}
						} else if selector, isSelector := stmt.Rhs[i].(*ast.SelectorExpr); isSelector {
							// 处理从其他包导入的变量，如 locales.User
							if xIdent, isXIdent := selector.X.(*ast.Ident); isXIdent {
								// 检查是否是已知的包导入
								if pkg := findPkgByID(f.Ast, xIdent.Name); pkg != "" {
									// 将变量信息添加到bundle中
									bundle := g.getBundleOrNew(selector.Sel.Name)
									bundle.AddVarDefine(ident.Name, pkg, f.FilePath)
								}
							}
						}
					}
				}
			}
		case *ast.GenDecl:
			// 处理 var user = i18n.Bundle("user") 这样的显式变量声明语句
			// 以及 var user = locales.User 这样的导入变量赋值
			if stmt.Tok == token.VAR {
				for _, spec := range stmt.Specs {
					if valueSpec, isValue := spec.(*ast.ValueSpec); isValue {
						for i, name := range valueSpec.Names {
							if len(valueSpec.Values) > i {
								if callExpr, isCall := valueSpec.Values[i].(*ast.CallExpr); isCall {
									if bundleName := extractBundleName(callExpr, f.I18nAlias); bundleName != "" {
										// 将变量信息添加到bundle中
										bundle := g.getBundleOrNew(bundleName)
										pkg, _ := g.getFullPkgPath(filepath.Dir(f.FilePath))
										bundle.AddVarDefine(name.Name, pkg, f.FilePath)
									}
								} else if selector, isSelector := valueSpec.Values[i].(*ast.SelectorExpr); isSelector {
									// 处理从其他包导入的变量
									if xIdent, isXIdent := selector.X.(*ast.Ident); isXIdent { // 检查是否是已知的包导入
										if pkg := findPkgByID(f.Ast, xIdent.Name); pkg != "" {
											// 将变量信息添加到bundle中
											bundle := g.getBundleOrNew(selector.Sel.Name)
											bundle.AddVarDefine(selector.Sel.Name, pkg, f.FilePath)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		return true
	})
}

func (g *Generator) collectBundleTranslations(f *ParsedFile) {
	ast.Inspect(f.Ast, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			// 处理选择器表达式（方法调用）
			if selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr); isSelector {
				methodName := selectorExpr.Sel.Name

				// 检查是否是 Str 或 Err 方法调用
				if methodName == "Str" || methodName == "Err" {
					// 检查是否是 i18n.Bundle().Str() 形式（直接链式调用）
					if funCall, isFunCall := selectorExpr.X.(*ast.CallExpr); isFunCall {
						if bundleName := extractBundleName(funCall, f.I18nAlias); bundleName != "" {
							bundle := g.getBundleOrNew(bundleName)
							g.addBundleStr(bundle, callExpr)
						}
						return true
					}

					// 检查是否是变量调用形式 bundleVar.Str()
					if ident, isIdent := selectorExpr.X.(*ast.Ident); isIdent {
						// 首先在当前包中查找变量
						if bundle, err := g.getBundleByVar(f.Pkg, ident.Name); err == nil {
							g.addBundleStr(bundle, callExpr)
						}
						return true
					}

					// 处理嵌套选择器，如 locales.User.Str()
					if selector, isSelector := selectorExpr.X.(*ast.SelectorExpr); isSelector {
						if xIdent, isXIdent := selector.X.(*ast.Ident); isXIdent {
							pkgPath := findPkgByID(f.Ast, xIdent.Name)
							if bundle, err := g.getBundleByVar(pkgPath, selector.Sel.Name); err == nil {
								g.addBundleStr(bundle, callExpr)
							}
						}
					}
					return true
				}

				// 检查是否是 NStr 或 NErr 方法调用
				if methodName == "NStr" || methodName == "NErr" {
					// 检查是否是 i18n.Bundle().NStr() 形式（直接链式调用）
					if funCall, isFunCall := selectorExpr.X.(*ast.CallExpr); isFunCall {
						if bundleName := extractBundleName(funCall, f.I18nAlias); bundleName != "" {
							bundle := g.getBundleOrNew(bundleName)
							g.addBundleNStrs(bundle, callExpr)
						}
						return true
					}

					// 检查是否是变量调用形式 bundleVar.NStr()
					if ident, isIdent := selectorExpr.X.(*ast.Ident); isIdent {
						// 首先在当前包中查找变量
						if bundle, err := g.getBundleByVar(f.Pkg, ident.Name); err == nil {
							g.addBundleNStrs(bundle, callExpr)
						}
						return true
					}

					// 处理嵌套选择器，如 locales.User.NStr()
					if selector, isSelector := selectorExpr.X.(*ast.SelectorExpr); isSelector {
						if xIdent, isXIdent := selector.X.(*ast.Ident); isXIdent {
							pkgPath := findPkgByID(f.Ast, xIdent.Name)
							if bundle, err := g.getBundleByVar(pkgPath, selector.Sel.Name); err == nil {
								g.addBundleNStrs(bundle, callExpr)
							}
						}
					}
					return true
				}
			}
		}
		return true
	})
}

func (g *Generator) addBundle(bundle *Bundle) {
	g.Bundles[bundle.Name] = bundle
}

func (g *Generator) getBundleOrNew(name string) *Bundle {
	bundle, ok := g.Bundles[name]
	if !ok {
		bundle = NewBundle(name)
		g.addBundle(bundle)
	}
	return bundle
}

func (g *Generator) getBundleByVar(packagePath string, varName string) (*Bundle, error) {
	for _, bundle := range g.Bundles {
		for _, bundleVar := range bundle.Vars {
			if bundleVar.Pkg == packagePath && bundleVar.Name == varName {
				return bundle, nil
			}
		}
	}

	return nil, fmt.Errorf("bundle not found for var %s in package %s", varName, packagePath)
}

func (g *Generator) getFullPkgPath(fileDir string) (string, error) {
	absFileDir, err := filepath.Abs(fileDir)
	if err != nil {
		return "", err
	}

	absModuleDir, err := filepath.Abs(g.ModuleDir)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(absModuleDir, absFileDir)
	if err != nil {
		return "", err
	}

	if rel == "." {
		return g.Module, nil
	}

	return filepath.Join(g.Module, rel), nil
}

func (g *Generator) parseFiles() (map[string]*ParsedFile, error) {
	// 收集所有需要处理的Go文件并预解析
	parsedFiles := map[string]*ParsedFile{}

	// 预先解析文件，后续无须重复解析文件，提升性能
	err := filepath.Walk(g.BaseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录和特定文件
		if info.IsDir() || strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
			return nil
		}

		// 只处理.go文件，忽略测试文件
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {

			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil // 跳过解析错误的文件
			}

			pkg, _ := g.getFullPkgPath(filepath.Dir(path))

			// 缓存解析结果
			parsedFiles[path] = &ParsedFile{
				Pkg:       pkg,
				FilePath:  path,
				Ast:       astFile,
				Fset:      fset,
				I18nAlias: findI18nImportAliases(astFile),
			}
		}

		return nil
	})

	return parsedFiles, err
}

func (g *Generator) addBundleStr(b *Bundle, callExpr *ast.CallExpr) {
	if transKey := getCallArgString(callExpr, 0); transKey != "" {
		b.AddTrans(transKey)
	}
}

func (g *Generator) addBundleNStrs(b *Bundle, callExpr *ast.CallExpr) {
	// singular
	if singular := getCallArgString(callExpr, 1); singular != "" {
		b.AddTrans(singular)
	}
	// plural
	if plural := getCallArgString(callExpr, 2); plural != "" {
		b.AddTrans(plural)
	}
}
