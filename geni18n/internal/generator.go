package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const I18nPackagePath = "github.com/epkgs/i18n"

type Generator struct {
	ModuleName string // module name
	ModuleDir  string // module dir
	BaseDir    string // generator work base path

	Bundles map[string]*Bundle // "package.varName" -> Bundle
}

func NewGenerator(searchPath string) *Generator {
	baseDir, _ := filepath.Abs(searchPath)
	moduleDir, moduleName, err := findModule(baseDir)
	if err != nil {
		panic(err)
	}

	return &Generator{
		BaseDir:    baseDir,
		ModuleDir:  moduleDir,
		ModuleName: moduleName,

		Bundles: make(map[string]*Bundle),
	}
}

func (g *Generator) GenerateTranslations(langs ...string) error {

	for _, bundle := range g.Bundles {
		if err := bundle.GenerateTranslationFile(g.BaseDir, langs...); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) CollectBundles() error {

	// First pass: collect all bundle variables
	err := filepath.Walk(g.BaseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the Go source file
		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, path, nil, 0)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Get the i18n package identifier
		i18nPkg := findPackageName(f, I18nPackagePath)

		// If no i18n package is imported, skip this file
		if i18nPkg == "" {
			return nil
		}

		// Collect bundle configurations
		g.collectBundleConfigs(f, i18nPkg, path)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory for bundle configs: %w", err)
	}

	// Second pass: collect format strings from Translate method calls
	err = filepath.Walk(g.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the Go source file
		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, path, nil, 0)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Collect format strings from Translate method calls
		g.collectDefinitions(f)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory for format strings: %w", err)
	}

	return nil
}

// findModule looks for go.mod file starting from searchPath and going up
// Returns the module path and the directory containing go.mod
func findModule(searchPath string) (moduleDir, moduleName string, err error) {
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
					moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
					return dir, moduleName, nil
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

// collectBundleConfigs collects bundle configurations from NewBundle calls
func (g *Generator) collectBundleConfigs(f *ast.File, i18nPkg string, filePath string) {
	// Look for variable declarations and assignments
	ast.Inspect(f, func(n ast.Node) bool {
		var assignStmt *ast.AssignStmt
		var ok bool

		// Handle regular assignments (var := value)
		if assignStmt, ok = n.(*ast.AssignStmt); ok {
			g.processAssignment(assignStmt, i18nPkg, filePath, f.Name.Name)
			return true
		}

		// Handle variable declarations (var name = value)
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok && len(valueSpec.Values) > 0 {
					// Create a temporary assignment statement for processing
					if callExpr, ok := valueSpec.Values[0].(*ast.CallExpr); ok {
						// Create a fake assignment statement for processing
						if len(valueSpec.Names) > 0 {
							fakeAssign := &ast.AssignStmt{
								Lhs: []ast.Expr{valueSpec.Names[0]},
								Tok: token.DEFINE, // or token.ASSIGN for var declarations
								Rhs: []ast.Expr{callExpr},
							}
							g.processAssignment(fakeAssign, i18nPkg, filePath, f.Name.Name)
						}
					}
				}
			}
			return true
		}

		return true
	})
}

