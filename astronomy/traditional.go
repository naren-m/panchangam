package astronomy

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TraditionalPeriods represents traditional Hindu time periods for a day
type TraditionalPeriods struct {
	RahuKalam     *TimePeriod `json:"rahu_kalam"`
	Yamagandam    *TimePeriod `json:"yamagandam"`
	GulikaKalam   *TimePeriod `json:"gulika_kalam"`
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

// CalculateTraditionalPeriods calculates all traditional time periods for a given location and date
func CalculateTraditionalPeriods(loc Location, date time.Time) (*TraditionalPeriods, error) {
	return CalculateTraditionalPeriodsWithContext(context.Background(), loc, date)
}

// CalculateTraditionalPeriodsWithContext calculates traditional periods with OpenTelemetry tracing
func CalculateTraditionalPeriodsWithContext(ctx context.Context, loc Location, date time.Time) (*TraditionalPeriods, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "CalculateTraditionalPeriods")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("timezone", date.Location().String()),
	)
	
	// First, get sunrise and sunset times
	ctx, sunTimesSpan := observer.CreateSpan(ctx, "getSunTimes")
	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	if err != nil {
		sunTimesSpan.RecordError(err)
		span.RecordError(err)
		return nil, err
	}
	sunTimesSpan.SetAttributes(
		attribute.String("sunrise", sunTimes.Sunrise.Format("15:04:05")),
		attribute.String("sunset", sunTimes.Sunset.Format("15:04:05")),
	)
	sunTimesSpan.End()
	
	// Calculate day length in minutes
	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
	dayLengthMinutes := int(dayLength.Minutes())
	
	span.SetAttributes(
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.Int("day_length_minutes", dayLengthMinutes),
	)
	
	// Calculate Rahu Kalam
	ctx, rahuSpan := observer.CreateSpan(ctx, "calculateRahuKalam")
	rahuKalam := calculateRahuKalam(ctx, sunTimes.Sunrise, sunTimes.Sunset, date)
	rahuSpan.SetAttributes(
		attribute.String("rahu_kalam_start", rahuKalam.Start.Format("15:04:05")),
		attribute.String("rahu_kalam_end", rahuKalam.End.Format("15:04:05")),
		attribute.Int("rahu_kalam_duration", rahuKalam.Duration),
	)
	rahuSpan.End()
	
	// Calculate Yamagandam
	ctx, yamaSpan := observer.CreateSpan(ctx, "calculateYamagandam")
	yamagandam := calculateYamagandam(ctx, sunTimes.Sunrise, sunTimes.Sunset, date)
	yamaSpan.SetAttributes(
		attribute.String("yamagandam_start", yamagandam.Start.Format("15:04:05")),
		attribute.String("yamagandam_end", yamagandam.End.Format("15:04:05")),
		attribute.Int("yamagandam_duration", yamagandam.Duration),
	)
	yamaSpan.End()
	
	// Calculate Gulika Kalam
	ctx, gulikaSpan := observer.CreateSpan(ctx, "calculateGulikaKalam")
	gulikaKalam := calculateGulikaKalam(ctx, sunTimes.Sunrise, sunTimes.Sunset, date)
	gulikaSpan.SetAttributes(
		attribute.String("gulika_kalam_start", gulikaKalam.Start.Format("15:04:05")),
		attribute.String("gulika_kalam_end", gulikaKalam.End.Format("15:04:05")),
		attribute.Int("gulika_kalam_duration", gulikaKalam.Duration),
	)
	gulikaSpan.End()
	
	// Calculate Abhijit Muhurta
	ctx, abhijitSpan := observer.CreateSpan(ctx, "calculateAbhijitMuhurta")
	abhijitMuhurta := calculateAbhijitMuhurta(ctx, sunTimes.Sunrise, sunTimes.Sunset)
	abhijitSpan.SetAttributes(
		attribute.String("abhijit_start", abhijitMuhurta.Start.Format("15:04:05")),
		attribute.String("abhijit_end", abhijitMuhurta.End.Format("15:04:05")),
		attribute.Int("abhijit_duration", abhijitMuhurta.Duration),
	)
	abhijitSpan.End()
	
	result := &TraditionalPeriods{
		RahuKalam:     rahuKalam,
		Yamagandam:    yamagandam,
		GulikaKalam:   gulikaKalam,
		AbhijitMuhurta: abhijitMuhurta,
	}
	
	span.AddEvent("Traditional periods calculated", trace.WithAttributes(
		attribute.String("rahu_kalam", rahuKalam.Start.Format("15:04:05")+" - "+rahuKalam.End.Format("15:04:05")),
		attribute.String("yamagandam", yamagandam.Start.Format("15:04:05")+" - "+yamagandam.End.Format("15:04:05")),
		attribute.String("gulika_kalam", gulikaKalam.Start.Format("15:04:05")+" - "+gulikaKalam.End.Format("15:04:05")),
		attribute.String("abhijit_muhurta", abhijitMuhurta.Start.Format("15:04:05")+" - "+abhijitMuhurta.End.Format("15:04:05")),
	))
	
	return result, nil
}

