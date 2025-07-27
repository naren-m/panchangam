package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// CalendarSystemPlugin handles different Hindu calendar systems
// Primarily Amanta (South Indian) vs Purnimanta (North Indian) month calculations
type CalendarSystemPlugin struct {
	enabled             bool
	config              map[string]interface{}
	tithiCalculator     *astronomy.TithiCalculator
	ephemerisManager    *ephemeris.Manager
}

// MonthInfo represents lunar month information
type MonthInfo struct {
	Name           string                `json:"name"`
	NameLocal      string                `json:"name_local"`
	Number         int                   `json:"number"`         // 1-12
	StartDate      time.Time             `json:"start_date"`
	EndDate        time.Time             `json:"end_date"`
	CalendarSystem api.CalendarSystem    `json:"calendar_system"`
	Region         api.Region            `json:"region"`
	Year           int                   `json:"year"`
	IsAdhikaMasa   bool                  `json:"is_adhika_masa"` // Intercalary month
	PrevailingTithi *astronomy.TithiInfo `json:"prevailing_tithi,omitempty"`
}

// NewCalendarSystemPlugin creates a new calendar system plugin
func NewCalendarSystemPlugin(ephemerisManager *ephemeris.Manager) *CalendarSystemPlugin {
	return &CalendarSystemPlugin{
		enabled:          false,
		config:           make(map[string]interface{}),
		tithiCalculator:  astronomy.NewTithiCalculator(ephemerisManager),
		ephemerisManager: ephemerisManager,
	}
}

// GetInfo returns plugin metadata
func (c *CalendarSystemPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "calendar_system_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Handles different Hindu calendar systems: Amanta (South Indian) vs Purnimanta (North Indian)",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
			string(api.CapabilityCalculation),
		},
		Dependencies: []string{"astronomy", "ephemeris"},
		Metadata: map[string]interface{}{
			"calendar_systems": []string{"amanta", "purnimanta", "lunar", "solar"},
			"regional_support": true,
			"month_calculation": "astronomical",
			"adhika_masa_support": true,
		},
	}
}

