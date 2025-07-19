package observability

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize observability for testing
	NewLocalObserver()
}

func TestErrorRecorder_RecordError(t *testing.T) {
	recorder := NewErrorRecorder()
	ctx := context.Background()

	t.Run("basic error recording", func(t *testing.T) {
		originalErr := errors.New("test error")
		errorCtx := ErrorContext{
			Severity:  SeverityHigh,
			Category:  CategoryCalculation,
			Operation: "test_operation",
			Component: "test_component",
			Additional: map[string]interface{}{
				"test_key": "test_value",
			},
			Retryable:   true,
			ExpectedErr: false,
		}

		enhancedErr := recorder.RecordError(ctx, originalErr, errorCtx)

		require.NotNil(t, enhancedErr)
		assert.Equal(t, originalErr, enhancedErr.OriginalError)
		assert.Equal(t, errorCtx.Severity, enhancedErr.Context.Severity)
		assert.Equal(t, errorCtx.Category, enhancedErr.Context.Category)
		assert.Equal(t, errorCtx.Operation, enhancedErr.Context.Operation)
		assert.NotEmpty(t, enhancedErr.CorrelationID)
		assert.NotEmpty(t, enhancedErr.StackTrace)
		assert.WithinDuration(t, time.Now(), enhancedErr.Timestamp, time.Second)
	})

	t.Run("nil error returns nil", func(t *testing.T) {
		errorCtx := ErrorContext{
			Severity: SeverityLow,
			Category: CategoryValidation,
		}

		enhancedErr := recorder.RecordError(ctx, nil, errorCtx)
		assert.Nil(t, enhancedErr)
	})

	t.Run("error with span context", func(t *testing.T) {
		// Create a span context
		observer := Observer()
		ctx, span := observer.CreateSpan(context.Background(), "test_span")
		defer span.End()

		originalErr := errors.New("span error")
		errorCtx := ErrorContext{
			Severity:  SeverityCritical,
			Category:  CategoryInternal,
			Operation: "span_test",
			Component: "test",
		}

		enhancedErr := recorder.RecordError(ctx, originalErr, errorCtx)

		require.NotNil(t, enhancedErr)
		assert.Equal(t, originalErr, enhancedErr.OriginalError)

		// Verify span is recording (though we can't easily test the internal state)
		assert.True(t, span.IsRecording())
	})
}

func TestErrorRecorder_RecordEvent(t *testing.T) {
	recorder := NewErrorRecorder()
	ctx := context.Background()

	t.Run("basic event recording", func(t *testing.T) {
		eventName := "test_event"
		attributes := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
			"key4": 3.14,
		}

		// This should not panic
		recorder.RecordEvent(ctx, eventName, attributes)
	})

	t.Run("event with span context", func(t *testing.T) {
		observer := Observer()
		ctx, span := observer.CreateSpan(context.Background(), "test_event_span")
		defer span.End()

		eventName := "span_event"
		attributes := map[string]interface{}{
			"event_type": "test",
			"count":      5,
		}

		recorder.RecordEvent(ctx, eventName, attributes)
		assert.True(t, span.IsRecording())
	})
}

func TestErrorRecorder_RecordCalculationStart(t *testing.T) {
	recorder := NewErrorRecorder()
	ctx := context.Background()

	operation := "sun_position_calculation"
	inputs := map[string]interface{}{
		"julian_day": 2451545.0,
		"latitude":   40.7128,
		"longitude":  -74.0060,
	}

	// Should not panic
	recorder.RecordCalculationStart(ctx, operation, inputs)
}

func TestErrorRecorder_RecordCalculationEnd(t *testing.T) {
	recorder := NewErrorRecorder()
	ctx := context.Background()

	operation := "moon_position_calculation"
	duration := 50 * time.Millisecond
	outputs := map[string]interface{}{
		"longitude": 252.018,
		"latitude":  -4.901,
		"distance":  390024.55,
	}

	t.Run("successful calculation", func(t *testing.T) {
		recorder.RecordCalculationEnd(ctx, operation, true, duration, outputs)
	})

	t.Run("failed calculation", func(t *testing.T) {
		recorder.RecordCalculationEnd(ctx, operation, false, duration, nil)
	})
}

func TestErrorRecorder_RecordValidationFailure(t *testing.T) {
	recorder := NewErrorRecorder()
	ctx := context.Background()

	field := "latitude"
	value := 200.0
	reason := "latitude must be between -90 and 90 degrees"

	recorder.RecordValidationFailure(ctx, field, value, reason)
}

func TestErrorRecorder_RecordRetryAttempt(t *testing.T) {
	recorder := NewErrorRecorder()
	ctx := context.Background()

	operation := "external_api_call"
	attempt := 2
	maxAttempts := 3
	lastError := errors.New("connection timeout")

	recorder.RecordRetryAttempt(ctx, operation, attempt, maxAttempts, lastError)
}

func TestEnhancedError_ErrorInterface(t *testing.T) {
	originalErr := errors.New("original error message")
	enhancedErr := &EnhancedError{
		OriginalError: originalErr,
		Context: ErrorContext{
			Severity: SeverityMedium,
			Category: CategoryNetwork,
		},
	}

	// Test Error() method
	assert.Equal(t, originalErr.Error(), enhancedErr.Error())

	// Test Unwrap() method
	assert.Equal(t, originalErr, enhancedErr.Unwrap())
}

