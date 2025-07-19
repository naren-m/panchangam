package api

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
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

	// Calculate Panchangam elements using plugins or fallback
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

// GetVersion returns the API version
func (api *CorePanchangamAPI) GetVersion() Version {
	return api.version
}

// GetSupportedRegions returns all supported regions
func (api *CorePanchangamAPI) GetSupportedRegions() []Region {
	return []Region{
		RegionGlobal,
		RegionNorthIndia,
		RegionSouthIndia,
		RegionTamilNadu,
		RegionKerala,
		RegionBengal,
		RegionGujarat,
		RegionMaha,
	}
}

// GetSupportedMethods returns all supported calculation methods
func (api *CorePanchangamAPI) GetSupportedMethods() []CalculationMethod {
	return []CalculationMethod{
		MethodDrik,
		MethodVakya,
		MethodAuto,
	}
}

// GetSupportedCalendars returns all supported calendar systems
func (api *CorePanchangamAPI) GetSupportedCalendars() []CalendarSystem {
	return []CalendarSystem{
		CalendarPurnimanta,
		CalendarAmanta,
		CalendarLunar,
		CalendarSolar,
	}
}

// GetPluginManager returns the plugin manager
func (api *CorePanchangamAPI) GetPluginManager() PluginManager {
	return api.pluginManager
}

// RegisterPlugin registers a plugin with the API
func (api *CorePanchangamAPI) RegisterPlugin(plugin Plugin) error {
	return api.pluginManager.RegisterPlugin(plugin)
}

// validateRequest validates the incoming request
func (api *CorePanchangamAPI) validateRequest(ctx context.Context, req PanchangamRequest) error {
	// Validate location
	if req.Location.Latitude < -90 || req.Location.Latitude > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", req.Location.Latitude)
	}
	if req.Location.Longitude < -180 || req.Location.Longitude > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", req.Location.Longitude)
	}

	// Validate date (basic check)
	if req.Date.IsZero() {
		return fmt.Errorf("date is required")
	}

	// Use validation plugins if available
	validationPlugins := api.pluginManager.GetPluginsByCapability(CapabilityValidation)
	for _, plugin := range validationPlugins {
		if validationPlugin, ok := plugin.(ValidationPlugin); ok && plugin.IsEnabled() {
			if err := validationPlugin.ValidateRequest(ctx, req); err != nil {
				return fmt.Errorf("validation plugin %s failed: %w", plugin.GetInfo().Name, err)
			}
		}
	}

	return nil
}

// calculateAstronomicalData calculates basic astronomical data
func (api *CorePanchangamAPI) calculateAstronomicalData(ctx context.Context, req PanchangamRequest, result *PanchangamData) error {
	// Calculate sun and moon times using existing astronomy package
	location := astronomy.Location{
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
	}

	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, req.Date)
	if err != nil {
		return fmt.Errorf("failed to calculate sun times: %w", err)
	}

	result.SunMoonTimes = SunMoonTimes{
		Sunrise:   sunTimes.Sunrise,
		Sunset:    sunTimes.Sunset,
		SolarNoon: sunTimes.Sunrise.Add(sunTimes.Sunset.Sub(sunTimes.Sunrise) / 2),
		DayLength: Duration{sunTimes.Sunset.Sub(sunTimes.Sunrise)},
	}

	// Calculate Julian Day
	result.JulianDay = float64(req.Date.Unix())/86400.0 + 2440587.5

	return nil
}

// calculatePanchangamElements calculates the five main Panchangam elements
func (api *CorePanchangamAPI) calculatePanchangamElements(ctx context.Context, req PanchangamRequest, result *PanchangamData) error {
	// Try calculation plugins first
	calculationPlugins := api.pluginManager.GetPluginsByCapability(CapabilityCalculation)

	for _, plugin := range calculationPlugins {
		if calcPlugin, ok := plugin.(CalculationPlugin); ok && plugin.IsEnabled() {
			// Check if plugin supports this method and region
			if api.pluginSupportsMethodAndRegion(calcPlugin, req.CalculationMethod, req.Region) {
				if err := api.calculateWithPlugin(ctx, calcPlugin, req, result); err == nil {
					return nil // Success with plugin
				} else {
					// Log plugin failure but continue with fallback
					observability.RecordError(ctx, err, observability.ErrorContext{
						Severity:  observability.SeverityMedium,
						Category:  observability.CategoryCalculation,
						Operation: "calculateWithPlugin",
						Component: "core_api",
						Additional: map[string]interface{}{
							"plugin": plugin.GetInfo().Name,
						},
						Retryable:   true,
						ExpectedErr: false,
					})
				}
			}
		}
	}

	// Fallback to default calculations
	return api.calculateWithDefaults(ctx, req, result)
}

