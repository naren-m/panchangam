package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
)

// MuhurtaPlugin provides comprehensive auspicious time calculations
type MuhurtaPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewMuhurtaPlugin creates a new muhurta calculation plugin
func NewMuhurtaPlugin() *MuhurtaPlugin {
	return &MuhurtaPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (m *MuhurtaPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "muhurta_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Comprehensive auspicious time calculations including Rahu Kalam, Yamagandam, and traditional muhurtas",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityMuhurta),
		},
		Dependencies: []string{"astronomy"},
		Metadata: map[string]interface{}{
			"muhurta_types":     []string{"rahu_kalam", "yamagandam", "gulikakalam", "abhijit", "brahma_muhurta"},
			"regional_support":  true,
			"calculation_based": "vedic_astronomy",
		},
	}
}

// Initialize sets up the plugin with configuration
func (m *MuhurtaPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	m.config = config
	m.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (m *MuhurtaPlugin) IsEnabled() bool {
	return m.enabled
}

// Shutdown cleans up plugin resources
func (m *MuhurtaPlugin) Shutdown(ctx context.Context) error {
	m.enabled = false
	return nil
}

// GetSupportedRegions returns regions this plugin supports
func (m *MuhurtaPlugin) GetSupportedRegions() []api.Region {
	return []api.Region{
		api.RegionGlobal,
		api.RegionNorthIndia,
		api.RegionSouthIndia,
		api.RegionTamilNadu,
		api.RegionKerala,
		api.RegionBengal,
		api.RegionGujarat,
		api.RegionMaha,
	}
}

// GetMuhurtas returns all muhurtas for a specific date and location
func (m *MuhurtaPlugin) GetMuhurtas(ctx context.Context, date time.Time, location api.Location, region api.Region) ([]api.Muhurta, error) {
	if !m.enabled {
		return nil, fmt.Errorf("muhurta plugin is not enabled")
	}

	var muhurtas []api.Muhurta

	// Calculate sunrise and sunset
	astroLocation := astronomy.Location{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}

	sunTimes, err := astronomy.CalculateSunTimes(astroLocation, date)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate sun times: %w", err)
	}

	if err := convertSunTimesToTimezone(sunTimes, location.Timezone); err != nil {
		return nil, err
	}

	// Calculate day length
	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)

	// Rahu Kalam
	rahuKalam := m.calculateRahuKalam(sunTimes.Sunrise, dayLength, date.Weekday())
	muhurtas = append(muhurtas, rahuKalam)

	// Yamagandam
	yamagandam := m.calculateYamagandam(sunTimes.Sunrise, dayLength, date.Weekday())
	muhurtas = append(muhurtas, yamagandam)

	// Gulika Kalam
	gulikaKalam := m.calculateGulikaKalam(sunTimes.Sunrise, dayLength, date.Weekday())
	muhurtas = append(muhurtas, gulikaKalam)

	// Abhijit Muhurta
	abhijitMuhurta := m.calculateAbhijitMuhurta(sunTimes.Sunrise, sunTimes.Sunset)
	muhurtas = append(muhurtas, abhijitMuhurta)

	// Brahma Muhurta
	brahmaMuhurta := m.calculateBrahmaMuhurta(sunTimes.Sunrise, date)
	muhurtas = append(muhurtas, brahmaMuhurta)

	// Regional specific muhurtas
	regionalMuhurtas := m.getRegionalMuhurtas(ctx, date, location, region, sunTimes, dayLength)
	muhurtas = append(muhurtas, regionalMuhurtas...)

	// Daily auspicious periods
	auspiciousPeriods := m.calculateDailyAuspiciousPeriods(sunTimes.Sunrise, sunTimes.Sunset, dayLength)
	muhurtas = append(muhurtas, auspiciousPeriods...)

	return muhurtas, nil
}

// FindAuspiciousTimes finds auspicious times for specific activities
func (m *MuhurtaPlugin) FindAuspiciousTimes(ctx context.Context, date time.Time, location api.Location, activities []string) ([]api.Muhurta, error) {
	allMuhurtas, err := m.GetMuhurtas(ctx, date, location, api.RegionGlobal)
	if err != nil {
		return nil, err
	}

	var auspiciousTimes []api.Muhurta

	for _, muhurta := range allMuhurtas {
		if muhurta.Quality == api.QualityAuspicious || muhurta.Quality == api.QualityHighly {
			// Check if this muhurta is suitable for the requested activities
			if m.isActivitySuitable(muhurta, activities) {
				auspiciousTimes = append(auspiciousTimes, muhurta)
			}
		}
	}

	return auspiciousTimes, nil
}

// IsTimeAuspicious checks if a specific time is auspicious for given activities
func (m *MuhurtaPlugin) IsTimeAuspicious(ctx context.Context, dateTime time.Time, location api.Location, activities []string) (bool, string, error) {
	date := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, dateTime.Location())

	muhurtas, err := m.GetMuhurtas(ctx, date, location, api.RegionGlobal)
	if err != nil {
		return false, "", err
	}

	for _, muhurta := range muhurtas {
		if dateTime.After(muhurta.StartTime) && dateTime.Before(muhurta.EndTime) {
			if muhurta.Quality == api.QualityInauspicious || muhurta.Quality == api.QualityAvoid {
				return false, fmt.Sprintf("Time falls in %s (%s)", muhurta.Name, muhurta.Quality), nil
			}

			if muhurta.Quality == api.QualityAuspicious || muhurta.Quality == api.QualityHighly {
				if m.isActivitySuitable(muhurta, activities) {
					return true, fmt.Sprintf("Time falls in %s (suitable for %v)", muhurta.Name, activities), nil
				}
			}
		}
	}

	return true, "Time is neutral - no specific restrictions", nil
}