// Initialize sets up the plugin with configuration
func (c *CalendarSystemPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	c.config = config
	c.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (c *CalendarSystemPlugin) IsEnabled() bool {
	return c.enabled
}

// Shutdown cleans up plugin resources
func (c *CalendarSystemPlugin) Shutdown(ctx context.Context) error {
	c.enabled = false
	return nil
}

// GetSupportedRegions returns regions this plugin supports
func (c *CalendarSystemPlugin) GetSupportedRegions() []api.Region {
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

// GetSupportedMethods returns calculation methods this plugin supports
func (c *CalendarSystemPlugin) GetSupportedMethods() []api.CalculationMethod {
	return []api.CalculationMethod{
		api.MethodDrik,
		api.MethodVakya,
		api.MethodAuto,
	}
}

// ApplyCalendarSystem adjusts Panchangam data based on the calendar system
func (c *CalendarSystemPlugin) ApplyCalendarSystem(ctx context.Context, data *api.PanchangamData) error {
	if !c.enabled {
		return fmt.Errorf("calendar system plugin is not enabled")
	}

	// Get current month information based on calendar system
	monthInfo, err := c.GetCurrentMonth(ctx, data.Date, data.Location, data.CalendarSystem, data.Region)
	if err != nil {
		return fmt.Errorf("failed to get month information: %w", err)
	}

	// Adjust month names and numbering based on calendar system
	if err := c.adjustMonthData(ctx, data, monthInfo); err != nil {
		return fmt.Errorf("failed to adjust month data: %w", err)
	}

	// Add calendar system specific metadata
	c.addCalendarSystemMetadata(data, monthInfo)

	return nil
}

// GetCurrentMonth returns current lunar month information based on calendar system
func (c *CalendarSystemPlugin) GetCurrentMonth(ctx context.Context, date time.Time, location api.Location, calendarSystem api.CalendarSystem, region api.Region) (*MonthInfo, error) {
	// Get Tithi for the date to determine month boundaries
	tithi, err := c.tithiCalculator.GetTithiForDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate tithi: %w", err)
	}

	var monthInfo *MonthInfo

	switch calendarSystem {
	case api.CalendarAmanta:
		monthInfo, err = c.calculateAmantaMonth(ctx, date, location, region, tithi)
	case api.CalendarPurnimanta:
		monthInfo, err = c.calculatePurnimantaMonth(ctx, date, location, region, tithi)
	case api.CalendarLunar:
		// Default to region-specific preference
		if c.isNorthIndianRegion(region) {
			monthInfo, err = c.calculatePurnimantaMonth(ctx, date, location, region, tithi)
		} else {
			monthInfo, err = c.calculateAmantaMonth(ctx, date, location, region, tithi)
		}
	case api.CalendarSolar:
		monthInfo, err = c.calculateSolarMonth(ctx, date, location, region)
	default:
		return nil, fmt.Errorf("unsupported calendar system: %s", calendarSystem)
	}

	if err != nil {
		return nil, err
	}

	return monthInfo, nil
}

// calculateAmantaMonth calculates month boundaries for Amanta system (month ends on Amavasya/New Moon)
func (c *CalendarSystemPlugin) calculateAmantaMonth(ctx context.Context, date time.Time, location api.Location, region api.Region, currentTithi *astronomy.TithiInfo) (*MonthInfo, error) {
	// In Amanta system, month begins the day after Amavasya and ends on next Amavasya
	
	// Find the most recent Amavasya (start of current month)
	monthStart, err := c.findPreviousAmavasya(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to find month start: %w", err)
	}

	// Find the next Amavasya (end of current month)
	monthEnd, err := c.findNextAmavasya(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to find month end: %w", err)
	}

	// Determine month number and name
	monthNumber, monthName := c.getAmantaMonthInfo(date, region)

	localName := c.getLocalMonthName(monthName, region)

	return &MonthInfo{
		Name:           monthName,
		NameLocal:      localName,
		Number:         monthNumber,
		StartDate:      monthStart,
		EndDate:        monthEnd,
		CalendarSystem: api.CalendarAmanta,
		Region:         region,
		Year:           date.Year(),
		IsAdhikaMasa:   c.isAdhikaMasa(ctx, monthStart, monthEnd),
		PrevailingTithi: currentTithi,
	}, nil
}

// calculatePurnimantaMonth calculates month boundaries for Purnimanta system (month ends on Purnima/Full Moon)
func (c *CalendarSystemPlugin) calculatePurnimantaMonth(ctx context.Context, date time.Time, location api.Location, region api.Region, currentTithi *astronomy.TithiInfo) (*MonthInfo, error) {
	// In Purnimanta system, month begins the day after Purnima and ends on next Purnima
	
	// Find the most recent Purnima (start of current month)
	monthStart, err := c.findPreviousPurnima(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to find month start: %w", err)
	}

	// Find the next Purnima (end of current month)
	monthEnd, err := c.findNextPurnima(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to find month end: %w", err)
	}

	// Determine month number and name
	monthNumber, monthName := c.getPurnimantaMonthInfo(date, region)

	localName := c.getLocalMonthName(monthName, region)

	return &MonthInfo{
		Name:           monthName,
		NameLocal:      localName,
		Number:         monthNumber,
		StartDate:      monthStart,
		EndDate:        monthEnd,
		CalendarSystem: api.CalendarPurnimanta,
		Region:         region,
		Year:           date.Year(),
		IsAdhikaMasa:   c.isAdhikaMasa(ctx, monthStart, monthEnd),
		PrevailingTithi: currentTithi,
	}, nil
}

// calculateSolarMonth calculates solar month (based on sun's position in zodiac)
func (c *CalendarSystemPlugin) calculateSolarMonth(ctx context.Context, date time.Time, location api.Location, region api.Region) (*MonthInfo, error) {
	// Solar months are based on sun's transit through zodiac signs
	// Each solar month begins when sun enters a new zodiac sign
	
	// Get sun's longitude
	jd := ephemeris.TimeToJulianDay(date)
	positions, err := c.ephemerisManager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		return nil, fmt.Errorf("failed to get planetary positions: %w", err)
	}

	sunLongitude := positions.Sun.Longitude

	// Determine solar month based on sun's longitude
	monthNumber, monthName := c.getSolarMonthInfo(sunLongitude, region)
	localName := c.getLocalMonthName(monthName, region)

	// Calculate approximate month boundaries (solar months are ~30 days)
	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)

	return &MonthInfo{
		Name:           monthName,
		NameLocal:      localName,
		Number:         monthNumber,
		StartDate:      monthStart,
		EndDate:        monthEnd,
		CalendarSystem: api.CalendarSolar,
		Region:         region,
		Year:           date.Year(),
		IsAdhikaMasa:   false,
	}, nil
}

