package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

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
		Name:            monthName,
		NameLocal:       localName,
		Number:          monthNumber,
		StartDate:       monthStart,
		EndDate:         monthEnd,
		CalendarSystem:  api.CalendarAmanta,
		Region:          region,
		Year:            date.Year(),
		IsAdhikaMasa:    c.isAdhikaMasa(ctx, monthStart, monthEnd),
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
		Name:            monthName,
		NameLocal:       localName,
		Number:          monthNumber,
		StartDate:       monthStart,
		EndDate:         monthEnd,
		CalendarSystem:  api.CalendarPurnimanta,
		Region:          region,
		Year:            date.Year(),
		IsAdhikaMasa:    c.isAdhikaMasa(ctx, monthStart, monthEnd),
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
