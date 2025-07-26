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
			string(api.CapabilityEvent),
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

	// Convert to local timezone if provided
	if location.Timezone != "" {
		if tz, err := time.LoadLocation(location.Timezone); err == nil {
			sunTimes.Sunrise = sunTimes.Sunrise.In(tz)
			sunTimes.Sunset = sunTimes.Sunset.In(tz)
		}
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

// calculateRahuKalam calculates Rahu Kalam based on weekday and sunrise
func (m *MuhurtaPlugin) calculateRahuKalam(sunrise time.Time, dayLength time.Duration, weekday time.Weekday) api.Muhurta {
	// Rahu Kalam calculation based on traditional Vedic astronomy
	// Each day is divided into 8 parts, Rahu Kalam occupies one part based on weekday
	
	partDuration := dayLength / 8
	var rahuPart int
	
	switch weekday {
	case time.Sunday:
		rahuPart = 4 // 5th part (12:00-1:30 PM approximately)
	case time.Monday:
		rahuPart = 1 // 2nd part (7:30-9:00 AM approximately)
	case time.Tuesday:
		rahuPart = 6 // 7th part (3:00-4:30 PM approximately)
	case time.Wednesday:
		rahuPart = 3 // 4th part (10:30 AM-12:00 PM approximately)
	case time.Thursday:
		rahuPart = 5 // 6th part (1:30-3:00 PM approximately)
	case time.Friday:
		rahuPart = 2 // 3rd part (9:00-10:30 AM approximately)
	case time.Saturday:
		rahuPart = 7 // 8th part (4:30-6:00 PM approximately)
	}
	
	startTime := sunrise.Add(time.Duration(rahuPart) * partDuration)
	endTime := startTime.Add(partDuration)
	
	return api.Muhurta{
		Name:         "Rahu Kalam",
		NameLocal:    "राहु काल",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityInauspicious,
		Purpose:      []string{}, // Not suitable for any auspicious activities
		Avoid:        []string{"all_auspicious_activities", "new_ventures", "travel", "ceremonies"},
		Significance: "Inauspicious period ruled by Rahu, avoid starting new activities",
		Metadata: map[string]interface{}{
			"planetary_ruler": "Rahu",
			"weekday_part":    rahuPart + 1,
			"total_parts":     8,
			"calculation":     "vedic_astronomy",
		},
	}
}

// calculateYamagandam calculates Yamagandam based on weekday and sunrise
func (m *MuhurtaPlugin) calculateYamagandam(sunrise time.Time, dayLength time.Duration, weekday time.Weekday) api.Muhurta {
	// Yamagandam calculation - another inauspicious period
	partDuration := dayLength / 8
	var yamaPart int
	
	switch weekday {
	case time.Sunday:
		yamaPart = 2 // 3rd part
	case time.Monday:
		yamaPart = 5 // 6th part
	case time.Tuesday:
		yamaPart = 0 // 1st part
	case time.Wednesday:
		yamaPart = 4 // 5th part
	case time.Thursday:
		yamaPart = 6 // 7th part
	case time.Friday:
		yamaPart = 3 // 4th part
	case time.Saturday:
		yamaPart = 1 // 2nd part
	}
	
	startTime := sunrise.Add(time.Duration(yamaPart) * partDuration)
	endTime := startTime.Add(partDuration)
	
	return api.Muhurta{
		Name:         "Yamagandam",
		NameLocal:    "यमगण्डम्",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityInauspicious,
		Purpose:      []string{},
		Avoid:        []string{"important_decisions", "travel", "business_deals", "ceremonies"},
		Significance: "Inauspicious period ruled by Yama, lord of death",
		Metadata: map[string]interface{}{
			"planetary_ruler": "Yama",
			"weekday_part":    yamaPart + 1,
			"total_parts":     8,
		},
	}
}

// calculateGulikaKalam calculates Gulika Kalam (similar to Rahu Kalam but different timing)
func (m *MuhurtaPlugin) calculateGulikaKalam(sunrise time.Time, dayLength time.Duration, weekday time.Weekday) api.Muhurta {
	partDuration := dayLength / 8
	var gulikaPart int
	
	switch weekday {
	case time.Sunday:
		gulikaPart = 6 // 7th part
	case time.Monday:
		gulikaPart = 3 // 4th part
	case time.Tuesday:
		gulikaPart = 4 // 5th part
	case time.Wednesday:
		gulikaPart = 5 // 6th part
	case time.Thursday:
		gulikaPart = 2 // 3rd part
	case time.Friday:
		gulikaPart = 7 // 8th part
	case time.Saturday:
		gulikaPart = 0 // 1st part
	}
	
	startTime := sunrise.Add(time.Duration(gulikaPart) * partDuration)
	endTime := startTime.Add(partDuration)
	
	return api.Muhurta{
		Name:         "Gulika Kalam",
		NameLocal:    "गुलिक काल",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityInauspicious,
		Purpose:      []string{},
		Avoid:        []string{"financial_transactions", "investments", "property_deals"},
		Significance: "Inauspicious period ruled by Gulika (son of Shani)",
		Metadata: map[string]interface{}{
			"planetary_ruler": "Gulika",
			"weekday_part":    gulikaPart + 1,
			"total_parts":     8,
		},
	}
}

// calculateAbhijitMuhurta calculates the auspicious Abhijit period around noon
func (m *MuhurtaPlugin) calculateAbhijitMuhurta(sunrise, sunset time.Time) api.Muhurta {
	// Abhijit is the 8th nakshatra period, approximately 24 minutes around solar noon
	dayLength := sunset.Sub(sunrise)
	solarNoon := sunrise.Add(dayLength / 2)
	
	// Abhijit duration is approximately 24 minutes (1/60th of a day)
	abhijitDuration := 24 * time.Minute
	startTime := solarNoon.Add(-abhijitDuration / 2)
	endTime := solarNoon.Add(abhijitDuration / 2)
	
	return api.Muhurta{
		Name:         "Abhijit Muhurta",
		NameLocal:    "अभिजित मुहूर्त",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityHighly,
		Purpose:      []string{"all_auspicious_activities", "business_ventures", "education", "spiritual_practices"},
		Avoid:        []string{},
		Significance: "Highly auspicious period around solar noon, victory and success assured",
		Metadata: map[string]interface{}{
			"nakshatra":    "Abhijit",
			"duration":     "24_minutes",
			"calculation":  "solar_noon_based",
			"deity":        "Brahma",
			"significance": "victory_success",
		},
	}
}

// calculateBrahmaMuhurta calculates the pre-dawn auspicious period
func (m *MuhurtaPlugin) calculateBrahmaMuhurta(sunrise time.Time, date time.Time) api.Muhurta {
	// Brahma Muhurta is approximately 1.5 hours before sunrise
	// It's considered the most auspicious time for spiritual practices
	
	brahmaDuration := 96 * time.Minute // 1 hour 36 minutes (1/15th of day-night cycle)
	endTime := sunrise
	startTime := endTime.Add(-brahmaDuration)
	
	// Ensure it's on the same date (early morning)
	if startTime.Day() != date.Day() {
		startTime = time.Date(date.Year(), date.Month(), date.Day(), 4, 0, 0, 0, date.Location())
		endTime = startTime.Add(brahmaDuration)
	}
	
	return api.Muhurta{
		Name:         "Brahma Muhurta",
		NameLocal:    "ब्रह्म मुहूर्त",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityHighly,
		Purpose:      []string{"meditation", "prayer", "study", "yoga", "spiritual_practices"},
		Avoid:        []string{},
		Significance: "Most auspicious pre-dawn period for spiritual activities",
		Metadata: map[string]interface{}{
			"duration":      "96_minutes",
			"calculation":   "pre_sunrise",
			"deity":         "Brahma",
			"time_of_day":   "pre_dawn",
			"spiritual_significance": "highest",
		},
	}
}

// getRegionalMuhurtas returns region-specific muhurtas
func (m *MuhurtaPlugin) getRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location, region api.Region, sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta
	
	switch region {
	case api.RegionTamilNadu:
		muhurtas = append(muhurtas, m.getTamilMuhurtas(sunTimes, dayLength)...)
	case api.RegionKerala:
		muhurtas = append(muhurtas, m.getKeralaMuhurtas(sunTimes, dayLength)...)
	case api.RegionBengal:
		muhurtas = append(muhurtas, m.getBengalMuhurtas(sunTimes, dayLength)...)
	}
	
	return muhurtas
}

