package observability

import (
	"context"
	"fmt"
	"sync"
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
func TestConcurrency(t *testing.T) {
	observer := NewLocalObserver()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := context.WithValue(context.Background(), "key", fmt.Sprintf("value%d", i))

			ctx, span := observer.Tracer("test").Start(ctx, "test")
			// assert.True(t, span.IsRecording())
			defer span.End()

			assert.NotNil(t, ctx)
			assert.NotNil(t, span)

			s := SpanFromContext(ctx)
			assert.NotNil(t, s)

			s.AddEvent("test")

		}(i)
	}

	wg.Wait()

	// Add assertions to check that all observations were made correctly
}


func BenchmarkObserver(b *testing.B) {
	observer := NewLocalObserver()

	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				ctx := context.WithValue(context.Background(), "key", fmt.Sprintf("value%d", i))

				ctx, span := observer.Tracer("test").Start(ctx, "test")
				defer span.End()

				s := SpanFromContext(ctx)
				s.AddEvent("test")
			}(i)
		}

		wg.Wait()
	}
}