package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TestCreateSpanWithGrpcContext(t *testing.T) {
	observer := newLocalObserverForTest(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "grpc.method", "/test.Service/TestMethod")

	ctx, span := observer.CreateSpan(ctx, "test_span")
	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
	defer span.End()
}

func TestContextPropagation(t *testing.T) {
	observer := newLocalObserverForTest(t)

	parentCtx, parentSpan := observer.CreateSpan(context.Background(), "parent_span")
	defer parentSpan.End()

	childCtx, childSpan := observer.CreateSpan(parentCtx, "child_span")
	defer childSpan.End()

	parentSpanFromCtx := SpanFromContext(parentCtx)
	childSpanFromCtx := SpanFromContext(childCtx)

	assert.NotNil(t, parentSpanFromCtx)
	assert.NotNil(t, childSpanFromCtx)
	assert.NotEqual(t, parentSpanFromCtx, childSpanFromCtx)
}

func TestSpanAttributesAndEvents(t *testing.T) {
	observer := newLocalObserverForTest(t)

	ctx, span := observer.CreateSpan(context.Background(), "test_span")
	defer span.End()

	span.SetAttributes(
		attribute.String("key1", "value1"),
		attribute.Int("key2", 42),
	)

	span.AddEvent("event1")
	span.AddEvent("event2", trace.WithAttributes(
		attribute.String("event_key", "event_value"),
	))

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
	assert.True(t, span.IsRecording())
}

func TestCreateSpanNotRecording(t *testing.T) {
	observer := newLocalObserverForTest(t)

	ctx, span := observer.CreateSpan(context.Background(), "test_span")
	defer span.End()

	span.AddEvent("test event")
	span.SetAttributes(attribute.String("key", "value"))

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}

func TestTracerEmptyName(t *testing.T) {
	observer := newLocalObserverForTest(t)
	tracer := observer.Tracer("")
	assert.NotNil(t, tracer)
}

func TestCreateSpanEmptyName(t *testing.T) {
	observer := newLocalObserverForTest(t)
	ctx, span := observer.CreateSpan(context.Background(), "")
	defer span.End()

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}

func TestNewLocalObserverSingleton(t *testing.T) {
	observer1 := newLocalObserverForTest(t)
	observer2 := NewLocalObserver()

	assert.Equal(t, observer1, observer2)

	ctx1, span1 := observer1.CreateSpan(context.Background(), "span1")
	ctx2, span2 := observer2.CreateSpan(context.Background(), "span2")

	assert.NotNil(t, ctx1)
	assert.NotNil(t, span1)
	assert.NotNil(t, ctx2)
	assert.NotNil(t, span2)

	span1.End()
	span2.End()
}
