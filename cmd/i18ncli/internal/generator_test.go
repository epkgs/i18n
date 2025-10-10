// generator_test.go
package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 创建模拟的 go.mod 文件
	goModContent := `module test/module
	
	go 1.21
	`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 测试正常情况
	gen := NewGenerator(tempDir)
	if gen == nil {
		t.Error("Expected generator to be created, got nil")
		return
	}

	if gen.BaseDir != tempDir {
		t.Errorf("Expected BaseDir %s, got %s", tempDir, gen.BaseDir)
	}

	if gen.Module != "test/module" {
		t.Errorf("Expected Module test/module, got %s", gen.Module)
	}
}

func TestGeneratorWalk(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 创建模拟的 go.mod 文件
	goModContent := `module test/module
	
	go 1.21
	`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建测试Go文件
	testGoFile := `
	package main
	
	import "github.com/epkgs/i18n"
	
	func main() {
		userBundle := i18n.Bundle("user")
		message := userBundle.Str("Hello, world!")
	}
	`

	err = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(testGoFile), 0644)
	if err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(tempDir)
	err = gen.Walk()
	if err != nil {
		t.Errorf("Walk failed: %v", err)
	}

	// 验证是否正确识别了bundle
	if len(gen.Bundles) == 0 {
		t.Error("Expected bundles to be collected, got none")
	}

	userBundle, exists := gen.Bundles["user"]
	if !exists {
		t.Error("Expected 'user' bundle to be created")
	}

	// 验证是否收集到了翻译文本
	if len(userBundle.Trans) == 0 {
		t.Error("Expected translations to be collected")
	}
}

func TestGeneratorGenerateTranslationFiles(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	resDir := "locales"

	// 创建模拟的 go.mod 文件
	goModContent := `module test/module
	
	go 1.21
	`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(tempDir)

	// 添加一些测试数据
	bundle := gen.getBundleOrNew("user")
	bundle.AddTrans("Hello, world!")
	bundle.AddTrans("Goodbye!")

	// 测试生成JSON格式的翻译文件
	err = gen.GenerateTranslationFiles("json", resDir, "en", "zh")
	if err != nil {
		t.Errorf("GenerateTranslationFiles failed: %v", err)
	}

	// 验证文件是否被创建
	enDir := filepath.Join(tempDir, resDir, "en")
	userJsonFile := filepath.Join(enDir, "user.json")

	if _, err := os.Stat(userJsonFile); os.IsNotExist(err) {
		t.Error("Expected user.json file to be created")
	}

	zhDir := filepath.Join(tempDir, resDir, "zh")
	userJsonFile = filepath.Join(zhDir, "user.json")

	if _, err := os.Stat(userJsonFile); os.IsNotExist(err) {
		t.Error("Expected user.json file to be created for zh locale")
	}
}

func TestGeneratorCollectBundles(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 创建模拟的 go.mod 文件
	goModContent := `module test/module
		
		go 1.21
		`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(tempDir)

	// 移除了未使用的 parsedFile 声明

	// TODO: 由于需要构建AST来进行完整的测试，这里只是基本结构检查
	// 在实际应用中，应该构建完整的AST来测试此方法的功能
	if gen.Bundles == nil {
		t.Error("Bundles should be initialized")
	}
}

func TestGetBundleOrNew(t *testing.T) {
	tempDir := t.TempDir()

	goModContent := `module test/module
		
		go 1.21
		`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(tempDir)

	// 获取新bundle
	bundle := gen.getBundleOrNew("test")
	if bundle == nil {
		t.Error("Expected new bundle to be created")
		return // 防止后续访问nil指针
	}

	if bundle.Name != "test" {
		t.Errorf("Expected bundle name 'test', got %s", bundle.Name)
	}

	// 再次获取相同名称的bundle，应返回同一个实例
	sameBundle := gen.getBundleOrNew("test")
	if bundle != sameBundle {
		t.Error("Expected same bundle instance")
	}
}

func TestAddBundleStr(t *testing.T) {
	tempDir := t.TempDir()

	goModContent := `module test/module
	
	go 1.21
	`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(tempDir)
	bundle := gen.getBundleOrNew("test")

	// TODO: 由于addBundleStr需要AST.CallExpr参数，此处仅为结构验证
	// 实际测试需要构建AST节点
	if len(bundle.Trans) == 0 {
		// 初始状态应该是空的
	} else {
		t.Error("Expected empty translations initially")
	}
}
