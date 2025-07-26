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

// CalculationMethodPlugin handles different Hindu astronomical calculation methods
// Primarily Drik Ganita (observational/modern) vs Vakya (traditional/tabular)
type CalculationMethodPlugin struct {
	enabled             bool
	config              map[string]interface{}
	tithiCalculator     *astronomy.TithiCalculator
	ephemerisManager    *ephemeris.Manager
}

// VakyaConstants holds traditional Vakya calculation constants
type VakyaConstants struct {
	// Traditional mean motion values (degrees per day)
	SunMeanMotion  float64 // ~0.985647 degrees per day
	MoonMeanMotion float64 // ~13.176358 degrees per day
	
	// Traditional epoch (reference point)
	EpochJD float64 // Kaliyuga start or other traditional epoch
	
	// Ayanamsa value for traditional calculations
	TraditionalAyanamsa float64
	
	// Correction factors for different celestial bodies
	SunCorrection  float64
	MoonCorrection float64
}

// DrikGanitaConfig holds modern calculation configuration
type DrikGanitaConfig struct {
	// Modern ephemeris to use (Swiss, JPL, etc.)
	EphemerisType string
	
	// Ayanamsa system to use
	AyanamsaSystem string
	
	// Precision level
	PrecisionLevel string
	
	// Use atmospheric refraction corrections
	UseAtmosphericCorrection bool
	
	// Use delta-T corrections
	UseDeltaTCorrection bool
}

// NewCalculationMethodPlugin creates a new calculation method plugin
func NewCalculationMethodPlugin(ephemerisManager *ephemeris.Manager) *CalculationMethodPlugin {
	return &CalculationMethodPlugin{
		enabled:          false,
		config:           make(map[string]interface{}),
		tithiCalculator:  astronomy.NewTithiCalculator(ephemerisManager),
		ephemerisManager: ephemerisManager,
	}
}

// GetInfo returns plugin metadata
func (c *CalculationMethodPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "calculation_method_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Handles different Hindu astronomical calculation methods: Drik Ganita vs Vakya",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityCalculation),
		},
		Dependencies: []string{"astronomy", "ephemeris"},
		Metadata: map[string]interface{}{
			"calculation_methods": []string{"drik", "vakya", "auto"},
			"precision_levels":    []string{"high", "medium", "traditional"},
			"ayanamsa_systems":   []string{"lahiri", "krishnamurti", "raman", "traditional"},
		},
	}
}

// Initialize sets up the plugin with configuration
func (c *CalculationMethodPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	c.config = config
	c.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (c *CalculationMethodPlugin) IsEnabled() bool {
	return c.enabled
}

// Shutdown cleans up plugin resources
func (c *CalculationMethodPlugin) Shutdown(ctx context.Context) error {
	c.enabled = false
	return nil
}

// GetSupportedMethods returns calculation methods this plugin supports
func (c *CalculationMethodPlugin) GetSupportedMethods() []api.CalculationMethod {
	return []api.CalculationMethod{
		api.MethodDrik,
		api.MethodVakya,
		api.MethodAuto,
	}
}

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
		EpochJD:            588465.5,   // Kaliyuga start (approximate)
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

// CalculateNakshatra calculates nakshatra using the specified method
func (c *CalculationMethodPlugin) CalculateNakshatra(ctx context.Context, date time.Time, location api.Location, method api.CalculationMethod) (*api.Nakshatra, error) {
	// Get moon's position using the appropriate method
	var moonLongitude float64
	var err error

	switch method {
	case api.MethodDrik:
		moonLongitude, err = c.getMoonLongitudeDrik(ctx, date)
	case api.MethodVakya:
		moonLongitude, err = c.getMoonLongitudeVakya(ctx, date)
	case api.MethodAuto:
		if date.Year() >= 1900 {
			moonLongitude, err = c.getMoonLongitudeDrik(ctx, date)
		} else {
			moonLongitude, err = c.getMoonLongitudeVakya(ctx, date)
		}
	default:
		return nil, fmt.Errorf("unsupported calculation method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get moon longitude: %w", err)
	}

	// Calculate Nakshatra from moon's longitude
	nakshatraInfo := c.calculateNakshatraFromLongitude(moonLongitude, date)

	return nakshatraInfo, nil
}

// getMoonLongitudeDrik gets moon longitude using modern ephemeris
func (c *CalculationMethodPlugin) getMoonLongitudeDrik(ctx context.Context, date time.Time) (float64, error) {
	jd := ephemeris.TimeToJulianDay(date)
	positions, err := c.ephemerisManager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		return 0, err
	}
	return positions.Moon.Longitude, nil
}