func TestErrorSeverity_Values(t *testing.T) {
	severities := []ErrorSeverity{
		SeverityLow,
		SeverityMedium,
		SeverityHigh,
		SeverityCritical,
	}

	for _, severity := range severities {
		assert.NotEmpty(t, string(severity))
	}
}

func TestErrorCategory_Values(t *testing.T) {
	categories := []ErrorCategory{
		CategoryValidation,
		CategoryAuthentication,
		CategoryAuthorization,
		CategoryNetwork,
		CategoryDatabase,
		CategoryExternal,
		CategoryCalculation,
		CategoryConfiguration,
		CategoryResource,
		CategoryInternal,
	}

	for _, category := range categories {
		assert.NotEmpty(t, string(category))
	}
}

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

func TestAttributeFromValue(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    interface{}
		expected string // We'll check the key, actual value testing is complex with otel attributes
	}{
		{"string value", "test_key", "test_value", "test_key"},
		{"int value", "test_int", 42, "test_int"},
		{"int64 value", "test_int64", int64(123), "test_int64"},
		{"float64 value", "test_float", 3.14, "test_float"},
		{"bool value", "test_bool", true, "test_bool"},
		{"complex value", "test_complex", map[string]string{"nested": "value"}, "test_complex"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := attributeFromValue(tt.key, tt.value)
			assert.Equal(t, tt.expected, string(attr.Key))
		})
	}
}

func TestGenerateCorrelationID(t *testing.T) {
	id1 := generateCorrelationID()
	
	// Sleep briefly to ensure different timestamp
	time.Sleep(1 * time.Nanosecond)
	id2 := generateCorrelationID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.Contains(t, id1, "err_")
	assert.Contains(t, id2, "err_")
	
	// Test that multiple IDs can be different (though not guaranteed due to timing)
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id := generateCorrelationID()
		ids[id] = true
		time.Sleep(1 * time.Nanosecond)
	}
	// Should have at least some unique IDs
	assert.GreaterOrEqual(t, len(ids), 5, "Should generate some unique correlation IDs")
}

func TestCaptureStackTrace(t *testing.T) {
	stackTrace := captureStackTrace(1)

	assert.NotEmpty(t, stackTrace)
	assert.Contains(t, stackTrace, "TestCaptureStackTrace")
	assert.Contains(t, stackTrace, "errors_test.go")
}

// Test recordToSpan function coverage with different severity levels
func TestRecordToSpanCoverage(t *testing.T) {
	NewLocalObserver()
	recorder := NewErrorRecorder()
	
	// Create a context with a span
	ctx, span := Observer().CreateSpan(context.Background(), "test_span")
	defer span.End()
	
	// Test with different severity levels to cover all branches
	severities := []ErrorSeverity{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
	}
	
	for _, severity := range severities {
		testErr := errors.New("test error for severity " + string(severity))
		errorCtx := ErrorContext{
			Severity:    severity,
			Category:    CategoryCalculation,
			Operation:   "test_operation",
			Component:   "test_component",
			UserID:      "user123",        // Test user context coverage
			RequestID:   "req456",         // Test request context coverage
			SessionID:   "session789",     // Test session context coverage
			Retryable:   true,
			ExpectedErr: false,
			Additional: map[string]interface{}{
				"extra_key": "extra_value",
			},
		}
		
		recorder.RecordError(ctx, testErr, errorCtx)
	}
}

// Test logStructuredError function coverage with different branches
func TestLogStructuredErrorCoverage(t *testing.T) {
	recorder := NewErrorRecorder()
	
	// Test with all severity levels to cover switch branches
	severities := []ErrorSeverity{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
		"unknown", // Test default case
	}
	
	for _, severity := range severities {
		testErr := errors.New("test error for logging")
		errorCtx := ErrorContext{
			Severity:  severity,
			Category:  CategoryValidation,
			Operation: "log_test",
			Additional: map[string]interface{}{
				"complex_struct": map[string]interface{}{
					"nested": "value",
				},
			},
		}
		
		// This will test the logStructuredError function
		recorder.RecordError(context.Background(), testErr, errorCtx)
	}
}

// Test edge cases for captureStackTrace
func TestCaptureStackTraceEdgeCases(t *testing.T) {
	// Test with high skip value
	stackTrace := captureStackTrace(100)
	assert.NotEmpty(t, stackTrace)
	
	// Test with zero skip
	stackTrace2 := captureStackTrace(0)
	assert.NotEmpty(t, stackTrace2)
	assert.Contains(t, stackTrace2, "captureStackTrace")
}

func BenchmarkErrorRecording(b *testing.B) {
	recorder := NewErrorRecorder()
	ctx := context.Background()
	err := errors.New("benchmark error")
	errorCtx := ErrorContext{
		Severity:  SeverityMedium,
		Category:  CategoryCalculation,
		Operation: "benchmark",
		Component: "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recorder.RecordError(ctx, err, errorCtx)
	}
}

func BenchmarkEventRecording(b *testing.B) {
	recorder := NewErrorRecorder()
	ctx := context.Background()
	attributes := map[string]interface{}{
		"benchmark": true,
		"iteration": 0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attributes["iteration"] = i
		recorder.RecordEvent(ctx, "benchmark_event", attributes)
	}
}