// Helper methods for finding lunar events

func (c *CalendarSystemPlugin) findPreviousAmavasya(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search backwards for the most recent Amavasya
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days back
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 30 { // Amavasya
			return tithi.StartTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, -1)
	}

	return time.Time{}, fmt.Errorf("could not find previous Amavasya")
}

func (c *CalendarSystemPlugin) findNextAmavasya(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search forwards for the next Amavasya
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days forward
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 30 { // Amavasya
			return tithi.EndTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, 1)
	}

	return time.Time{}, fmt.Errorf("could not find next Amavasya")
}

func (c *CalendarSystemPlugin) findPreviousPurnima(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search backwards for the most recent Purnima
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days back
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 15 { // Purnima
			return tithi.StartTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, -1)
	}

	return time.Time{}, fmt.Errorf("could not find previous Purnima")
}

func (c *CalendarSystemPlugin) findNextPurnima(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search forwards for the next Purnima
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days forward
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 15 { // Purnima
			return tithi.EndTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, 1)
	}

	return time.Time{}, fmt.Errorf("could not find next Purnima")
}

// Month name and numbering logic

func (c *CalendarSystemPlugin) getAmantaMonthInfo(date time.Time, region api.Region) (int, string) {
	// Amanta month names based on when the month's Amavasya falls
	// This is a simplified implementation
	monthNames := []string{
		"Chaitra", "Vaisakha", "Jyeshtha", "Ashadha",
		"Shravana", "Bhadrapada", "Ashwin", "Kartik",
		"Margashirsha", "Pausha", "Magha", "Phalgun",
	}

	// Simplified mapping based on Gregorian calendar
	// In actual implementation, this would be based on precise lunar calculations
	monthIndex := (int(date.Month()) + 10) % 12 // Rough approximation
	return monthIndex + 1, monthNames[monthIndex]
}

func (c *CalendarSystemPlugin) getPurnimantaMonthInfo(date time.Time, region api.Region) (int, string) {
	// Purnimanta month names based on when the month's Purnima falls
	// The key difference is that Purnimanta months start ~15 days later than Amanta
	monthNames := []string{
		"Chaitra", "Vaisakha", "Jyeshtha", "Ashadha",
		"Shravana", "Bhadrapada", "Ashwin", "Kartik",
		"Margashirsha", "Pausha", "Magha", "Phalgun",
	}

	// For Purnimanta, adjust by approximately 15 days
	adjustedDate := date.AddDate(0, 0, -15)
	monthIndex := (int(adjustedDate.Month()) + 10) % 12
	return monthIndex + 1, monthNames[monthIndex]
}

