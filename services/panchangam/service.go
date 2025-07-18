// server.go
package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var logger = log.Logger()

type PanchangamServer struct {
	observer observability.ObserverInterface
	ppb.UnimplementedPanchangamServer
}

func NewPanchangamServer() *PanchangamServer {
	return &PanchangamServer{
		observer: observability.Observer(),
	}
}

// Helper functions for tracing
func traceAttribute(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

func traceAttributes(keyValues ...string) []trace.EventOption {
	if len(keyValues)%2 != 0 {
		return nil
	}
	
	attrs := make([]attribute.KeyValue, 0, len(keyValues)/2)
	for i := 0; i < len(keyValues); i += 2 {
		attrs = append(attrs, attribute.String(keyValues[i], keyValues[i+1]))
	}
	
	return []trace.EventOption{trace.WithAttributes(attrs...)}
}

func (s *PanchangamServer) Get(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.GetPanchangamResponse, error) {
	ctx, span := s.observer.CreateSpan(ctx, "Get")
	defer span.End()
	
	// Log request with comprehensive context
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
	
	// Add span attributes for better tracing
	span.SetAttributes(
		traceAttribute("request.date", req.Date),
		traceAttribute("request.latitude", fmt.Sprintf("%.4f", req.Latitude)),
		traceAttribute("request.longitude", fmt.Sprintf("%.4f", req.Longitude)),
		traceAttribute("request.timezone", req.Timezone),
		traceAttribute("request.region", req.Region),
		traceAttribute("request.calculation_method", req.CalculationMethod),
		traceAttribute("request.locale", req.Locale),
	)
	
	// Validate request parameters
	logger.DebugContext(ctx, "Validating request parameters")
	if req.Latitude < -90 || req.Latitude > 90 {
		err := status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
		logger.WarnContext(ctx, "Invalid latitude parameter", 
			"latitude", req.Latitude,
			"error", err)
		span.RecordError(err)
		return nil, err
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		err := status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
		logger.WarnContext(ctx, "Invalid longitude parameter", 
			"longitude", req.Longitude,
			"error", err)
		span.RecordError(err)
		return nil, err
	}
	
	logger.DebugContext(ctx, "Request parameters validated successfully")
	
	// Fetch panchangam data
	logger.InfoContext(ctx, "Fetching panchangam data")
	d, err := s.fetchPanchangamData(ctx, req)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch panchangam data", 
			"error", err,
			"operation", "fetchPanchangamData")
		span.RecordError(err)
		return nil, err
	}
	
	// Prepare response
	logger.DebugContext(ctx, "Building response object")
	response := &ppb.GetPanchangamResponse{
		PanchangamData: d,
	}
	
	// Add processing delay simulation
	time.Sleep(100 * time.Millisecond)
	
	// Log successful response
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

