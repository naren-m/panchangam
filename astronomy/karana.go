package astronomy

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// KaranaType represents the category of Karana
type KaranaType string

const (
	KaranaTypeMovable KaranaType = "Movable"  // Chara Karanas (7 types)
	KaranaTypeFixed   KaranaType = "Fixed"    // Sthira Karanas (4 types)
)

// KaranaInfo represents a Karana with its properties
type KaranaInfo struct {
	Number      int        `json:"number"`        // 1-11 (in the cycle)
	Name        string     `json:"name"`          // Sanskrit name
	Type        KaranaType `json:"type"`          // Movable or Fixed
	Description string     `json:"description"`   // Meaning and effects
	IsVishti    bool       `json:"is_vishti"`     // Special flag for Vishti (Bhadra) karana
	StartTime   time.Time  `json:"start_time"`    // When this Karana begins
	EndTime     time.Time  `json:"end_time"`      // When this Karana ends
	Duration    float64    `json:"duration"`      // Duration in hours
	MoonSunDiff float64    `json:"moon_sun_diff"` // Moon longitude - Sun longitude in degrees
	TithiNumber int        `json:"tithi_number"`  // Associated Tithi number
	HalfTithi   int        `json:"half_tithi"`    // Which half of the Tithi (1 or 2)
}

// KaranaCalculator handles Karana calculations
type KaranaCalculator struct {
	tithiCalculator *TithiCalculator
	observer        observability.ObserverInterface
}

// NewKaranaCalculator creates a new KaranaCalculator
func NewKaranaCalculator(ephemerisManager *ephemeris.Manager) *KaranaCalculator {
	return &KaranaCalculator{
		tithiCalculator: NewTithiCalculator(ephemerisManager),
		observer:        observability.Observer(),
	}
}

// KaranaData contains detailed information about each Karana
// The 11 Karanas cycle through the lunar month in a specific pattern
var KaranaData = map[int]struct {
	Name        string
	Type        KaranaType
	Description string
	IsVishti    bool
}{
	1:  {"Kintughna", KaranaTypeMovable, "Destroyer of insects, good for destroying enemies", false},
	2:  {"Bava", KaranaTypeMovable, "Child-like, good for creative and joyful activities", false},
	3:  {"Balava", KaranaTypeMovable, "Strong and powerful, good for strength-based activities", false},
	4:  {"Kaulava", KaranaTypeMovable, "Of the family, good for family-related activities", false},
	5:  {"Taitila", KaranaTypeMovable, "Sesame seed, good for detailed work", false},
	6:  {"Gara", KaranaTypeMovable, "Poison, avoid important activities", false},
	7:  {"Vanija", KaranaTypeMovable, "Merchant, good for business and trade", false},
	8:  {"Vishti", KaranaTypeMovable, "Obstruction, very inauspicious - avoid all important work", true}, // Special Karana
	9:  {"Shakuni", KaranaTypeFixed, "Bird of ill omen, inauspicious", false},
	10: {"Chatushpada", KaranaTypeFixed, "Four-footed, stable and grounding", false},
	11: {"Naga", KaranaTypeFixed, "Serpent, mysterious and transformative", false},
	// Note: Kimstughna (12th) is the same as Kintughna (1st) - it's a special case for the last half of Amavasya
}

// GetKaranaForDate calculates the Karana for a given date
func (kc *KaranaCalculator) GetKaranaForDate(ctx context.Context, date time.Time) (*KaranaInfo, error) {
	ctx, span := kc.observer.CreateSpan(ctx, "KaranaCalculator.GetKaranaForDate")
	defer span.End()

	span.SetAttributes(
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("timezone", date.Location().String()),
	)

	// Get Tithi first using the existing calculator
	tithi, err := kc.tithiCalculator.GetTithiForDate(ctx, date)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get tithi for karana calculation: %w", err)
	}

	span.SetAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.Float64("moon_sun_diff", tithi.MoonSunDiff),
	)

	// Calculate Karana from Tithi information
	karana, err := kc.calculateKaranaFromTithi(ctx, tithi, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("karana_number", karana.Number),
		attribute.String("karana_name", karana.Name),
		attribute.String("karana_type", string(karana.Type)),
		attribute.Bool("is_vishti", karana.IsVishti),
		attribute.Int("tithi_number", karana.TithiNumber),
		attribute.Int("half_tithi", karana.HalfTithi),
		attribute.Float64("moon_sun_diff", karana.MoonSunDiff),
	)

	span.AddEvent("Karana calculated", trace.WithAttributes(
		attribute.Int("karana_number", karana.Number),
		attribute.String("karana_name", karana.Name),
		attribute.Bool("is_vishti", karana.IsVishti),
	))

	return karana, nil
}

