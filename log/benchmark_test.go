package log

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func BenchmarkLogging(b *testing.B) {
	o, _ := observability.NewObserver("")
	assert.NotNil(b, o)
	ctxWithSpanAndRecording, span := o.CreateSpan(context.Background(), "test")
	defer span.End()

	ctxWithSpanAndNotRecording, spanNotRecording := o.CreateSpan(context.Background(), "Context for non recording span")
	spanNotRecording.End()

	tests := []struct {
		name  string
		count int
		ctx   context.Context
		span  trace.Span
	}{
		{"no span", 100, context.Background(), nil},
		{"span", 100, ctxWithSpanAndRecording, observability.SpanFromContext(ctxWithSpanAndRecording)},
		{"span recording", 100, ctxWithSpanAndRecording, observability.SpanFromContext(ctxWithSpanAndRecording)},
		{"span not recording", 100, ctxWithSpanAndNotRecording, observability.SpanFromContext(ctxWithSpanAndNotRecording)},
		{"context is nil", 100, nil, nil},
	}
	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	log := slog.New(h)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tt := range tests {
			for j := 0; j < tt.count; j++ {
				log.InfoContext(tt.ctx, tt.name)
				log.DebugContext(tt.ctx, tt.name)
				log.WarnContext(tt.ctx, tt.name)
				log.ErrorContext(tt.ctx, tt.name)
			}
			buf.Reset()
		}
	}
}
