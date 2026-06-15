package log

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
)

// TestLoggerWithoutSpan tests that logging works normally without a span.
func TestLoggerWithoutSpan(t *testing.T) {
	observability.NewLocalObserver()

	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	log := slog.New(h)
	ctx := context.Background()

	tests := []struct {
		name     string
		logLevel slog.Level
		message  string
	}{
		{"Info without span", slog.LevelInfo, "Info message without span"},
		{"Debug without span", slog.LevelDebug, "Debug message without span"},
		{"Warn without span", slog.LevelWarn, "Warn message without span"},
		{"Error without span", slog.LevelError, "Error message without span"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			switch tt.logLevel {
			case slog.LevelInfo:
				log.InfoContext(ctx, tt.message)
			case slog.LevelDebug:
				log.DebugContext(ctx, tt.message)
			case slog.LevelWarn:
				log.WarnContext(ctx, tt.message)
			case slog.LevelError:
				log.ErrorContext(ctx, tt.message)
			}

			assert.Contains(t, buf.String(), tt.message)

			t.Logf("Test: %s", tt.name)
			t.Logf("Message: %s", tt.message)
			t.Logf("Log output: %s", buf.String())
		})
	}
}

// TestLoggerWithNilContext tests that logging works with nil context.
func TestLoggerWithNilContext(t *testing.T) {
	observability.NewLocalObserver()

	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	log := slog.New(h)
	var ctx context.Context

	tests := []struct {
		name     string
		logLevel slog.Level
		message  string
	}{
		{"Info with nil context", slog.LevelInfo, "Info message with nil context"},
		{"Error with nil context", slog.LevelError, "Error message with nil context"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			switch tt.logLevel {
			case slog.LevelInfo:
				log.InfoContext(ctx, tt.message)
			case slog.LevelError:
				log.ErrorContext(ctx, tt.message)
			}

			assert.Contains(t, buf.String(), tt.message)

			t.Logf("Test: %s", tt.name)
			t.Logf("Message: %s", tt.message)
			t.Logf("Successfully handled nil context")
		})
	}
}
