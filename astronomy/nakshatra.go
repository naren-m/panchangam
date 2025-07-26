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

// NakshatraInfo represents a Nakshatra with its properties
type NakshatraInfo struct {
	Number       int       `json:"number"`         // 1-27
	Name         string    `json:"name"`           // Sanskrit name
	Deity        string    `json:"deity"`          // Ruling deity
	PlanetaryLord string   `json:"planetary_lord"` // Ruling planet
	Symbol       string    `json:"symbol"`         // Traditional symbol
	Pada         int       `json:"pada"`           // Current pada (1-4)
	StartTime    time.Time `json:"start_time"`     // When this Nakshatra begins
	EndTime      time.Time `json:"end_time"`       // When this Nakshatra ends
	Duration     float64   `json:"duration"`       // Duration in hours
	MoonLongitude float64  `json:"moon_longitude"` // Moon's longitude in degrees
}

// NakshatraCalculator handles Nakshatra calculations
type NakshatraCalculator struct {
	ephemerisManager *ephemeris.Manager
	observer         observability.ObserverInterface
}

// NewNakshatraCalculator creates a new NakshatraCalculator
func NewNakshatraCalculator(ephemerisManager *ephemeris.Manager) *NakshatraCalculator {
	return &NakshatraCalculator{
		ephemerisManager: ephemerisManager,
		observer:         observability.Observer(),
	}
}

// NakshatraData contains detailed information about each Nakshatra
// Sources: 
// - "Hindu Astronomy" by W.E. van Wijk (1930)
// - "Surya Siddhanta" - Ancient Sanskrit astronomical text
// - "Brihat Parashara Hora Shastra" by Sage Parashara
// - "Muhurta Chintamani" by Daivagya Ramachandra
var NakshatraData = map[int]struct {
	Name          string
	Deity         string
	PlanetaryLord string
	Symbol        string
}{
	1:  {"Ashwini", "Ashwini Kumaras", "Ketu", "Horse's Head"},
	2:  {"Bharani", "Yama", "Venus", "Yoni (Vagina)"},
	3:  {"Krittika", "Agni", "Sun", "Razor/Knife"},
	4:  {"Rohini", "Brahma", "Moon", "Cart/Chariot"},
	5:  {"Mrigashira", "Soma", "Mars", "Deer's Head"},
	6:  {"Ardra", "Rudra", "Rahu", "Teardrop/Diamond"},
	7:  {"Punarvasu", "Aditi", "Jupiter", "Bow and Quiver"},
	8:  {"Pushya", "Brihaspati", "Saturn", "Cow's Udder"},
	9:  {"Ashlesha", "Nagas", "Mercury", "Serpent"},
	10: {"Magha", "Pitrs (Ancestors)", "Ketu", "Throne"},
	11: {"Purva Phalguni", "Bhaga", "Venus", "Front Legs of Bed"},
	12: {"Uttara Phalguni", "Aryaman", "Sun", "Back Legs of Bed"},
	13: {"Hasta", "Savitar", "Moon", "Hand"},
	14: {"Chitra", "Tvashtar", "Mars", "Bright Jewel"},
	15: {"Swati", "Vayu", "Rahu", "Young Shoot of Plant"},
	16: {"Vishakha", "Indra-Agni", "Jupiter", "Triumphal Arch"},
	17: {"Anuradha", "Mitra", "Saturn", "Lotus"},
	18: {"Jyeshtha", "Indra", "Mercury", "Circular Amulet"},
	19: {"Mula", "Nirriti", "Ketu", "Bunch of Roots"},
	20: {"Purva Ashadha", "Apas", "Venus", "Elephant Tusk"},
	21: {"Uttara Ashadha", "Vishve Devas", "Sun", "Elephant Tusk"},
	22: {"Shravana", "Vishnu", "Moon", "Ear/Three Footprints"},
	23: {"Dhanishta", "Vasus", "Mars", "Drum"},
	24: {"Shatabhisha", "Varuna", "Rahu", "Empty Circle"},
	25: {"Purva Bhadrapada", "Aja Ekapada", "Jupiter", "Front Legs of Funeral Cot"},
	26: {"Uttara Bhadrapada", "Ahir Budhnya", "Saturn", "Back Legs of Funeral Cot"},
	27: {"Revati", "Pushan", "Mercury", "Fish/Pair of Fish"},
}