// calculateRahuKalam calculates Rahu Kalam based on the day of the week
// Rahu Kalam is considered inauspicious and varies by day of the week
func calculateRahuKalam(ctx context.Context, sunrise, sunset time.Time, date time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateRahuKalam")
	defer span.End()
	
	dayOfWeek := date.Weekday()
	span.SetAttributes(attribute.String("day_of_week", dayOfWeek.String()))
	
	// Calculate day length and divide into 8 parts
	dayLength := sunset.Sub(sunrise)
	partDuration := dayLength / 8
	
	span.SetAttributes(
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.Float64("part_duration_minutes", partDuration.Minutes()),
	)
	
	// Rahu Kalam timing based on day of the week (traditional calculation)
	var rahuPart int
	switch dayOfWeek {
	case time.Sunday:
		rahuPart = 7 // 7th part (4:30 PM - 6:00 PM traditionally)
	case time.Monday:
		rahuPart = 1 // 1st part (7:30 AM - 9:00 AM traditionally)
	case time.Tuesday:
		rahuPart = 6 // 6th part (3:00 PM - 4:30 PM traditionally)
	case time.Wednesday:
		rahuPart = 4 // 4th part (12:00 PM - 1:30 PM traditionally)
	case time.Thursday:
		rahuPart = 3 // 3rd part (10:30 AM - 12:00 PM traditionally)
	case time.Friday:
		rahuPart = 2 // 2nd part (9:00 AM - 10:30 AM traditionally)
	case time.Saturday:
		rahuPart = 5 // 5th part (1:30 PM - 3:00 PM traditionally)
	}
	
	// Calculate start and end times
	start := sunrise.Add(time.Duration(rahuPart-1) * partDuration)
	end := sunrise.Add(time.Duration(rahuPart) * partDuration)
	
	span.SetAttributes(
		attribute.Int("rahu_part", rahuPart),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)
	
	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(partDuration.Minutes()),
		Description: "Rahu Kalam - Inauspicious period ruled by Rahu",
		Auspicious:  false,
	}
	
	span.AddEvent("Rahu Kalam calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
	))
	
	return result
}

// calculateYamagandam calculates Yamagandam based on the day of the week
// Yamagandam is another inauspicious period
func calculateYamagandam(ctx context.Context, sunrise, sunset time.Time, date time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateYamagandam")
	defer span.End()
	
	dayOfWeek := date.Weekday()
	span.SetAttributes(attribute.String("day_of_week", dayOfWeek.String()))
	
	// Calculate day length and divide into 8 parts
	dayLength := sunset.Sub(sunrise)
	partDuration := dayLength / 8
	
	// Yamagandam timing based on day of the week
	var yamaPart int
	switch dayOfWeek {
	case time.Sunday:
		yamaPart = 4 // 4th part
	case time.Monday:
		yamaPart = 7 // 7th part
	case time.Tuesday:
		yamaPart = 2 // 2nd part
	case time.Wednesday:
		yamaPart = 5 // 5th part
	case time.Thursday:
		yamaPart = 8 // 8th part
	case time.Friday:
		yamaPart = 6 // 6th part
	case time.Saturday:
		yamaPart = 3 // 3rd part
	}
	
	// Calculate start and end times
	start := sunrise.Add(time.Duration(yamaPart-1) * partDuration)
	end := sunrise.Add(time.Duration(yamaPart) * partDuration)
	
	span.SetAttributes(
		attribute.Int("yama_part", yamaPart),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)
	
	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(partDuration.Minutes()),
		Description: "Yamagandam - Inauspicious period ruled by Yama",
		Auspicious:  false,
	}
	
	span.AddEvent("Yamagandam calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
	))
	
	return result
}

// calculateGulikaKalam calculates Gulika Kalam based on the day of the week
// Gulika Kalam is also considered inauspicious
func calculateGulikaKalam(ctx context.Context, sunrise, sunset time.Time, date time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateGulikaKalam")
	defer span.End()
	
	dayOfWeek := date.Weekday()
	span.SetAttributes(attribute.String("day_of_week", dayOfWeek.String()))
	
	// Calculate day length and divide into 8 parts
	dayLength := sunset.Sub(sunrise)
	partDuration := dayLength / 8
	
	// Gulika Kalam timing based on day of the week
	var gulikaPart int
	switch dayOfWeek {
	case time.Sunday:
		gulikaPart = 6 // 6th part
	case time.Monday:
		gulikaPart = 8 // 8th part
	case time.Tuesday:
		gulikaPart = 4 // 4th part
	case time.Wednesday:
		gulikaPart = 7 // 7th part
	case time.Thursday:
		gulikaPart = 2 // 2nd part
	case time.Friday:
		gulikaPart = 5 // 5th part
	case time.Saturday:
		gulikaPart = 1 // 1st part
	}
	
	// Calculate start and end times
	start := sunrise.Add(time.Duration(gulikaPart-1) * partDuration)
	end := sunrise.Add(time.Duration(gulikaPart) * partDuration)
	
	span.SetAttributes(
		attribute.Int("gulika_part", gulikaPart),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)
	
	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(partDuration.Minutes()),
		Description: "Gulika Kalam - Inauspicious period ruled by Gulika",
		Auspicious:  false,
	}
	
	span.AddEvent("Gulika Kalam calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
	))
	
	return result
}

