package observability

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

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
