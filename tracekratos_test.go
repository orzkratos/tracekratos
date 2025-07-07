package tracekratos

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
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

func TestExtractArgs(t *testing.T) {
	type reqType struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	req := &reqType{}
	must.Done(gofakeit.Struct(req))
	t.Log(ExtractArgs(req))
}
