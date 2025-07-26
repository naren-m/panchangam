package astronomy

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TithiType represents the categorization of Tithi
type TithiType string

const (
	TithiTypeNanda  TithiType = "Nanda"  // 1, 6, 11 (Joyful)
	TithiTypeBhadra TithiType = "Bhadra" // 2, 7, 12 (Auspicious)
	TithiTypeJaya   TithiType = "Jaya"   // 3, 8, 13 (Victorious)
	TithiTypeRikta  TithiType = "Rikta"  // 4, 9, 14 (Empty)
	TithiTypePurna  TithiType = "Purna"  // 5, 10, 15 (Full/Complete)
)

// TithiInfo represents a Tithi with its properties
type TithiInfo struct {
	Number      int       `json:"number"`        // 1-15 for Shukla Paksha, 16-30 for Krishna Paksha
	Name        string    `json:"name"`          // Sanskrit name of the Tithi
	Type        TithiType `json:"type"`          // Category (Nanda, Bhadra, Jaya, Rikta, Purna)
	StartTime   time.Time `json:"start_time"`    // When this Tithi begins
	EndTime     time.Time `json:"end_time"`      // When this Tithi ends
	Duration    float64   `json:"duration"`      // Duration in hours
	IsShukla    bool      `json:"is_shukla"`     // true for Shukla Paksha, false for Krishna Paksha
	MoonSunDiff float64   `json:"moon_sun_diff"` // Moon longitude - Sun longitude in degrees
}

// TithiCalculator handles Tithi calculations
type TithiCalculator struct {
	ephemerisManager *ephemeris.Manager
	observer         observability.ObserverInterface
}

// NewTithiCalculator creates a new TithiCalculator
func NewTithiCalculator(ephemerisManager *ephemeris.Manager) *TithiCalculator {
	return &TithiCalculator{
		ephemerisManager: ephemerisManager,
		observer:         observability.Observer(),
	}
}

// TithiNames maps Tithi numbers to their Sanskrit names
var TithiNames = map[int]string{
	1: "Pratipada", 2: "Dwitiya", 3: "Tritiya", 4: "Chaturthi", 5: "Panchami",
	6: "Shashthi", 7: "Saptami", 8: "Ashtami", 9: "Navami", 10: "Dashami",
	11: "Ekadashi", 12: "Dwadashi", 13: "Trayodashi", 14: "Chaturdashi", 15: "Purnima",
	16: "Pratipada", 17: "Dwitiya", 18: "Tritiya", 19: "Chaturthi", 20: "Panchami",
	21: "Shashthi", 22: "Saptami", 23: "Ashtami", 24: "Navami", 25: "Dashami",
	26: "Ekadashi", 27: "Dwadashi", 28: "Trayodashi", 29: "Chaturdashi", 30: "Amavasya",
}

// GetTithiForDate calculates the Tithi for a given date
func (tc *TithiCalculator) GetTithiForDate(ctx context.Context, date time.Time) (*TithiInfo, error) {
	ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.GetTithiForDate")
	defer span.End()

	span.SetAttributes(
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("timezone", date.Location().String()),
	)

	// Convert to Julian day (using noon for calculation)
	noonDate := time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location())
	jd := ephemeris.TimeToJulianDay(noonDate)

	span.SetAttributes(attribute.Float64("julian_day", float64(jd)))

	// Get planetary positions
	ctx, posSpan := tc.observer.CreateSpan(ctx, "getTithiPositions")
	positions, err := tc.ephemerisManager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		posSpan.RecordError(err)
		posSpan.End()
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get planetary positions: %w", err)
	}

	sunLong := positions.Sun.Longitude
	moonLong := positions.Moon.Longitude

	posSpan.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
	)
	posSpan.End()

	// Calculate Tithi
	tithi, err := tc.calculateTithiFromLongitudes(ctx, sunLong, moonLong, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.String("tithi_type", string(tithi.Type)),
		attribute.Bool("is_shukla", tithi.IsShukla),
		attribute.Float64("moon_sun_diff", tithi.MoonSunDiff),
	)

	span.AddEvent("Tithi calculated", trace.WithAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.String("tithi_type", string(tithi.Type)),
	))

	return tithi, nil
}

// calculateTithiFromLongitudes calculates Tithi from Sun and Moon longitudes
func (tc *TithiCalculator) calculateTithiFromLongitudes(ctx context.Context, sunLong, moonLong float64, referenceDate time.Time) (*TithiInfo, error) {
	ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.calculateTithiFromLongitudes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
	)

	// Calculate the difference (Moon longitude - Sun longitude)
	moonSunDiff := moonLong - sunLong

	// Normalize to 0-360 degrees
	if moonSunDiff < 0 {
		moonSunDiff += 360
	}
	if moonSunDiff >= 360 {
		moonSunDiff -= 360
	}

	span.SetAttributes(attribute.Float64("normalized_moon_sun_diff", moonSunDiff))

	// Calculate Tithi number (each Tithi is 12 degrees)
	tithiFloat := moonSunDiff / 12.0
	tithiNumber := int(tithiFloat) + 1

	// Ensure Tithi number is in valid range (1-30)
	if tithiNumber > 30 {
		tithiNumber = 30
	}
	if tithiNumber < 1 {
		tithiNumber = 1
	}

	span.SetAttributes(
		attribute.Float64("tithi_float", tithiFloat),
		attribute.Int("tithi_number", tithiNumber),
	)

	// Determine if it's Shukla Paksha (waxing) or Krishna Paksha (waning)
	isShukla := tithiNumber <= 15

	// Get Tithi name
	tithiName := TithiNames[tithiNumber]

	// Determine Tithi type
	tithiType := getTithiType(tithiNumber)

	// Calculate approximate start and end times
	// Each Tithi spans 12 degrees, so we can estimate timing
	startTime, endTime := tc.calculateTithiTimes(ctx, tithiFloat, referenceDate)

	span.SetAttributes(
		attribute.String("tithi_name", tithiName),
		attribute.String("tithi_type", string(tithiType)),
		attribute.Bool("is_shukla", isShukla),
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
	)

	duration := endTime.Sub(startTime).Hours()

	tithi := &TithiInfo{
		Number:      tithiNumber,
		Name:        tithiName,
		Type:        tithiType,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    duration,
		IsShukla:    isShukla,
		MoonSunDiff: moonSunDiff,
	}

	span.AddEvent("Tithi calculation completed", trace.WithAttributes(
		attribute.Int("tithi_number", tithiNumber),
		attribute.String("tithi_name", tithiName),
		attribute.Float64("duration_hours", duration),
	))

	return tithi, nil
}

