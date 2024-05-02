// logger.go
package logging

import (
    "context"
    "os"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"

    "github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
    // Create a new Logger instance
    Logger = logrus.New()

    // Set the logger output to stdout
    Logger.SetOutput(os.Stdout)

    // Set log level to Info by default
    Logger.SetLevel(logrus.InfoLevel)

    // Ensure logrus behaves like TTY is disabled
    Logger.SetFormatter(&logrus.TextFormatter{
        DisableColors: true,
        FullTimestamp: true,
    })
}

// LogrusFields creates a logrus.Fields map containing trace and span identifiers
// that can be used to enrich log entries with tracing information.
func LogrusFields(span trace.Span) logrus.Fields {
    return logrus.Fields{
        "span_id":  span.SpanContext().SpanID().String(),
        "trace_id": span.SpanContext().TraceID().String(),
    }
}

// CreateSpan creates a new span with the given name and returns the span and the
// context containing the span. This function should be called at the beginning of
// the operations you want to trace.
func CreateSpan(spanName string) (context.Context, trace.Span) {
    tracer := otel.GetTracerProvider().Tracer("")

	ctx := context.Background()
    // Start a span and return it along with the newly derived context containing the span.
    ctx, span := tracer.Start(ctx, spanName)

    // Optionally, set attributes to the span.
    span.SetAttributes(attribute.String("example_key", "example_value"))

    return ctx, span
}

// You can also add a function to log messages with trace context.
// LogWithTrace logs a message with trace context using the provided log level, message and fields.
func LogWithTrace(span trace.Span, level logrus.Level, msg string, fields logrus.Fields) {
    if span != nil {
        fields = mergeFields(LogrusFields(span), fields)
    }

    entry := Logger.WithFields(fields)
    switch level {
    case logrus.DebugLevel:
        entry.Debug(msg)
    case logrus.InfoLevel:
        entry.Info(msg)
    case logrus.WarnLevel:
        entry.Warn(msg)
    case logrus.ErrorLevel:
        entry.Error(msg)
    case logrus.FatalLevel:
        entry.Fatal(msg)
    case logrus.PanicLevel:
        entry.Panic(msg)
    default:
        entry.Info(msg)
    }
}

// mergeFields merges two logrus.Fields maps.
func mergeFields(a, b logrus.Fields) logrus.Fields {
    for k, v := range b {
        a[k] = v
    }
    return a
}