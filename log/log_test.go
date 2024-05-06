package log

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

type wrappingHandler struct {
	h slog.Handler
	l []slog.Record
}

func (h wrappingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}
func (h wrappingHandler) WithGroup(name string) slog.Handler    { return h.h.WithGroup(name) }
func (h wrappingHandler) WithAttrs(as []slog.Attr) slog.Handler { return h.h.WithAttrs(as) }
func (h wrappingHandler) Handle(ctx context.Context, r slog.Record) error {
	h.l = append(h.l, r)
	return h.h.Handle(ctx, r)
}

func TestHandler(t *testing.T) {
	h := NewHandler(slog.LevelDebug, slog.NewTextHandler(os.Stdout, nil))

	if h.Handler() != h.handler {
		t.Errorf("Handler() = %v, want %v", h.Handler(), h.handler)
	}
}

func TestWithAttrs(t *testing.T) {
	h := NewHandler(slog.LevelDebug, slog.NewTextHandler(os.Stdout, nil))
	attrs := []slog.Attr{slog.String("key", "value")}

	newHandler := h.WithAttrs(attrs)

	if _, ok := newHandler.(*Handler); !ok {
		t.Errorf("WithAttrs() should return a Handler")
	}
}

func TestWithGroup(t *testing.T) {
	h := NewHandler(slog.LevelDebug, slog.NewTextHandler(os.Stdout, nil))
	groupName := "testGroup"

	newHandler := h.WithGroup(groupName)

	if _, ok := newHandler.(*Handler); !ok {
		t.Errorf("WithGroup() should return a Handler")
	}
}

// Test logging with and wihout span. If the context does not have span,
// the log should be written to the handler. But should not fail.
func TestHandleWithoutSpan(t *testing.T) {
	lh := wrappingHandler{
		h: NewHandler(slog.LevelDebug, slog.NewTextHandler(os.Stdout, nil)),
		l: []slog.Record{}}

	log := slog.New(lh)

	ctxWithSpanAndRecording, span := observability.NewObserver("").CreateSpan(context.Background(), "test")
	defer span.End()

	ctxWithSpanAndNotRecording, spanNotRecording := observability.NewObserver("").CreateSpan(context.Background(), "test")
	spanNotRecording.End()

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{"no span", context.Background()},
		{"span", ctxWithSpanAndRecording},
		{"span not recording", ctxWithSpanAndNotRecording},
		{"span recording", ctxWithSpanAndRecording},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.InfoContext(tt.ctx, tt.name)
			log.ErrorContext(tt.ctx, tt.name)
			log.DebugContext(tt.ctx, tt.name)
			log.WarnContext(tt.ctx, tt.name)

			if len(lh.l) != 4 {
				t.Errorf("Handle() should have been called 4 times, got %d", len(lh.l))
			}
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