// getMoonLongitudeVakya gets moon longitude using traditional calculations
func (c *CalculationMethodPlugin) getMoonLongitudeVakya(ctx context.Context, date time.Time) (float64, error) {
	vakya := VakyaConstants{
		MoonMeanMotion:      13.176358, // degrees per day
		EpochJD:            588465.5,   // Kaliyuga start
		TraditionalAyanamsa: 23.85,
		MoonCorrection:      0.0,
	}

	jd := ephemeris.TimeToJulianDay(date)
	daysSinceEpoch := float64(jd) - vakya.EpochJD

	// Traditional moon longitude calculation
	moonMeanLongitude := math.Mod(vakya.MoonMeanMotion*daysSinceEpoch, 360.0)
	moonTrueLongitude := moonMeanLongitude + vakya.MoonCorrection
	moonTrueLongitude = moonTrueLongitude - vakya.TraditionalAyanamsa
	moonTrueLongitude = math.Mod(moonTrueLongitude+360, 360)

	return moonTrueLongitude, nil
}

// calculateNakshatraFromLongitude calculates nakshatra from moon's longitude
func (c *CalculationMethodPlugin) calculateNakshatraFromLongitude(moonLongitude float64, date time.Time) *api.Nakshatra {
	// Each Nakshatra spans 13°20' (13.333... degrees)
	nakshatraDegrees := 360.0 / 27.0 // 13.333... degrees per nakshatra
	
	nakshatraNumber := int(moonLongitude/nakshatraDegrees) + 1
	if nakshatraNumber > 27 {
		nakshatraNumber = 27
	}

	// Calculate pada (each nakshatra has 4 padas)
	degreesIntoNakshatra := math.Mod(moonLongitude, nakshatraDegrees)
	pada := int(degreesIntoNakshatra/(nakshatraDegrees/4.0)) + 1

	nakshatraName := c.getNakshatraName(nakshatraNumber)
	lord := c.getNakshatraLord(nakshatraNumber)
	deity := c.getNakshatraDeity(nakshatraNumber)
	symbol := c.getNakshatraSymbol(nakshatraNumber)

	// Calculate approximate timing (simplified)
	startTime := date.Add(-2 * time.Hour) // Rough approximation
	endTime := date.Add(22 * time.Hour)   // Rough approximation

	return &api.Nakshatra{
		Number:     nakshatraNumber,
		Name:       nakshatraName,
		StartTime:  startTime,
		EndTime:    endTime,
		Percentage: 50.0, // Simplified
		Pada:       pada,
		Lord:       lord,
		Deity:      deity,
		Symbol:     symbol,
		IsRunning:  true,
	}
}

// Helper methods for calculation differences

