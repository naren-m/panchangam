package observability

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var resource *sdkresource.Resource
var initResourcesOnce sync.Once
var initObserverOnce sync.Once

// Wrappers for OpenTelemetry trace package
var WithAttributes = trace.WithAttributes
var SpanFromContext = trace.SpanFromContext
var NewServerHandler = otelgrpc.NewServerHandler

// https://github.com/wavefrontHQ/opentelemetry-examples/blob/master/go-example/manual-instrumentation/main.go
// https://github.com/wavefrontHQ/opentelemetry-examples/blob/master/go-example/manual-instrumentation/README.md
// https://opentelemetry.io/docs/demo/services/checkout/

type ObserverInterface interface {
	Shutdown(ctx context.Context)
	Tracer(name string) trace.Tracer
}
type observer struct {
	tp *sdktrace.TracerProvider
}

var oi *observer

// NewObserver creates a new Observer instance.
func Observer(address string) ObserverInterface {
	// Initialize the TracerProvider and Tracer.
	initObserverOnce.Do(func() {
		tp, _ := initTracerProvider("")
		oi = &observer{
			tp: tp,
		}
	})

	return oi
}

// Shutdown stops the observer.
func (o *observer) Shutdown(ctx context.Context) {
	o.tp.Shutdown(ctx)
}

// Tracer returns the tracer.
func (o *observer) Tracer(name string) trace.Tracer {
	return o.tp.Tracer(name)
}

// Now you can use observability.TracerProvider the same way as sdktrace.TracerProvider.
func initResource() *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithHost(),
			sdkresource.WithAttributes(
				attribute.String("application", "panchangam"),
				attribute.String("service.name", "panchangam"),
				attribute.String("service.namespace", "observability"),
				attribute.String("application.version", "0.0.1"),
			),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

func initTracerProvider(address string) (*sdktrace.TracerProvider, error) {
	if address == "" {
		address = "localhost:4317"
	}
	conn, err := grpc.NewClient(address,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// Set up a trace exporter
	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func InitMeterProvider() *sdkmetric.MeterProvider {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(fmt.Sprintf("new otlp metric grpc exporter failed: %v", err))
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource()),
	)
	otel.SetMeterProvider(mp)

	return mp
}
