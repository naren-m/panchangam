package log

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
)

// TestErrorCorrelationAndAggregation tests error correlation across multiple log entries.
func TestErrorCorrelationAndAggregation(t *testing.T) {
	observability.NewLocalObserver()

	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	log := slog.New(h)
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "error-correlation-test")
	defer span.End()

	errors := []struct {
		message string
		level   slog.Level
		attrs   []slog.Attr
	}{
		{
			message: "Database connection failed",
			level:   slog.LevelError,
			attrs: []slog.Attr{
				slog.String("error", "connection timeout"),
				slog.String("operation", "db_connect"),
				slog.Int("attempt", 1),
			},
		},
		{
			message: "Retrying database connection",
			level:   slog.LevelWarn,
			attrs: []slog.Attr{
				slog.String("operation", "db_connect"),
				slog.Int("attempt", 2),
			},
		},
		{
			message: "Database connection retry failed",
			level:   slog.LevelError,
			attrs: []slog.Attr{
				slog.String("error", "max retries exceeded"),
				slog.String("operation", "db_connect"),
				slog.Int("attempt", 2),
			},
		},
		{
			message: "Fallback to read-only mode",
			level:   slog.LevelInfo,
			attrs: []slog.Attr{
				slog.String("operation", "db_connect"),
				slog.String("mode", "read-only"),
			},
		},
	}

	for i, err := range errors {
		t.Run(fmt.Sprintf("Error_%d", i+1), func(t *testing.T) {
			buf.Reset()

			args := make([]any, 0, len(err.attrs)*2)
			for _, attr := range err.attrs {
				args = append(args, attr.Key, attr.Value.Any())
			}

			switch err.level {
			case slog.LevelError:
				log.ErrorContext(ctx, err.message, args...)
			case slog.LevelWarn:
				log.WarnContext(ctx, err.message, args...)
			case slog.LevelInfo:
				log.InfoContext(ctx, err.message, args...)
			}

			assert.Contains(t, buf.String(), err.message)
			assert.True(t, span.IsRecording())

			t.Logf("Step %d: %s", i+1, err.message)
			t.Logf("Level: %s", err.level.String())
			t.Logf("Attributes: %d", len(err.attrs))
		})
	}

	t.Logf("All related errors logged in span: %s", "error-correlation-test")
}
