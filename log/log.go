package log

import (
	"context"
	"fmt"
	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"os"
	"sync"
	"time"
)

var logger *slog.Logger
var initOnce sync.Once
var spanEnabled = true

func init() {
	initOnce.Do(func() {
		logger = slog.New(NewHandler(slog.LevelDebug,
			slog.NewTextHandler(os.Stdout, nil)))
	})
}

func Logger() *slog.Logger {
	return logger
}

// A Handler wraps a Handler with an Enabled method
// that returns false for levels below a minimum.
type Handler struct {
	level   slog.Leveler
	handler slog.Handler
}

// NewHandler returns a LevelHandler with the given level.
// All methods except Enabled delegate to h.
func NewHandler(level slog.Leveler, h slog.Handler) *Handler {
	// Optimization: avoid chains of LevelHandlers.
	if lh, ok := h.(*Handler); ok {
		h = lh.Handler()
	}
	return &Handler{level, h}
}

// Enabled implements Handler.Enabled by reporting whether
// level is at least as large as h's level.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle implements Handler.Handle.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil && spanEnabled {
		span := observability.SpanFromContext(ctx)
		if !span.IsRecording() {
			return h.handler.Handle(ctx, r)
		}
		span.AddEvent(r.Message)

	}

	return h.handler.Handle(ctx, r)
}

func ConvertSlogAttrToSpanAttr(key string, attr slog.Value) (attribute.KeyValue, error) {
	var kv attribute.KeyValue
	switch attr.Kind() {
	case slog.KindString:
		kv = attribute.String(key, attr.Any().(string))
	case slog.KindBool:
		kv = attribute.Bool(key, attr.Any().(bool))
	case slog.KindInt64:
		kv = attribute.Int64(key, attr.Any().(int64))
	case slog.KindUint64:
		// OpenTelemetry does not support Uint64 directly, convert to Int64
		kv = attribute.Int64(key, int64(attr.Any().(uint64)))
	case slog.KindFloat64:
		kv = attribute.Float64(key, attr.Any().(float64))
	case slog.KindDuration:
		kv = attribute.String(key, attr.Any().(time.Duration).String())
	case slog.KindTime:
		kv = attribute.String(key, attr.Any().(time.Time).String())
	default:
		// For unsupported types, or in case of any errors, encode as a string
		kv = attribute.String(key, fmt.Sprint(attr.Any()))
	}

	if !kv.Valid() {
		return kv, fmt.Errorf("invalid attribute.KeyValue: %v", kv)
	}

	return kv, nil
}

// Handler returns the Handler wrapped by h.
func (h *Handler) Handler() slog.Handler {
	return h.handler
}

// WithAttrs implements Handler.WithAttrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewHandler(h.level, h.handler.WithAttrs(attrs))
}

// WithGroup implements Handler.WithGroup.
func (h *Handler) WithGroup(name string) slog.Handler {
	return NewHandler(h.level, h.handler.WithGroup(name))
}