// calculateWithPlugin uses a plugin for calculations
func (api *CorePanchangamAPI) calculateWithPlugin(ctx context.Context, plugin CalculationPlugin, req PanchangamRequest, result *PanchangamData) error {
	var err error

	// Calculate each element
	tithiPtr, err := plugin.CalculateTithi(ctx, req.Date, req.Location, req.CalculationMethod)
	if err != nil {
		return fmt.Errorf("plugin tithi calculation failed: %w", err)
	}
	if tithiPtr != nil {
		result.Tithi = *tithiPtr
	}

	nakshatraPtr, err := plugin.CalculateNakshatra(ctx, req.Date, req.Location, req.CalculationMethod)
	if err != nil {
		return fmt.Errorf("plugin nakshatra calculation failed: %w", err)
	}
	if nakshatraPtr != nil {
		result.Nakshatra = *nakshatraPtr
	}

	yogaPtr, err := plugin.CalculateYoga(ctx, req.Date, req.Location, req.CalculationMethod)
	if err != nil {
		return fmt.Errorf("plugin yoga calculation failed: %w", err)
	}
	if yogaPtr != nil {
		result.Yoga = *yogaPtr
	}

	karanaPtr, err := plugin.CalculateKarana(ctx, req.Date, req.Location, req.CalculationMethod)
	if err != nil {
		return fmt.Errorf("plugin karana calculation failed: %w", err)
	}
	if karanaPtr != nil {
		result.Karana = *karanaPtr
	}

	// Set Vara (weekday)
	result.Vara = Vara{
		Number:    int(req.Date.Weekday()),
		Name:      req.Date.Weekday().String(),
		NameLocal: req.Date.Weekday().String(), // Will be localized later
	}

	return nil
}

// calculateWithDefaults provides fallback calculations
func (api *CorePanchangamAPI) calculateWithDefaults(ctx context.Context, req PanchangamRequest, result *PanchangamData) error {
	// Simple fallback calculations (these would be replaced with actual algorithms)
	weekday := req.Date.Weekday()

	result.Tithi = Tithi{
		Number:     1, // Placeholder
		Name:       "Pratipada",
		StartTime:  req.Date,
		EndTime:    req.Date.Add(24 * time.Hour),
		Percentage: 50.0,
		IsRunning:  true,
	}

	result.Nakshatra = Nakshatra{
		Number:     1, // Placeholder
		Name:       "Ashwini",
		StartTime:  req.Date,
		EndTime:    req.Date.Add(24 * time.Hour),
		Percentage: 50.0,
		Pada:       1,
		IsRunning:  true,
	}

	result.Yoga = Yoga{
		Number:     1, // Placeholder
		Name:       "Vishkambha",
		StartTime:  req.Date,
		EndTime:    req.Date.Add(24 * time.Hour),
		Percentage: 50.0,
		IsRunning:  true,
	}

	result.Karana = Karana{
		Number:     1, // Placeholder
		Name:       "Bava",
		StartTime:  req.Date,
		EndTime:    req.Date.Add(12 * time.Hour),
		Percentage: 50.0,
		Type:       "movable",
		IsRunning:  true,
	}

	result.Vara = Vara{
		Number:    int(weekday),
		Name:      weekday.String(),
		NameLocal: weekday.String(),
	}

	return nil
}

// calculateEvents calculates events using event plugins
func (api *CorePanchangamAPI) calculateEvents(ctx context.Context, req PanchangamRequest, result *PanchangamData) error {
	eventPlugins := api.pluginManager.GetPluginsByCapability(CapabilityEvent)

	for _, plugin := range eventPlugins {
		if eventPlugin, ok := plugin.(EventPlugin); ok && plugin.IsEnabled() {
			events, err := eventPlugin.GetEvents(ctx, req.Date, req.Location, req.Region)
			if err != nil {
				return fmt.Errorf("event plugin %s failed: %w", plugin.GetInfo().Name, err)
			}
			result.Events = append(result.Events, events...)
		}
	}

	return nil
}