// calculateKaranaFromTithi calculates Karana from Tithi information
func (kc *KaranaCalculator) calculateKaranaFromTithi(ctx context.Context, tithi *TithiInfo, referenceDate time.Time) (*KaranaInfo, error) {
	ctx, span := kc.observer.CreateSpan(ctx, "KaranaCalculator.calculateKaranaFromTithi")
	defer span.End()

	span.SetAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.Float64("moon_sun_diff", tithi.MoonSunDiff),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	// Each Tithi is divided into 2 Karanas (6 degrees each)
	// Determine which half of the Tithi we're in
	positionInTithi := tithi.MoonSunDiff - (float64(tithi.Number-1) * 12.0)
	halfTithi := 1
	if positionInTithi >= 6.0 {
		halfTithi = 2
	}

	span.SetAttributes(
		attribute.Float64("position_in_tithi", positionInTithi),
		attribute.Int("half_tithi", halfTithi),
	)

	// Calculate Karana number using the traditional cycle
	karanaNumber := kc.calculateKaranaNumber(tithi.Number, halfTithi)

	span.SetAttributes(attribute.Int("karana_number", karanaNumber))

	// Get Karana details
	karanaDetails := KaranaData[karanaNumber]

	// Calculate approximate start and end times based on Tithi times
	// Each Karana is half a Tithi, so approximately 12.395 hours
	startTime, endTime := kc.calculateKaranaTimesFromTithi(ctx, tithi, halfTithi)

	span.SetAttributes(
		attribute.String("karana_name", karanaDetails.Name),
		attribute.String("karana_type", string(karanaDetails.Type)),
		attribute.String("karana_description", karanaDetails.Description),
		attribute.Bool("is_vishti", karanaDetails.IsVishti),
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
	)

	duration := endTime.Sub(startTime).Hours()

	karana := &KaranaInfo{
		Number:      karanaNumber,
		Name:        karanaDetails.Name,
		Type:        karanaDetails.Type,
		Description: karanaDetails.Description,
		IsVishti:    karanaDetails.IsVishti,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    duration,
		MoonSunDiff: tithi.MoonSunDiff,
		TithiNumber: tithi.Number,
		HalfTithi:   halfTithi,
	}

	span.AddEvent("Karana calculation completed", trace.WithAttributes(
		attribute.Int("karana_number", karanaNumber),
		attribute.String("karana_name", karanaDetails.Name),
		attribute.Bool("is_vishti", karanaDetails.IsVishti),
		attribute.Float64("duration_hours", duration),
	))

	return karana, nil
}

// calculateKaranaNumber determines the Karana number based on Tithi and half
func (kc *KaranaCalculator) calculateKaranaNumber(tithiNumber, halfTithi int) int {
	// Special cases for fixed Karanas
	if tithiNumber == 30 && halfTithi == 1 {
		// First half of Amavasya (30th Tithi) - Shakuni
		return 9
	}
	if tithiNumber == 30 && halfTithi == 2 {
		// Second half of Amavasya - Chatushpada
		return 10
	}
	if tithiNumber == 1 && halfTithi == 1 {
		// First half of Pratipada (1st Tithi) - Naga
		return 11
	}
	if tithiNumber == 1 && halfTithi == 2 {
		// Second half of Pratipada - Kintughna (same as 1st but special name Kimstughna)
		return 1
	}

	// For movable Karanas (Tithis 2-29), use the 8-Karana cycle
	// Starting from the second half of Pratipada, cycle through 8 movable Karanas
	
	// Calculate the position in the movable cycle
	// Tithis 2-29 = 28 Tithis × 2 halves = 56 half-Tithis
	// Starting from Tithi 2, half 1
	
	var karanaPosition int
	if tithiNumber >= 2 && tithiNumber <= 29 {
		// Position in the movable cycle (0-55)
		karanaPosition = ((tithiNumber - 2) * 2) + (halfTithi - 1)
		// Cycle through the 8 movable Karanas (1-8)
		karanaNumber := (karanaPosition % 8) + 1
		return karanaNumber
	}

	// Fallback (should not happen with proper input)
	return 1
}

// calculateKaranaTimesFromTithi estimates the start and end times of a Karana based on Tithi times
func (kc *KaranaCalculator) calculateKaranaTimesFromTithi(ctx context.Context, tithi *TithiInfo, halfTithi int) (startTime, endTime time.Time) {
	ctx, span := kc.observer.CreateSpan(ctx, "KaranaCalculator.calculateKaranaTimesFromTithi")
	defer span.End()

	span.SetAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.Int("half_tithi", halfTithi),
		attribute.Float64("tithi_duration_hours", tithi.Duration),
	)

	// Each Karana is half a Tithi
	karanaDuration := time.Duration(tithi.Duration/2.0) * time.Hour

	if halfTithi == 1 {
		// First half of Tithi
		startTime = tithi.StartTime
		endTime = tithi.StartTime.Add(karanaDuration)
	} else {
		// Second half of Tithi
		startTime = tithi.StartTime.Add(karanaDuration)
		endTime = tithi.EndTime
	}

	span.SetAttributes(
		attribute.Float64("karana_duration_hours", karanaDuration.Hours()),
		attribute.String("calculated_start_time", startTime.Format("2006-01-02 15:04:05")),
		attribute.String("calculated_end_time", endTime.Format("2006-01-02 15:04:05")),
	)

	span.AddEvent("Karana times calculated", trace.WithAttributes(
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
		attribute.Float64("duration_hours", endTime.Sub(startTime).Hours()),
	))

	return startTime, endTime
}

