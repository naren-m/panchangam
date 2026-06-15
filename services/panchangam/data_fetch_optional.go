package panchangam

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/astronomy"
)

type optionalPanchangamData struct {
	lunarTimes         *astronomy.LunarTimes
	lunarPhase         *astronomy.LunarPhase
	traditionalPeriods *astronomy.TraditionalPeriods
	festivals          []string
}

func calculateOptionalPanchangamData(ctx context.Context, date time.Time, location astronomy.Location, tithiNumber int) optionalPanchangamData {
	logger.DebugContext(ctx, "Calculating lunar times")
	lunarTimes, err := astronomy.CalculateLunarTimesWithContext(ctx, location, date)
	if err != nil {
		logger.WarnContext(ctx, "Failed to calculate lunar times", "error", err)
	}

	lunarPhase, err := astronomy.CalculateLunarPhaseWithContext(ctx, date)
	if err != nil {
		logger.WarnContext(ctx, "Failed to calculate lunar phase", "error", err)
	}

	logger.DebugContext(ctx, "Calculating traditional periods")
	traditionalPeriods, err := astronomy.CalculateTraditionalPeriodsWithContext(ctx, location, date)
	if err != nil {
		logger.WarnContext(ctx, "Failed to calculate traditional periods", "error", err)
	}

	logger.DebugContext(ctx, "Calculating festivals for date")
	festivals, err := astronomy.GetFestivalNamesForDate(ctx, date, tithiNumber)
	if err != nil {
		logger.WarnContext(ctx, "Failed to calculate festivals", "error", err)
		festivals = []string{}
	}

	return optionalPanchangamData{
		lunarTimes:         lunarTimes,
		lunarPhase:         lunarPhase,
		traditionalPeriods: traditionalPeriods,
		festivals:          festivals,
	}
}
