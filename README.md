# 🌍 i18n - Internationalization library for Go

A simple yet powerful internationalization library for Go applications with support for translation and localized error handling.

## 📌 Features

- ✅ Simple API for string translation
- ✅ Automatic language detection from context
- ✅ Support for parameterized translations
- ✅ JSON-based translation files
- ✅ Gin middleware for HTTP applications
- ✅ Internationalized error handling
- ✅ Thread-safe bundle caching
- ✅ CLI tool for extracting translation keys from source code

## 🚀 Quick Start
### 1. Define translation bundles
Create a bundle for your translations:

```go
// locales/user.go
package locales

import "github.com/epkgs/i18n"

var User = i18n.Bundle("user")
```

### 2. Create translation files
Create JSON translation files in your resources directory:

```bash
locales/
├── en/
│   └── user.json
└── zh-CN/
    └── user.json
```

Example `locales/en/user.json`:
```json
{
  "User %s not exist": "User %s does not exist",
  "%d item found": "%d item found",
  "%d items found": "%d items found"
}
```

Example `locales/zh-CN/user.json`:
```json
{
  "User %s not exist": "用户 %s 不存在",
  "%d item found": "找到 %d 个项目",
  "%d items found": "找到 %d 个项目"
}
```

### 3. Use translations in your code
```go
package main

import (
    "context"
    "fmt"
    
    "github.com/epkgs/i18n"
    "path/to/locales"
)

func main() {
    // Create a context with language preference
    ctx := i18n.WithAcceptLanguages(context.Background(), "zh-CN")
    
    // Create a translatable string
    user := "alice"
    message := locales.User.Str("User %s not exist", user)
    
    // Get default string
    fmt.Printf("Default: %s\n", message)
    
    // Get translated string
    fmt.Printf("Translated: %s\n", message.T(ctx))
    
    // Using singular/plural forms
    itemCount := 1
    pluralMessage := locales.User.NStr(itemCount, "%d item found", "%d items found", itemCount)
    fmt.Printf("Plural Translated: %s\n", pluralMessage.T(ctx))
}
```

### 4. Use with Gin web framework

```go
package main

import (
    "github.com/epkgs/i18n"
    "github.com/gin-gonic/gin"
    "golang.org/x/text/language"
)

func main() {
    r := gin.Default()
    
    // Add i18n middleware
    r.Use(i18n.GinMiddleware(language.AmericanEnglish.String()))
    
    r.GET("/api/user", func(c *gin.Context) {
        // The context now contains language preferences
        // based on Accept-Language header, query params, or cookies
        message := locales.User.Str("User not found")
        c.JSON(404, gin.H{
            "error": message.T(c.Request.Context()),
        })
    })
    
    r.Run(":8080")
}
```

### 5. Internationalized errors
```go
func someHandler(c *gin.Context) {
    err := locales.User.Err("User %s not exist", "alice")
    // err implements error interface and can be translated
    response.Fail(c, err)
    
    // Using singular/plural forms for errors
    itemCount := 0
    pluralErr := locales.User.NErr(itemCount, "%d item found", "%d items found", itemCount)
    response.Fail(c, pluralErr)
}
```

## 🛠️ API Reference
### Bundle
The main component for managing translations.
```go
// Create a translatable string
str := bundle.Str("Hello %s", "world")

// Create an internationalized error
err := bundle.Err("Something went wrong: %s", details)

// Create a translatable string with singular/plural forms
itemCount := 1
str := bundle.NStr(itemCount, "%d item found", "%d items found", itemCount)

// Create an internationalized error with singular/plural forms
import "github.com/epkgs/i18n/plural"

itemCount := 5
err := bundle.NErr(plural.IsOne(itemCount), "%d item found", "%d items found", itemCount)
```

### Context Integration
```go
// Set language preferences in context
ctx := i18n.WithAcceptLanguages(context.Background(), "zh-CN", "zh", "en")

// Get language preferences from context
langs := i18n.GetAcceptLanguages(ctx)
```

### Gin Middleware
```go
// Use the middleware to automatically detect language preferences
r.Use(i18n.GinMiddleware("en")) // "en" is the fallback language
```

The middleware checks for language preferences in this order:

  1. Query parameter lang
  2. Cookie lang
  3. Accept-Language header
  4. Default language

## 📁 Translation File Structure

Translation files are organized by language directories:

```bash
locales/
├── en/
│   ├── user.json
│   └── common.json
├── zh-CN/
│   ├── user.json
│   └── common.json
└── es/
    ├── user.json
    └── common.json
```

Each JSON file contains key-value pairs where the key is the original string and the value is the translation:
```json
{
  "Welcome %s": "欢迎 %s",
  "User not found": "用户未找到",
  "Invalid input": "输入无效",
  "%d item found": "找到 %d 个项目",
  "%d items found": "找到 %d 个项目"
}
```

## 🧰 i18n CLI Tool
The project includes a CLI tool to help extract translation keys from your source code:

```bash
# Install the CLI tool
go install github.com/epkgs/i18n/cmd/i18ncli@latest

# Extract translation keys from your project
i18ncli extract

# You can also use go generate, as shown in the examples
//go:generate i18ncli extract
```

This tool scans your Go source files for `Str`, `NStr`, `Err` and `NErr` calls of `i18n.Bundle`, 
extracts the format strings, and automatically creates or updates the translation files.

## 📄 License
This project is licensed under the MIT License.