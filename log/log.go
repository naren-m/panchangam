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

func init() {
	initOnce.Do(func() {
		logger = slog.New(NewHandler(slog.NewTextHandler(os.Stdout, nil)))
	})
}

func Logger() *slog.Logger {
	return logger
}

// A Handler wraps a Handler with an Enabled method
// that returns false for levels below a minimum.
type Handler struct {
	handler slog.Handler
}

// NewHandler returns a LevelHandler with the given level.
// All methods except Enabled delegate to h.
func NewHandler(h slog.Handler) *Handler {
	// Optimization: avoid chains of LevelHandlers.
	if lh, ok := h.(*Handler); ok {
		h = lh.Handler()
	}
	return &Handler{h}
}

// Enabled implements Handler.Enabled by reporting whether
// level is at least as large as h's level.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements Handler.Handle.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		span := observability.SpanFromContext(ctx)
		if span != nil && span.IsRecording() {
			// Convert slog attributes to OpenTelemetry attributes
			var spanAttrs []attribute.KeyValue
			r.Attrs(func(attr slog.Attr) bool {
				if spanAttr, err := convertSlogAttrToSpanAttr(attr.Key, attr.Value); err == nil {
					spanAttrs = append(spanAttrs, spanAttr)
				}
				return true
			})
			
			// Add log level as span attribute
			spanAttrs = append(spanAttrs, attribute.String("log.level", r.Level.String()))
			
			// Create span event with attributes
			eventName := fmt.Sprintf("log.%s", r.Level.String())
			span.AddEvent(eventName, observability.WithAttributes(spanAttrs...))
			
			// For errors, also record the error on the span
			if r.Level >= slog.LevelError {
				// Try to extract error from attributes
				var errorAttr slog.Attr
				r.Attrs(func(attr slog.Attr) bool {
					if attr.Key == "error" {
						errorAttr = attr
						return false
					}
					return true
				})
				
				if errorAttr.Key != "" {
					if err, ok := errorAttr.Value.Any().(error); ok {
						span.RecordError(err)
					} else {
						// Create a synthetic error from the error attribute
						span.RecordError(fmt.Errorf("%v", errorAttr.Value.Any()))
					}
				} else {
					// Create a synthetic error from the log message
					span.RecordError(fmt.Errorf("%s", r.Message))
				}
			}
		}
	}

	return h.handler.Handle(ctx, r)
}

func convertSlogAttrToSpanAttr(key string, attr slog.Value) (attribute.KeyValue, error) {
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
func (h *Handler) Handler() slog.Handler { return h.handler }

// WithAttrs implements Handler.WithAttrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewHandler(h.handler.WithAttrs(attrs))
}

// WithGroup implements Handler.WithGroup.
func (h *Handler) WithGroup(name string) slog.Handler {
	return NewHandler(h.handler.WithGroup(name))
}
