package astronomy

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func estimateTithiTimes(tithiFloat float64, referenceDate time.Time) (startTime, endTime time.Time) {
	tithiProgress := tithiFloat - math.Floor(tithiFloat)
	avgTithiDuration := time.Duration(24.79 * float64(time.Hour))
	timeIntoTithi := time.Duration(tithiProgress * float64(avgTithiDuration))

	startTime = referenceDate.Add(-timeIntoTithi)
	endTime = startTime.Add(avgTithiDuration)
	return startTime, endTime
}

func (tc *TithiCalculator) calculateTithiTimes(ctx context.Context, tithiNumber int, referenceDate time.Time) (startTime, endTime time.Time, err error) {
	_, span := tc.observer.CreateSpan(ctx, "TithiCalculator.calculateTithiTimes")
	defer span.End()

	span.SetAttributes(
		attribute.Int("tithi_number", tithiNumber),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	startTime, err = tc.findTithiBoundary(ctx, referenceDate, tithiNumber, -1)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime, err = tc.findTithiBoundary(ctx, referenceDate, tithiNumber, 1)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	span.SetAttributes(
		attribute.String("calculated_start_time", startTime.Format("2006-01-02 15:04:05")),
		attribute.String("calculated_end_time", endTime.Format("2006-01-02 15:04:05")),
	)

	span.AddEvent("Tithi times calculated", trace.WithAttributes(
		attribute.String("start_time", startTime.Format("15:04:05")),
		attribute.String("end_time", endTime.Format("15:04:05")),
		attribute.Float64("duration_hours", endTime.Sub(startTime).Hours()),
	))

	return startTime, endTime, nil
}

func (tc *TithiCalculator) findTithiBoundary(ctx context.Context, referenceDate time.Time, tithiNumber int, direction int) (time.Time, error) {
	step := 6 * time.Hour
	if direction > 0 {
		lo := referenceDate
		hi := referenceDate.Add(step)
		for i := 0; i < 8; i++ {
			number, err := tc.tithiNumberAt(ctx, hi)
			if err != nil {
				return time.Time{}, err
			}
			if number != tithiNumber {
				return tc.refineTithiBoundary(ctx, lo, hi, tithiNumber, direction)
			}
			lo = hi
			hi = hi.Add(step)
		}
		return time.Time{}, fmt.Errorf("failed to find next tithi boundary near %s", referenceDate.Format(time.RFC3339))
	}

	hi := referenceDate
	lo := referenceDate.Add(-step)
	for i := 0; i < 8; i++ {
		number, err := tc.tithiNumberAt(ctx, lo)
		if err != nil {
			return time.Time{}, err
		}
		if number != tithiNumber {
			return tc.refineTithiBoundary(ctx, lo, hi, tithiNumber, direction)
		}
		hi = lo
		lo = lo.Add(-step)
	}
	return time.Time{}, fmt.Errorf("failed to find previous tithi boundary near %s", referenceDate.Format(time.RFC3339))
}

func (tc *TithiCalculator) refineTithiBoundary(ctx context.Context, lo, hi time.Time, tithiNumber int, direction int) (time.Time, error) {
	for i := 0; i < 40; i++ {
		mid := lo.Add(hi.Sub(lo) / 2)
		number, err := tc.tithiNumberAt(ctx, mid)
		if err != nil {
			return time.Time{}, err
		}
		if direction > 0 {
			if number == tithiNumber {
				lo = mid
			} else {
				hi = mid
			}
		} else {
			if number == tithiNumber {
				hi = mid
			} else {
				lo = mid
			}
		}
	}
	return hi, nil
}

func (tc *TithiCalculator) tithiNumberAt(ctx context.Context, dateTime time.Time) (int, error) {
	jd := ephemeris.TimeToJulianDay(dateTime)
	positions, err := tc.ephemerisManager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		return 0, fmt.Errorf("failed to get planetary positions: %w", err)
	}
	return tithiNumberFromLongitudes(positions.Sun.Longitude, positions.Moon.Longitude), nil
}