// calculateAbhijitMuhurta calculates Abhijit Muhurta - the most auspicious time of the day
// Abhijit Muhurta is approximately the 8th muhurta of the day (around midday)
func calculateAbhijitMuhurta(ctx context.Context, sunrise, sunset time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateAbhijitMuhurta")
	defer span.End()
	
	// Calculate day length and divide into 30 muhurtas (each muhurta is 1/30 of day length)
	dayLength := sunset.Sub(sunrise)
	muhurtaDuration := dayLength / 30
	
	span.SetAttributes(
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.Float64("muhurta_duration_minutes", muhurtaDuration.Minutes()),
	)
	
	// Abhijit Muhurta is the 8th muhurta (7 * muhurta_duration after sunrise)
	// This is typically around midday and is considered very auspicious
	start := sunrise.Add(7 * muhurtaDuration)
	end := sunrise.Add(8 * muhurtaDuration)
	
	// However, according to some traditions, if the start time is after 12:30 PM,
	// Abhijit Muhurta is not considered valid for the day
	midday := time.Date(sunrise.Year(), sunrise.Month(), sunrise.Day(), 12, 30, 0, 0, sunrise.Location())
	if start.After(midday) {
		span.SetAttributes(attribute.Bool("abhijit_valid", false))
		span.AddEvent("Abhijit Muhurta not valid - starts after 12:30 PM")
		
		// Return a neutral period indicating it's not valid
		result := &TimePeriod{
			Start:       start,
			End:         end,
			Duration:    int(muhurtaDuration.Minutes()),
			Description: "Abhijit Muhurta - Not valid today (starts after 12:30 PM)",
			Auspicious:  false,
		}
		return result
	}
	
	span.SetAttributes(
		attribute.Bool("abhijit_valid", true),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)
	
	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(muhurtaDuration.Minutes()),
		Description: "Abhijit Muhurta - Most auspicious period of the day",
		Auspicious:  true,
	}
	
	span.AddEvent("Abhijit Muhurta calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
		attribute.Bool("is_valid", true),
	))
	
	return result
}

// GetRahuKalam returns just the Rahu Kalam period for a location and date
func GetRahuKalam(loc Location, date time.Time) (*TimePeriod, error) {
	return GetRahuKalamWithContext(context.Background(), loc, date)
}

// GetRahuKalamWithContext returns Rahu Kalam with tracing
func GetRahuKalamWithContext(ctx context.Context, loc Location, date time.Time) (*TimePeriod, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetRahuKalam")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)
	
	periods, err := CalculateTraditionalPeriodsWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	
	span.SetAttributes(
		attribute.String("rahu_kalam", periods.RahuKalam.Start.Format("15:04:05")+" - "+periods.RahuKalam.End.Format("15:04:05")),
	)
	span.AddEvent("Rahu Kalam extracted", trace.WithAttributes(
		attribute.String("rahu_kalam", periods.RahuKalam.Start.Format("15:04:05")+" - "+periods.RahuKalam.End.Format("15:04:05")),
	))
	
	return periods.RahuKalam, nil
}

// GetAbhijitMuhurta returns just the Abhijit Muhurta period for a location and date
func GetAbhijitMuhurta(loc Location, date time.Time) (*TimePeriod, error) {
	return GetAbhijitMuhurtaWithContext(context.Background(), loc, date)
}

// GetAbhijitMuhurtaWithContext returns Abhijit Muhurta with tracing
func GetAbhijitMuhurtaWithContext(ctx context.Context, loc Location, date time.Time) (*TimePeriod, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetAbhijitMuhurta")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)
	
	periods, err := CalculateTraditionalPeriodsWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	
	span.SetAttributes(
		attribute.String("abhijit_muhurta", periods.AbhijitMuhurta.Start.Format("15:04:05")+" - "+periods.AbhijitMuhurta.End.Format("15:04:05")),
		attribute.Bool("abhijit_auspicious", periods.AbhijitMuhurta.Auspicious),
	)
	span.AddEvent("Abhijit Muhurta extracted", trace.WithAttributes(
		attribute.String("abhijit_muhurta", periods.AbhijitMuhurta.Start.Format("15:04:05")+" - "+periods.AbhijitMuhurta.End.Format("15:04:05")),
	))
	
	return periods.AbhijitMuhurta, nil
}

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