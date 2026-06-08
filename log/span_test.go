package log

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
)

// TestSpanEventCreation tests that span events are created with correct attributes.
func TestSpanEventCreation(t *testing.T) {
	observability.NewLocalObserver()

	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	log := slog.New(h)
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "test-span")
	defer span.End()

	tests := []struct {
		name     string
		logLevel slog.Level
		logFunc  func(ctx context.Context, msg string, args ...any)
		message  string
		attrs    []slog.Attr
	}{
		{
			name:     "Info log with span event",
			logLevel: slog.LevelInfo,
			logFunc:  log.InfoContext,
			message:  "Test info message",
			attrs:    []slog.Attr{slog.String("key", "value")},
		},
		{
			name:     "Debug log with span event",
			logLevel: slog.LevelDebug,
			logFunc:  log.DebugContext,
			message:  "Test debug message",
			attrs:    []slog.Attr{slog.Int("count", 42)},
		},
		{
			name:     "Warn log with span event",
			logLevel: slog.LevelWarn,
			logFunc:  log.WarnContext,
			message:  "Test warn message",
			attrs:    []slog.Attr{slog.Float64("value", 3.14)},
		},
		{
			name:     "Error log with span event",
			logLevel: slog.LevelError,
			logFunc:  log.ErrorContext,
			message:  "Test error message",
			attrs:    []slog.Attr{slog.Bool("critical", true)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			args := make([]any, 0, len(tt.attrs)*2)
			for _, attr := range tt.attrs {
				args = append(args, attr.Key, attr.Value.Any())
			}
			tt.logFunc(ctx, tt.message, args...)

			assert.Contains(t, buf.String(), tt.message)
			assert.True(t, span.IsRecording())

			t.Logf("Log level: %s", tt.logLevel.String())
			t.Logf("Message: %s", tt.message)
			t.Logf("Span is recording: %t", span.IsRecording())
		})
	}
}

// TestErrorRecordingOnSpan tests that errors are properly recorded on spans.
func TestErrorRecordingOnSpan(t *testing.T) {
	observability.NewLocalObserver()

	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	log := slog.New(h)
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "test-error-span")
	defer span.End()

	tests := []struct {
		name        string
		logLevel    slog.Level
		message     string
		errorAttr   slog.Attr
		expectError bool
	}{
		{
			name:        "Error log with error attribute",
			logLevel:    slog.LevelError,
			message:     "Database connection failed",
			errorAttr:   slog.String("error", "connection timeout"),
			expectError: true,
		},
		{
			name:        "Error log with actual error",
			logLevel:    slog.LevelError,
			message:     "Processing failed",
			errorAttr:   slog.Any("error", assert.AnError),
			expectError: true,
		},
		{
			name:        "Error log without error attribute",
			logLevel:    slog.LevelError,
			message:     "Something went wrong",
			errorAttr:   slog.String("other", "value"),
			expectError: true,
		},
		{
			name:        "Warning log should not record error",
			logLevel:    slog.LevelWarn,
			message:     "Warning message",
			errorAttr:   slog.String("warning", "non-critical"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			testCtx, testSpan := observer.CreateSpan(ctx, "test-"+tt.name)
			defer testSpan.End()

			switch tt.logLevel {
			case slog.LevelError:
				log.ErrorContext(testCtx, tt.message, tt.errorAttr)
			case slog.LevelWarn:
				log.WarnContext(testCtx, tt.message, tt.errorAttr)
			default:
				log.InfoContext(testCtx, tt.message, tt.errorAttr)
			}

			assert.Contains(t, buf.String(), tt.message)
			assert.True(t, testSpan.IsRecording())

			t.Logf("Test: %s", tt.name)
			t.Logf("Log level: %s", tt.logLevel.String())
			t.Logf("Message: %s", tt.message)
			t.Logf("Should record error: %t", tt.expectError)
		})
	}
}

// TestSpanAttributeConversion tests conversion of slog attributes to span attributes.
func TestSpanAttributeConversion(t *testing.T) {
	observability.NewLocalObserver()

	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	log := slog.New(h)
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "test-attributes")
	defer span.End()

	tests := []struct {
		name  string
		attrs []slog.Attr
	}{
		{
			name: "String attributes",
			attrs: []slog.Attr{
				slog.String("service", "panchangam"),
				slog.String("version", "1.0.0"),
			},
		},
		{
			name: "Numeric attributes",
			attrs: []slog.Attr{
				slog.Int("count", 42),
				slog.Int64("id", 123456789),
				slog.Float64("ratio", 0.75),
			},
		},
		{
			name: "Boolean and time attributes",
			attrs: []slog.Attr{
				slog.Bool("enabled", true),
				slog.Duration("elapsed", time.Second*5),
				slog.Time("timestamp", time.Now()),
			},
		},
		{
			name: "Mixed attributes",
			attrs: []slog.Attr{
				slog.String("operation", "calculate"),
				slog.Int("retries", 3),
				slog.Bool("success", false),
				slog.Float64("latitude", 40.7128),
				slog.Float64("longitude", -74.0060),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			args := make([]any, 0, len(tt.attrs)*2)
			for _, attr := range tt.attrs {
				args = append(args, attr.Key, attr.Value.Any())
			}
			log.InfoContext(ctx, "Test message with attributes", args...)

			assert.Contains(t, buf.String(), "Test message with attributes")
			assert.True(t, span.IsRecording())

			t.Logf("Test: %s", tt.name)
			t.Logf("Number of attributes: %d", len(tt.attrs))

			for _, attr := range tt.attrs {
				spanAttr, err := convertSlogAttrToSpanAttr(attr.Key, attr.Value)
				assert.NoError(t, err, "Failed to convert attribute %s", attr.Key)
				assert.True(t, spanAttr.Valid(), "Span attribute should be valid for %s", attr.Key)

				t.Logf("  %s: %v -> %v", attr.Key, attr.Value.Any(), spanAttr.Value.AsString())
			}
		})
	}
}
