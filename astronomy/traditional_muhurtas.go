package astronomy

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CalculateMuhurtas calculates all 30 muhurtas of the day with their qualities
func CalculateMuhurtas(loc Location, date time.Time) ([]*MuhurtaInfo, error) {
	return CalculateMuhurtasWithContext(context.Background(), loc, date)
}

// CalculateMuhurtasWithContext calculates all muhurtas with OpenTelemetry tracing
func CalculateMuhurtasWithContext(ctx context.Context, loc Location, date time.Time) ([]*MuhurtaInfo, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "CalculateMuhurtas")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)

	// Get sunrise and sunset times
	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Calculate day length and muhurta duration
	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
	muhurtaDuration := dayLength / 30

	span.SetAttributes(
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.Float64("muhurta_duration_minutes", muhurtaDuration.Minutes()),
	)

	// Traditional muhurta names (simplified list of 30)
	muhurtaNames := []string{
		"Rudra", "Ahi", "Mitra", "Pritas", "Vasu", "Varaha", "Vishve", "Abhijit",
		"Savitri", "Aditya", "Sadhya", "Ganga", "Brahma", "Yamya", "Vasu", "Varuna",
		"Aryama", "Bhaga", "Girish", "Dhanvantari", "Ananda", "Rakshasa", "Sarpa", "Shakra",
		"Indra", "Vayu", "Dhruva", "Vaidhriti", "Ketu", "Shubha",
	}

	muhurtas := make([]*MuhurtaInfo, 30)

	for i := 0; i < 30; i++ {
		start := sunTimes.Sunrise.Add(time.Duration(i) * muhurtaDuration)
		end := sunTimes.Sunrise.Add(time.Duration(i+1) * muhurtaDuration)

		// Determine quality based on traditional knowledge
		var quality string
		var recommended []string
		var avoid []string
		var auspicious bool

		// Abhijit Muhurta (8th) is the most auspicious
		if i == 7 {
			quality = "excellent"
			auspicious = true
			recommended = []string{"All activities", "Important ceremonies", "Travel", "Business"}
			avoid = []string{}
		} else if i >= 5 && i <= 9 { // Middle of the day muhurtas
			quality = "good"
			auspicious = true
			recommended = []string{"General activities", "Work", "Study"}
			avoid = []string{}
		} else if i >= 15 && i <= 20 { // Late afternoon muhurtas
			quality = "neutral"
			auspicious = true
			recommended = []string{"Routine work"}
			avoid = []string{"Important ceremonies"}
		} else {
			quality = "neutral"
			auspicious = true
			recommended = []string{"General activities"}
			avoid = []string{}
		}

		muhurtas[i] = &MuhurtaInfo{
			Name: muhurtaNames[i],
			Period: &TimePeriod{
				Start:       start,
				End:         end,
				Duration:    int(muhurtaDuration.Minutes()),
				Description: "Muhurta " + muhurtaNames[i],
				Auspicious:  auspicious,
			},
			Quality:     quality,
			Recommended: recommended,
			Avoid:       avoid,
		}
	}

	span.SetAttributes(attribute.Int("total_muhurtas", len(muhurtas)))
	span.AddEvent("All muhurtas calculated", trace.WithAttributes(
		attribute.Int("total_muhurtas", len(muhurtas)),
		attribute.Float64("muhurta_duration_minutes", muhurtaDuration.Minutes()),
	))

	return muhurtas, nil
}
