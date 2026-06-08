package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
)

// GetEvents returns inauspicious events for a specific date and location
func (i *InauspiciousEventsPlugin) GetEvents(ctx context.Context, date time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !i.enabled {
		return nil, fmt.Errorf("inauspicious events plugin is not enabled")
	}

	var events []api.Event

	// Calculate sunrise and sunset for the location
	astroLocation := astronomy.Location{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}

	sunTimes, err := astronomy.CalculateSunTimes(astroLocation, date)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate sun times: %w", err)
	}

	if err := convertSunTimesToTimezone(sunTimes, location.Timezone); err != nil {
		return nil, err
	}

	// Calculate day length
	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
	weekday := date.Weekday()

	// Generate Rahu Kalam event
	rahuKalam := i.calculateRahuKalamEvent(sunTimes.Sunrise, dayLength, weekday, region)
	events = append(events, rahuKalam)

	// Generate Yamagandam event
	yamagandam := i.calculateYamgandamEvent(sunTimes.Sunrise, dayLength, weekday, region)
	events = append(events, yamagandam)

	// Generate Gulika Kalam event
	gulikaKalam := i.calculateGulikaKalamEvent(sunTimes.Sunrise, dayLength, weekday, region)
	events = append(events, gulikaKalam)

	return events, nil
}

// GetEventsInRange returns inauspicious events for a date range
func (i *InauspiciousEventsPlugin) GetEventsInRange(ctx context.Context, start, end time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !i.enabled {
		return nil, fmt.Errorf("inauspicious events plugin is not enabled")
	}

	if err := validateEventDateRange(start, end); err != nil {
		return nil, err
	}

	var allEvents []api.Event

	current := start
	for current.Before(end) || current.Equal(end) {
		dayEvents, err := i.GetEvents(ctx, current, location, region)
		if err != nil {
			return nil, fmt.Errorf("failed to get events for %s: %w", current.Format("2006-01-02"), err)
		}
		allEvents = append(allEvents, dayEvents...)
		current = current.AddDate(0, 0, 1)
	}

	return allEvents, nil
}
