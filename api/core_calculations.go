package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
)

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
	method := req.CalculationMethod
	if method == "" {
		method = result.CalculationMethod
	}
	region := req.Region
	if region == "" {
		region = result.Region
	}
	req.CalculationMethod = method
	req.Region = region

	calculationPlugins := api.pluginManager.GetPluginsByCapability(CapabilityCalculation)
	var failedPlugins []string

	for _, plugin := range calculationPlugins {
		if calcPlugin, ok := plugin.(CalculationPlugin); ok && plugin.IsEnabled() {
			if api.pluginSupportsMethodAndRegion(calcPlugin, method, region) {
				if err := api.calculateWithPlugin(ctx, calcPlugin, req, result); err == nil {
					return nil
				} else {
					failedPlugins = append(failedPlugins, plugin.GetInfo().Name)
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

	if len(failedPlugins) > 0 {
		return fmt.Errorf(
			"no enabled calculation plugin completed for method %q and region %q; failed plugins: %s",
			method,
			region,
			strings.Join(failedPlugins, ", "),
		)
	}

	return fmt.Errorf("no enabled calculation plugin for method %q and region %q", method, region)
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
