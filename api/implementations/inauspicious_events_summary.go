package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

// Helper function to format duration for display
func (i *InauspiciousEventsPlugin) formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d hour(s) %d minute(s)", hours, minutes)
	}
	return fmt.Sprintf("%d minute(s)", minutes)
}

// GetEventSummary provides a summary of all inauspicious periods for a day
func (i *InauspiciousEventsPlugin) GetEventSummary(ctx context.Context, date time.Time, location api.Location, region api.Region) (map[string]interface{}, error) {
	events, err := i.GetEvents(ctx, date, location, region)
	if err != nil {
		return nil, err
	}

	eventSummaries := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		eventSummary := map[string]interface{}{
			"name":       event.Name,
			"name_local": event.NameLocal,
			"type":       string(event.Type),
			"start_time": event.StartTime.Format("15:04"),
			"end_time":   event.EndTime.Format("15:04"),
			"duration":   i.formatDuration(event.EndTime.Sub(event.StartTime)),
		}
		eventSummaries = append(eventSummaries, eventSummary)
	}

	return map[string]interface{}{
		"date":         date.Format("2006-01-02"),
		"weekday":      date.Weekday().String(),
		"location":     location.Name,
		"region":       string(region),
		"total_events": len(events),
		"events":       eventSummaries,
	}, nil
}
