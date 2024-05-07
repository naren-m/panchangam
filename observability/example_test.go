package observability_test

import (
	"context"
	"fmt"

	"github.com/naren-m/panchangam/observability"
)

func ExampleObserver_Tracer() {
	// Create a new observer
	observer := observability.NewLocalObserver()

	// Get a tracer from the observer
	tracer := observer.Tracer("test")

	// Output: {<nil>}
	fmt.Println(tracer)
}

func ExampleObserver_CreateSpan() {
	// Create a new observer
	observer := observability.NewLocalObserver()

	// Create a new span using the observer
	_, span := observer.CreateSpan(context.Background(), "test")
	defer span.End()
	span.AddEvent("test event")
	// if span.IsRecording() {
	// 	fmt.Println("Span is recording")
	// } else {
	// 	fmt.Println("Span is not recording")
	// }
	// Output: Successfully created span using observer
	fmt.Println("Successfully created span using observer")
}