func (c *CalendarSystemPlugin) getSolarMonthInfo(sunLongitude float64, region api.Region) (int, string) {
	// Solar months based on sun's zodiac position
	solarMonths := []string{
		"Mesha", "Vrishabha", "Mithuna", "Karkataka",
		"Simha", "Kanya", "Tula", "Vrishchika",
		"Dhanus", "Makara", "Kumbha", "Meena",
	}

	// Each zodiac sign is 30 degrees
	monthIndex := int(sunLongitude / 30.0)
	if monthIndex >= 12 {
		monthIndex = 11
	}

	return monthIndex + 1, solarMonths[monthIndex]
}

func (c *CalendarSystemPlugin) getLocalMonthName(sanskritName string, region api.Region) string {
	// Regional month name mappings
	monthMappings := map[api.Region]map[string]string{
		api.RegionTamilNadu: {
			"Chaitra":     "சித்திரை",
			"Vaisakha":    "வைகாசி",
			"Jyeshtha":    "ஜெயிஷ்டா",
			"Ashadha":     "ஆஷாட",
			"Shravana":    "ஸ்ராவண",
			"Bhadrapada":  "பத்ரபத",
			"Ashwin":      "ஆஸ்வின",
			"Kartik":      "கார்த்திக",
			"Margashirsha": "மார்கழி",
			"Pausha":      "பௌஷ",
			"Magha":       "மாக",
			"Phalgun":     "பால்குன",
		},
		api.RegionBengal: {
			"Chaitra":     "চৈত্র",
			"Vaisakha":    "বৈশাখ",
			"Jyeshtha":    "জ্যৈষ্ঠ",
			"Ashadha":     "আষাঢ়",
			"Shravana":    "শ্রাবণ",
			"Bhadrapada":  "ভাদ্র",
			"Ashwin":      "আশ্বিন",
			"Kartik":      "কার্তিক",
			"Margashirsha": "অগ্রহায়ণ",
			"Pausha":      "পৌষ",
			"Magha":       "মাঘ",
			"Phalgun":     "ফাল্গুন",
		},
		// Add more regional mappings as needed
	}

	if regionalMap, exists := monthMappings[region]; exists {
		if localName, exists := regionalMap[sanskritName]; exists {
			return localName
		}
	}

	return sanskritName // Fallback to Sanskrit name
}

// Helper methods

func (c *CalendarSystemPlugin) isNorthIndianRegion(region api.Region) bool {
	northIndianRegions := map[api.Region]bool{
		api.RegionNorthIndia: true,
		api.RegionGujarat:    true,
		api.RegionMaha:       true,
	}
	return northIndianRegions[region]
}

func (c *CalendarSystemPlugin) adjustMonthData(ctx context.Context, data *api.PanchangamData, monthInfo *MonthInfo) error {
	// Add month-specific adjustments to the Panchangam data
	// Note: Since PanchangamData doesn't have a Metadata field, we'll add this information
	// to the Events or Muhurtas as appropriate, or extend the structure if needed
	
	// For now, we'll create a special event to carry month information
	monthEvent := api.Event{
		Name:         "Lunar Month",
		NameLocal:    monthInfo.NameLocal,
		Type:         api.EventTypeLunar,
		StartTime:    monthInfo.StartDate,
		EndTime:      monthInfo.EndDate,
		Significance: fmt.Sprintf("%s month in %s calendar system", monthInfo.Name, monthInfo.CalendarSystem),
		Region:       data.Region,
		Metadata: map[string]interface{}{
			"type":            "month_info",
			"name":            monthInfo.Name,
			"name_local":      monthInfo.NameLocal,
			"number":          monthInfo.Number,
			"calendar_system": string(monthInfo.CalendarSystem),
			"is_adhika_masa":  monthInfo.IsAdhikaMasa,
		},
	}
	
	data.Events = append(data.Events, monthEvent)
	return nil
}

