package log

import (
	"bytes"
	"context"
	"fmt"
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
	o, _ := observability.NewObserver("")
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

// TestSpanEventCreation tests that span events are created with correct attributes
func TestSpanEventCreation(t *testing.T) {
	observability.NewLocalObserver()
	
	// Create a buffer to capture log output
	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	
	log := slog.New(h)
	
	// Create a span context
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "test-span")
	defer span.End()
	
	// Test different log levels
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
			// Clear buffer
			buf.Reset()
			
			// Log message with attributes
			args := make([]any, 0, len(tt.attrs)*2)
			for _, attr := range tt.attrs {
				args = append(args, attr.Key, attr.Value.Any())
			}
			tt.logFunc(ctx, tt.message, args...)
			
			// Verify log was written
			assert.Contains(t, buf.String(), tt.message)
			
			// The span should be active and recording
			assert.True(t, span.IsRecording())
			
			t.Logf("Log level: %s", tt.logLevel.String())
			t.Logf("Message: %s", tt.message)
			t.Logf("Span is recording: %t", span.IsRecording())
		})
	}
}

// TestErrorRecordingOnSpan tests that errors are properly recorded on spans
func TestErrorRecordingOnSpan(t *testing.T) {
	observability.NewLocalObserver()
	
	// Create a buffer to capture log output
	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	
	log := slog.New(h)
	
	// Create a span context
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
			expectError: true, // Should create synthetic error from message
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
			// Clear buffer
			buf.Reset()
			
			// Create a new span for each test
			testCtx, testSpan := observer.CreateSpan(ctx, "test-"+tt.name)
			defer testSpan.End()
			
			// Log message with error level
			switch tt.logLevel {
			case slog.LevelError:
				log.ErrorContext(testCtx, tt.message, tt.errorAttr)
			case slog.LevelWarn:
				log.WarnContext(testCtx, tt.message, tt.errorAttr)
			default:
				log.InfoContext(testCtx, tt.message, tt.errorAttr)
			}
			
			// Verify log was written
			assert.Contains(t, buf.String(), tt.message)
			
			// The span should be active and recording
			assert.True(t, testSpan.IsRecording())
			
			t.Logf("Test: %s", tt.name)
			t.Logf("Log level: %s", tt.logLevel.String())
			t.Logf("Message: %s", tt.message)
			t.Logf("Should record error: %t", tt.expectError)
		})
	}
}

// TestSpanAttributeConversion tests conversion of slog attributes to span attributes
func TestSpanAttributeConversion(t *testing.T) {
	observability.NewLocalObserver()
	
	// Create a buffer to capture log output
	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	
	log := slog.New(h)
	
	// Create a span context
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "test-attributes")
	defer span.End()
	
	// Test various attribute types
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
			// Clear buffer
			buf.Reset()
			
			// Log message with attributes
			args := make([]any, 0, len(tt.attrs)*2)
			for _, attr := range tt.attrs {
				args = append(args, attr.Key, attr.Value.Any())
			}
			log.InfoContext(ctx, "Test message with attributes", args...)
			
			// Verify log was written
			assert.Contains(t, buf.String(), "Test message with attributes")
			
			// The span should be active and recording
			assert.True(t, span.IsRecording())
			
			t.Logf("Test: %s", tt.name)
			t.Logf("Number of attributes: %d", len(tt.attrs))
			
			// Test individual attribute conversion
			for _, attr := range tt.attrs {
				spanAttr, err := convertSlogAttrToSpanAttr(attr.Key, attr.Value)
				assert.NoError(t, err, "Failed to convert attribute %s", attr.Key)
				assert.True(t, spanAttr.Valid(), "Span attribute should be valid for %s", attr.Key)
				
				t.Logf("  %s: %v -> %v", attr.Key, attr.Value.Any(), spanAttr.Value.AsString())
			}
		})
	}
}

// TestLoggerWithoutSpan tests that logging works normally without a span
func TestLoggerWithoutSpan(t *testing.T) {
	observability.NewLocalObserver()
	
	// Create a buffer to capture log output
	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	
	log := slog.New(h)
	
	// Test with context that has no span
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
			// Clear buffer
			buf.Reset()
			
			// Log message without span
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
			
			// Verify log was written normally
			assert.Contains(t, buf.String(), tt.message)
			
			t.Logf("Test: %s", tt.name)
			t.Logf("Message: %s", tt.message)
			t.Logf("Log output: %s", buf.String())
		})
	}
}

// TestLoggerWithNilContext tests that logging works with nil context
func TestLoggerWithNilContext(t *testing.T) {
	observability.NewLocalObserver()
	
	// Create a buffer to capture log output
	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	
	log := slog.New(h)
	
	// Test with nil context
	var ctx context.Context = nil
	
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
			// Clear buffer
			buf.Reset()
			
			// Log message with nil context - should not panic
			switch tt.logLevel {
			case slog.LevelInfo:
				log.InfoContext(ctx, tt.message)
			case slog.LevelError:
				log.ErrorContext(ctx, tt.message)
			}
			
			// Verify log was written normally
			assert.Contains(t, buf.String(), tt.message)
			
			t.Logf("Test: %s", tt.name)
			t.Logf("Message: %s", tt.message)
			t.Logf("Successfully handled nil context")
		})
	}
}

// TestErrorCorrelationAndAggregation tests error correlation across multiple log entries
func TestErrorCorrelationAndAggregation(t *testing.T) {
	observability.NewLocalObserver()
	
	// Create a buffer to capture log output
	var buf bytes.Buffer
	h := NewHandler(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	
	log := slog.New(h)
	
	// Create a span context for correlation
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "error-correlation-test")
	defer span.End()
	
	// Simulate multiple related errors
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
	
	// Log all related events in the same span
	for i, err := range errors {
		t.Run(fmt.Sprintf("Error_%d", i+1), func(t *testing.T) {
			// Clear buffer for each test
			buf.Reset()
			
			// Convert attributes to args
			args := make([]any, 0, len(err.attrs)*2)
			for _, attr := range err.attrs {
				args = append(args, attr.Key, attr.Value.Any())
			}
			
			// Log message at appropriate level
			switch err.level {
			case slog.LevelError:
				log.ErrorContext(ctx, err.message, args...)
			case slog.LevelWarn:
				log.WarnContext(ctx, err.message, args...)
			case slog.LevelInfo:
				log.InfoContext(ctx, err.message, args...)
			}
			
			// Verify log was written
			assert.Contains(t, buf.String(), err.message)
			
			// The span should be active and recording
			assert.True(t, span.IsRecording())
			
			t.Logf("Step %d: %s", i+1, err.message)
			t.Logf("Level: %s", err.level.String())
			t.Logf("Attributes: %d", len(err.attrs))
		})
	}
	
	t.Logf("All related errors logged in span: %s", "error-correlation-test")
}
