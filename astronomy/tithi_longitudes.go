package astronomy

import (
	"context"
	"math"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// calculateTithiFromLongitudes calculates Tithi from Sun and Moon longitudes
func (tc *TithiCalculator) calculateTithiFromLongitudes(ctx context.Context, sunLong, moonLong float64, referenceDate time.Time, calendarSystem string) (*TithiInfo, error) {
	return tc.calculateTithiFromLongitudesWithTiming(ctx, sunLong, moonLong, referenceDate, calendarSystem, false)
}

func (tc *TithiCalculator) calculateExactTithiFromLongitudes(ctx context.Context, sunLong, moonLong float64, referenceDate time.Time, calendarSystem string) (*TithiInfo, error) {
	return tc.calculateTithiFromLongitudesWithTiming(ctx, sunLong, moonLong, referenceDate, calendarSystem, true)
}

func (tc *TithiCalculator) calculateTithiFromLongitudesWithTiming(ctx context.Context, sunLong, moonLong float64, referenceDate time.Time, calendarSystem string, exactTimes bool) (*TithiInfo, error) {
	ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.calculateTithiFromLongitudes")
	defer span.End()

	calendarSystem = normalizeTithiCalendarSystem(calendarSystem)

	span.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("calendar_system", calendarSystem),
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
	baseTithiNumber := int(tithiFloat) + 1

	// Ensure base Tithi number is in valid range (1-30)
	if baseTithiNumber > 30 {
		baseTithiNumber = 30
	}
	if baseTithiNumber < 1 {
		baseTithiNumber = 1
	}

	span.SetAttributes(
		attribute.Float64("tithi_float", tithiFloat),
		attribute.Int("base_tithi_number", baseTithiNumber),
	)

	// Calculate paksha information and adjust for calendar system
	var tithiNumber, pakshaDay int
	var paksha string
	var isShukla bool
	var traditionalName string

	if calendarSystem == "Amanta" {
		// In Amanta system, Krishna paksha comes first (1-15), then Shukla paksha (1-15)
		// But astronomically, we still use 1-30 numbering where 1-15 is Shukla, 16-30 is Krishna

		if baseTithiNumber <= 15 {
			// Shukla Paksha (waxing moon)
			isShukla = true
			paksha = "Shukla"
			pakshaDay = baseTithiNumber
			tithiNumber = baseTithiNumber // Keep original numbering for internal calculations
		} else {
			// Krishna Paksha (waning moon) - adjust numbering for Amanta
			isShukla = false
			paksha = "Krishna"
			pakshaDay = baseTithiNumber - 15 // 16 becomes 1, 17 becomes 2, etc.
			tithiNumber = baseTithiNumber    // Keep original for internal calculations
		}

		// Get traditional name based on paksha day
		if pakshaDay == 15 && !isShukla {
			traditionalName = "Amavasya" // Special case for new moon
		} else {
			traditionalName = PakshaNames[pakshaDay]
		}
	} else {
		// Purnimanta system (standard)
		if baseTithiNumber <= 15 {
			// Shukla Paksha (waxing moon)
			isShukla = true
			paksha = "Shukla"
			pakshaDay = baseTithiNumber
		} else {
			// Krishna Paksha (waning moon)
			isShukla = false
			paksha = "Krishna"
			pakshaDay = baseTithiNumber - 15 // 16 becomes 1, 17 becomes 2, etc.
		}
		tithiNumber = baseTithiNumber
		traditionalName = TraditionalTithiNames[baseTithiNumber]
	}

	// Get standard name and traditional name
	tithiName := TithiNames[baseTithiNumber]

	// Determine Tithi type based on paksha day
	tithiType := getTithiType(pakshaDay)

	var startTime, endTime time.Time
	var err error
	if exactTimes {
		startTime, endTime, err = tc.calculateTithiTimes(ctx, baseTithiNumber, referenceDate)
		if err != nil {
			return nil, err
		}
	} else {
		startTime, endTime = estimateTithiTimes(tithiFloat, referenceDate)
	}

	span.SetAttributes(
		attribute.String("tithi_name", tithiName),
		attribute.String("traditional_name", traditionalName),
		attribute.String("paksha", paksha),
		attribute.Int("paksha_day", pakshaDay),
		attribute.String("tithi_type", string(tithiType)),
		attribute.Bool("is_shukla", isShukla),
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
	)

	duration := endTime.Sub(startTime).Hours()

	tithi := &TithiInfo{
		Number:          tithiNumber,
		Name:            tithiName,
		Type:            tithiType,
		StartTime:       startTime,
		EndTime:         endTime,
		Duration:        duration,
		IsShukla:        isShukla,
		Paksha:          paksha,
		PakshaDay:       pakshaDay,
		TraditionalName: traditionalName,
		MoonSunDiff:     moonSunDiff,
		CalendarSystem:  calendarSystem,
	}

	span.AddEvent("Tithi calculation completed", trace.WithAttributes(
		attribute.Int("tithi_number", tithiNumber),
		attribute.String("tithi_name", tithiName),
		attribute.String("traditional_name", traditionalName),
		attribute.String("paksha", paksha),
		attribute.Int("paksha_day", pakshaDay),
		attribute.Float64("duration_hours", duration),
	))

	return tithi, nil
}

func tithiNumberFromLongitudes(sunLong, moonLong float64) int {
	diff := moonLong - sunLong
	for diff < 0 {
		diff += 360
	}
	for diff >= 360 {
		diff -= 360
	}

	number := int(math.Floor(diff/12.0)) + 1
	if number < 1 {
		return 1
	}
	if number > 30 {
		return 30
	}
	return number
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

func normalizeTithiCalendarSystem(calendarSystem string) string {
	switch strings.ToLower(strings.TrimSpace(calendarSystem)) {
	case "", "traditional", "purnimanta":
		return "Purnimanta"
	case "amanta":
		return "Amanta"
	default:
		return calendarSystem
	}
}

// GetTithiFromLongitudes is a convenience function for direct longitude input with default Purnimanta system
func (tc *TithiCalculator) GetTithiFromLongitudes(ctx context.Context, sunLong, moonLong float64, date time.Time) (*TithiInfo, error) {
	return tc.GetTithiFromLongitudesWithCalendarSystem(ctx, sunLong, moonLong, date, "Purnimanta")
}

// GetTithiFromLongitudesWithCalendarSystem is a convenience function for direct longitude input with specified calendar system
func (tc *TithiCalculator) GetTithiFromLongitudesWithCalendarSystem(ctx context.Context, sunLong, moonLong float64, date time.Time, calendarSystem string) (*TithiInfo, error) {
	ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.GetTithiFromLongitudesWithCalendarSystem")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("sun_longitude", sunLong),
		attribute.Float64("moon_longitude", moonLong),
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("calendar_system", calendarSystem),
	)

	return tc.calculateTithiFromLongitudes(ctx, sunLong, moonLong, date, calendarSystem)
}
