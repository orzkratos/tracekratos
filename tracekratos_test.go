package tracekratos

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig("X-B3-TRACE-ID")
	t.Log(config.NewTraceID(context.Background()))

	type reqType struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	req := &reqType{}
	require.NoError(t, gofakeit.Struct(req))
	t.Log(config.FormatArgs(req))
}

func TestNewConfigWithOptions(t *testing.T) {
	config := NewConfig("X-Trace-ID",
		WithLogLevel(log.LevelDebug),
		WithNewTraceID(func(ctx context.Context) string {
			return "custom-trace-id"
		}),
	)
	require.Equal(t, log.LevelDebug, config.LogLevel)
	require.Equal(t, "custom-trace-id", config.NewTraceID(context.Background()))
}

func TestExtractArgs(t *testing.T) {
	type reqType struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	req := &reqType{}
	require.NoError(t, gofakeit.Struct(req))
	t.Log(ExtractArgs(req))
}

func TestGetTraceID(t *testing.T) {
	ctx := context.Background()
	require.Empty(t, GetTraceID(ctx))
	require.Empty(t, GetTraceIDFromContext(ctx))

	ctx = context.WithValue(ctx, traceCtxKey{}, "test-trace-id")
	require.Equal(t, "test-trace-id", GetTraceID(ctx))
	require.Equal(t, "test-trace-id", GetTraceIDFromContext(ctx))
}

func TestWithLogReply(t *testing.T) {
	config := NewConfig("X-Trace-ID", WithLogReply(true))
	require.True(t, config.LogReply)

	config2 := NewConfig("X-Trace-ID")
	require.False(t, config2.LogReply)
}
