package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewObserver(t *testing.T) {
	observer := NewObserver("localhost:4317")
	assert.NotNil(t, observer)
}

func TestObserver(t *testing.T) {
	observer := Observer()
	assert.NotNil(t, observer)
}
func TestObserverSingleton(t *testing.T) {
	// Create two observers using NewObserver
	observer1 := NewObserver("localhost:4317")
	observer2 := NewObserver("localhost:4317")

	// Create two observers using Observer
	observer3 := Observer()
	observer4 := Observer()

	// Check that the observers created by NewObserver are not the same instance
	assert.Equal(t, observer1, observer2)

	// Check that the observers created by Observer are the same instance
	assert.Equal(t, observer3, observer4)
}

// TODO: Test is not verifying.
func TestShutdown(t *testing.T) {
	observer := NewObserver("localhost:4317")
	err := observer.Shutdown(context.Background())
	assert.Nil(t, err)
}

func TestTracer(t *testing.T) {
	observer := NewObserver("localhost:4317")
	tracer := observer.Tracer("test")
	assert.NotNil(t, tracer)
}

func TestCreateSpan(t *testing.T) {
	observer := NewObserver("localhost:4317")
	ctx, span := observer.CreateSpan(context.Background(), "test")
	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}