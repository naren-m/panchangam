package observability

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalFunctions(t *testing.T) {
	ctx := context.Background()

	t.Run("global RecordError", func(t *testing.T) {
		err := errors.New("global test error")
		errorCtx := ErrorContext{
			Severity:  SeverityMedium,
			Category:  CategoryValidation,
			Operation: "global_test",
		}

		enhancedErr := RecordError(ctx, err, errorCtx)
		require.NotNil(t, enhancedErr)
		assert.Equal(t, err, enhancedErr.OriginalError)
	})

	t.Run("global RecordEvent", func(t *testing.T) {
		attributes := map[string]interface{}{
			"global": true,
		}
		RecordEvent(ctx, "global_event", attributes)
	})

	t.Run("global RecordCalculationStart", func(t *testing.T) {
		inputs := map[string]interface{}{
			"test": "input",
		}
		RecordCalculationStart(ctx, "global_calculation", inputs)
	})

	t.Run("global RecordCalculationEnd", func(t *testing.T) {
		outputs := map[string]interface{}{
			"test": "output",
		}
		RecordCalculationEnd(ctx, "global_calculation", true, time.Millisecond, outputs)
	})

	t.Run("global RecordValidationFailure", func(t *testing.T) {
		RecordValidationFailure(ctx, "global_field", "invalid_value", "test validation failure")
	})

	t.Run("global RecordRetryAttempt", func(t *testing.T) {
		lastErr := errors.New("retry test error")
		RecordRetryAttempt(ctx, "global_retry_operation", 1, 3, lastErr)
	})
}

func TestGlobalErrorRecorderConcurrentAccessDoesNotRace(t *testing.T) {
	globalErrorRecorder = nil
	globalErrorRecorderOnce = sync.Once{}

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			RecordEvent(context.Background(), "global_recorder_race_check", nil)
		}()
	}

	close(start)
	wg.Wait()
}