// GetNakshatraForDate calculates the Nakshatra for a given date
func (nc *NakshatraCalculator) GetNakshatraForDate(ctx context.Context, date time.Time) (*NakshatraInfo, error) {
	ctx, span := nc.observer.CreateSpan(ctx, "NakshatraCalculator.GetNakshatraForDate")
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
	ctx, posSpan := nc.observer.CreateSpan(ctx, "getNakshatraPositions")
	positions, err := nc.ephemerisManager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		posSpan.RecordError(err)
		posSpan.End()
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get planetary positions: %w", err)
	}

	moonLong := positions.Moon.Longitude

	posSpan.SetAttributes(attribute.Float64("moon_longitude", moonLong))
	posSpan.End()

	// Calculate Nakshatra
	nakshatra, err := nc.calculateNakshatraFromLongitude(ctx, moonLong, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("nakshatra_number", nakshatra.Number),
		attribute.String("nakshatra_name", nakshatra.Name),
		attribute.String("deity", nakshatra.Deity),
		attribute.String("planetary_lord", nakshatra.PlanetaryLord),
		attribute.Int("pada", nakshatra.Pada),
		attribute.Float64("moon_longitude", nakshatra.MoonLongitude),
	)

	span.AddEvent("Nakshatra calculated", trace.WithAttributes(
		attribute.Int("nakshatra_number", nakshatra.Number),
		attribute.String("nakshatra_name", nakshatra.Name),
		attribute.Int("pada", nakshatra.Pada),
	))

	return nakshatra, nil
}

// calculateNakshatraFromLongitude calculates Nakshatra from Moon's longitude
func (nc *NakshatraCalculator) calculateNakshatraFromLongitude(ctx context.Context, moonLong float64, referenceDate time.Time) (*NakshatraInfo, error) {
	ctx, span := nc.observer.CreateSpan(ctx, "NakshatraCalculator.calculateNakshatraFromLongitude")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	// Normalize longitude to 0-360 degrees
	normalizedLong := moonLong
	if normalizedLong < 0 {
		normalizedLong += 360
	}
	if normalizedLong >= 360 {
		normalizedLong -= 360
	}

	span.SetAttributes(attribute.Float64("normalized_moon_longitude", normalizedLong))

	// Each Nakshatra spans 13°20' (13.333... degrees)
	// There are 27 Nakshatras covering the full 360° zodiac
	nakshatraSpan := 360.0 / 27.0 // 13.333... degrees
	
	// Calculate Nakshatra number (1-27)
	nakshatraFloat := normalizedLong / nakshatraSpan
	nakshatraNumber := int(nakshatraFloat) + 1

	// Ensure Nakshatra number is in valid range (1-27)
	if nakshatraNumber > 27 {
		nakshatraNumber = 27
	}
	if nakshatraNumber < 1 {
		nakshatraNumber = 1
	}

	// Calculate Pada (1-4) - each Nakshatra is divided into 4 equal parts
	// Each Pada spans 3°20' (3.333... degrees)
	padaSpan := nakshatraSpan / 4.0 // 3.333... degrees
	positionInNakshatra := normalizedLong - (float64(nakshatraNumber-1) * nakshatraSpan)
	pada := int(positionInNakshatra/padaSpan) + 1

	// Ensure Pada is in valid range (1-4)
	if pada > 4 {
		pada = 4
	}
	if pada < 1 {
		pada = 1
	}

	span.SetAttributes(
		attribute.Float64("nakshatra_span", nakshatraSpan),
		attribute.Float64("nakshatra_float", nakshatraFloat),
		attribute.Int("nakshatra_number", nakshatraNumber),
		attribute.Float64("position_in_nakshatra", positionInNakshatra),
		attribute.Float64("pada_span", padaSpan),
		attribute.Int("pada", pada),
	)

	// Get Nakshatra details
	nakshatraDetails := NakshatraData[nakshatraNumber]

	// Calculate approximate start and end times
	// Average Nakshatra duration is approximately 24.79 hours (lunar month / 27)
	startTime, endTime := nc.calculateNakshatraTimes(ctx, nakshatraFloat, referenceDate)

	span.SetAttributes(
		attribute.String("nakshatra_name", nakshatraDetails.Name),
		attribute.String("deity", nakshatraDetails.Deity),
		attribute.String("planetary_lord", nakshatraDetails.PlanetaryLord),
		attribute.String("symbol", nakshatraDetails.Symbol),
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
	)

	duration := endTime.Sub(startTime).Hours()

	nakshatra := &NakshatraInfo{
		Number:        nakshatraNumber,
		Name:          nakshatraDetails.Name,
		Deity:         nakshatraDetails.Deity,
		PlanetaryLord: nakshatraDetails.PlanetaryLord,
		Symbol:        nakshatraDetails.Symbol,
		Pada:          pada,
		StartTime:     startTime,
		EndTime:       endTime,
		Duration:      duration,
		MoonLongitude: normalizedLong,
	}

	span.AddEvent("Nakshatra calculation completed", trace.WithAttributes(
		attribute.Int("nakshatra_number", nakshatraNumber),
		attribute.String("nakshatra_name", nakshatraDetails.Name),
		attribute.Int("pada", pada),
		attribute.Float64("duration_hours", duration),
	))

	return nakshatra, nil
}

