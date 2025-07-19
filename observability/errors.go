package observability

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	SeverityLow      ErrorSeverity = "low"
	SeverityMedium   ErrorSeverity = "medium"
	SeverityHigh     ErrorSeverity = "high"
	SeverityCritical ErrorSeverity = "critical"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	CategoryValidation     ErrorCategory = "validation"
	CategoryAuthentication ErrorCategory = "authentication"
	CategoryAuthorization  ErrorCategory = "authorization"
	CategoryNetwork        ErrorCategory = "network"
	CategoryDatabase       ErrorCategory = "database"
	CategoryExternal       ErrorCategory = "external_service"
	CategoryCalculation    ErrorCategory = "calculation"
	CategoryConfiguration  ErrorCategory = "configuration"
	CategoryResource       ErrorCategory = "resource"
	CategoryInternal       ErrorCategory = "internal"
)

// ErrorContext contains additional context for error reporting
type ErrorContext struct {
	Severity    ErrorSeverity
	Category    ErrorCategory
	Operation   string
	Component   string
	UserID      string
	RequestID   string
	SessionID   string
	Additional  map[string]interface{}
	Retryable   bool
	ExpectedErr bool
}

// EnhancedError wraps an error with additional observability context
type EnhancedError struct {
	OriginalError error
	Context       ErrorContext
	Timestamp     time.Time
	StackTrace    string
	CorrelationID string
}

// Error implements the error interface
func (e *EnhancedError) Error() string {
	return e.OriginalError.Error()
}

// Unwrap returns the original error
func (e *EnhancedError) Unwrap() error {
	return e.OriginalError
}

// ErrorRecorder provides enhanced error recording capabilities
type ErrorRecorder struct {
	observer ObserverInterface
}

// NewErrorRecorder creates a new error recorder
func NewErrorRecorder() *ErrorRecorder {
	return &ErrorRecorder{
		observer: Observer(),
	}
}

// RecordError records an error with comprehensive context and span events
func (er *ErrorRecorder) RecordError(ctx context.Context, err error, errorCtx ErrorContext) *EnhancedError {
	if err == nil {
		return nil
	}

	// Generate correlation ID
	correlationID := generateCorrelationID()

	// Capture stack trace
	stackTrace := captureStackTrace(2) // Skip this function and the caller

	// Create enhanced error
	enhancedErr := &EnhancedError{
		OriginalError: err,
		Context:       errorCtx,
		Timestamp:     time.Now(),
		StackTrace:    stackTrace,
		CorrelationID: correlationID,
	}

	// Record to span if available
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		er.recordToSpan(span, enhancedErr)
	}

	// Log structured error
	er.logStructuredError(ctx, enhancedErr)

	return enhancedErr
}

// RecordEvent records an important event with span events
func (er *ErrorRecorder) RecordEvent(ctx context.Context, eventName string, attributes map[string]interface{}) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		// Convert attributes to OpenTelemetry attributes
		otelAttrs := make([]attribute.KeyValue, 0, len(attributes))
		for key, value := range attributes {
			otelAttrs = append(otelAttrs, attributeFromValue(key, value))
		}

		span.AddEvent(eventName, trace.WithAttributes(otelAttrs...))
	}

	// Log the event
	slog.InfoContext(ctx, "Important event recorded",
		"event_name", eventName,
		"attributes", attributes,
		"timestamp", time.Now().Format(time.RFC3339),
	)
}

// RecordCalculationStart records the start of a calculation operation
func (er *ErrorRecorder) RecordCalculationStart(ctx context.Context, operation string, inputs map[string]interface{}) {
	attributes := map[string]interface{}{
		"operation_type": "calculation_start",
		"operation":      operation,
		"inputs":         inputs,
		"start_time":     time.Now().Format(time.RFC3339),
	}
	er.RecordEvent(ctx, fmt.Sprintf("Calculation started: %s", operation), attributes)
}

// RecordCalculationEnd records the end of a calculation operation
func (er *ErrorRecorder) RecordCalculationEnd(ctx context.Context, operation string, success bool, duration time.Duration, outputs map[string]interface{}) {
	attributes := map[string]interface{}{
		"operation_type": "calculation_end",
		"operation":      operation,
		"success":        success,
		"duration_ms":    duration.Milliseconds(),
		"outputs":        outputs,
		"end_time":       time.Now().Format(time.RFC3339),
	}
	er.RecordEvent(ctx, fmt.Sprintf("Calculation completed: %s", operation), attributes)
}

