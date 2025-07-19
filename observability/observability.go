package observability

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
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
	Shutdown(ctx context.Context) error
	Tracer(name string) trace.Tracer
	CreateSpan(ctx context.Context, name string) (context.Context, trace.Span)
}
type observer struct {
	tp *sdktrace.TracerProvider
}

var oi *observer

func NewLocalObserver() ObserverInterface {
	// Initialize the TracerProvider and Tracer.
	initObserverOnce.Do(func() {
		tp, _ := initStdoutProvider()
		oi = &observer{
			tp: tp,
		}
	})

	return oi
}

// NewObserver creates a new Observer instance.
func NewObserver(address string) (ObserverInterface, error) {
	// Initialize the TracerProvider and Tracer.
	var tp *sdktrace.TracerProvider
	var err error
	initObserverOnce.Do(func() {
		if address == "" {
			tp, err = initStdoutProvider()
			oi = &observer{
				tp: tp,
			}
		} else {
			tp, err = initTracerProvider(address)
			oi = &observer{
				tp: tp,
			}
		}
	})

	return oi, err
}

// Observer returns the observer instance. 
// If no observer has been initialized, it will create a local observer with stdout output.
func Observer() ObserverInterface {
	if oi == nil {
		// Auto-initialize with local observer if not already initialized
		// This provides a safe default instead of panicking
		return NewLocalObserver()
	}

	return oi
}

// Shutdown stops the observer.
func (o *observer) Shutdown(ctx context.Context) error {
	return o.tp.Shutdown(ctx)
}

// Tracer returns the tracer.
func (o *observer) Tracer(name string) trace.Tracer {
	return o.tp.Tracer(name)
}

// CreateSpan starts a new span. Getting the RPC method name from the context.
func (o *observer) CreateSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	fullMethod, ok := grpc.Method(ctx)
	if !ok {
		fullMethod = "unknown"
	}
	tracer := otel.GetTracerProvider().Tracer(fullMethod)
	return tracer.Start(ctx, name)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		slog.Info("Entering observability interceptor")
		tracer := Observer().Tracer(fmt.Sprintf("ParentSpan %s", info.FullMethod))
		ctx, oSpan := tracer.Start(ctx, info.FullMethod)
		defer oSpan.End()
		resp, err := handler(ctx, req)
		if err != nil {
			slog.ErrorContext(ctx, "Request failed.", "error", err)
		}
		if oSpan.IsRecording() {
			// If oSpan is recoding, record the oSpan.
			slog.Info("Recording Span.")
			if err != nil {
				oSpan.AddEvent("Request failed.", trace.WithAttributes(attribute.String("error", err.Error())))
				oSpan.RecordError(err)
				oSpan.SetStatus(codes.Error, err.Error())
			} else {
				oSpan.AddEvent("Request completed successfully.")
				oSpan.SetStatus(codes.Ok, "OK")
			}
		}

		slog.Info("Leaving observability interceptor")
		return resp, err
	}
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

func initStdoutProvider() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize stdouttrace export pipeline: %v", err))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func initTracerProvider(address string) (*sdktrace.TracerProvider, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
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
