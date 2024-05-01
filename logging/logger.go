// logger.go
package logging

import (
	"context"
	"os"

	oteltrace "go.opentelemetry.io/otel/trace"

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
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func LogrusFields(span oteltrace.Span) logrus.Fields {
    return logrus.Fields{
        "span_id": span.SpanContext().SpanID().String(),
        "trace_id": span.SpanContext().TraceID().String(),
    }
}

func CreateSpan(ctx context.Context, spanName string) (context.Context, oteltrace.Span) {

	return oteltrace.Tracer.Start(ctx, spanName, ...)
	
}