// getTamilMuhurtas returns Tamil Nadu specific muhurtas
func (m *MuhurtaPlugin) getTamilMuhurtas(sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta
	
	// Naazhikai-based calculations (Tamil time units)
	// 1 Naazhikai = 24 minutes, 60 Naazhikai = 1 day
	naazhikaiDuration := 24 * time.Minute
	
	// Shubha Muhurta (auspicious period in Tamil tradition)
	// Typically 2 Naazhikai in the morning
	shubhaStart := sunTimes.Sunrise.Add(2 * naazhikaiDuration)
	shubhaEnd := shubhaStart.Add(2 * naazhikaiDuration)
	
	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Shubha Muhurta",
		NameLocal:    "ஶுப முஹூர்த்த",
		StartTime:    shubhaStart,
		EndTime:      shubhaEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"business", "education", "ceremonies"},
		Significance: "Tamil traditional auspicious period",
		Region:       api.RegionTamilNadu,
		Metadata: map[string]interface{}{
			"duration_naazhikai": 2,
			"calculation_system": "tamil_naazhikai",
		},
	})
	
	return muhurtas
}

// getKeralaMuhurtas returns Kerala specific muhurtas
func (m *MuhurtaPlugin) getKeralaMuhurtas(sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta
	
	// Kerala specific Malyalam calendar muhurtas
	// Uchcha Kalam (auspicious time)
	ucchaStart := sunTimes.Sunrise.Add(dayLength / 4)
	ucchaEnd := ucchaStart.Add(90 * time.Minute)
	
	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Uchcha Kalam",
		NameLocal:    "ഉച്ച കാലം",
		StartTime:    ucchaStart,
		EndTime:      ucchaEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"religious_ceremonies", "important_decisions"},
		Significance: "Kerala traditional auspicious period",
		Region:       api.RegionKerala,
		Metadata: map[string]interface{}{
			"calculation_system": "malayalam_calendar",
		},
	})
	
	return muhurtas
}