// calculateNakshatraTimes estimates the start and end times of a Nakshatra
func (nc *NakshatraCalculator) calculateNakshatraTimes(ctx context.Context, nakshatraFloat float64, referenceDate time.Time) (startTime, endTime time.Time) {
	_, span := nc.observer.CreateSpan(ctx, "NakshatraCalculator.calculateNakshatraTimes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("nakshatra_float", nakshatraFloat),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	// Calculate how far into the current Nakshatra we are
	nakshatraProgress := nakshatraFloat - math.Floor(nakshatraFloat)

	// Average Nakshatra duration is approximately 27.32 hours (sidereal month / 27)
	avgNakshatraDuration := time.Duration(27.32 * float64(time.Hour))

	// Estimate when this Nakshatra started and will end
	timeIntoNakshatra := time.Duration(nakshatraProgress * float64(avgNakshatraDuration))

	// Start time is reference time minus how far we are into the Nakshatra
	noonRef := time.Date(referenceDate.Year(), referenceDate.Month(), referenceDate.Day(), 12, 0, 0, 0, referenceDate.Location())
	startTime = noonRef.Add(-timeIntoNakshatra)
	endTime = startTime.Add(avgNakshatraDuration)

	span.SetAttributes(
		attribute.Float64("nakshatra_progress", nakshatraProgress),
		attribute.Float64("avg_nakshatra_duration_hours", avgNakshatraDuration.Hours()),
		attribute.Float64("time_into_nakshatra_hours", timeIntoNakshatra.Hours()),
		attribute.String("calculated_start_time", startTime.Format("2006-01-02 15:04:05")),
		attribute.String("calculated_end_time", endTime.Format("2006-01-02 15:04:05")),
	)

	span.AddEvent("Nakshatra times calculated", trace.WithAttributes(
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
		attribute.Float64("duration_hours", endTime.Sub(startTime).Hours()),
	))

	return startTime, endTime
}

// GetNakshatraFromLongitude is a convenience function for direct longitude input
func (nc *NakshatraCalculator) GetNakshatraFromLongitude(ctx context.Context, moonLong float64, date time.Time) (*NakshatraInfo, error) {
	ctx, span := nc.observer.CreateSpan(ctx, "NakshatraCalculator.GetNakshatraFromLongitude")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("date", date.Format("2006-01-02")),
	)

	return nc.calculateNakshatraFromLongitude(ctx, moonLong, date)
}

// GetPadaDescription returns a description of the Pada
func GetPadaDescription(nakshatraNumber, pada int) string {
	// Each Nakshatra's 4 padas have specific meanings and associations
	// This is a simplified version - in practice, each Nakshatra has unique pada characteristics
	switch pada {
	case 1:
		return "First pada - represents new beginnings and initiation"
	case 2:
		return "Second pada - represents growth and development"
	case 3:
		return "Third pada - represents maturity and stability"
	case 4:
		return "Fourth pada - represents completion and transformation"
	default:
		return "Unknown pada"
	}
}

// ValidateNakshatraCalculation validates a Nakshatra calculation result
func ValidateNakshatraCalculation(nakshatra *NakshatraInfo) error {
	if nakshatra == nil {
		return fmt.Errorf("nakshatra cannot be nil")
	}

	if nakshatra.Number < 1 || nakshatra.Number > 27 {
		return fmt.Errorf("invalid nakshatra number: %d, must be between 1 and 27", nakshatra.Number)
	}

	if nakshatra.Pada < 1 || nakshatra.Pada > 4 {
		return fmt.Errorf("invalid pada: %d, must be between 1 and 4", nakshatra.Pada)
	}

	if nakshatra.MoonLongitude < 0 || nakshatra.MoonLongitude >= 360 {
		return fmt.Errorf("invalid moon longitude: %f, must be between 0 and 360 degrees", nakshatra.MoonLongitude)
	}

	if nakshatra.Duration <= 0 || nakshatra.Duration > 48 {
		return fmt.Errorf("invalid nakshatra duration: %f hours, must be positive and reasonable", nakshatra.Duration)
	}

	if nakshatra.EndTime.Before(nakshatra.StartTime) {
		return fmt.Errorf("nakshatra end time cannot be before start time")
	}

	if nakshatra.Name == "" {
		return fmt.Errorf("nakshatra name cannot be empty")
	}

	return nil
}