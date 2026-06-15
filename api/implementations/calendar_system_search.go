package implementations

import (
	"context"
	"fmt"
	"time"
)

// Helper methods for finding lunar events

func (c *CalendarSystemPlugin) findPreviousAmavasya(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search backwards for the most recent Amavasya
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days back
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 30 { // Amavasya
			return tithi.StartTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, -1)
	}

	return time.Time{}, fmt.Errorf("could not find previous Amavasya")
}

func (c *CalendarSystemPlugin) findNextAmavasya(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search forwards for the next Amavasya
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days forward
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 30 { // Amavasya
			return tithi.EndTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, 1)
	}

	return time.Time{}, fmt.Errorf("could not find next Amavasya")
}

func (c *CalendarSystemPlugin) findPreviousPurnima(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search backwards for the most recent Purnima
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days back
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 15 { // Purnima
			return tithi.StartTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, -1)
	}

	return time.Time{}, fmt.Errorf("could not find previous Purnima")
}

func (c *CalendarSystemPlugin) findNextPurnima(ctx context.Context, fromDate time.Time) (time.Time, error) {
	// Search forwards for the next Purnima
	searchDate := fromDate
	for i := 0; i < 45; i++ { // Search up to 45 days forward
		tithi, err := c.tithiCalculator.GetTithiForDate(ctx, searchDate)
		if err != nil {
			return time.Time{}, err
		}

		if tithi.Number == 15 { // Purnima
			return tithi.EndTime, nil
		}

		searchDate = searchDate.AddDate(0, 0, 1)
	}

	return time.Time{}, fmt.Errorf("could not find next Purnima")
}
