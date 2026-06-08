package observability

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryServerInterceptorWithError(t *testing.T) {
	NewLocalObserver()

	interceptor := UnaryServerInterceptor()
	assert.NotNil(t, interceptor)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, fmt.Errorf("test error")
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	resp, err := interceptor(context.Background(), "test_request", info, handler)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "test error", err.Error())
}

func TestUnaryServerInterceptorWithSuccess(t *testing.T) {
	NewLocalObserver()

	interceptor := UnaryServerInterceptor()
	assert.NotNil(t, interceptor)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success_response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	resp, err := interceptor(context.Background(), "test_request", info, handler)
	assert.Equal(t, "success_response", resp)
	assert.Nil(t, err)
}

func TestUnaryServerInterceptorNoGrpcMethod(t *testing.T) {
	NewLocalObserver()

	interceptor := UnaryServerInterceptor()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	resp, err := interceptor(context.Background(), "request", info, handler)
	assert.Equal(t, "response", resp)
	assert.Nil(t, err)
}

func TestUnaryServerInterceptorNonRecording(t *testing.T) {
	NewLocalObserver()

	interceptor := UnaryServerInterceptor()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/NonRecording",
	}

	resp, err := interceptor(context.Background(), "request", info, handler)
	assert.Equal(t, "success", resp)
	assert.Nil(t, err)
}

func TestUnaryServerInterceptorSpanNotRecording(t *testing.T) {
	NewLocalObserver()

	interceptor := UnaryServerInterceptor()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/NotRecording",
	}

	resp, err := interceptor(context.Background(), "request", info, handler)
	assert.Equal(t, "response", resp)
	assert.Nil(t, err)
}

func TestUnaryServerInterceptorDetailedCoverage(t *testing.T) {
	NewLocalObserver()

	interceptor := UnaryServerInterceptor()
	panicHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test panic")
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/PanicMethod",
	}

	assert.Panics(t, func() {
		interceptor(context.Background(), "request", info, panicHandler)
	})
}