// RecordValidationFailure records validation failures with detailed context
func (er *ErrorRecorder) RecordValidationFailure(ctx context.Context, field string, value interface{}, reason string) {
	errorCtx := ErrorContext{
		Severity:  SeverityMedium,
		Category:  CategoryValidation,
		Operation: "validation",
		Component: "input_validator",
		Additional: map[string]interface{}{
			"field":  field,
			"value":  value,
			"reason": reason,
		},
		Retryable:   false,
		ExpectedErr: true,
	}

	err := fmt.Errorf("validation failed for field '%s': %s", field, reason)
	er.RecordError(ctx, err, errorCtx)
}

// RecordRetryAttempt records retry attempts with context
func (er *ErrorRecorder) RecordRetryAttempt(ctx context.Context, operation string, attempt int, maxAttempts int, lastError error) {
	attributes := map[string]interface{}{
		"operation":    operation,
		"attempt":      attempt,
		"max_attempts": maxAttempts,
		"last_error":   lastError.Error(),
		"retry_time":   time.Now().Format(time.RFC3339),
	}
	er.RecordEvent(ctx, fmt.Sprintf("Retry attempt %d/%d for %s", attempt, maxAttempts, operation), attributes)
}

// recordToSpan records error information to the current span
func (er *ErrorRecorder) recordToSpan(span trace.Span, enhancedErr *EnhancedError) {
	// Record the error
	span.RecordError(enhancedErr.OriginalError)

	// Set span status
	var statusCode codes.Code
	switch enhancedErr.Context.Severity {
	case SeverityCritical, SeverityHigh:
		statusCode = codes.Error
	case SeverityMedium:
		statusCode = codes.Error
	case SeverityLow:
		statusCode = codes.Ok // Low severity might not be considered an error
	default:
		statusCode = codes.Error
	}
	span.SetStatus(statusCode, enhancedErr.OriginalError.Error())

	// Add comprehensive attributes
	attributes := []attribute.KeyValue{
		attribute.String("error.type", string(enhancedErr.Context.Category)),
		attribute.String("error.severity", string(enhancedErr.Context.Severity)),
		attribute.String("error.operation", enhancedErr.Context.Operation),
		attribute.String("error.component", enhancedErr.Context.Component),
		attribute.String("error.correlation_id", enhancedErr.CorrelationID),
		attribute.Bool("error.retryable", enhancedErr.Context.Retryable),
		attribute.Bool("error.expected", enhancedErr.Context.ExpectedErr),
		attribute.String("error.timestamp", enhancedErr.Timestamp.Format(time.RFC3339)),
	}

	// Add user context if available
	if enhancedErr.Context.UserID != "" {
		attributes = append(attributes, attribute.String("user.id", enhancedErr.Context.UserID))
	}
	if enhancedErr.Context.RequestID != "" {
		attributes = append(attributes, attribute.String("request.id", enhancedErr.Context.RequestID))
	}
	if enhancedErr.Context.SessionID != "" {
		attributes = append(attributes, attribute.String("session.id", enhancedErr.Context.SessionID))
	}

	// Add additional context
	for key, value := range enhancedErr.Context.Additional {
		attributes = append(attributes, attributeFromValue(fmt.Sprintf("error.%s", key), value))
	}

	span.SetAttributes(attributes...)

	// Add span event
	eventName := fmt.Sprintf("Error recorded: %s", enhancedErr.Context.Category)
	eventAttributes := []attribute.KeyValue{
		attribute.String("error.message", enhancedErr.OriginalError.Error()),
		attribute.String("error.severity", string(enhancedErr.Context.Severity)),
		attribute.String("error.correlation_id", enhancedErr.CorrelationID),
	}

	span.AddEvent(eventName, trace.WithAttributes(eventAttributes...))
}

