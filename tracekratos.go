package tracekratos

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport"
)

// Config trace middleware config
// 追踪中间件配置
type Config struct {
	TraceKeyName string                       // HTTP head name to get trace ID // 请求头键名
	NewTraceID   func(context.Context) string // Generate new trace ID // 生成新的 trace ID
	FormatArgs   func(req any) string         // Format request args // 格式化请求参数
	LogLevel     log.Level                    // Log config // 日志级别
	LogReply     bool                         // Log response and cost // 记录响应和耗时
}

// Option is a function to set Config options
// Config 的函数式选项
type Option func(*Config)

// WithLogLevel set log config
// 设置日志级别
func WithLogLevel(level log.Level) Option {
	return func(c *Config) {
		c.LogLevel = level
	}
}

// WithNewTraceID set custom trace ID creation function
// 设置自定义 trace ID 生成函数
func WithNewTraceID(fn func(context.Context) string) Option {
	return func(c *Config) {
		c.NewTraceID = fn
	}
}

// WithFormatArgs set custom args format function
// 设置自定义参数格式化函数
func WithFormatArgs(fn func(req any) string) Option {
	return func(c *Config) {
		c.FormatArgs = fn
	}
}

// WithLogReply enable response and cost logging
// 启用响应日志和耗时记录
func WithLogReply(enable bool) Option {
	return func(c *Config) {
		c.LogReply = enable
	}
}

// NewConfig create config with trace name and options
// 使用追踪键名和选项创建配置
func NewConfig(keyName string, opts ...Option) *Config {
	cfg := &Config{
		TraceKeyName: keyName,
		NewTraceID: func(ctx context.Context) string {
			return strconv.FormatInt(time.Now().UnixNano(), 10)
		},
		FormatArgs: func(req any) string {
			return "((" + strings.ReplaceAll(ExtractArgs(req), `"`, `'`) + "))"
		},
		LogLevel: log.LevelInfo,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// NewTraceMiddleware create trace middleware
// 创建追踪中间件
func NewTraceMiddleware(config *Config, logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if tsp, ok := transport.FromServerContext(ctx); ok {
				traceID := tsp.RequestHeader().Get(config.TraceKeyName)
				if traceID == "" {
					traceID = config.NewTraceID(ctx)
					tsp.RequestHeader().Set(config.TraceKeyName, traceID)
				}

				// Store trace ID in context
				// 将 trace ID 存入 context
				ctx = context.WithValue(ctx, traceCtxKey{}, traceID)

				kind := tsp.Kind().String()
				operation := tsp.Operation()

				startTime := time.Now()
				// cp from: https://github.com/go-kratos/kratos/blob/15dd2f638e3d53d059913ca83818f5843d67a277/middleware/logging/logging.go#L47
				// cp from: https://github.com/go-kratos/kratos/blob/15dd2f638e3d53d059913ca83818f5843d67a277/middleware/logging/logging.go#L87
				log.NewHelper(log.WithContext(ctx, logger)).Log(config.LogLevel,
					"kind", "server",
					"component", kind,
					"operation", operation,
					"args", config.FormatArgs(req),
					"trace", traceID,
					"startTime", startTime.Format(time.RFC3339Nano),
				)

				reply, err = handler(ctx, req)

				// Log response if enabled
				// 如果启用则记录响应
				if config.LogReply {
					latency := time.Since(startTime)
					code := 0
					if erk := errors.FromError(err); erk != nil {
						code = int(erk.Code)
					}
					log.NewHelper(log.WithContext(ctx, logger)).Log(config.LogLevel,
						"kind", "server",
						"component", kind,
						"operation", operation,
						"trace", traceID,
						"code", code,
						"latency", latency.String(),
					)
				}
				return reply, err
			}
			return handler(ctx, req)
		}
	}
}

// traceCtxKey is the context storage name to get trace ID
// context 存储 trace ID 的类型
type traceCtxKey struct{}

// GetTraceIDFromContext get trace ID from context
// 从 context 中获取 trace ID
func GetTraceIDFromContext(ctx context.Context) string {
	if v := ctx.Value(traceCtxKey{}); v != nil {
		if traceID, ok := v.(string); ok {
			return traceID
		}
	}
	return ""
}

// GetTraceID alias to GetTraceIDFromContext
// GetTraceIDFromContext 的简写
func GetTraceID(ctx context.Context) string {
	return GetTraceIDFromContext(ctx)
}

// ExtractArgs returns the string of the req-param
// cp from: https://github.com/go-kratos/kratos/blob/15dd2f638e3d53d059913ca83818f5843d67a277/middleware/logging/logging.go#L102
// 提取请求参数的字符串表示
func ExtractArgs(req any) string {
	if redacter, ok := req.(logging.Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
