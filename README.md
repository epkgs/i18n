下面是一个完整的 [README.md](README.md) 文档，适用于你的 i18n 错误系统项目。该文档从**项目简介**、**核心功能**、**使用示例**、**结构说明**、**接口定义**、**错误包装器**、**国际化支持**、**测试建议**等维度进行详细说明，适用于团队内部文档或开源项目说明。

---

# 🌍 i18n 错误系统（支持国际化、错误包装、上下文注入）

一个支持 **国际化（i18n）**、**错误定义**、**错误包装器（Wrapper）** 和 **上下文注入** 的 Go 错误系统。

## 📌 项目简介

该项目提供了一个完整的错误定义与包装系统，支持：

- ✅ 多语言翻译支持（基于 context）
- ✅ 自定义错误模板
- ✅ 错误包装器（如错误码、HTTP 状态码、TraceID）
- ✅ 支持 Gin 等 Web 框架集成
- ✅ 链式构建错误定义
- ✅ 可扩展的错误系统

## 🧱 模块结构

```bash
i18n/
├── i18n.go                # i18n 翻译器主结构
├── i18n_item.go           # 翻译项定义
├── context.go             # 上下文操作（用于注入 Accept-Language）
├── error.go               # 核心错误类型定义
├── error_definition.go    # 错误定义构造器
└── error_extras.go        # 错误包装器定义
```

## 🧩 核心功能

### ✅ 1. 国际化（i18n）支持

通过 [I18n](i18n_item.go#L11-L11) 和 [Item](i18n_item.go#L10-L16) 实现多语言翻译支持，支持：

- 多语言映射
- 自动语言匹配
- 支持 fallback 到默认语言
- 支持上下文注入语言

### ✅ 2. 错误定义系统

通过 [errorDefinition](error_definition.go#L12-L16) 定义错误模板，支持：

- 定义错误格式（如 `"user %s not found"`）
- 支持链式构建错误定义
- 支持运行时注入参数（如 `err.New(ctx, "alice")`）

### ✅ 3. 错误包装器（Error Wrapper）

支持链式包装错误，例如：

- 错误码（Code）
- HTTP 状态码（HTTP Status）
- TraceID（追踪 ID）
- 自定义包装器（如日志、监控）

### ✅ 4. 上下文注入语言

通过 `context` 支持注入 `Accept-Language`，用于动态选择语言：

```go
ctx := i18n.WithAcceptLanguages(context.Background(), "zh-CN", "zh")
```

### ✅ 5. 与 Gin 集成（可选）

支持与 Gin Web 框架无缝集成，可作为 Gin 的 `context.Error()` 使用。

## 🧱 使用示例

### 1. 初始化 i18n 翻译器

```go
i18n := i18n.NewCatalog("user")
i18n.AddTrans("zh", "user %s not found", "用户 %s 未找到")
i18n.AddTrans("en", "user %s not found", "User %s not found")
```

### 2. 定义错误模板

```go
errDef := i18n.DefineError("user %s not found").WithStatus(404, 404).WithTraceID("abc123")
```

### 3. 根据 ctx 构建错误

```go
err := errDef.New(ctx, "alice")
```

### 4. 获取错误信息

```go
fmt.Println(err.Error()) // 输出："用户 alice 未找到"
```

### 5. 提取错误信息

```go
if statusErr, ok := err.(interface {
	Code() int
	HttpStatus() int
}); ok {
	fmt.Println("Code:", statusErr.Code())
	fmt.Println("HTTP Status:", statusErr.HttpStatus())
}
```



## 🧰 错误包装器

| 包装器 | 说明 | 接口方法 |
|--------|------|----------|
| `WithStatus` | 添加错误码和 HTTP 状态码 | `Code()`, `HttpStatus()` |
| `WithTraceID` | 添加 TraceID | `TraceID()` |


## 📁 配置文件结构（示例）

```
locales/
└── zh/
    └── user.json
```

`user.json` 示例：

```json
{
  "user %s not found": "用户 %s 未找到"
}
```


## 🧪 与 Gin 集成示例

```go
func someHandler(c *gin.Context) {
	err := errDef.New(c.Request.Context(), "alice")
	c.AbortWithError(404, err)
}
```



## 📄 License

MIT License