// calculateMuhurtas calculates muhurtas using muhurta plugins
func (api *CorePanchangamAPI) calculateMuhurtas(ctx context.Context, req PanchangamRequest, result *PanchangamData) error {
	muhurtaPlugins := api.pluginManager.GetPluginsByCapability(CapabilityMuhurta)

	for _, plugin := range muhurtaPlugins {
		if muhurtaPlugin, ok := plugin.(MuhurtaPlugin); ok && plugin.IsEnabled() {
			muhurtas, err := muhurtaPlugin.GetMuhurtas(ctx, req.Date, req.Location, req.Region)
			if err != nil {
				return fmt.Errorf("muhurta plugin %s failed: %w", plugin.GetInfo().Name, err)
			}
			result.Muhurtas = append(result.Muhurtas, muhurtas...)
		}
	}

	return nil
}

// applyRegionalExtensions applies region-specific modifications
func (api *CorePanchangamAPI) applyRegionalExtensions(ctx context.Context, result *PanchangamData) error {
	regionalPlugins := api.pluginManager.GetPluginsByCapability(CapabilityRegional)

	for _, plugin := range regionalPlugins {
		if regionalPlugin, ok := plugin.(RegionalExtension); ok && plugin.IsEnabled() {
			if regionalPlugin.GetRegion() == result.Region {
				if err := regionalPlugin.ApplyRegionalRules(ctx, result); err != nil {
					return fmt.Errorf("regional plugin %s failed: %w", plugin.GetInfo().Name, err)
				}
			}
		}
	}

	return nil
}

// applyLocalization applies localization using localization plugins
func (api *CorePanchangamAPI) applyLocalization(ctx context.Context, result *PanchangamData) error {
	if result.Locale == "" {
		return nil // No localization needed
	}

	localizationPlugins := api.pluginManager.GetPluginsByCapability(CapabilityLocalization)

	for _, plugin := range localizationPlugins {
		if locPlugin, ok := plugin.(LocalizationPlugin); ok && plugin.IsEnabled() {
			// Localize each element
			if err := locPlugin.LocalizeTithi(&result.Tithi, result.Locale, result.Region); err != nil {
				return fmt.Errorf("localization plugin %s failed for tithi: %w", plugin.GetInfo().Name, err)
			}
			if err := locPlugin.LocalizeNakshatra(&result.Nakshatra, result.Locale, result.Region); err != nil {
				return fmt.Errorf("localization plugin %s failed for nakshatra: %w", plugin.GetInfo().Name, err)
			}
			if err := locPlugin.LocalizeYoga(&result.Yoga, result.Locale, result.Region); err != nil {
				return fmt.Errorf("localization plugin %s failed for yoga: %w", plugin.GetInfo().Name, err)
			}
			if err := locPlugin.LocalizeKarana(&result.Karana, result.Locale, result.Region); err != nil {
				return fmt.Errorf("localization plugin %s failed for karana: %w", plugin.GetInfo().Name, err)
			}

			// Localize events and muhurtas
			for i := range result.Events {
				if err := locPlugin.LocalizeEvent(&result.Events[i], result.Locale, result.Region); err != nil {
					return fmt.Errorf("localization plugin %s failed for event: %w", plugin.GetInfo().Name, err)
				}
			}
			for i := range result.Muhurtas {
				if err := locPlugin.LocalizeMuhurta(&result.Muhurtas[i], result.Locale, result.Region); err != nil {
					return fmt.Errorf("localization plugin %s failed for muhurta: %w", plugin.GetInfo().Name, err)
				}
			}
		}
	}

	return nil
}

// pluginSupportsMethodAndRegion checks if a plugin supports given method and region
func (api *CorePanchangamAPI) pluginSupportsMethodAndRegion(plugin CalculationPlugin, method CalculationMethod, region Region) bool {
	supportedMethods := plugin.GetSupportedMethods()
	supportedRegions := plugin.GetSupportedRegions()

	methodSupported := false
	for _, m := range supportedMethods {
		if m == method {
			methodSupported = true
			break
		}
	}

	regionSupported := false
	for _, r := range supportedRegions {
		if r == region || r == RegionGlobal {
			regionSupported = true
			break
		}
	}

	return methodSupported && regionSupported
}