func (c *CalendarSystemPlugin) addCalendarSystemMetadata(data *api.PanchangamData, monthInfo *MonthInfo) {
	// Add calendar system information as a special event since PanchangamData doesn't have Metadata field
	calendarSystemEvent := api.Event{
		Name:         "Calendar System Info",
		NameLocal:    "",
		Type:         api.EventTypeLunar,
		StartTime:    data.Date,
		EndTime:      data.Date.Add(24 * time.Hour),
		Significance: c.getCalendarSystemDescription(data.CalendarSystem),
		Region:       data.Region,
		Metadata: map[string]interface{}{
			"type":                  "calendar_system_info",
			"system":                string(data.CalendarSystem),
			"description":           c.getCalendarSystemDescription(data.CalendarSystem),
			"month_boundary_rule":   c.getMonthBoundaryRule(data.CalendarSystem),
			"regional_preference":   c.getRegionalPreference(data.Region),
			"calculation_precision": "astronomical",
		},
	}
	
	data.Events = append(data.Events, calendarSystemEvent)
}

func (c *CalendarSystemPlugin) getCalendarSystemDescription(system api.CalendarSystem) string {
	descriptions := map[api.CalendarSystem]string{
		api.CalendarAmanta:     "South Indian system where lunar months end on Amavasya (New Moon)",
		api.CalendarPurnimanta: "North Indian system where lunar months end on Purnima (Full Moon)",
		api.CalendarLunar:      "Pure lunar calendar following regional preferences",
		api.CalendarSolar:      "Solar calendar based on sun's movement through zodiac signs",
	}
	return descriptions[system]
}

func (c *CalendarSystemPlugin) getMonthBoundaryRule(system api.CalendarSystem) string {
	rules := map[api.CalendarSystem]string{
		api.CalendarAmanta:     "Month begins day after Amavasya, ends on next Amavasya",
		api.CalendarPurnimanta: "Month begins day after Purnima, ends on next Purnima",
		api.CalendarLunar:      "Follows regional calendar system preference",
		api.CalendarSolar:      "Month begins when sun enters new zodiac sign",
	}
	return rules[system]
}

func (c *CalendarSystemPlugin) getRegionalPreference(region api.Region) string {
	preferences := map[api.Region]string{
		api.RegionNorthIndia: "Purnimanta",
		api.RegionSouthIndia: "Amanta",
		api.RegionTamilNadu:  "Amanta",
		api.RegionKerala:     "Amanta",
		api.RegionBengal:     "Amanta",
		api.RegionGujarat:    "Purnimanta",
		api.RegionMaha:       "Purnimanta",
		api.RegionGlobal:     "Mixed (region-dependent)",
	}
	return preferences[region]
}

// isAdhikaMasa determines if a lunar month is an intercalary month (Adhika Masa)
// A lunar month is considered Adhika Masa if it doesn't contain a solar month transition (Sankranti)
func (c *CalendarSystemPlugin) isAdhikaMasa(ctx context.Context, monthStart, monthEnd time.Time) bool {
	// Get sun's longitude at month start and end
	startJD := ephemeris.TimeToJulianDay(monthStart)
	endJD := ephemeris.TimeToJulianDay(monthEnd)
	
	startPositions, err := c.ephemerisManager.GetPlanetaryPositions(ctx, startJD)
	if err != nil {
		// If we can't get planetary positions, assume it's not Adhika Masa
		return false
	}
	
	endPositions, err := c.ephemerisManager.GetPlanetaryPositions(ctx, endJD)
	if err != nil {
		return false
	}
	
	// Get sun's longitude (in degrees) at start and end of month
	startSunLong := startPositions.Sun.Longitude
	endSunLong := endPositions.Sun.Longitude
	
	// Handle longitude wraparound (360 degrees)
	if endSunLong < startSunLong {
		endSunLong += 360
	}
	
	// Calculate how many zodiac signs (30-degree segments) the sun has traversed
	startSign := int(startSunLong / 30)
	endSign := int(endSunLong / 30)
	
	// If no solar month transition occurred (no Sankranti), it's an Adhika Masa
	// This means the sun remained in the same zodiac sign throughout the lunar month
	return startSign == endSign
}