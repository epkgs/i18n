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

	bundleCache map[string]*Bundle // name => bundle
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

		bundleCache: map[string]*Bundle{},
	}
}

func (g *Generator) AddBundle(bundle *Bundle) {
	g.bundleCache[bundle.Name] = bundle
}

func (g *Generator) Bundle(name string) *Bundle {
	if bundle, ok := g.bundleCache[name]; ok {
		return bundle
	}

	bundle := newBundle(name)
	g.AddBundle(bundle)
	return bundle
}

func (g *Generator) GetBundleByVar(packagePath string, varName string) *Bundle {
	for _, bundle := range g.bundleCache {
		for _, bundleVar := range bundle.vars {
			if bundleVar.PackagePath == packagePath && bundleVar.Name == varName {
				return bundle
			}
		}
	}

	return nil
}

func (g *Generator) GenerateTranslations(langs ...string) error {

	for _, bundle := range g.bundleCache {
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
		g.collectFunctions(f, i18nPkg, path)

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

// collectFunctions collects bundle configurations from NewBundle calls
func (g *Generator) collectFunctions(f *ast.File, i18nPkg string, filePath string) {
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

	// Check if it's a call to function
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
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

	// Check if the function name is "SetDefaultLanguage"
	if selExpr.Sel.Name == "SetDefaultLanguage" {
		// Check if it has at least 1 argument
		if len(callExpr.Args) < 1 {
			return
		}

		// Extract the bundle name (first argument)
		lit, ok := callExpr.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		// Clean the string literal (remove quotes)
		g_DefaultLanguage = Unquote(lit.Value)
		return
	}

	// Check if the function name is "SetResourcesDir"
	if selExpr.Sel.Name == "SetResourcesDir" {
		// Check if it has at least 1 argument
		if len(callExpr.Args) < 1 {
			return
		}

		// Extract the bundle name (first argument)
		lit, ok := callExpr.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		// Clean the string literal (remove quotes)
		g_ResourcesDir = Unquote(lit.Value)
		return
	}

	// Check if the function name is "New"
	if selExpr.Sel.Name == "New" {
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

		bundle := g.Bundle(bundleName)
		bundle.AddVarDefine(
			ident.Name,
			filePath,
			packageName,
			g.formatPackagePath(filepath.Dir(filePath), packageName),
		)
		return
	}

	// Check if the function name is "Str" or "Err"
	if selExpr.Sel.Name == "Str" || selExpr.Sel.Name == "Err" {
		// Check if it has at least 2 argument
		if len(callExpr.Args) < 2 {
			return
		}

		// Extract the bundle name (first argument)
		lit, ok := callExpr.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		// Clean the string literal (remove quotes)
		bundleName := Unquote(lit.Value)

		bundle := g.Bundle(bundleName)

		// Extract the bundle name (first argument)
		lit, ok = callExpr.Args[1].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		// Clean the string literal (remove quotes)
		transKey := Unquote(lit.Value)
		bundle.AddTrans(transKey)
		return
	}

}

// collectDefinitions collects format strings from Bundle.Str and Bundle.Err method calls
func (g *Generator) collectDefinitions(f *ast.File) {
	// Collect format strings from Bundle.Str and Bundle.Err method calls
	ast.Inspect(f, func(n ast.Node) bool {
		// Look for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a method call on a selector expression
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Check if the method name is "Str" or "Err"
		if selExpr.Sel.Name != "Str" && selExpr.Sel.Name != "Err" {
			return true
		}

		// Check if it has at least 1 argument (the format string)
		if len(callExpr.Args) < 1 {
			return true
		}

		// Get the target bundle from the selector expression
		target, ok := selExpr.X.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Get the bundle variable name
		bundleVarName := target.Sel.Name

		// Find the package path for the selector's package
		packagePath := findPackagePath(f, target.X.(*ast.Ident).Name)
		if packagePath == "" {
			// Assume it's in the same package
			packagePath = g.formatPackagePath(f.Name.Name, f.Name.Name)
		}

		// Find bundle by var definition
		bundle := g.GetBundleByVar(packagePath, bundleVarName)
		if bundle == nil {
			return true
		}

		// Get the first argument (format string)
		firstArg := callExpr.Args[0]
		if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			// Argument is a string literal
			// Remove quotes and add to bundle definitions
			tranKey := Unquote(lit.Value)
			bundle.AddTrans(tranKey)
		}

		return true
	})
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
