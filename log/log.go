package log

import (
	"context"
	"log/slog"
	"os"
	"github.com/vincentfree/opentelemetry/otelslog"
)
var th = slog.NewTextHandler(os.Stdout, nil)

var Logger = slog.New(NewHandler(slog.LevelDebug, th))


var AddTracingContext = otelslog.AddTracingContext

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
	slog.Info(r.Message, "level", r.Level, "name", "naren")
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewHandler(h.level, h.handler.WithAttrs(attrs))
}

// WithGroup implements Handler.WithGroup.
func (h *Handler) WithGroup(name string) slog.Handler {
	return NewHandler(h.level, h.handler.WithGroup(name))
}

// Handler returns the Handler wrapped by h.
func (h *Handler) Handler() slog.Handler {
	return h.handler
}
