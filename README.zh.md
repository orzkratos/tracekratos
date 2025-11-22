[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/orzkratos/tracekratos/release.yml?branch=main&label=BUILD)](https://github.com/orzkratos/tracekratos/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/orzkratos/tracekratos)](https://pkg.go.dev/github.com/orzkratos/tracekratos)
[![Coverage Status](https://img.shields.io/coveralls/github/orzkratos/tracekratos/main.svg)](https://coveralls.io/github/orzkratos/tracekratos?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/orzkratos/tracekratos.svg)](https://github.com/orzkratos/tracekratos/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/orzkratos/tracekratos)](https://goreportcard.com/report/github.com/orzkratos/tracekratos)

# tracekratos

Trace ID 中间件，在 Kratos 框架的日志中显示 trace ID，提供请求追踪和调试功能。

---

## 英文文档

[ENGLISH README](README.md)

## 主要功能

- 🔍 Trace ID 日志 - 在请求日志中显示 trace ID
- 🚀 自动生成 - 未提供时自动生成 trace ID
- ⚙️ 灵活配置 - 使用自定义选项构建配置
- 📊 响应日志 - 可选的响应和耗时日志
- 🎯 Context 访问 - 在业务代码中从 context 获取 trace ID

## 安装

```bash
go get github.com/orzkratos/tracekratos
```

## 快速开始

```go
package main

import (
    "context"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/google/uuid"
    "github.com/orzkratos/tracekratos"
)

func main() {
    // 创建追踪配置，使用所有选项
    config := tracekratos.NewConfig("X-Trace-ID",
        tracekratos.WithLogLevel(log.LevelDebug),
        tracekratos.WithLogReply(true),
        tracekratos.WithNewTraceID(func(ctx context.Context) string {
            return uuid.New().String()
        }),
        tracekratos.WithFormatArgs(func(req any) string {
            return fmt.Sprintf("%+v", req)
        }),
    )

    // 创建追踪中间件
    middleware := tracekratos.NewTraceMiddleware(config, log.DefaultLogger)

    // 在 Kratos 服务中使用
    // httpSrv := http.NewServer(
    //     http.Middleware(middleware),
    // )
}
```

## 高级用法

### 自定义 Trace ID 生成

```go
tracekratos.WithNewTraceID(func(ctx context.Context) string {
    return uuid.New().String()
})
```

### 启用响应日志

```go
tracekratos.WithLogReply(true)
```

### 自定义日志配置

```go
tracekratos.WithLogLevel(log.LevelDebug)
```

### 自定义参数格式化

```go
tracekratos.WithFormatArgs(func(req any) string {
    return fmt.Sprintf("%+v", req)
})
```

### 在业务代码中获取 Trace ID

```go
func (s *Service) DoSomething(ctx context.Context) {
    traceID := tracekratos.GetTraceID(ctx)
    // 或者
    traceID := tracekratos.GetTraceIDFromContext(ctx)
}
```

## 完整示例

查看 [tracekratos-demos](https://github.com/orzkratos/tracekratos-demos) 查看在实际 Kratos 项目中的完整集成：

- **[demo1kratos](https://github.com/orzkratos/tracekratos-demos/tree/main/demo1kratos)** - HTTP 和 gRPC 基础集成
- **[demo2kratos](https://github.com/orzkratos/tracekratos-demos/tree/main/demo2kratos)** - 使用 Wire DI 的高级用法

## API 参考

### Config

追踪中间件的配置结构体。

```go
type Config struct {
    TraceKeyName string                       // 获取 trace ID 的 HTTP 头名称
    NewTraceID   func(context.Context) string // 生成新的 trace ID
    FormatArgs   func(req any) string         // 格式化请求参数
    LogLevel     log.Level                    // 日志配置
    LogReply     bool                         // 记录响应和耗时
}
```

### 选项函数

```go
func WithLogLevel(level log.Level) Option      // 设置日志配置
func WithNewTraceID(fn func(context.Context) string) Option // 设置自定义 trace ID 生成函数
func WithFormatArgs(fn func(req any) string) Option // 设置自定义参数格式化函数
func WithLogReply(enable bool) Option          // 启用响应和耗时日志
```

### 函数

```go
func NewConfig(keyName string, opts ...Option) *Config // 使用追踪键名和选项创建配置
func NewTraceMiddleware(config *Config, logger log.Logger) middleware.Middleware // 创建追踪中间件
func GetTraceID(ctx context.Context) string // 从 context 获取 trace ID（简写）
func GetTraceIDFromContext(ctx context.Context) string // 从 context 获取 trace ID
func ExtractArgs(req any) string // 提取请求参数的字符串表示
```

## 许可证

MIT License - 查看 [LICENSE](LICENSE) 文件

---

## GitHub Stars

[![Stargazers](https://starchart.cc/orzkratos/tracekratos.svg?variant=adaptive)](https://starchart.cc/orzkratos/tracekratos)
