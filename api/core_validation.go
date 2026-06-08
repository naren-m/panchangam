package api

import (
	"context"
	"fmt"
)

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