// processAssignment processes an assignment statement to check if it's a NewBundle call
func (g *Generator) processAssignment(assignStmt *ast.AssignStmt, i18nPkg string, filePath string, packageName string) {
	// Check for assignments with a single left-hand side (e.g., var := ...)
	if len(assignStmt.Lhs) != 1 {
		return
	}

	// Check if the right-hand side is a call expression
	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return
	}

	// Check if it's a call to NewBundle function
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// Check if the function name is "NewBundle" and it's from the i18n package
	if selExpr.Sel.Name != "NewBundle" {
		return
	}

	// Check if the selector expression is from the i18n package
	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return
	}

	// Check if this identifier refers to the i18n package
	if ident.Name != i18nPkg {
		return
	}

	// Check if it has at least 1 argument (bundle name)
	if len(callExpr.Args) < 1 {
		return
	}

	// Extract the bundle name (first argument)
	lit, ok := callExpr.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	// Clean the string literal (remove quotes)
	bundleName := Unquote(lit.Value)

	// Get the variable name
	ident, ok = assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return
	}

	varName := ident.Name
	packagePath := g.formatPackagePath(filepath.Dir(filePath), packageName)
	// Record the variable association with package name as prefix
	key := packagePath + "." + varName
	var bundle *Bundle
	if b, ok := g.Bundles[key]; ok {
		bundle = b
	} else {
		bundle = newBundle()
	}
	bundle.Key = key
	bundle.Name = bundleName
	bundle.VarName = varName
	bundle.FilePath = filePath
	bundle.PackageName = packageName
	bundle.PackagePath = packagePath

	// Check all arguments for function literals (OptionsFunc)
	for _, arg := range callExpr.Args[1:] {
		// Try to extract options from function literal
		if funLit, ok := arg.(*ast.FuncLit); ok {
			// Get the parameter name from the function literal
			paramName := getOptionsParamName(funLit)
			if paramName != "" {
				// Look for assignments to opts fields
				ast.Inspect(funLit, func(n ast.Node) bool {
					// Look for assignment statements
					assignStmt, ok := n.(*ast.AssignStmt)
					if !ok {
						return true
					}

					// Check for assignments to fields of opts
					for i, lhs := range assignStmt.Lhs {
						// Check if this is a selector expression (e.g., opts.DefaultLang)
						selExpr, ok := lhs.(*ast.SelectorExpr)
						if !ok {
							continue
						}

						// Check if the receiver is the parameter name
						ident, ok := selExpr.X.(*ast.Ident)
						if !ok || ident.Name != paramName {
							continue
						}

						// Extract the value being assigned
						rhs := assignStmt.Rhs[i]
						lit, ok := rhs.(*ast.BasicLit)
						if !ok || lit.Kind != token.STRING {
							continue
						}

						// Clean the string literal (remove quotes)
						value := lit.Value
						if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
							// Unquote the string
							unquoted, err := strconv.Unquote(value)
							if err != nil {
								// If unquoting fails, just remove the quotes
								value = value[1 : len(value)-1]
							} else {
								value = unquoted
							}
						}

						// Set the appropriate field based on the selector
						switch selExpr.Sel.Name {
						case "DefaultLang":
							bundle.opts.DefaultLang = value
						case "ResourcesPath":
							bundle.opts.ResourcesPath = value
						}
					}

					return true
				})
			}
		}
	}

	// Record the bundle
	g.Bundles[key] = bundle
}

