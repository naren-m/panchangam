package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type panchangamRequestTime struct {
	date             time.Time
	calculationTime  time.Time
	responseDate     string
	timezoneLocation *time.Location
	timezoneString   string
	timezoneInfo     TimezoneInfo
	hasClockTime     bool
}

func parseRequestTime(ctx context.Context, req *ppb.GetPanchangamRequest, span trace.Span) (*panchangamRequestTime, error) {
	logger.DebugContext(ctx, "Parsing date", "date", req.Date)
	parsedDate, parsedDateTime, hasClockTime, err := parsePanchangamDateInput(req.Date)
	if err != nil {
		grpcErr := status.Error(codes.InvalidArgument, fmt.Sprintf("invalid date format: %v", err))

		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityMedium,
			Category:  observability.CategoryValidation,
			Operation: "date_parsing",
			Component: "panchangam_service",
			Additional: map[string]interface{}{
				"date_input":      req.Date,
				"expected_format": "2006-01-02 or RFC3339",
				"parse_error":     err.Error(),
			},
			Retryable:   false,
			ExpectedErr: true,
		})

		observability.RecordEvent(ctx, "Date parsing failed", map[string]interface{}{
			"date":            req.Date,
			"expected_format": "2006-01-02 or RFC3339",
			"error_type":      "invalid_format",
			"parse_error":     err.Error(),
		})

		logger.WarnContext(ctx, "Date parsing failed",
			"date", req.Date,
			"error", grpcErr,
			"expected_format", "2006-01-02 or RFC3339")
		span.RecordError(grpcErr)
		return nil, grpcErr
	}
	logger.DebugContext(ctx, "Date parsed successfully", "has_clock_time", hasClockTime)

	tzParser := NewTimezoneParser()
	tzString := req.Timezone
	if tzString == "" {
		tzString = "UTC"
		logger.DebugContext(ctx, "No timezone specified, using UTC default")
	}

	logger.DebugContext(ctx, "Processing timezone", "timezone", tzString)
	loc, err := tzParser.ParseTimezone(tzString)
	if err != nil {
		grpcErr := status.Error(codes.InvalidArgument, fmt.Sprintf("invalid timezone: %v", err))

		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityMedium,
			Category:  observability.CategoryValidation,
			Operation: "timezone_parsing",
			Component: "panchangam_service",
			Additional: map[string]interface{}{
				"timezone_input": tzString,
				"parse_error":    err.Error(),
			},
			Retryable:   false,
			ExpectedErr: true,
		})

		logger.WarnContext(ctx, "Timezone parsing failed",
			"timezone", tzString,
			"error", grpcErr)
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	logger.DebugContext(ctx, "Timezone parsed successfully",
		"timezone", tzString,
		"location", loc.String())

	var calculationTime time.Time
	var date time.Time
	if hasClockTime {
		calculationTime = parsedDateTime.In(loc)
		date = time.Date(calculationTime.Year(), calculationTime.Month(), calculationTime.Day(), 0, 0, 0, 0, loc)
	} else {
		date = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, loc)
		calculationTime = date
	}

	isValid, warning := tzParser.ValidateTimezoneForLocation(loc, req.Latitude, req.Longitude, date)
	if !isValid {
		logger.WarnContext(ctx, "Timezone validation warning",
			"timezone", tzString,
			"latitude", req.Latitude,
			"longitude", req.Longitude,
			"warning", warning)
		span.AddEvent("Timezone validation warning", traceAttributes(
			"warning", warning,
		)...)
	}

	logger.DebugContext(ctx, "Date adjusted to timezone",
		"adjusted_date", date.Format("2006-01-02 15:04:05 MST"),
		"calculation_time", calculationTime.Format(time.RFC3339),
		"timezone", loc.String())

	return &panchangamRequestTime{
		date:             date,
		calculationTime:  calculationTime,
		responseDate:     date.Format("2006-01-02"),
		timezoneLocation: loc,
		timezoneString:   tzString,
		timezoneInfo:     tzParser.GetTimezoneInfo(loc, date),
		hasClockTime:     hasClockTime,
	}, nil
}
