package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
)

func TestServiceErrorHandling(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()
	ctx := context.Background()

	resp, err := server.Get(ctx, nil)
	assert.Error(t, err, "Nil request should cause error")
	assert.Nil(t, resp, "Response should be nil for nil request")

	resp, err = server.Get(ctx, &ppb.GetPanchangamRequest{})
	assert.Error(t, err, "Empty request should cause error")
	assert.Nil(t, resp, "Response should be nil for empty request")

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	_, _ = server.Get(cancelCtx, req)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), -time.Nanosecond)
	defer cancel()

	_, _ = server.Get(timeoutCtx, req)
}
