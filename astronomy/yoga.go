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

// YogaQuality represents the auspicious nature of a Yoga
type YogaQuality string

const (
	YogaQualityAuspicious    YogaQuality = "Auspicious"
	YogaQualityInauspicious  YogaQuality = "Inauspicious"
	YogaQualityMixed         YogaQuality = "Mixed"
	YogaQualityNeutral       YogaQuality = "Neutral"
)

// YogaInfo represents a Yoga with its properties
type YogaInfo struct {
	Number        int         `json:"number"`          // 1-27
	Name          string      `json:"name"`            // Sanskrit name
	Quality       YogaQuality `json:"quality"`         // Auspicious nature
	Description   string      `json:"description"`     // Meaning and effects
	StartTime     time.Time   `json:"start_time"`      // When this Yoga begins
	EndTime       time.Time   `json:"end_time"`        // When this Yoga ends
	Duration      float64     `json:"duration"`        // Duration in hours
	SunLongitude  float64     `json:"sun_longitude"`   // Sun's longitude in degrees
	MoonLongitude float64     `json:"moon_longitude"`  // Moon's longitude in degrees
	CombinedValue float64     `json:"combined_value"`  // Sum of Sun and Moon longitudes
}

// YogaCalculator handles Yoga calculations
type YogaCalculator struct {
	ephemerisManager *ephemeris.Manager
	observer         observability.ObserverInterface
}

// NewYogaCalculator creates a new YogaCalculator
func NewYogaCalculator(ephemerisManager *ephemeris.Manager) *YogaCalculator {
	return &YogaCalculator{
		ephemerisManager: ephemerisManager,
		observer:         observability.Observer(),
	}
}

// YogaData contains detailed information about each Yoga
// Sources:
// - "Brihat Parashara Hora Shastra" by Sage Parashara
// - "Muhurta Chintamani" by Daivagya Ramachandra
// - "Hindu Astronomy" by W.E. van Wijk (1930)
// - "Surya Siddhanta" - Ancient Sanskrit astronomical text
var YogaData = map[int]struct {
	Name        string
	Quality     YogaQuality
	Description string
}{
	1:  {"Vishkambha", YogaQualityInauspicious, "Obstructive, delays and obstacles"},
	2:  {"Priti", YogaQualityAuspicious, "Love and affection, good for relationships"},
	3:  {"Ayushman", YogaQualityAuspicious, "Longevity, health and vitality"},
	4:  {"Saubhagya", YogaQualityAuspicious, "Good fortune, prosperity and happiness"},
	5:  {"Shobhana", YogaQualityAuspicious, "Beauty, auspicious for ceremonies"},
	6:  {"Atiganda", YogaQualityInauspicious, "Great danger, avoid important work"},
	7:  {"Sukarma", YogaQualityAuspicious, "Good deeds, meritorious actions"},
	8:  {"Dhriti", YogaQualityAuspicious, "Determination, steadfastness"},
	9:  {"Shula", YogaQualityInauspicious, "Pain and suffering, inauspicious"},
	10: {"Ganda", YogaQualityInauspicious, "Danger, avoid travel and new ventures"},
	11: {"Vriddhi", YogaQualityAuspicious, "Growth and prosperity"},
	12: {"Dhruva", YogaQualityAuspicious, "Stability, permanent gains"},
	13: {"Vyaghata", YogaQualityInauspicious, "Destruction, avoid important work"},
	14: {"Harshana", YogaQualityAuspicious, "Joy and happiness"},
	15: {"Vajra", YogaQualityMixed, "Diamond-like strength, can be harsh"},
	16: {"Siddhi", YogaQualityAuspicious, "Success and achievement"},
	17: {"Vyatipata", YogaQualityInauspicious, "Great calamity, very inauspicious"},
	18: {"Variyana", YogaQualityMixed, "Choice and selection, mixed results"},
	19: {"Parigha", YogaQualityInauspicious, "Iron rod, obstacles and delays"},
	20: {"Shiva", YogaQualityAuspicious, "Auspicious, beneficial for all activities"},
	21: {"Siddha", YogaQualityAuspicious, "Accomplished, success assured"},
	22: {"Sadhya", YogaQualityAuspicious, "Achievable, goals can be accomplished"},
	23: {"Shubha", YogaQualityAuspicious, "Pure and auspicious"},
	24: {"Shukla", YogaQualityAuspicious, "Bright and pure"},
	25: {"Brahma", YogaQualityAuspicious, "Divine, highly auspicious"},
	26: {"Indra", YogaQualityAuspicious, "Royal, powerful and prosperous"},
	27: {"Vaidhriti", YogaQualityInauspicious, "Separation, avoid joint ventures"},
}

