package observability

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newObserver(t *testing.T) (ObserverInterface, error) {
	observer, err := NewObserver("localhost:4317")
	assert.NotNil(t, observer)
	assert.Nil(t, err)
	return observer, err
}

// NOTE: This must be the first test, As it tests if call of Observer() panics if not initialized.
func TestObserver(t *testing.T) {

	t.Run("Test for panic", func(t *testing.T) {
		assert.Panics(t, func() { Observer() }, "The code did not panic")
	})

	t.Run("Test for Observer", func(t *testing.T) {
		newObserver(t)
		observer := Observer()
		assert.NotNil(t, observer)
	})

}

func TestNewObserver(t *testing.T) {
	newObserver(t)
}

func TestObserverSingleton(t *testing.T) {
	// Create two observers using NewObserver
	observer1, _ := newObserver(t)
	observer2, _ := newObserver(t)

	// Create two observers using Observer
	observer3 := Observer()
	observer4 := Observer()

	assert.Equal(t, observer1, observer2)
	assert.Equal(t, observer3, observer4)
	assert.Equal(t, observer1, observer3)
}

// TODO: Test is not verifying.
func TestShutdown(t *testing.T) {
	observer, err := newObserver(t)
	err = observer.Shutdown(context.Background())
	assert.Nil(t, err)
}

func TestTracer(t *testing.T) {
	observer, _ := newObserver(t)
	tracer := observer.Tracer("test")
	assert.NotNil(t, tracer)
}

func TestCreateSpan(t *testing.T) {
	observer, _ := newObserver(t)
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
