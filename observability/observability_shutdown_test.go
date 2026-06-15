package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShutdown(t *testing.T) {
	observer := newLocalObserverForTest(t)
	assert.NotNil(t, observer)

	err := observer.Shutdown(context.Background())
	assert.Nil(t, err)

	err = observer.Shutdown(context.Background())
	assert.Nil(t, err)
}

func TestShutdownWithCancelledContext(t *testing.T) {
	observer := newLocalObserverForTest(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := observer.Shutdown(ctx)
	_ = err
}

func TestShutdownWithTimeout(t *testing.T) {
	observer := newLocalObserverForTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 1)
	defer cancel()

	err := observer.Shutdown(ctx)
	_ = err
}
