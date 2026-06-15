package implementations

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// CalculateTithi calculates tithi using the specified method
func (c *CalculationMethodPlugin) CalculateTithi(ctx context.Context, date time.Time, location api.Location, method api.CalculationMethod) (*api.Tithi, error) {
	if !c.enabled {
		return nil, fmt.Errorf("calculation method plugin is not enabled")
	}

	var tithiInfo *astronomy.TithiInfo
	var err error

	switch method {
	case api.MethodDrik:
		tithiInfo, err = c.calculateTithiDrikGanita(ctx, date, location)
	case api.MethodVakya:
		tithiInfo, err = c.calculateTithiVakya(ctx, date, location)
	case api.MethodAuto:
		// Auto method chooses based on date and region
		tithiInfo, err = c.calculateTithiAuto(ctx, date, location)
	default:
		return nil, fmt.Errorf("unsupported calculation method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to calculate tithi using %s method: %w", method, err)
	}

	// Convert to API format
	apiTithi := &api.Tithi{
		Number:     tithiInfo.Number,
		Name:       tithiInfo.Name,
		StartTime:  tithiInfo.StartTime,
		EndTime:    tithiInfo.EndTime,
		Percentage: c.calculateTithiPercentage(tithiInfo, date),
		IsRunning:  c.isTithiRunning(tithiInfo, date),
		Lord:       c.getTithiLord(tithiInfo.Number),
		Quality:    c.getTithiQuality(tithiInfo.Type),
	}

	return apiTithi, nil
}

// calculateTithiDrikGanita uses modern astronomical calculations
func (c *CalculationMethodPlugin) calculateTithiDrikGanita(ctx context.Context, date time.Time, location api.Location) (*astronomy.TithiInfo, error) {
	// Use modern ephemeris for precise planetary positions
	return c.tithiCalculator.GetTithiForDate(ctx, date)
}

// calculateTithiVakya uses traditional Vakya (tabular) calculations
func (c *CalculationMethodPlugin) calculateTithiVakya(ctx context.Context, date time.Time, location api.Location) (*astronomy.TithiInfo, error) {
	// Traditional Vakya constants (simplified implementation)
	vakya := VakyaConstants{
		SunMeanMotion:       0.985647,  // degrees per day
		MoonMeanMotion:      13.176358, // degrees per day
		EpochJD:             588465.5,  // Kaliyuga start (approximate)
		TraditionalAyanamsa: 23.85,     // Traditional ayanamsa value
		SunCorrection:       0.0,       // Traditional corrections
		MoonCorrection:      0.0,
	}

	// Calculate Julian Day
	jd := ephemeris.TimeToJulianDay(date)
	daysSinceEpoch := float64(jd) - vakya.EpochJD

	// Calculate mean longitudes using traditional mean motions
	sunMeanLongitude := math.Mod(vakya.SunMeanMotion*daysSinceEpoch, 360.0)
	moonMeanLongitude := math.Mod(vakya.MoonMeanMotion*daysSinceEpoch, 360.0)

	// Apply traditional corrections (simplified)
	sunTrueLongitude := sunMeanLongitude + vakya.SunCorrection
	moonTrueLongitude := moonMeanLongitude + vakya.MoonCorrection

	// Apply traditional ayanamsa correction
	sunTrueLongitude = sunTrueLongitude - vakya.TraditionalAyanamsa
	moonTrueLongitude = moonTrueLongitude - vakya.TraditionalAyanamsa

	// Normalize to 0-360 range
	sunTrueLongitude = math.Mod(sunTrueLongitude+360, 360)
	moonTrueLongitude = math.Mod(moonTrueLongitude+360, 360)

	// Calculate Tithi from longitudes using traditional method
	return c.tithiCalculator.GetTithiFromLongitudes(ctx, sunTrueLongitude, moonTrueLongitude, date)
}

// calculateTithiAuto automatically chooses the best method
func (c *CalculationMethodPlugin) calculateTithiAuto(ctx context.Context, date time.Time, location api.Location) (*astronomy.TithiInfo, error) {
	// Decision logic for auto method:
	// - Use Drik Ganita for modern dates (after 1900)
	// - Use Vakya for historical dates
	// - Consider regional preferences

	if date.Year() >= 1900 {
		return c.calculateTithiDrikGanita(ctx, date, location)
	} else {
		return c.calculateTithiVakya(ctx, date, location)
	}
}

func (c *CalculationMethodPlugin) calculateTithiPercentage(tithi *astronomy.TithiInfo, currentTime time.Time) float64 {
	if currentTime.Before(tithi.StartTime) || currentTime.After(tithi.EndTime) {
		return 0.0
	}

	totalDuration := tithi.EndTime.Sub(tithi.StartTime)
	elapsed := currentTime.Sub(tithi.StartTime)

	percentage := (elapsed.Seconds() / totalDuration.Seconds()) * 100.0
	if percentage > 100.0 {
		percentage = 100.0
	}
	if percentage < 0.0 {
		percentage = 0.0
	}

	return percentage
}

func (c *CalculationMethodPlugin) isTithiRunning(tithi *astronomy.TithiInfo, currentTime time.Time) bool {
	return !currentTime.Before(tithi.StartTime) && !currentTime.After(tithi.EndTime)
}

func (c *CalculationMethodPlugin) getTithiLord(tithiNumber int) string {
	// Tithi lords based on traditional astronomy
	lords := []string{
		"Agni", "Vayu", "Surya", "Vishnu", "Chandra",
		"Kartikeya", "Indra", "Vasu", "Sarpa", "Dharma",
		"Rudra", "Aditya", "Vishvedeva", "Shiva", "Brahma",
	}

	normalizedNumber := tithiNumber
	if normalizedNumber > 15 {
		normalizedNumber = normalizedNumber - 15
	}

	if normalizedNumber >= 1 && normalizedNumber <= 15 {
		return lords[normalizedNumber-1]
	}

	return ""
}

func (c *CalculationMethodPlugin) getTithiQuality(tithiType astronomy.TithiType) string {
	switch tithiType {
	case astronomy.TithiTypeNanda:
		return "joyful"
	case astronomy.TithiTypeBhadra:
		return "auspicious"
	case astronomy.TithiTypeJaya:
		return "victorious"
	case astronomy.TithiTypeRikta:
		return "empty"
	case astronomy.TithiTypePurna:
		return "complete"
	default:
		return ""
	}
}
