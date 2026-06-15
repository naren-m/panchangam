package panchangam

import (
	"context"
	"fmt"
	"strings"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PanchangamServer) Get(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.GetPanchangamResponse, error) {
	ctx, span := s.observer.CreateSpan(ctx, "Get")
	defer span.End()

	if req == nil {
		err := status.Error(codes.InvalidArgument, "request cannot be nil")
		span.RecordError(err)
		return nil, err
	}

	logger.InfoContext(ctx, "Panchangam request received",
		"operation", "Get",
		"date", req.Date,
		"latitude", req.Latitude,
		"longitude", req.Longitude,
		"timezone", req.Timezone,
		"region", req.Region,
		"calculation_method", req.CalculationMethod,
		"locale", req.Locale,
	)

	span.SetAttributes(
		traceAttribute("request.date", req.Date),
		traceAttribute("request.latitude", fmt.Sprintf("%.4f", req.Latitude)),
		traceAttribute("request.longitude", fmt.Sprintf("%.4f", req.Longitude)),
		traceAttribute("request.timezone", req.Timezone),
		traceAttribute("request.region", req.Region),
		traceAttribute("request.calculation_method", req.CalculationMethod),
		traceAttribute("request.locale", req.Locale),
	)

	logger.DebugContext(ctx, "Validating request parameters")

	date := strings.TrimSpace(req.Date)
	if date == "" {
		err := status.Error(codes.InvalidArgument, "date parameter is required")
		observability.RecordValidationFailure(ctx, "date", date, "date parameter cannot be empty")
		observability.RecordEvent(ctx, "Validation failure detected", map[string]interface{}{
			"field":      "date",
			"value":      date,
			"error_type": "missing_required",
		})
		span.RecordError(err)
		return nil, err
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		err := status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")

		observability.RecordValidationFailure(ctx, "latitude", req.Latitude, "latitude must be between -90 and 90 degrees")
		observability.RecordEvent(ctx, "Validation failure detected", map[string]interface{}{
			"field":       "latitude",
			"value":       req.Latitude,
			"valid_range": "[-90, 90]",
			"error_type":  "out_of_range",
		})

		span.RecordError(err)
		return nil, err
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		err := status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")

		observability.RecordValidationFailure(ctx, "longitude", req.Longitude, "longitude must be between -180 and 180 degrees")
		observability.RecordEvent(ctx, "Validation failure detected", map[string]interface{}{
			"field":       "longitude",
			"value":       req.Longitude,
			"valid_range": "[-180, 180]",
			"error_type":  "out_of_range",
		})

		span.RecordError(err)
		return nil, err
	}

	logger.DebugContext(ctx, "Request parameters validated successfully")
	logger.InfoContext(ctx, "Fetching panchangam data")

	observability.RecordEvent(ctx, "Panchangam data fetch started", map[string]interface{}{
		"operation": "fetchPanchangamData",
		"date":      req.Date,
		"location":  fmt.Sprintf("%.4f,%.4f", req.Latitude, req.Longitude),
		"timezone":  req.Timezone,
		"region":    req.Region,
	})

	d, err := s.fetchPanchangamData(ctx, req)
	if err != nil {
		observability.RecordError(ctx, err, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryInternal,
			Operation: "fetchPanchangamData",
			Component: "panchangam_service",
			Additional: map[string]interface{}{
				"request_date":      req.Date,
				"request_latitude":  req.Latitude,
				"request_longitude": req.Longitude,
				"request_timezone":  req.Timezone,
				"request_region":    req.Region,
			},
			Retryable:   true,
			ExpectedErr: false,
		})

		observability.RecordEvent(ctx, "Panchangam data fetch failed", map[string]interface{}{
			"operation":  "fetchPanchangamData",
			"error_type": "data_fetch_failure",
			"date":       req.Date,
			"location":   fmt.Sprintf("%.4f,%.4f", req.Latitude, req.Longitude),
		})

		logger.ErrorContext(ctx, "Failed to fetch panchangam data",
			"error", err,
			"operation", "fetchPanchangamData")
		span.RecordError(err)
		return nil, err
	}

	observability.RecordEvent(ctx, "Panchangam data fetch completed", map[string]interface{}{
		"operation":    "fetchPanchangamData",
		"date":         d.Date,
		"events_count": len(d.Events),
		"success":      true,
	})

	logger.DebugContext(ctx, "Building response object")
	response := &ppb.GetPanchangamResponse{
		PanchangamData: d,
	}

	logger.InfoContext(ctx, "Panchangam response prepared successfully",
		"operation", "Get",
		"date", d.Date,
		"tithi", d.Tithi,
		"nakshatra", d.Nakshatra,
		"yoga", d.Yoga,
		"karana", d.Karana,
		"sunrise", d.SunriseTime,
		"sunset", d.SunsetTime,
		"events_count", len(d.Events),
	)

	span.AddEvent("Response prepared", traceAttributes(
		"response.date", d.Date,
		"response.events_count", fmt.Sprintf("%d", len(d.Events)),
	)...)

	return response, nil
}
