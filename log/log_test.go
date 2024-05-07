package log

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type wrappingHandler struct {
	h slog.Handler
	C *int16
}

func (h wrappingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}
func (h wrappingHandler) WithGroup(name string) slog.Handler    { return h.h.WithGroup(name) }
func (h wrappingHandler) WithAttrs(as []slog.Attr) slog.Handler { return h.h.WithAttrs(as) }
func (h wrappingHandler) Handle(ctx context.Context, r slog.Record) error {
	*h.C++
	return h.h.Handle(ctx, r)
}

var opts = &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true,
}

func TestHandler(t *testing.T) {
	h := NewHandler(slog.NewTextHandler(os.Stdout, opts))

	if h.Handler() != h.handler {
		t.Errorf("Handler() = %v, want %v", h.Handler(), h.handler)
	}
}

func TestWithAttrs(t *testing.T) {
	h := NewHandler(slog.NewTextHandler(os.Stdout, opts))
	attrs := []slog.Attr{slog.String("key", "value")}

	newHandler := h.WithAttrs(attrs)

	if _, ok := newHandler.(*Handler); !ok {
		t.Errorf("WithAttrs() should return a Handler")
	}
}

func TestWithGroup(t *testing.T) {
	h := NewHandler(slog.NewTextHandler(os.Stdout, opts))
	groupName := "testGroup"

	newHandler := h.WithGroup(groupName)

	if _, ok := newHandler.(*Handler); !ok {
		t.Errorf("WithGroup() should return a Handler")
	}
}

// Test logging with and without span. If the context does not have span,
// the log should be written to the handler. But should not fail.
func TestLogWithSpan(t *testing.T) {
	i := int16(0)
	lh := wrappingHandler{
		h: NewHandler(slog.NewTextHandler(os.Stdout, opts)),
		C: &i,
	}

	o, _ := observability.NewObserver("")

	log := slog.New(lh)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctxWithSpanAndRecording, span := o.CreateSpan(context.Background(), "test")
	defer span.End()

	ctxWithSpanAndNotRecording, spanNotRecording := o.CreateSpan(context.Background(), "test")
	spanNotRecording.End()

	tests := []struct {
		name string
		ctx  context.Context
		span trace.Span
	}{
		{"no span", context.Background(), nil},
		{"span", ctxWithSpanAndRecording, observability.SpanFromContext(ctxWithSpanAndRecording)},
		{"span recording", ctxWithSpanAndRecording, observability.SpanFromContext(ctxWithSpanAndRecording)},
		{"span not recording", ctxWithSpanAndNotRecording, observability.SpanFromContext(ctxWithSpanAndNotRecording)},
		{"context is nil", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.InfoContext(tt.ctx, tt.name)
			log.DebugContext(tt.ctx, tt.name)
			log.WarnContext(tt.ctx, tt.name)
			log.ErrorContext(tt.ctx, tt.name)

			if *lh.C != 4 {
				t.Errorf("Handle() should have been called 4 times, got %v", lh.C)
			}
			i = 0
		})
	}

}
func TestMultiRoutines(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	o, _ :=  observability.NewObserver("")
	ctxWithSpanAndRecording, span := o.CreateSpan(context.Background(), "test")
	defer span.End()

	ctxWithSpanAndNotRecording, spanNotRecording := o.CreateSpan(context.Background(), "test")
	spanNotRecording.End()

	tests := []struct {
		name  string
		count int
		ctx   context.Context
		span  trace.Span
	}{
		// Existing test cases...
		{"no span", 100, context.Background(), nil},
		{"span", 100, ctxWithSpanAndRecording, observability.SpanFromContext(ctxWithSpanAndRecording)},
		{"span recording", 100, ctxWithSpanAndRecording, observability.SpanFromContext(ctxWithSpanAndRecording)},
		{"span not recording", 100, ctxWithSpanAndNotRecording, observability.SpanFromContext(ctxWithSpanAndNotRecording)},
		{"context is nil", 100, nil, nil},
	}
	h := NewHandler(slog.NewTextHandler(os.Stdout, opts))
	log := slog.New(h)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(tt.count)

			for i := 0; i < tt.count; i++ {
				go func() {
					defer wg.Done()
					log.InfoContext(tt.ctx, tt.name)
					log.DebugContext(tt.ctx, tt.name)
					log.WarnContext(tt.ctx, tt.name)
					log.ErrorContext(tt.ctx, tt.name)
				}()
			}

			wg.Wait()
		})
	}
}

func TestConvertSlogAttrToSpanAttr(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		attr    slog.Attr
		want    attribute.KeyValue
		wantErr bool
	}{
		{"string", "key", slog.String("key", "value"), attribute.String("key", "value"), false},
		{"bool", "key", slog.Bool("key", true), attribute.Bool("key", true), false},
		{"int64", "key", slog.Int64("key", 123), attribute.Int64("key", 123), false},
		{"uint64", "key", slog.Uint64("key", 123), attribute.Int64("key", 123), false},
		{"float64", "key", slog.Float64("key", 1.23), attribute.Float64("key", 1.23), false},
		{"duration", "key", slog.Duration("key", time.Second), attribute.String("key", "1s"), false},
		// {"time", "key", slog.Time("key", time.Unix(0, 0)), attribute.String("key", "1970-01-01 00:00:00 +0000 UTC"), false},
		// {"unsupported", "key", slog.Bytes("key", []byte("value")), attribute.String("key", "[118 97 108 117 101]"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertSlogAttrToSpanAttr(tt.key, tt.attr.Value)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertSlogAttrToSpanAttr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != tt.want.Value {
				t.Errorf("convertSlogAttrToSpanAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// Existing test cases...
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

	b.ResetTimer() // Reset the timer to exclude setup time

	for i := 0; i < b.N; i++ { // b.N is provided by the testing framework
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
