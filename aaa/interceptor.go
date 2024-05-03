package aaa

import (
	"context"
	"math/rand"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/naren-m/panchangam/observability"
	"github.com/naren-m/panchangam/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
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

// Tracing middleware
func (i *Interceptor) TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Start a new span
	ctx, span := i.tracer.Start(ctx, info.FullMethod)
	defer span.End()

    if err != nil {
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.SetAttributes(statusCodeAttr(s.Code()))
	  } else {
		span.SetAttributes(statusCodeAttr(grpc_codes.OK))
	  }
	  i, err := handler(ctx, req)
	// Call the next interceptor or handler
	return i, errd
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
func (i *Interceptor) AuthorizeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error){
	
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