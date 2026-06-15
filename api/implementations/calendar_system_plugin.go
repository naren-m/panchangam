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
	enabled          bool
	config           map[string]interface{}
	tithiCalculator  *astronomy.TithiCalculator
	ephemerisManager *ephemeris.Manager
}

// MonthInfo represents lunar month information
type MonthInfo struct {
	Name            string               `json:"name"`
	NameLocal       string               `json:"name_local"`
	Number          int                  `json:"number"` // 1-12
	StartDate       time.Time            `json:"start_date"`
	EndDate         time.Time            `json:"end_date"`
	CalendarSystem  api.CalendarSystem   `json:"calendar_system"`
	Region          api.Region           `json:"region"`
	Year            int                  `json:"year"`
	IsAdhikaMasa    bool                 `json:"is_adhika_masa"` // Intercalary month
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
		Name:         "calendar_system_plugin",
		Version:      api.Version{Major: 1, Minor: 0, Patch: 0},
		Description:  "Handles different Hindu calendar systems: Amanta (South Indian) vs Purnimanta (North Indian)",
		Author:       "Panchangam Team",
		Capabilities: []string{},
		Dependencies: []string{"astronomy", "ephemeris"},
		Metadata: map[string]interface{}{
			"calendar_systems":    []string{"amanta", "purnimanta", "lunar", "solar"},
			"regional_support":    true,
			"month_calculation":   "astronomical",
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
	if !c.enabled {
		return nil, fmt.Errorf("calendar system plugin is not enabled")
	}

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