// collectDefinitions collects format strings from Translate method calls
func (g *Generator) collectDefinitions(f *ast.File) {
	// Collect format strings from Translate method calls
	ast.Inspect(f, func(n ast.Node) bool {
		// Look for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a method call on an identifier with name "Definef"
		// Handle both regular calls and generic calls
		var selExpr *ast.SelectorExpr

		// Case 1: Regular function call like i18n.Define(...)
		if se, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			selExpr = se
		} else if idxExpr, ok := callExpr.Fun.(*ast.IndexExpr); ok {
			// Case 2: Generic function call like i18n.Definef[string](...)
			if se, ok := idxExpr.X.(*ast.SelectorExpr); ok {
				selExpr = se
			}
		}

		// If we don't have a selector expression, skip
		if selExpr == nil {
			return true
		}

		// Check if the method name is "Definef", "Define"
		if selExpr.Sel.Name != "Define" && selExpr.Sel.Name != "Definef" &&
			selExpr.Sel.Name != "DefineError" && selExpr.Sel.Name != "DefineErrorf" {
			return true
		}

		// Check if it has at least 2 arguments (bundle and format string)
		if len(callExpr.Args) < 2 {
			return true
		}

		// 获取第一个参数
		//   - 参数应为变量
		//   - 如果是当前包的变量，则格式化 key 为当前package path 全路径 + 变量名
		//   - 如果是外部包的变量，则格式化 key 为外部包的全路径 + 变量名
		//   - 通过 key 查找语言包
		//   - 如果未找到，则跳过
		firstArg := callExpr.Args[0]
		var bundle *Bundle

		switch arg := firstArg.(type) {
		case *ast.Ident:
			// Direct identifier like "bundle"
			// Create key using current package
			key := g.formatPackagePath(f.Name.Name, f.Name.Name) + "." + arg.Name
			var exists bool
			bundle, exists = g.Bundles[key]
			if !exists {
				return true
			}

		case *ast.SelectorExpr:
			// Selector expression like "locales.User"
			if ident, ok := arg.X.(*ast.Ident); ok {
				// Find the package path for the selector's package
				packagePath := findPackagePath(f, ident.Name)
				if packagePath != "" {
					// Create key using the found package path
					key := packagePath + "." + arg.Sel.Name
					var exists bool
					bundle, exists = g.Bundles[key]
					if !exists {
						return true
					}
				} else {
					// Assume it's in the same package
					key := g.formatPackagePath(f.Name.Name, f.Name.Name) + "." + arg.Sel.Name
					var exists bool
					bundle, exists = g.Bundles[key]
					if !exists {
						return true
					}
				}
			} else {
				return true
			}

		default:
			// Unsupported argument type
			return true
		}

		// 获取第二个参数
		//   - 参数应为字符串
		//   - 将其去除双引号后添加进语言包的definition
		secondArg := callExpr.Args[1]
		if lit, ok := secondArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			// 参数是字符串
			// 去除双引号后添加进语言包的definition
			definition := Unquote(lit.Value)
			if bundle != nil {
				bundle.AddDefinition(definition)
			}
		}

		// 尝试获取第三个参数（如有）
		//   - 尝试获取字符串参数
		//   - 如果是字符串，则将其去除双引号后添加进语言包的definition
		if len(callExpr.Args) >= 3 {
			thirdArg := callExpr.Args[2]
			if lit, ok := thirdArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				// 参数是字符串
				// 去除双引号后添加进语言包的definition
				definition := Unquote(lit.Value)
				if bundle != nil {
					bundle.AddDefinition(definition)
				}
			}
		}

		return true
	})
}

// getOptionsParamName extracts the parameter name from a function literal
// that matches the OptionsFunc pattern func(opts *i18n.Options)
func getOptionsParamName(funLit *ast.FuncLit) string {
	// Check if the function has exactly one parameter
	if funLit.Type.Params == nil || len(funLit.Type.Params.List) != 1 {
		return ""
	}

	// Get the first parameter
	param := funLit.Type.Params.List[0]

	// Check if it has a name
	if len(param.Names) == 0 {
		return ""
	}

	// Return the parameter name
	return param.Names[0].Name
}

// formatPackagePath constructs the full package path from file path, module directory and module path
func (g *Generator) formatPackagePath(fileDir, packageName string) string {

	absFileDir := fileDir
	if !filepath.IsAbs(fileDir) {
		absFileDir = filepath.Join(g.BaseDir, fileDir)
	}

	absModuleDir, err := filepath.Abs(g.ModuleDir)
	if err != nil {
		return packageName
	}

	// Get the relative path from module directory
	relPath, err := filepath.Rel(absModuleDir, absFileDir)
	if err != nil {
		return packageName
	}

	// If the file is in the module root, return just the module path
	if relPath == "." {
		return g.ModuleName
	}

	// Otherwise, combine module path with relative path
	return filepath.Join(g.ModuleName, relPath)
}

// findPackagePath finds the full package path from the file's import statements
func findPackagePath(f *ast.File, packageName string) string {
	// Look through the imports for the package name
	for _, imp := range f.Imports {
		// Remove quotes from the import path
		importPath := Unquote(imp.Path.Value)

		if imp.Name != nil {
			// Check if there's an alias that matches the package name
			if imp.Name.Name == packageName {
				return importPath
			}
		} else {
			// Check if the import path ends with the package name
			if filepath.Base(importPath) == packageName {
				return importPath
			}
		}

	}

	return packageName // Not found
}

func findPackageName(f *ast.File, packagePath string) string {
	// Look through the imports for github.com/epkgs/i18n
	for _, imp := range f.Imports {
		if Unquote(imp.Path.Value) == packagePath {
			// Check if it has an alias
			if imp.Name != nil {
				return imp.Name.Name
			}
			return filepath.Base(packagePath)
		}
	}
	return ""
}

func Unquote(str string) string {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
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
