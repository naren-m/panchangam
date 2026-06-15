package astronomy

import (
	"context"
	"math"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CalculateLunarPhase calculates detailed lunar phase information
func CalculateLunarPhase(date time.Time) (*LunarPhase, error) {
	return CalculateLunarPhaseWithContext(context.Background(), date)
}

// CalculateLunarPhaseWithContext calculates lunar phase with OpenTelemetry tracing
func CalculateLunarPhaseWithContext(ctx context.Context, date time.Time) (*LunarPhase, error) {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "CalculateLunarPhase")
	defer span.End()

	span.SetAttributes(attribute.String("date", date.Format("2006-01-02")))

	year, month, day := date.Date()
	jd := julianDayNumber(year, int(month), day)

	// Calculate days since known new moon (January 6, 2000)
	knownNewMoon := julianDayNumber(2000, 1, 6) // JD 2451549.5
	daysSinceKnownNew := jd - knownNewMoon

	// Calculate the current lunation number
	lunationNumber := daysSinceKnownNew / LunarSynodicMonth
	currentLunation := math.Floor(lunationNumber)

	// Days into current lunar cycle
	daysIntoCycle := (lunationNumber - currentLunation) * LunarSynodicMonth

	// Calculate phase (0.0 = new moon, 0.5 = full moon)
	phase := daysIntoCycle / LunarSynodicMonth

	// Calculate illumination percentage
	phaseAngle := phase * 2 * math.Pi
	illumination := (1 - math.Cos(phaseAngle)) / 2 * 100

	// Determine phase name
	var phaseName string
	if phase < 0.125 || phase >= 0.875 {
		phaseName = "New Moon"
	} else if phase < 0.375 {
		phaseName = "Waxing Crescent"
	} else if phase < 0.625 {
		phaseName = "Full Moon"
	} else {
		phaseName = "Waning Crescent"
	}

	// More precise phase names
	if phase >= 0.125 && phase < 0.25 {
		phaseName = "Waxing Crescent"
	} else if phase >= 0.25 && phase < 0.375 {
		phaseName = "First Quarter"
	} else if phase >= 0.375 && phase < 0.5 {
		phaseName = "Waxing Gibbous"
	} else if phase >= 0.5 && phase < 0.625 {
		phaseName = "Full Moon"
	} else if phase >= 0.625 && phase < 0.75 {
		phaseName = "Waning Gibbous"
	} else if phase >= 0.75 && phase < 0.875 {
		phaseName = "Last Quarter"
	}

	// Calculate next major phase
	var nextPhaseJD float64
	var nextPhaseType float64

	if phase < 0.25 {
		nextPhaseType = 0.25 // First Quarter
		nextPhaseJD = knownNewMoon + (currentLunation * LunarSynodicMonth) + (0.25 * LunarSynodicMonth)
	} else if phase < 0.5 {
		nextPhaseType = 0.5 // Full Moon
		nextPhaseJD = knownNewMoon + (currentLunation * LunarSynodicMonth) + (0.5 * LunarSynodicMonth)
	} else if phase < 0.75 {
		nextPhaseType = 0.75 // Last Quarter
		nextPhaseJD = knownNewMoon + (currentLunation * LunarSynodicMonth) + (0.75 * LunarSynodicMonth)
	} else {
		nextPhaseType = 0.0 // New Moon (next cycle)
		nextPhaseJD = knownNewMoon + ((currentLunation + 1) * LunarSynodicMonth)
	}

	// Convert Julian day to time
	nextPhaseTime := julianDayToTime(nextPhaseJD, date.Location())

	result := &LunarPhase{
		Phase:        phase,
		Illumination: illumination,
		Name:         phaseName,
		Age:          daysIntoCycle,
		NextPhase:    nextPhaseTime,
	}

	span.SetAttributes(
		attribute.Float64("julian_day", jd),
		attribute.Float64("days_since_known_new", daysSinceKnownNew),
		attribute.Float64("lunation_number", lunationNumber),
		attribute.Float64("current_lunation", currentLunation),
		attribute.Float64("days_into_cycle", daysIntoCycle),
		attribute.Float64("phase", phase),
		attribute.Float64("illumination", illumination),
		attribute.String("phase_name", phaseName),
		attribute.Float64("age_days", daysIntoCycle),
		attribute.Float64("next_phase_type", nextPhaseType),
		attribute.String("next_phase_time", nextPhaseTime.Format("2006-01-02 15:04:05")),
	)

	span.AddEvent("Lunar phase calculated", trace.WithAttributes(
		attribute.Float64("phase", phase),
		attribute.Float64("illumination", illumination),
		attribute.String("phase_name", phaseName),
		attribute.Float64("age_days", daysIntoCycle),
	))

	return result, nil
}

// julianDayToTime converts Julian day to time.Time
func julianDayToTime(jd float64, loc *time.Location) time.Time {
	// Add 0.5 to adjust for noon-based Julian day
	jd += 0.5

	// Integer part is the number of days since Julian epoch
	z := int(jd)
	f := jd - float64(z)

	var a int
	if z < 2299161 {
		a = z
	} else {
		alpha := int((float64(z) - 1867216.25) / 36524.25)
		a = z + 1 + alpha - alpha/4
	}

	b := a + 1524
	c := int((float64(b) - 122.1) / 365.25)
	d := int(365.25 * float64(c))
	e := int(float64(b-d) / 30.6001)

	// Day of month
	day := b - d - int(30.6001*float64(e))

	// Month
	var month int
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}

	// Year
	var year int
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}

	// Hours, minutes, seconds from fractional part
	dayFraction := f
	hours := int(dayFraction * 24)
	minuteFraction := (dayFraction*24 - float64(hours)) * 60
	minutes := int(minuteFraction)
	seconds := int((minuteFraction - float64(minutes)) * 60)

	return time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, loc)
}
