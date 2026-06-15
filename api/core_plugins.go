package api

import (
	"context"
	"fmt"
	"slices"
)

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

	return slices.Contains(supportedMethods, method) &&
		(slices.Contains(supportedRegions, region) || slices.Contains(supportedRegions, RegionGlobal))
}