// GetYogaForDate calculates the Yoga for a given date
func (yc *YogaCalculator) GetYogaForDate(ctx context.Context, date time.Time) (*YogaInfo, error) {
	ctx, span := yc.observer.CreateSpan(ctx, "YogaCalculator.GetYogaForDate")
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
	ctx, posSpan := yc.observer.CreateSpan(ctx, "getYogaPositions")
	positions, err := yc.ephemerisManager.GetPlanetaryPositions(ctx, jd)
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

	// Calculate Yoga
	yoga, err := yc.calculateYogaFromLongitudes(ctx, sunLong, moonLong, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("yoga_number", yoga.Number),
		attribute.String("yoga_name", yoga.Name),
		attribute.String("yoga_quality", string(yoga.Quality)),
		attribute.Float64("combined_value", yoga.CombinedValue),
	)

	span.AddEvent("Yoga calculated", trace.WithAttributes(
		attribute.Int("yoga_number", yoga.Number),
		attribute.String("yoga_name", yoga.Name),
		attribute.String("yoga_quality", string(yoga.Quality)),
	))

	return yoga, nil
}

// calculateYogaFromLongitudes calculates Yoga from Sun and Moon longitudes
func (yc *YogaCalculator) calculateYogaFromLongitudes(ctx context.Context, sunLong, moonLong float64, referenceDate time.Time) (*YogaInfo, error) {
	ctx, span := yc.observer.CreateSpan(ctx, "YogaCalculator.calculateYogaFromLongitudes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	// Normalize longitudes to 0-360 degrees
	normalizedSunLong := normalizeLongitude(sunLong)
	normalizedMoonLong := normalizeLongitude(moonLong)

	// Calculate the sum of Sun and Moon longitudes
	combinedValue := normalizedSunLong + normalizedMoonLong

	// Normalize to 0-360 degrees
	if combinedValue >= 360 {
		combinedValue -= 360
	}

	span.SetAttributes(
		attribute.Float64("normalized_sun_longitude", normalizedSunLong),
		attribute.Float64("normalized_moon_longitude", normalizedMoonLong),
		attribute.Float64("combined_value", combinedValue),
	)

	// Each Yoga spans 13°20' (13.333... degrees), same as Nakshatra
	// There are 27 Yogas covering the full 360° zodiac
	yogaSpan := 360.0 / 27.0 // 13.333... degrees
	
	// Calculate Yoga number (1-27)
	yogaFloat := combinedValue / yogaSpan
	yogaNumber := int(yogaFloat) + 1

	// Ensure Yoga number is in valid range (1-27)
	if yogaNumber > 27 {
		yogaNumber = 27
	}
	if yogaNumber < 1 {
		yogaNumber = 1
	}

	span.SetAttributes(
		attribute.Float64("yoga_span", yogaSpan),
		attribute.Float64("yoga_float", yogaFloat),
		attribute.Int("yoga_number", yogaNumber),
	)

	// Get Yoga details
	yogaDetails := YogaData[yogaNumber]

	// Calculate approximate start and end times
	// Average Yoga duration is approximately 24.79 hours (similar to Tithi)
	startTime, endTime := yc.calculateYogaTimes(ctx, yogaFloat, referenceDate)

	span.SetAttributes(
		attribute.String("yoga_name", yogaDetails.Name),
		attribute.String("yoga_quality", string(yogaDetails.Quality)),
		attribute.String("yoga_description", yogaDetails.Description),
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
	)

	duration := endTime.Sub(startTime).Hours()

	yoga := &YogaInfo{
		Number:        yogaNumber,
		Name:          yogaDetails.Name,
		Quality:       yogaDetails.Quality,
		Description:   yogaDetails.Description,
		StartTime:     startTime,
		EndTime:       endTime,
		Duration:      duration,
		SunLongitude:  normalizedSunLong,
		MoonLongitude: normalizedMoonLong,
		CombinedValue: combinedValue,
	}

	span.AddEvent("Yoga calculation completed", trace.WithAttributes(
		attribute.Int("yoga_number", yogaNumber),
		attribute.String("yoga_name", yogaDetails.Name),
		attribute.String("yoga_quality", string(yogaDetails.Quality)),
		attribute.Float64("duration_hours", duration),
	))

	return yoga, nil
}

// calculateYogaTimes estimates the start and end times of a Yoga
func (yc *YogaCalculator) calculateYogaTimes(ctx context.Context, yogaFloat float64, referenceDate time.Time) (startTime, endTime time.Time) {
	_, span := yc.observer.CreateSpan(ctx, "YogaCalculator.calculateYogaTimes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("yoga_float", yogaFloat),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	// Calculate how far into the current Yoga we are
	yogaProgress := yogaFloat - math.Floor(yogaFloat)

	// Average Yoga duration varies based on the combined motion of Sun and Moon
	// Approximately 24.79 hours (similar to average Tithi duration)
	avgYogaDuration := time.Duration(24.79 * float64(time.Hour))

	// Estimate when this Yoga started and will end
	timeIntoYoga := time.Duration(yogaProgress * float64(avgYogaDuration))

	// Start time is reference time minus how far we are into the Yoga
	noonRef := time.Date(referenceDate.Year(), referenceDate.Month(), referenceDate.Day(), 12, 0, 0, 0, referenceDate.Location())
	startTime = noonRef.Add(-timeIntoYoga)
	endTime = startTime.Add(avgYogaDuration)

	span.SetAttributes(
		attribute.Float64("yoga_progress", yogaProgress),
		attribute.Float64("avg_yoga_duration_hours", avgYogaDuration.Hours()),
		attribute.Float64("time_into_yoga_hours", timeIntoYoga.Hours()),
		attribute.String("calculated_start_time", startTime.Format("2006-01-02 15:04:05")),
		attribute.String("calculated_end_time", endTime.Format("2006-01-02 15:04:05")),
	)

	span.AddEvent("Yoga times calculated", trace.WithAttributes(
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
		attribute.Float64("duration_hours", endTime.Sub(startTime).Hours()),
	))

	return startTime, endTime
}

// GetYogaFromLongitudes is a convenience function for direct longitude input
func (yc *YogaCalculator) GetYogaFromLongitudes(ctx context.Context, sunLong, moonLong float64, date time.Time) (*YogaInfo, error) {
	ctx, span := yc.observer.CreateSpan(ctx, "YogaCalculator.GetYogaFromLongitudes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("date", date.Format("2006-01-02")),
	)

	return yc.calculateYogaFromLongitudes(ctx, sunLong, moonLong, date)
}

// IsAuspiciousYoga returns true if the Yoga is considered auspicious
func IsAuspiciousYoga(yoga *YogaInfo) bool {
	return yoga.Quality == YogaQualityAuspicious
}

// IsInauspiciousYoga returns true if the Yoga is considered inauspicious
func IsInauspiciousYoga(yoga *YogaInfo) bool {
	return yoga.Quality == YogaQualityInauspicious
}

// GetYogaQualityDescription returns a detailed description of the Yoga quality
func GetYogaQualityDescription(quality YogaQuality) string {
	switch quality {
	case YogaQualityAuspicious:
		return "Favorable for all activities, brings good fortune and success"
	case YogaQualityInauspicious:
		return "Unfavorable, avoid important activities and new ventures"
	case YogaQualityMixed:
		return "Mixed results, proceed with caution and careful planning"
	case YogaQualityNeutral:
		return "Neutral influence, neither particularly favorable nor unfavorable"
	default:
		return "Unknown yoga quality"
	}
}

// normalizeLongitude normalizes a longitude value to 0-360 degrees
func normalizeLongitude(longitude float64) float64 {
	normalized := longitude
	for normalized < 0 {
		normalized += 360
	}
	for normalized >= 360 {
		normalized -= 360
	}
	return normalized
}

// ValidateYogaCalculation validates a Yoga calculation result
func ValidateYogaCalculation(yoga *YogaInfo) error {
	if yoga == nil {
		return fmt.Errorf("yoga cannot be nil")
	}

	if yoga.Number < 1 || yoga.Number > 27 {
		return fmt.Errorf("invalid yoga number: %d, must be between 1 and 27", yoga.Number)
	}

	if yoga.SunLongitude < 0 || yoga.SunLongitude >= 360 {
		return fmt.Errorf("invalid sun longitude: %f, must be between 0 and 360 degrees", yoga.SunLongitude)
	}

	if yoga.MoonLongitude < 0 || yoga.MoonLongitude >= 360 {
		return fmt.Errorf("invalid moon longitude: %f, must be between 0 and 360 degrees", yoga.MoonLongitude)
	}

	if yoga.CombinedValue < 0 || yoga.CombinedValue >= 360 {
		return fmt.Errorf("invalid combined value: %f, must be between 0 and 360 degrees", yoga.CombinedValue)
	}

	if yoga.Duration <= 0 || yoga.Duration > 48 {
		return fmt.Errorf("invalid yoga duration: %f hours, must be positive and reasonable", yoga.Duration)
	}

	if yoga.EndTime.Before(yoga.StartTime) {
		return fmt.Errorf("yoga end time cannot be before start time")
	}

	if yoga.Name == "" {
		return fmt.Errorf("yoga name cannot be empty")
	}

	// Validate quality is one of the defined values
	switch yoga.Quality {
	case YogaQualityAuspicious, YogaQualityInauspicious, YogaQualityMixed, YogaQualityNeutral:
		// Valid quality
	default:
		return fmt.Errorf("invalid yoga quality: %s", yoga.Quality)
	}

	return nil
}