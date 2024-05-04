package aaa

import (
	"context"
	"math/rand"
	"time"

	"github.com/naren-m/panchangam/observability"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	tracer observability.Tracer
}

func NewInterceptor(t observability.Tracer) *Interceptor {
	i := &Interceptor{
		tracer: *observability.NewTracer("Panchangam-server"),
	}

	return i
}

// statusCodeAttr assumes to return an appropriate OpenTelemetry attribute based on the gRPC status code.
func statusCodeAttr(code codes.Code) trace.Attribute {
	return trace.StringAttribute("grpc.status_code", code.String())
}

// TraceInterceptor traces a gRPC request by starting a span and recording the status.
func (i *Interceptor) TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Start a new span
	ctx, span := i.tracer.Start(ctx, info.FullMethod)
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.AddEvent("error", trace.WithAttributes(semconv.ExceptionMessageKey.String(err.Error())))
		}
	}()

	// Call the handler
	resp, err := handler(ctx, req)

	if err != nil {
		// Convert the error to a gRPC status, then set the span status and attributes accordingly
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		// Set the span status to OK
		span.SetStatus(codes.Ok, "")
		span.SetAttributes(statusCodeAttr(codes.Ok))
	}

	// Return the response from the handler
	return resp, err
}

// Authentication middleware
func (i *Interceptor) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, span := i.tracer.Start(ctx, info.FullMethod)
	defer span.End()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		span.RecordError()
		return nil, status.Error(codes.Unauthenticated, "metadata missing from request")
	}

	if len(md["authorization"]) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token missing")
	}

	return ctx, nil
}

// Authorization middleware
func (i *Interceptor) AuthorizeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if rand.Intn(3-1+1) == 2 {
		return nil, status.Error(codes.PermissionDenied, "Authorization failed")
	}
	// Simulate accounting by sleeping for a random duration between 1 and 3 seconds
	sleepDuration := time.Duration(rand.Intn(3-1+1)+1) * time.Second
	time.Sleep(sleepDuration)

	logrus.Infof("User authorized")

	return ctx, nil
}

// Accounting middleware
func AccountInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Simulate accounting by sleeping for a random duration between 1 and 3 seconds
	sleepDuration := time.Duration(rand.Intn(3-1+1)+1) * time.Second
	time.Sleep(sleepDuration)

	elapsed := time.Since(start)
	logrus.Infof("Accounting completed in %v", elapsed)

	return handler(ctx, req)
}
