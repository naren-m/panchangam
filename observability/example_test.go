package observability_test

import (
	"context"
	"fmt"

	"github.com/naren-m/panchangam/observability"
)

func ExampleNewLocalObserver() {
	// Create a new observer
	observer := observability.NewLocalObserver()

	// Get a tracer from the observer
	_ = observer.Tracer("test")

	// Output: Successfully created tracer
	fmt.Println("Successfully created tracer")
}

func ExampleObserverInterface_CreateSpan() {
	// Create a new observer
	observer := observability.NewLocalObserver()

	// Create a new span using the observer
	_, span := observer.CreateSpan(context.Background(), "test")
	defer span.End()
	span.AddEvent("test event")
	// Output: Successfully created span using observer
	fmt.Println("Successfully created span using observer")
}
