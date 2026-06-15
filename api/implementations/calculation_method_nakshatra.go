package implementations

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// CalculateNakshatra calculates nakshatra using the specified method
func (c *CalculationMethodPlugin) CalculateNakshatra(ctx context.Context, date time.Time, location api.Location, method api.CalculationMethod) (*api.Nakshatra, error) {
	if !c.enabled {
		return nil, fmt.Errorf("calculation method plugin is not enabled")
	}

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
		EpochJD:             588465.5,  // Kaliyuga start
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
