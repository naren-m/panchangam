package implementations

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

// calculateDailyAuspiciousPeriods calculates general auspicious periods throughout the day
func (m *MuhurtaPlugin) calculateDailyAuspiciousPeriods(sunrise, sunset time.Time, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta

	// Pratah Kaal (Morning period) - First 3 hours after sunrise
	pratahEnd := sunrise.Add(3 * time.Hour)
	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Pratah Kaal",
		NameLocal:    "प्रातः काल",
		StartTime:    sunrise,
		EndTime:      pratahEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"daily_activities", "exercise", "study", "work"},
		Significance: "Morning period, good for daily activities",
		Metadata: map[string]interface{}{
			"period":         "morning",
			"duration_hours": 3,
		},
	})

	// Sandhya Kaal (Evening twilight) - Last hour before sunset
	sandhyaStart := sunset.Add(-1 * time.Hour)
	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Sandhya Kaal",
		NameLocal:    "संध्या काल",
		StartTime:    sandhyaStart,
		EndTime:      sunset,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"prayer", "meditation", "spiritual_practices"},
		Significance: "Evening twilight, good for spiritual activities",
		Metadata: map[string]interface{}{
			"period":         "evening_twilight",
			"duration_hours": 1,
		},
	})

	return muhurtas
}
