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

// TithiNames maps Tithi numbers to their standard Sanskrit names
var TithiNames = map[int]string{
	1: "Pratipada", 2: "Dwitiya", 3: "Tritiya", 4: "Chaturthi", 5: "Panchami",
	6: "Shashthi", 7: "Saptami", 8: "Ashtami", 9: "Navami", 10: "Dashami",
	11: "Ekadashi", 12: "Dwadashi", 13: "Trayodashi", 14: "Chaturdashi", 15: "Purnima",
	16: "Pratipada", 17: "Dwitiya", 18: "Tritiya", 19: "Chaturthi", 20: "Panchami",
	21: "Shashthi", 22: "Saptami", 23: "Ashtami", 24: "Navami", 25: "Dashami",
	26: "Ekadashi", 27: "Dwadashi", 28: "Trayodashi", 29: "Chaturdashi", 30: "Amavasya",
}

// TraditionalTithiNames maps Tithi numbers to traditional Sanskrit names with preferred spellings
var TraditionalTithiNames = map[int]string{
	1: "Pratipada", 2: "Dvithiya", 3: "Thuthiya", 4: "Chathurthi", 5: "Panchami",
	6: "Shashthi", 7: "Sapthami", 8: "Ashtami", 9: "Navami", 10: "Dashami",
	11: "Ekadashi", 12: "Dvadashi", 13: "Thrayodashi", 14: "Chathurdashi", 15: "Pournima",
	16: "Pratipada", 17: "Dvithiya", 18: "Thuthiya", 19: "Chathurthi", 20: "Panchami",
	21: "Shashthi", 22: "Sapthami", 23: "Ashtami", 24: "Navami", 25: "Dashami",
	26: "Ekadashi", 27: "Dvadashi", 28: "Thrayodashi", 29: "Chathurdashi", 30: "Amavasya",
}

// PakshaNames maps paksha day numbers (1-15) to their traditional names
var PakshaNames = map[int]string{
	1: "Pratipada", 2: "Dvithiya", 3: "Thuthiya", 4: "Chathurthi", 5: "Panchami",
	6: "Shashthi", 7: "Sapthami", 8: "Ashtami", 9: "Navami", 10: "Dashami",
	11: "Ekadashi", 12: "Dvadashi", 13: "Thrayodashi", 14: "Chathurdashi", 15: "Pournima",
}

// GetTithiForDate calculates the Tithi for a given date with default Purnimanta system
func (tc *TithiCalculator) GetTithiForDate(ctx context.Context, date time.Time) (*TithiInfo, error) {
	return tc.GetTithiForDateWithCalendarSystem(ctx, date, "Purnimanta")
}

// GetTithiForDateWithCalendarSystem calculates the Tithi for a given date with specified calendar system
func (tc *TithiCalculator) GetTithiForDateWithCalendarSystem(ctx context.Context, date time.Time, calendarSystem string) (*TithiInfo, error) {
	ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.GetTithiForDateWithCalendarSystem")
	defer span.End()

	span.SetAttributes(
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("timezone", date.Location().String()),
		attribute.String("calendar_system", calendarSystem),
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
	tithi, err := tc.calculateExactTithiFromLongitudes(ctx, sunLong, moonLong, noonDate, calendarSystem)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.String("paksha", tithi.Paksha),
		attribute.Int("paksha_day", tithi.PakshaDay),
		attribute.String("traditional_name", tithi.TraditionalName),
		attribute.String("tithi_type", string(tithi.Type)),
		attribute.Bool("is_shukla", tithi.IsShukla),
		attribute.Float64("moon_sun_diff", tithi.MoonSunDiff),
		attribute.String("calendar_system", tithi.CalendarSystem),
	)

	span.AddEvent("Tithi calculated", trace.WithAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.String("paksha", tithi.Paksha),
		attribute.String("traditional_name", tithi.TraditionalName),
		attribute.String("tithi_type", string(tithi.Type)),
	))

	return tithi, nil
}

// GetTithiForTimeWithCalendarSystem calculates the Tithi at an exact time.
func (tc *TithiCalculator) GetTithiForTimeWithCalendarSystem(ctx context.Context, dateTime time.Time, calendarSystem string) (*TithiInfo, error) {
	ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.GetTithiForTimeWithCalendarSystem")
	defer span.End()

	span.SetAttributes(
		attribute.String("date_time", dateTime.Format(time.RFC3339)),
		attribute.String("timezone", dateTime.Location().String()),
		attribute.String("calendar_system", calendarSystem),
	)

	jd := ephemeris.TimeToJulianDay(dateTime)
	span.SetAttributes(attribute.Float64("julian_day", float64(jd)))

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

	tithi, err := tc.calculateExactTithiFromLongitudes(ctx, sunLong, moonLong, dateTime, calendarSystem)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("tithi_number", tithi.Number),
		attribute.String("tithi_name", tithi.Name),
		attribute.String("paksha", tithi.Paksha),
		attribute.Int("paksha_day", tithi.PakshaDay),
	)

	return tithi, nil
}