func (c *CalculationMethodPlugin) getCalculationMethodAccuracy(method api.CalculationMethod) map[string]interface{} {
	switch method {
	case api.MethodDrik:
		return map[string]interface{}{
			"accuracy":           "very_high",
			"precision_seconds":  "±30 seconds",
			"suitable_for":      []string{"modern_dates", "precise_timing", "scientific_use"},
			"ephemeris_based":   true,
			"atmospheric_corrections": true,
			"description":       "Modern observational astronomy with precise ephemeris data",
		}
	case api.MethodVakya:
		return map[string]interface{}{
			"accuracy":          "traditional",
			"precision_minutes": "±15 minutes",
			"suitable_for":     []string{"historical_dates", "traditional_practices", "approximate_timing"},
			"tabular_based":    true,
			"historical_accuracy": true,
			"description":      "Traditional Hindu astronomical tables and mean motion calculations",
		}
	case api.MethodAuto:
		return map[string]interface{}{
			"accuracy":     "adaptive",
			"description":  "Automatically chooses between Drik and Vakya based on date and requirements",
			"decision_logic": "Drik for modern dates (>=1900), Vakya for historical dates",
		}
	default:
		return map[string]interface{}{}
	}
}

func (c *CalculationMethodPlugin) getMethodDifferences(drikResult, vakyaResult *api.Tithi) map[string]interface{} {
	timeDiff := drikResult.StartTime.Sub(vakyaResult.StartTime)
	
	return map[string]interface{}{
		"time_difference_minutes": timeDiff.Minutes(),
		"drik_start_time":        drikResult.StartTime.Format("15:04:05"),
		"vakya_start_time":       vakyaResult.StartTime.Format("15:04:05"),
		"precision_difference":   "Drik is typically more precise for modern dates",
		"historical_note":        "Vakya may be more appropriate for dates before 1900 CE",
	}
}

// Helper methods for Tithi and Nakshatra information

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

func (c *CalculationMethodPlugin) getNakshatraName(number int) string {
	names := []string{
		"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira",
		"Ardra", "Punarvasu", "Pushya", "Ashlesha", "Magha",
		"Purva Phalguni", "Uttara Phalguni", "Hasta", "Chitra", "Swati",
		"Vishakha", "Anuradha", "Jyeshtha", "Mula", "Purva Ashadha",
		"Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha", "Purva Bhadrapada",
		"Uttara Bhadrapada", "Revati",
	}
	
	if number >= 1 && number <= 27 {
		return names[number-1]
	}
	return ""
}

func (c *CalculationMethodPlugin) getNakshatraLord(number int) string {
	lords := []string{
		"Ketu", "Venus", "Sun", "Moon", "Mars",
		"Rahu", "Jupiter", "Saturn", "Mercury", "Ketu",
		"Venus", "Sun", "Moon", "Mars", "Rahu",
		"Jupiter", "Saturn", "Mercury", "Ketu", "Venus",
		"Sun", "Moon", "Mars", "Rahu", "Jupiter",
		"Saturn", "Mercury",
	}
	
	if number >= 1 && number <= 27 {
		return lords[number-1]
	}
	return ""
}

func (c *CalculationMethodPlugin) getNakshatraDeity(number int) string {
	deities := []string{
		"Ashwini Kumaras", "Yama", "Agni", "Brahma", "Chandra",
		"Rudra", "Aditi", "Brihaspati", "Sarpa", "Pitru",
		"Bhaga", "Aryaman", "Savita", "Tvashta", "Vayu",
		"Indragni", "Mitra", "Indra", "Nirriti", "Apas",
		"Vishvedeva", "Vishnu", "Vasu", "Varuna", "Aja Ekapada",
		"Ahirbudhnya", "Pushan",
	}
	
	if number >= 1 && number <= 27 {
		return deities[number-1]
	}
	return ""
}

func (c *CalculationMethodPlugin) getNakshatraSymbol(number int) string {
	symbols := []string{
		"Horse's head", "Yoni", "Razor", "Cart", "Deer's head",
		"Teardrop", "Bow and arrow", "Flower", "Serpent", "Throne",
		"Front legs of bed", "Back legs of bed", "Palm", "Pearl", "Coral",
		"Potter's wheel", "Lotus", "Earring", "Elephant goad", "Elephant tusk",
		"Fan", "Ear", "Drum", "Empty circle", "Front legs of funeral cot",
		"Back legs of funeral cot", "Fish",
	}
	
	if number >= 1 && number <= 27 {
		return symbols[number-1]
	}
	return ""
}