func (s *PanchangamServer) fetchPanchangamData(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.PanchangamData, error) {
	ctx, span := s.observer.CreateSpan(ctx, "fetchPanchangamData")
	defer span.End()

	logger.InfoContext(ctx, "Starting panchangam data fetch", 
		"operation", "fetchPanchangamData",
		"date", req.Date,
		"location", fmt.Sprintf("%.4f,%.4f", req.Latitude, req.Longitude),
	)
	
	// Add span attributes for detailed tracing
	span.SetAttributes(
		traceAttribute("operation", "fetchPanchangamData"),
		traceAttribute("date", req.Date),
		traceAttribute("location", fmt.Sprintf("%.4f,%.4f", req.Latitude, req.Longitude)),
	)
	
	// Simulate data fetching delay
	logger.DebugContext(ctx, "Simulating data fetch delay")
	time.Sleep(29 * time.Millisecond)

	// Parse and validate the date
	logger.DebugContext(ctx, "Parsing date", "date", req.Date)
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("invalid date format: %v", err))
		logger.WarnContext(ctx, "Date parsing failed", 
			"date", req.Date,
			"error", err,
			"expected_format", "2006-01-02")
		span.RecordError(err)
		return nil, err
	}
	logger.DebugContext(ctx, "Date parsed successfully", "parsed_date", date.Format("2006-01-02"))

	// Handle timezone configuration
	loc := time.Local // default to local timezone
	logger.DebugContext(ctx, "Processing timezone", "timezone", req.Timezone)
	if req.Timezone != "" {
		parsedLoc, err := time.LoadLocation(req.Timezone)
		if err != nil {
			logger.WarnContext(ctx, "Failed to load timezone, falling back to local", 
				"requested_timezone", req.Timezone, 
				"error", err,
				"fallback_timezone", "local")
			span.AddEvent("Timezone fallback", traceAttributes(
				"requested_timezone", req.Timezone,
				"fallback_timezone", "local",
			)...)
		} else {
			loc = parsedLoc
			logger.DebugContext(ctx, "Timezone loaded successfully", "timezone", req.Timezone)
		}
	} else {
		logger.DebugContext(ctx, "No timezone specified, using local timezone")
	}
	
	// Adjust date to the requested timezone
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	logger.DebugContext(ctx, "Date adjusted to timezone", 
		"adjusted_date", date.Format("2006-01-02 15:04:05 MST"),
		"timezone", loc.String())

	// Prepare location for astronomical calculations
	location := astronomy.Location{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	logger.DebugContext(ctx, "Starting astronomical calculations", 
		"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude))
	
	// Calculate sunrise and sunset times
	logger.InfoContext(ctx, "Calculating sun times", 
		"operation", "CalculateSunTimes",
		"date", date.Format("2006-01-02"),
		"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude))
	sunTimes, err := astronomy.CalculateSunTimes(location, date)
	if err != nil {
		err = status.Error(codes.Internal, fmt.Sprintf("failed to calculate sun times: %v", err))
		logger.ErrorContext(ctx, "Astronomical calculation failed", 
			"operation", "CalculateSunTimes",
			"error", err,
			"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude),
			"date", date.Format("2006-01-02"))
		span.RecordError(err)
		return nil, err
	}
	logger.DebugContext(ctx, "Sun times calculated successfully", 
		"sunrise", sunTimes.Sunrise.Format("15:04:05"),
		"sunset", sunTimes.Sunset.Format("15:04:05"))

	// Simulate random error for testing (as in original code)
	if rand.Intn(10)%2 == 0 {
		err := status.Error(codes.Internal, "failed to fetch panchangam data")
		logger.ErrorContext(ctx, "Simulated random error occurred", 
			"operation", "fetchPanchangamData",
			"error", err,
			"note", "This is a simulated error for testing purposes")
		span.RecordError(err)
		span.AddEvent("Random error simulated")
		return nil, err
	}
	
	// Build panchangam data response
	logger.InfoContext(ctx, "Building panchangam data response", 
		"operation", "buildResponse",
		"date", req.Date,
		"sunrise", sunTimes.Sunrise.Format("15:04:05"),
		"sunset", sunTimes.Sunset.Format("15:04:05"))
	
	data := &ppb.PanchangamData{
		Date:        req.Date,
		Tithi:       "Some Tithi",
		Nakshatra:   "Some Nakshatra",
		Yoga:        "Some Yoga",
		Karana:      "Some Karana",
		SunriseTime: sunTimes.Sunrise.Format("15:04:05"),
		SunsetTime:  sunTimes.Sunset.Format("15:04:05"),
		Events: []*ppb.PanchangamEvent{
			{Name: "Some Event 1", Time: "08:00:00", EventType: "RAHU_KALAM"},
			{Name: "Some Event 2", Time: "12:00:00", EventType: "FESTIVAL"},
		},
	}
	
	logger.InfoContext(ctx, "Panchangam data fetched successfully", 
		"operation", "fetchPanchangamData",
		"date", data.Date,
		"tithi", data.Tithi,
		"nakshatra", data.Nakshatra,
		"events_count", len(data.Events))
	
	span.AddEvent("Data fetch completed", traceAttributes(
		"success", "true",
		"events_count", fmt.Sprintf("%d", len(data.Events)),
	)...)
	
	return data, nil
}