// logStructuredError logs the error with structured logging
func (er *ErrorRecorder) logStructuredError(ctx context.Context, enhancedErr *EnhancedError) {
	// Determine log level based on severity
	logArgs := []interface{}{
		"error", enhancedErr.OriginalError.Error(),
		"error_type", enhancedErr.Context.Category,
		"severity", enhancedErr.Context.Severity,
		"operation", enhancedErr.Context.Operation,
		"component", enhancedErr.Context.Component,
		"correlation_id", enhancedErr.CorrelationID,
		"retryable", enhancedErr.Context.Retryable,
		"expected", enhancedErr.Context.ExpectedErr,
		"timestamp", enhancedErr.Timestamp.Format(time.RFC3339),
	}

	// Add user context
	if enhancedErr.Context.UserID != "" {
		logArgs = append(logArgs, "user_id", enhancedErr.Context.UserID)
	}
	if enhancedErr.Context.RequestID != "" {
		logArgs = append(logArgs, "request_id", enhancedErr.Context.RequestID)
	}
	if enhancedErr.Context.SessionID != "" {
		logArgs = append(logArgs, "session_id", enhancedErr.Context.SessionID)
	}

	// Add additional context
	for key, value := range enhancedErr.Context.Additional {
		logArgs = append(logArgs, key, value)
	}

	// Add stack trace for high severity errors
	if enhancedErr.Context.Severity == SeverityHigh || enhancedErr.Context.Severity == SeverityCritical {
		logArgs = append(logArgs, "stack_trace", enhancedErr.StackTrace)
	}

	// Log based on severity
	switch enhancedErr.Context.Severity {
	case SeverityCritical:
		slog.ErrorContext(ctx, "CRITICAL ERROR occurred", logArgs...)
	case SeverityHigh:
		slog.ErrorContext(ctx, "HIGH severity error occurred", logArgs...)
	case SeverityMedium:
		slog.WarnContext(ctx, "MEDIUM severity error occurred", logArgs...)
	case SeverityLow:
		slog.InfoContext(ctx, "LOW severity error occurred", logArgs...)
	default:
		slog.ErrorContext(ctx, "Error occurred", logArgs...)
	}
}

// attributeFromValue creates an OpenTelemetry attribute from a value
func attributeFromValue(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case int:
		return attribute.Int(key, v)
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	default:
		return attribute.String(key, fmt.Sprintf("%v", v))
	}
}

// generateCorrelationID generates a unique correlation ID for error tracking
func generateCorrelationID() string {
	return fmt.Sprintf("err_%d_%d", time.Now().UnixNano(), runtime.NumGoroutine())
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) string {
	const maxStackSize = 50
	pc := make([]uintptr, maxStackSize)
	n := runtime.Callers(skip, pc)

	if n == 0 {
		return "no stack trace available"
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	var stackTrace string
	for {
		frame, more := frames.Next()
		stackTrace += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}

	return stackTrace
}

// GetGlobalErrorRecorder returns a global instance of ErrorRecorder
var globalErrorRecorder *ErrorRecorder

func getGlobalErrorRecorder() *ErrorRecorder {
	if globalErrorRecorder == nil {
		globalErrorRecorder = NewErrorRecorder()
	}
	return globalErrorRecorder
}

// RecordError provides a convenient global function for error recording
func RecordError(ctx context.Context, err error, errorCtx ErrorContext) *EnhancedError {
	return getGlobalErrorRecorder().RecordError(ctx, err, errorCtx)
}

// RecordEvent provides a convenient global function for event recording
func RecordEvent(ctx context.Context, eventName string, attributes map[string]interface{}) {
	getGlobalErrorRecorder().RecordEvent(ctx, eventName, attributes)
}

// RecordCalculationStart provides a convenient global function
func RecordCalculationStart(ctx context.Context, operation string, inputs map[string]interface{}) {
	getGlobalErrorRecorder().RecordCalculationStart(ctx, operation, inputs)
}

// RecordCalculationEnd provides a convenient global function
func RecordCalculationEnd(ctx context.Context, operation string, success bool, duration time.Duration, outputs map[string]interface{}) {
	getGlobalErrorRecorder().RecordCalculationEnd(ctx, operation, success, duration, outputs)
}

// RecordValidationFailure provides a convenient global function
func RecordValidationFailure(ctx context.Context, field string, value interface{}, reason string) {
	getGlobalErrorRecorder().RecordValidationFailure(ctx, field, value, reason)
}

// RecordRetryAttempt provides a convenient global function
func RecordRetryAttempt(ctx context.Context, operation string, attempt int, maxAttempts int, lastError error) {
	getGlobalErrorRecorder().RecordRetryAttempt(ctx, operation, attempt, maxAttempts, lastError)
}
