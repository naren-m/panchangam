package api

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/observability"
)

// CorePanchangamAPI implements the main PanchangamAPI interface
type CorePanchangamAPI struct {
	pluginManager PluginManager
	version       Version
	observer      observability.ObserverInterface
	logger        *observability.ErrorRecorder
}

// NewCorePanchangamAPI creates a new core API instance
func NewCorePanchangamAPI(observer observability.ObserverInterface) *CorePanchangamAPI {
	return &CorePanchangamAPI{
		pluginManager: NewPluginManager(),
		version: Version{
			Major: 1,
			Minor: 0,
			Patch: 0,
			Pre:   "alpha",
		},
		observer: observer,
		logger:   observability.NewErrorRecorder(),
	}
}

// GetPanchangam returns Panchangam data for a specific request
func (api *CorePanchangamAPI) GetPanchangam(ctx context.Context, req PanchangamRequest) (*PanchangamData, error) {
	ctx, span := api.observer.CreateSpan(ctx, "CorePanchangamAPI.GetPanchangam")
	defer span.End()

	// Record operation start
	observability.RecordEvent(ctx, "API request started", map[string]interface{}{
		"operation": "GetPanchangam",
		"date":      req.Date.Format("2006-01-02"),
		"location":  fmt.Sprintf("%.4f,%.4f", req.Location.Latitude, req.Location.Longitude),
		"region":    string(req.Region),
		"method":    string(req.CalculationMethod),
	})

	// Validate request
	if err := api.validateRequest(ctx, req); err != nil {
		observability.RecordError(ctx, err, observability.ErrorContext{
			Severity:  observability.SeverityMedium,
			Category:  observability.CategoryValidation,
			Operation: "validateRequest",
			Component: "core_api",
			Additional: map[string]interface{}{
				"request": req,
			},
			Retryable:   false,
			ExpectedErr: true,
		})
		return nil, err
	}

	// Initialize result with basic information
	result := &PanchangamData{
		Date:              req.Date,
		Location:          req.Location,
		Region:            req.Region,
		CalendarSystem:    req.CalendarSystem,
		CalculationMethod: req.CalculationMethod,
		Locale:            req.Locale,
		Version:           api.version,
		GeneratedAt:       time.Now(),
	}

	// Set defaults if not specified
	if result.Region == "" {
		result.Region = RegionGlobal
	}
	if result.CalendarSystem == "" {
		result.CalendarSystem = CalendarPurnimanta
	}
	if result.CalculationMethod == "" {
		result.CalculationMethod = MethodDrik
	}

	// Calculate basic astronomical data
	if err := api.calculateAstronomicalData(ctx, req, result); err != nil {
		observability.RecordError(ctx, err, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "calculateAstronomicalData",
			Component: "core_api",
			Additional: map[string]interface{}{
				"request": req,
			},
			Retryable:   true,
			ExpectedErr: false,
		})
		return nil, err
	}

	// Calculate Panchangam elements using enabled calculation plugins.
	if err := api.calculatePanchangamElements(ctx, req, result); err != nil {
		observability.RecordError(ctx, err, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "calculatePanchangamElements",
			Component: "core_api",
			Additional: map[string]interface{}{
				"request": req,
			},
			Retryable:   true,
			ExpectedErr: false,
		})
		return nil, err
	}

	// Add events if requested
	if req.IncludeEvents {
		if err := api.calculateEvents(ctx, req, result); err != nil {
			// Log error but don't fail the entire request
			observability.RecordError(ctx, err, observability.ErrorContext{
				Severity:  observability.SeverityMedium,
				Category:  observability.CategoryCalculation,
				Operation: "calculateEvents",
				Component: "core_api",
				Additional: map[string]interface{}{
					"request": req,
				},
				Retryable:   true,
				ExpectedErr: false,
			})
		}
	}

	// Add muhurtas if requested
	if req.IncludeMuhurtas {
		if err := api.calculateMuhurtas(ctx, req, result); err != nil {
			// Log error but don't fail the entire request
			observability.RecordError(ctx, err, observability.ErrorContext{
				Severity:  observability.SeverityMedium,
				Category:  observability.CategoryCalculation,
				Operation: "calculateMuhurtas",
				Component: "core_api",
				Additional: map[string]interface{}{
					"request": req,
				},
				Retryable:   true,
				ExpectedErr: false,
			})
		}
	}

	// Apply regional extensions
	if err := api.applyRegionalExtensions(ctx, result); err != nil {
		// Log error but don't fail the entire request
		observability.RecordError(ctx, err, observability.ErrorContext{
			Severity:  observability.SeverityLow,
			Category:  observability.CategoryCalculation,
			Operation: "applyRegionalExtensions",
			Component: "core_api",
			Additional: map[string]interface{}{
				"region": result.Region,
			},
			Retryable:   true,
			ExpectedErr: false,
		})
	}

	// Apply localization
	if err := api.applyLocalization(ctx, result); err != nil {
		// Log error but don't fail the entire request
		observability.RecordError(ctx, err, observability.ErrorContext{
			Severity:  observability.SeverityLow,
			Category:  observability.CategoryCalculation,
			Operation: "applyLocalization",
			Component: "core_api",
			Additional: map[string]interface{}{
				"locale": result.Locale,
				"region": result.Region,
			},
			Retryable:   true,
			ExpectedErr: false,
		})
	}

	// Record successful completion
	observability.RecordEvent(ctx, "API request completed", map[string]interface{}{
		"operation":     "GetPanchangam",
		"date":          result.Date.Format("2006-01-02"),
		"events_count":  len(result.Events),
		"muhurta_count": len(result.Muhurtas),
		"success":       true,
	})

	return result, nil
}

// GetDateRange returns Panchangam data for a range of dates
func (api *CorePanchangamAPI) GetDateRange(ctx context.Context, start, end time.Time, location Location, options ...RequestOption) ([]*PanchangamData, error) {
	ctx, span := api.observer.CreateSpan(ctx, "CorePanchangamAPI.GetDateRange")
	defer span.End()

	if end.Before(start) {
		return nil, fmt.Errorf("end date %s is before start date %s", end.Format("2006-01-02"), start.Format("2006-01-02"))
	}

	// Create base request
	req := PanchangamRequest{
		Location: location,
	}

	// Apply options
	for _, option := range options {
		option(&req)
	}

	var results []*PanchangamData
	current := start

	for current.Before(end) || current.Equal(end) {
		req.Date = current

		data, err := api.GetPanchangam(ctx, req)
		if err != nil {
			observability.RecordError(ctx, err, observability.ErrorContext{
				Severity:  observability.SeverityMedium,
				Category:  observability.CategoryCalculation,
				Operation: "GetDateRange",
				Component: "core_api",
				Additional: map[string]interface{}{
					"date":  current.Format("2006-01-02"),
					"start": start.Format("2006-01-02"),
					"end":   end.Format("2006-01-02"),
				},
				Retryable:   true,
				ExpectedErr: false,
			})
			return nil, fmt.Errorf("failed to get Panchangam for %s: %w", current.Format("2006-01-02"), err)
		}

		results = append(results, data)
		current = current.AddDate(0, 0, 1)
	}

	return results, nil
}
