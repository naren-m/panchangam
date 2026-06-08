package panchangam

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func traceAttribute(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

func traceAttributes(keyValues ...string) []trace.EventOption {
	if len(keyValues)%2 != 0 {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(keyValues)/2)
	for i := 0; i < len(keyValues); i += 2 {
		attrs = append(attrs, attribute.String(keyValues[i], keyValues[i+1]))
	}

	return []trace.EventOption{trace.WithAttributes(attrs...)}
}