// GetKaranaFromLongitudes is a convenience function for direct longitude input
func (kc *KaranaCalculator) GetKaranaFromLongitudes(ctx context.Context, sunLong, moonLong float64, date time.Time) (*KaranaInfo, error) {
	ctx, span := kc.observer.CreateSpan(ctx, "KaranaCalculator.GetKaranaFromLongitudes")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("date", date.Format("2006-01-02")),
	)

	// Get Tithi from longitudes first
	tithi, err := kc.tithiCalculator.GetTithiFromLongitudes(ctx, sunLong, moonLong, date)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get tithi for karana calculation: %w", err)
	}

	return kc.calculateKaranaFromTithi(ctx, tithi, date)
}

// IsAuspiciousKarana returns true if the Karana is considered auspicious
func IsAuspiciousKarana(karana *KaranaInfo) bool {
	// Vishti (Bhadra) and Gara are generally considered inauspicious
	if karana.IsVishti || karana.Name == "Gara" || karana.Name == "Shakuni" {
		return false
	}
	return true
}

// GetKaranaTypeDescription returns a description of the Karana type
func GetKaranaTypeDescription(karanaType KaranaType) string {
	switch karanaType {
	case KaranaTypeMovable:
		return "Movable Karana - cycles through the lunar month, each has specific qualities"
	case KaranaTypeFixed:
		return "Fixed Karana - appears in specific positions during new moon and first tithi"
	default:
		return "Unknown karana type"
	}
}

// GetKaranaRecommendations returns recommendations based on the Karana
func GetKaranaRecommendations(karana *KaranaInfo) string {
	switch karana.Name {
	case "Vishti":
		return "Avoid all important activities, travel, and new ventures. Time for rest and introspection."
	case "Gara":
		return "Avoid important activities. Not favorable for new beginnings."
	case "Shakuni":
		return "Inauspicious time. Avoid important decisions and activities."
	case "Bava":
		return "Good time for creative activities, learning, and joyful pursuits."
	case "Balava":
		return "Favorable for activities requiring strength and determination."
	case "Kaulava":
		return "Good for family-related activities and domestic affairs."
	case "Vanija":
		return "Excellent time for business, trade, and commercial activities."
	case "Taitila":
		return "Good for detailed work, craftsmanship, and precision tasks."
	case "Kintughna":
		return "Favorable for activities that require removing obstacles or enemies."
	case "Chatushpada":
		return "Stable and grounding energy. Good for foundational work."
	case "Naga":
		return "Mysterious energy. Good for spiritual practices and transformation."
	default:
		return "General karana with moderate influence."
	}
}

// ValidateKaranaCalculation validates a Karana calculation result
func ValidateKaranaCalculation(karana *KaranaInfo) error {
	if karana == nil {
		return fmt.Errorf("karana cannot be nil")
	}

	if karana.Number < 1 || karana.Number > 11 {
		return fmt.Errorf("invalid karana number: %d, must be between 1 and 11", karana.Number)
	}

	if karana.TithiNumber < 1 || karana.TithiNumber > 30 {
		return fmt.Errorf("invalid tithi number: %d, must be between 1 and 30", karana.TithiNumber)
	}

	if karana.HalfTithi < 1 || karana.HalfTithi > 2 {
		return fmt.Errorf("invalid half tithi: %d, must be 1 or 2", karana.HalfTithi)
	}

	if karana.MoonSunDiff < 0 || karana.MoonSunDiff >= 360 {
		return fmt.Errorf("invalid moon-sun difference: %f, must be between 0 and 360 degrees", karana.MoonSunDiff)
	}

	if karana.Duration <= 0 || karana.Duration > 24 {
		return fmt.Errorf("invalid karana duration: %f hours, must be positive and reasonable", karana.Duration)
	}

	if karana.EndTime.Before(karana.StartTime) {
		return fmt.Errorf("karana end time cannot be before start time")
	}

	if karana.Name == "" {
		return fmt.Errorf("karana name cannot be empty")
	}

	// Validate type is one of the defined values
	switch karana.Type {
	case KaranaTypeMovable, KaranaTypeFixed:
		// Valid type
	default:
		return fmt.Errorf("invalid karana type: %s", karana.Type)
	}

	return nil
}