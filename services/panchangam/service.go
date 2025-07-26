// server.go
package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = log.Logger()

type PanchangamServer struct {
	observer         observability.ObserverInterface
	ephemerisManager *ephemeris.Manager
	tithiCalc        *astronomy.TithiCalculator
	nakshatraCalc    *astronomy.NakshatraCalculator
	yogaCalc         *astronomy.YogaCalculator
	karanaCalc       *astronomy.KaranaCalculator
	varaCalc         *astronomy.VaraCalculator
	ppb.UnimplementedPanchangamServer
}

func NewPanchangamServer() *PanchangamServer {
	// Initialize ephemeris providers and cache
	jplProvider := ephemeris.NewJPLProvider()
	swissProvider := ephemeris.NewSwissProvider()
	cache := ephemeris.NewMemoryCache(1000, 24*time.Hour) // 1000 entries, 24h TTL

	// Create ephemeris manager
	ephemerisManager := ephemeris.NewManager(jplProvider, swissProvider, cache)

	// Initialize calculators
	tithiCalc := astronomy.NewTithiCalculator(ephemerisManager)
	nakshatraCalc := astronomy.NewNakshatraCalculator(ephemerisManager)
	yogaCalc := astronomy.NewYogaCalculator(ephemerisManager)
	karanaCalc := astronomy.NewKaranaCalculator(ephemerisManager)
	varaCalc := astronomy.NewVaraCalculator()

	return &PanchangamServer{
		observer:         observability.Observer(),
		ephemerisManager: ephemerisManager,
		tithiCalc:        tithiCalc,
		nakshatraCalc:    nakshatraCalc,
		yogaCalc:         yogaCalc,
		karanaCalc:       karanaCalc,
		varaCalc:         varaCalc,
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
	// Enhanced validation with comprehensive error recording
	if req.Latitude < -90 || req.Latitude > 90 {
		err := status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")

		// Use enhanced error recording
		observability.RecordValidationFailure(ctx, "latitude", req.Latitude, "latitude must be between -90 and 90 degrees")

		// Record as an important event
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

		// Use enhanced error recording
		observability.RecordValidationFailure(ctx, "longitude", req.Longitude, "longitude must be between -180 and 180 degrees")

		// Record as an important event
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

	// Fetch panchangam data
	logger.InfoContext(ctx, "Fetching panchangam data")

	// Record operation start
	observability.RecordEvent(ctx, "Panchangam data fetch started", map[string]interface{}{
		"operation": "fetchPanchangamData",
		"date":      req.Date,
		"location":  fmt.Sprintf("%.4f,%.4f", req.Latitude, req.Longitude),
		"timezone":  req.Timezone,
		"region":    req.Region,
	})

	d, err := s.fetchPanchangamData(ctx, req)
	if err != nil {
		// Use enhanced error recording
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

		// Record as an important event
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

	// Record successful data fetch
	observability.RecordEvent(ctx, "Panchangam data fetch completed", map[string]interface{}{
		"operation":    "fetchPanchangamData",
		"date":         d.Date,
		"events_count": len(d.Events),
		"success":      true,
	})

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
		grpcErr := status.Error(codes.InvalidArgument, fmt.Sprintf("invalid date format: %v", err))

		// Use enhanced error recording
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityMedium,
			Category:  observability.CategoryValidation,
			Operation: "date_parsing",
			Component: "panchangam_service",
			Additional: map[string]interface{}{
				"date_input":      req.Date,
				"expected_format": "2006-01-02",
				"parse_error":     err.Error(),
			},
			Retryable:   false,
			ExpectedErr: true,
		})

		// Record as an important event
		observability.RecordEvent(ctx, "Date parsing failed", map[string]interface{}{
			"date":            req.Date,
			"expected_format": "2006-01-02",
			"error_type":      "invalid_format",
			"parse_error":     err.Error(),
		})

		logger.WarnContext(ctx, "Date parsing failed",
			"date", req.Date,
			"error", grpcErr,
			"expected_format", "2006-01-02")
		span.RecordError(grpcErr)
		return nil, grpcErr
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
		"operation", "CalculateSunTimesWithContext",
		"date", date.Format("2006-01-02"),
		"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude))

	// Record calculation start
	observability.RecordCalculationStart(ctx, "sun_times_calculation", map[string]interface{}{
		"date":      date.Format("2006-01-02"),
		"latitude":  location.Latitude,
		"longitude": location.Longitude,
		"timezone":  loc.String(),
	})

	calcStart := time.Now()
	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, date)
	calcDuration := time.Since(calcStart)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate sun times: %v", err))

		// Record calculation end with failure
		observability.RecordCalculationEnd(ctx, "sun_times_calculation", false, calcDuration, nil)

		// Use enhanced error recording
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "sun_times_calculation",
			Component: "astronomy_service",
			Additional: map[string]interface{}{
				"calculation_type": "sun_times",
				"latitude":         location.Latitude,
				"longitude":        location.Longitude,
				"date":             date.Format("2006-01-02"),
				"duration_ms":      calcDuration.Milliseconds(),
				"original_error":   err.Error(),
			},
			Retryable:   true,
			ExpectedErr: false,
		})

		// Record as an important event
		observability.RecordEvent(ctx, "Astronomical calculation failed", map[string]interface{}{
			"calculation_type": "sun_times",
			"error_type":       "calculation_failure",
			"location":         fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude),
			"date":             date.Format("2006-01-02"),
			"duration_ms":      calcDuration.Milliseconds(),
		})

		logger.ErrorContext(ctx, "Astronomical calculation failed",
			"operation", "CalculateSunTimesWithContext",
			"error", grpcErr,
			"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude),
			"date", date.Format("2006-01-02"))
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	// Record successful calculation
	observability.RecordCalculationEnd(ctx, "sun_times_calculation", true, calcDuration, map[string]interface{}{
		"sunrise_time": sunTimes.Sunrise.Format("15:04:05"),
		"sunset_time":  sunTimes.Sunset.Format("15:04:05"),
	})
	logger.DebugContext(ctx, "Sun times calculated successfully",
		"sunrise", sunTimes.Sunrise.Format("15:04:05"),
		"sunset", sunTimes.Sunset.Format("15:04:05"))

	// Simulate random error for testing (as in original code)
	if rand.Intn(10)%2 == 0 {
		grpcErr := status.Error(codes.Internal, "failed to fetch panchangam data")

		// Use enhanced error recording
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityMedium,
			Category:  observability.CategoryInternal,
			Operation: "fetchPanchangamData",
			Component: "panchangam_service",
			Additional: map[string]interface{}{
				"error_type":    "simulated_random_error",
				"testing_mode":  true,
				"probability":   "50%",
				"error_purpose": "testing_error_handling",
			},
			Retryable:   true,
			ExpectedErr: true, // This is expected in testing
		})

		// Record as an important event
		observability.RecordEvent(ctx, "Simulated error triggered", map[string]interface{}{
			"error_type":        "random_testing_error",
			"operation":         "fetchPanchangamData",
			"testing_mode":      true,
			"retry_recommended": true,
		})

		logger.ErrorContext(ctx, "Simulated random error occurred",
			"operation", "fetchPanchangamData",
			"error", grpcErr,
			"note", "This is a simulated error for testing purposes")
		span.RecordError(grpcErr)
		span.AddEvent("Random error simulated")
		return nil, grpcErr
	}

	// Calculate all Panchangam elements
	logger.InfoContext(ctx, "Calculating Panchangam elements",
		"operation", "calculatePanchangamElements",
		"date", req.Date,
		"sunrise", sunTimes.Sunrise.Format("15:04:05"),
		"sunset", sunTimes.Sunset.Format("15:04:05"))

	// Calculate Tithi
	tithi, err := s.tithiCalc.GetTithiForDate(ctx, date)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate tithi: %v", err))
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "tithi_calculation",
			Component: "panchangam_service",
		})
		logger.ErrorContext(ctx, "Failed to calculate tithi", "error", grpcErr)
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	// Calculate Nakshatra
	nakshatra, err := s.nakshatraCalc.GetNakshatraForDate(ctx, date)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate nakshatra: %v", err))
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "nakshatra_calculation",
			Component: "panchangam_service",
		})
		logger.ErrorContext(ctx, "Failed to calculate nakshatra", "error", grpcErr)
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	// Calculate Yoga
	yoga, err := s.yogaCalc.GetYogaForDate(ctx, date)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate yoga: %v", err))
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "yoga_calculation",
			Component: "panchangam_service",
		})
		logger.ErrorContext(ctx, "Failed to calculate yoga", "error", grpcErr)
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	// Calculate Karana
	karana, err := s.karanaCalc.GetKaranaForDate(ctx, date)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate karana: %v", err))
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "karana_calculation",
			Component: "panchangam_service",
		})
		logger.ErrorContext(ctx, "Failed to calculate karana", "error", grpcErr)
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	// Calculate Vara
	vara, err := s.varaCalc.GetVaraForDate(ctx, date, location)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate vara: %v", err))
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "vara_calculation",
			Component: "panchangam_service",
		})
		logger.ErrorContext(ctx, "Failed to calculate vara", "error", grpcErr)
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	// Build panchangam data response with real calculations
	logger.InfoContext(ctx, "Building panchangam data response with real calculations",
		"operation", "buildResponse",
		"date", req.Date,
		"tithi", tithi.Name,
		"nakshatra", nakshatra.Name,
		"yoga", yoga.Name,
		"karana", karana.Name,
		"vara", vara.Name,
		"sunrise", sunTimes.Sunrise.Format("15:04:05"),
		"sunset", sunTimes.Sunset.Format("15:04:05"))

	data := &ppb.PanchangamData{
		Date:        req.Date,
		Tithi:       fmt.Sprintf("%s (%d)", tithi.Name, tithi.Number),
		Nakshatra:   fmt.Sprintf("%s (%d)", nakshatra.Name, nakshatra.Number),
		Yoga:        fmt.Sprintf("%s (%d)", yoga.Name, yoga.Number),
		Karana:      fmt.Sprintf("%s (%d)", karana.Name, karana.Number),
		SunriseTime: sunTimes.Sunrise.Format("15:04:05"),
		SunsetTime:  sunTimes.Sunset.Format("15:04:05"),
		Events: []*ppb.PanchangamEvent{
			{Name: fmt.Sprintf("Tithi: %s", tithi.Name), Time: tithi.StartTime.Format("15:04:05"), EventType: "TITHI"},
			{Name: fmt.Sprintf("Nakshatra: %s", nakshatra.Name), Time: "00:00:00", EventType: "NAKSHATRA"},
			{Name: fmt.Sprintf("Yoga: %s", yoga.Name), Time: "00:00:00", EventType: "YOGA"},
			{Name: fmt.Sprintf("Karana: %s", karana.Name), Time: "00:00:00", EventType: "KARANA"},
			{Name: fmt.Sprintf("Vara: %s", vara.Name), Time: "00:00:00", EventType: "VARA"},
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
