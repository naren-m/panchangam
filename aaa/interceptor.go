package aaa

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

var logger = log.Logger()

type Auth struct {
	observer observability.ObserverInterface
}

func NewAuth() *Auth {
	o := observability.Observer()
	return &Auth{
		observer: o,
	}
}

func (a *Auth) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, span := a.observer.Tracer(info.FullMethod).Start(ctx, "aaa.AuthInterceptor")
		logger.Info("Successfully authenticated.", "rpc", info.FullMethod)
		span.AddEvent("authenticated")
		// Continue the handler chain.
		time.Sleep(100 * time.Millisecond)
		span.End()
		return handler(ctx, req)
	}
}

func AccountingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Continue the handler chain.
		resp, err := handler(ctx, req)

		// Log the call details.
		logger.Info("Method: %s, Duration: %s\n", info.FullMethod, time.Since(startTime))

		return resp, err
	}
}
