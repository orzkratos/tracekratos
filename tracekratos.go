package tracekratos

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport"
)

type Config struct {
	TraceKeyName string
	NewTraceID   func(context.Context) string
	FormatArgs   func(req any) string
}

func NewConfig(keyName string) *Config {
	return &Config{
		TraceKeyName: keyName,
		NewTraceID: func(ctx context.Context) string {
			return strconv.FormatInt(time.Now().UnixNano(), 10)
		},
		FormatArgs: func(req any) string {
			return "((" + strings.ReplaceAll(ExtractArgs(req), `"`, `'`) + "))"
		},
	}
}

func NewTraceMiddleware(config *Config, logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if tsp, ok := transport.FromServerContext(ctx); ok {
				traceID := tsp.RequestHeader().Get(config.TraceKeyName)
				if traceID == "" {
					traceID = config.NewTraceID(ctx)
					tsp.RequestHeader().Set(config.TraceKeyName, traceID)
				}

				kind := tsp.Kind().String()
				operation := tsp.Operation()

				startTime := time.Now()
				//cp from: https://github.com/go-kratos/kratos/blob/15dd2f638e3d53d059913ca83818f5843d67a277/middleware/logging/logging.go#L47
				//cp from: https://github.com/go-kratos/kratos/blob/15dd2f638e3d53d059913ca83818f5843d67a277/middleware/logging/logging.go#L87
				log.NewHelper(log.WithContext(ctx, logger)).Log(log.LevelInfo,
					"kind", "server",
					"component", kind,
					"operation", operation,
					"args", config.FormatArgs(req),
					"trace", traceID,
					"startTime", startTime.Format(time.RFC3339Nano),
				)
			}
			return handler(ctx, req)
		}
	}
}

// ExtractArgs returns the string of the req-param.
// cp from: https://github.com/go-kratos/kratos/blob/15dd2f638e3d53d059913ca83818f5843d67a277/middleware/logging/logging.go#L102
func ExtractArgs(req any) string {
	if redacter, ok := req.(logging.Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
