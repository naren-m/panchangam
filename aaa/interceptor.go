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

		c, span := a.observer.Tracer(info.FullMethod).Start(ctx, "aaa.AuthInterceptor")
		logger.InfoContext(c, "Successfully authenticated.", "rpc", info.FullMethod)
		time.Sleep(100 * time.Millisecond)
		span.End()

		return handler(ctx, req)
	}
}

func (a *Auth) AccountingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		c, span := a.observer.Tracer(info.FullMethod).Start(ctx, "aaa.AccountingInterceptor")
		startTime := time.Now()
		time.Sleep(30 * time.Millisecond)
		logger.InfoContext(c, "Accounting successful", "Method", info.FullMethod, "timetook", time.Since(startTime))
		span.End()

		// Continue the handler chain.
		return handler(ctx, req)
	}
}
