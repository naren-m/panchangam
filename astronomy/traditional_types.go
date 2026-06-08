package astronomy

import "time"

// TraditionalPeriods represents traditional Hindu time periods for a day
type TraditionalPeriods struct {
	RahuKalam      *TimePeriod `json:"rahu_kalam"`
	Yamagandam     *TimePeriod `json:"yamagandam"`
	GulikaKalam    *TimePeriod `json:"gulika_kalam"`
	AbhijitMuhurta *TimePeriod `json:"abhijit_muhurta"`
}

// TimePeriod represents a time period with start and end times
type TimePeriod struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Duration    int       `json:"duration_minutes"`
	Description string    `json:"description"`
	Auspicious  bool      `json:"auspicious"`
}

// MuhurtaInfo represents muhurta (auspicious time) information
type MuhurtaInfo struct {
	Name        string      `json:"name"`
	Period      *TimePeriod `json:"period"`
	Quality     string      `json:"quality"` // "good", "neutral", "avoid"
	Recommended []string    `json:"recommended_activities"`
	Avoid       []string    `json:"avoid_activities"`
}