// calculateTithiTimes estimates the start and end times of a Tithi
func (tc *TithiCalculator) calculateTithiTimes(ctx context.Context, tithiFloat float64, referenceDate time.Time) (startTime, endTime time.Time) {
	_, span := tc.observer.CreateSpan(ctx, "TithiCalculator.calculateTithiTimes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("tithi_float", tithiFloat),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	// Calculate how far into the current Tithi we are
	tithiProgress := tithiFloat - math.Floor(tithiFloat)

	// Average Tithi duration is approximately 24.79 hours (lunar month / 30)
	avgTithiDuration := time.Duration(24.79 * float64(time.Hour))

	// Estimate when this Tithi started and will end
	timeIntoTithi := time.Duration(tithiProgress * float64(avgTithiDuration))

	// Start time is reference time minus how far we are into the Tithi
	noonRef := time.Date(referenceDate.Year(), referenceDate.Month(), referenceDate.Day(), 12, 0, 0, 0, referenceDate.Location())
	startTime = noonRef.Add(-timeIntoTithi)
	endTime = startTime.Add(avgTithiDuration)

	span.SetAttributes(
		attribute.Float64("tithi_progress", tithiProgress),
		attribute.Float64("avg_tithi_duration_hours", avgTithiDuration.Hours()),
		attribute.Float64("time_into_tithi_hours", timeIntoTithi.Hours()),
		attribute.String("calculated_start_time", startTime.Format("2006-01-02 15:04:05")),
		attribute.String("calculated_end_time", endTime.Format("2006-01-02 15:04:05")),
	)

	span.AddEvent("Tithi times calculated", trace.WithAttributes(
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
		attribute.Float64("duration_hours", endTime.Sub(startTime).Hours()),
	))

	return startTime, endTime
}

// getTithiType returns the type/category of a Tithi
func getTithiType(tithiNumber int) TithiType {
	// Normalize to 1-15 range for type calculation
	normalizedTithi := tithiNumber
	if normalizedTithi > 15 {
		normalizedTithi = normalizedTithi - 15
	}

	switch normalizedTithi {
	case 1, 6, 11:
		return TithiTypeNanda // Joyful
	case 2, 7, 12:
		return TithiTypeBhadra // Auspicious
	case 3, 8, 13:
		return TithiTypeJaya // Victorious
	case 4, 9, 14:
		return TithiTypeRikta // Empty
	case 5, 10, 15:
		return TithiTypePurna // Full/Complete
	default:
		return TithiTypeNanda // Default fallback
	}
}

// GetTithiFromLongitudes is a convenience function for direct longitude input
func (tc *TithiCalculator) GetTithiFromLongitudes(ctx context.Context, sunLong, moonLong float64, date time.Time) (*TithiInfo, error) {
	ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.GetTithiFromLongitudes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("date", date.Format("2006-01-02")),
	)

	return tc.calculateTithiFromLongitudes(ctx, sunLong, moonLong, date)
}

// GetTithiTypeDescription returns a description of the Tithi type
func GetTithiTypeDescription(tithiType TithiType) string {
	switch tithiType {
	case TithiTypeNanda:
		return "Joyful, good for celebrations and new beginnings"
	case TithiTypeBhadra:
		return "Auspicious, good for all activities"
	case TithiTypeJaya:
		return "Victorious, good for achieving success"
	case TithiTypeRikta:
		return "Empty, avoid starting new ventures"
	case TithiTypePurna:
		return "Complete, excellent for completion of tasks"
	default:
		return "Unknown Tithi type"
	}
}

// ValidateTithiCalculation validates a Tithi calculation result
func ValidateTithiCalculation(tithi *TithiInfo) error {
	if tithi == nil {
		return fmt.Errorf("tithi cannot be nil")
	}

	if tithi.Number < 1 || tithi.Number > 30 {
		return fmt.Errorf("invalid tithi number: %d, must be between 1 and 30", tithi.Number)
	}

	if tithi.MoonSunDiff < 0 || tithi.MoonSunDiff >= 360 {
		return fmt.Errorf("invalid moon-sun difference: %f, must be between 0 and 360 degrees", tithi.MoonSunDiff)
	}

	if tithi.Duration <= 0 || tithi.Duration > 48 {
		return fmt.Errorf("invalid tithi duration: %f hours, must be positive and reasonable", tithi.Duration)
	}

	if tithi.EndTime.Before(tithi.StartTime) {
		return fmt.Errorf("tithi end time cannot be before start time")
	}

	return nil
}