// getBengalMuhurtas returns Bengal specific muhurtas
func (m *MuhurtaPlugin) getBengalMuhurtas(sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta
	
	// Labha Kaal (beneficial time in Bengali tradition)
	labhaStart := sunTimes.Sunrise.Add(dayLength / 3)
	labhaEnd := labhaStart.Add(72 * time.Minute)
	
	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Labha Kaal",
		NameLocal:    "লাভ কাল",
		StartTime:    labhaStart,
		EndTime:      labhaEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"business", "financial_transactions", "investments"},
		Significance: "Bengali traditional beneficial period for gains",
		Region:       api.RegionBengal,
		Metadata: map[string]interface{}{
			"calculation_system": "bengali_calendar",
		},
	})
	
	return muhurtas
}

// calculateDailyAuspiciousPeriods calculates general auspicious periods throughout the day
func (m *MuhurtaPlugin) calculateDailyAuspiciousPeriods(sunrise, sunset time.Time, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta
	
	// Pratah Kaal (Morning period) - First 3 hours after sunrise
	pratahEnd := sunrise.Add(3 * time.Hour)
	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Pratah Kaal",
		NameLocal:    "प्रातः काल",
		StartTime:    sunrise,
		EndTime:      pratahEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"daily_activities", "exercise", "study", "work"},
		Significance: "Morning period, good for daily activities",
		Metadata: map[string]interface{}{
			"period": "morning",
			"duration_hours": 3,
		},
	})
	
	// Sandhya Kaal (Evening twilight) - Last hour before sunset
	sandhyaStart := sunset.Add(-1 * time.Hour)
	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Sandhya Kaal",
		NameLocal:    "संध्या काल",
		StartTime:    sandhyaStart,
		EndTime:      sunset,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"prayer", "meditation", "spiritual_practices"},
		Significance: "Evening twilight, good for spiritual activities",
		Metadata: map[string]interface{}{
			"period": "evening_twilight",
			"duration_hours": 1,
		},
	})
	
	return muhurtas
}

// isActivitySuitable checks if a muhurta is suitable for given activities
func (m *MuhurtaPlugin) isActivitySuitable(muhurta api.Muhurta, activities []string) bool {
	if len(muhurta.Purpose) == 0 {
		return false
	}
	
	// Check if any of the requested activities match the muhurta's purposes
	for _, activity := range activities {
		for _, purpose := range muhurta.Purpose {
			if activity == purpose || purpose == "all_auspicious_activities" {
				return true
			}
		}
	}
	
	// Check if any activities are in the avoid list
	for _, activity := range activities {
		for _, avoid := range muhurta.Avoid {
			if activity == avoid || avoid == "all_auspicious_activities" {
				return false
			}
		}
	}
	
	return false
}