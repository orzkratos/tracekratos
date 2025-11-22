[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/orzkratos/tracekratos/release.yml?branch=main&label=BUILD)](https://github.com/orzkratos/tracekratos/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/orzkratos/tracekratos)](https://pkg.go.dev/github.com/orzkratos/tracekratos)
[![Coverage Status](https://img.shields.io/coveralls/github/orzkratos/tracekratos/main.svg)](https://coveralls.io/github/orzkratos/tracekratos?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/orzkratos/tracekratos.svg)](https://github.com/orzkratos/tracekratos/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/orzkratos/tracekratos)](https://goreportcard.com/report/github.com/orzkratos/tracekratos)

# tracekratos

Trace ID middleware that shows trace ID in logs using Kratos framework, providing request tracking and debugging.

---

## CHINESE README

[中文说明](README.zh.md)

## Main Features

- 🔍 Trace ID Logging - Show trace ID in request logs
- 🚀 Auto Generation - Auto generate trace ID when not provided
- ⚙️ Flexible Config - Build config with custom options
- 📊 Response Logging - Response and cost logging on demand
- 🎯 Context Access - Get trace ID from context in business code

## Installation

```bash
go get github.com/orzkratos/tracekratos
```

## Quick Start

```go
package main

import (
    "context"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/google/uuid"
    "github.com/orzkratos/tracekratos"
)

func main() {
    // Create trace config with all options
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

    // Create trace middleware
    middleware := tracekratos.NewTraceMiddleware(config, log.DefaultLogger)

    // Use in Kratos server
    // httpSrv := http.NewServer(
    //     http.Middleware(middleware),
    // )
}
```

## Advanced Usage

### Custom Trace ID Generation

```go
tracekratos.WithNewTraceID(func(ctx context.Context) string {
    return uuid.New().String()
})
```

### Enable Response Logging

```go
tracekratos.WithLogReply(true)
```

### Custom Log Config

```go
tracekratos.WithLogLevel(log.LevelDebug)
```

### Custom Args Format

```go
tracekratos.WithFormatArgs(func(req any) string {
    return fmt.Sprintf("%+v", req)
})
```

### Get Trace ID in Business Code

```go
func (s *Service) DoSomething(ctx context.Context) {
    traceID := tracekratos.GetTraceID(ctx)
    // or
    traceID := tracekratos.GetTraceIDFromContext(ctx)
}
```

## Complete Examples

See [tracekratos-demos](https://github.com/orzkratos/tracekratos-demos) to view complete integration in Kratos projects:

- **[demo1kratos](https://github.com/orzkratos/tracekratos-demos/tree/main/demo1kratos)** - Basic integration with HTTP and gRPC
- **[demo2kratos](https://github.com/orzkratos/tracekratos-demos/tree/main/demo2kratos)** - Advanced usage with Wire DI

## API Reference

### Config

Config struct to trace middleware.

```go
type Config struct {
    TraceKeyName string                       // HTTP head name to get trace ID
    NewTraceID   func(context.Context) string // Generate new trace ID
    FormatArgs   func(req any) string         // Format request args
    LogLevel     log.Level                    // Log config
    LogReply     bool                         // Log response and cost
}
```

### Options

```go
func WithLogLevel(level log.Level) Option      // Set log config
func WithNewTraceID(fn func(context.Context) string) Option // Set custom trace ID creation function
func WithFormatArgs(fn func(req any) string) Option // Set custom args format function
func WithLogReply(enable bool) Option          // Enable response and cost logging
```

### Functions

```go
func NewConfig(keyName string, opts ...Option) *Config // Create config with trace name and options
func NewTraceMiddleware(config *Config, logger log.Logger) middleware.Middleware // Create trace middleware
func GetTraceID(ctx context.Context) string // Get trace ID from context (short name)
func GetTraceIDFromContext(ctx context.Context) string // Get trace ID from context
func ExtractArgs(req any) string // Extract request args to string
```

## License

MIT License - see [LICENSE](LICENSE) file

---

## GitHub Stars

[![Stargazers](https://starchart.cc/orzkratos/tracekratos.svg?variant=adaptive)](https://starchart.cc/orzkratos/tracekratos)
