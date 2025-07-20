package astronomy

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// VaraInfo represents a Vara (weekday) with its properties
type VaraInfo struct {
	Number          int       `json:"number"`           // 1-7 (Sunday=1, Monday=2, etc.)
	Name            string    `json:"name"`             // Sanskrit name
	PlanetaryLord   string    `json:"planetary_lord"`   // Ruling planet
	Quality         string    `json:"quality"`          // General quality/nature
	Color           string    `json:"color"`            // Associated color
	Deity           string    `json:"deity"`            // Presiding deity
	StartTime       time.Time `json:"start_time"`       // Sunrise time when Vara begins
	EndTime         time.Time `json:"end_time"`         // Next sunrise when Vara ends
	Duration        float64   `json:"duration"`         // Duration in hours
	GregorianDay    string    `json:"gregorian_day"`    // English weekday name
	IsAuspicious    bool      `json:"is_auspicious"`    // General auspiciousness
	CurrentHora     int       `json:"current_hora"`     // Current hora (1-24)
	HoraPlanet      string    `json:"hora_planet"`      // Planet ruling current hora
}

// VaraCalculator handles Vara calculations
type VaraCalculator struct {
	observer observability.ObserverInterface
}

// NewVaraCalculator creates a new VaraCalculator
func NewVaraCalculator() *VaraCalculator {
	return &VaraCalculator{
		observer: observability.Observer(),
	}
}

// VaraData contains detailed information about each Vara
// Sources:
// - "Brihat Parashara Hora Shastra" by Sage Parashara
// - "Muhurta Chintamani" by Daivagya Ramachandra
// - "Hindu Astronomy" by W.E. van Wijk (1930)
// - "Surya Siddhanta" - Ancient Sanskrit astronomical text
var VaraData = map[int]struct {
	Name          string
	PlanetaryLord string
	Quality       string
	Color         string
	Deity         string
	GregorianDay  string
	IsAuspicious  bool
}{
	1: {"Ravivara", "Sun", "Fierce and authoritative", "Red", "Surya", "Sunday", true},
	2: {"Somavara", "Moon", "Gentle and nurturing", "White", "Chandra", "Monday", true},
	3: {"Mangalavara", "Mars", "Energetic and aggressive", "Red", "Mangala", "Tuesday", false},
	4: {"Budhavara", "Mercury", "Intellectual and communicative", "Green", "Budha", "Wednesday", true},
	5: {"Guruvara", "Jupiter", "Wise and benevolent", "Yellow", "Brihaspati", "Thursday", true},
	6: {"Shukravara", "Venus", "Artistic and luxurious", "White", "Shukra", "Friday", true},
	7: {"Shanivara", "Saturn", "Disciplined and restrictive", "Black", "Shani", "Saturday", false},
}

// HoraPlanets defines the sequence of planets ruling each hora
// Starting with the day's planetary lord, then following the traditional sequence
var HoraPlanets = []string{"Sun", "Venus", "Mercury", "Moon", "Saturn", "Jupiter", "Mars"}

// GetVaraForDate calculates the Vara for a given date and location
func (vc *VaraCalculator) GetVaraForDate(ctx context.Context, date time.Time, location Location) (*VaraInfo, error) {
	ctx, span := vc.observer.CreateSpan(ctx, "VaraCalculator.GetVaraForDate")
	defer span.End()

	span.SetAttributes(
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("timezone", date.Location().String()),
		attribute.Float64("location.latitude", location.Latitude),
		attribute.Float64("location.longitude", location.Longitude),
	)

	// Calculate sunrise times for current and next day
	ctx, sunriseSpan := vc.observer.CreateSpan(ctx, "calculateSunriseTimes")
	
	// Current day sunrise
	currentSunTimes, err := CalculateSunTimesWithContext(ctx, location, date)
	if err != nil {
		sunriseSpan.RecordError(err)
		sunriseSpan.End()
		span.RecordError(err)
		return nil, fmt.Errorf("failed to calculate current day sunrise: %w", err)
	}

	// Next day sunrise
	nextDay := date.Add(24 * time.Hour)
	nextSunTimes, err := CalculateSunTimesWithContext(ctx, location, nextDay)
	if err != nil {
		sunriseSpan.RecordError(err)
		sunriseSpan.End()
		span.RecordError(err)
		return nil, fmt.Errorf("failed to calculate next day sunrise: %w", err)
	}

	sunriseSpan.SetAttributes(
		attribute.String("current_sunrise", currentSunTimes.Sunrise.Format("15:04:05")),
		attribute.String("next_sunrise", nextSunTimes.Sunrise.Format("15:04:05")),
	)
	sunriseSpan.End()

	// Determine the Vara based on sunrise
	vara, err := vc.calculateVaraFromSunrise(ctx, currentSunTimes.Sunrise, nextSunTimes.Sunrise, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("vara_number", vara.Number),
		attribute.String("vara_name", vara.Name),
		attribute.String("planetary_lord", vara.PlanetaryLord),
		attribute.String("gregorian_day", vara.GregorianDay),
		attribute.Bool("is_auspicious", vara.IsAuspicious),
		attribute.Int("current_hora", vara.CurrentHora),
		attribute.String("hora_planet", vara.HoraPlanet),
	)

	span.AddEvent("Vara calculated", trace.WithAttributes(
		attribute.Int("vara_number", vara.Number),
		attribute.String("vara_name", vara.Name),
		attribute.String("planetary_lord", vara.PlanetaryLord),
	))

	return vara, nil
}

// calculateVaraFromSunrise calculates Vara based on sunrise times
func (vc *VaraCalculator) calculateVaraFromSunrise(ctx context.Context, currentSunrise, nextSunrise time.Time, referenceDate time.Time) (*VaraInfo, error) {
	ctx, span := vc.observer.CreateSpan(ctx, "VaraCalculator.calculateVaraFromSunrise")
	defer span.End()

	span.SetAttributes(
		attribute.String("current_sunrise", currentSunrise.Format("2006-01-02 15:04:05")),
		attribute.String("next_sunrise", nextSunrise.Format("2006-01-02 15:04:05")),
		attribute.String("reference_date", referenceDate.Format("2006-01-02")),
	)

	// In Hindu calendar, the day changes at sunrise, not midnight
	// So we need to determine which Vara (weekday) is active based on sunrise
	
	// Get the Gregorian weekday for the sunrise date
	// Note: We use the sunrise date to determine the Vara
	sunriseDate := currentSunrise
	
	// Calculate Vara number (1-7, Sunday=1)
	// Go's time.Weekday() returns Sunday=0, Monday=1, etc.
	// We convert to traditional Hindu numbering where Sunday=1
	gregorianWeekday := sunriseDate.Weekday()
	varaNumber := int(gregorianWeekday) + 1
	if varaNumber > 7 {
		varaNumber = 1
	}

	span.SetAttributes(
		attribute.Int("gregorian_weekday", int(gregorianWeekday)),
		attribute.Int("vara_number", varaNumber),
	)

	// Get Vara details
	varaDetails := VaraData[varaNumber]

	// Calculate current hora and its ruling planet
	currentHora, horaPlanet := vc.calculateCurrentHora(ctx, currentSunrise, nextSunrise, referenceDate, varaNumber)

	span.SetAttributes(
		attribute.String("vara_name", varaDetails.Name),
		attribute.String("planetary_lord", varaDetails.PlanetaryLord),
		attribute.String("quality", varaDetails.Quality),
		attribute.String("color", varaDetails.Color),
		attribute.String("deity", varaDetails.Deity),
		attribute.String("gregorian_day", varaDetails.GregorianDay),
		attribute.Bool("is_auspicious", varaDetails.IsAuspicious),
		attribute.Int("current_hora", currentHora),
		attribute.String("hora_planet", horaPlanet),
	)

	duration := nextSunrise.Sub(currentSunrise).Hours()

	vara := &VaraInfo{
		Number:          varaNumber,
		Name:            varaDetails.Name,
		PlanetaryLord:   varaDetails.PlanetaryLord,
		Quality:         varaDetails.Quality,
		Color:           varaDetails.Color,
		Deity:           varaDetails.Deity,
		StartTime:       currentSunrise,
		EndTime:         nextSunrise,
		Duration:        duration,
		GregorianDay:    varaDetails.GregorianDay,
		IsAuspicious:    varaDetails.IsAuspicious,
		CurrentHora:     currentHora,
		HoraPlanet:      horaPlanet,
	}

	span.AddEvent("Vara calculation completed", trace.WithAttributes(
		attribute.Int("vara_number", varaNumber),
		attribute.String("vara_name", varaDetails.Name),
		attribute.Float64("duration_hours", duration),
		attribute.Int("current_hora", currentHora),
	))

	return vara, nil
}

// calculateCurrentHora calculates the current hora and its ruling planet
func (vc *VaraCalculator) calculateCurrentHora(ctx context.Context, currentSunrise, nextSunrise time.Time, referenceTime time.Time, varaNumber int) (int, string) {
	ctx, span := vc.observer.CreateSpan(ctx, "VaraCalculator.calculateCurrentHora")
	defer span.End()

	span.SetAttributes(
		attribute.String("current_sunrise", currentSunrise.Format("15:04:05")),
		attribute.String("next_sunrise", nextSunrise.Format("15:04:05")),
		attribute.String("reference_time", referenceTime.Format("15:04:05")),
		attribute.Int("vara_number", varaNumber),
	)

	// Calculate the total daylight duration
	totalDuration := nextSunrise.Sub(currentSunrise)
	
	// Each day is divided into 24 horas (planetary hours)
	// Each hora is 1/24 of the total day duration
	horaDuration := totalDuration / 24

	// Find which hora we're currently in
	var timeFromSunrise time.Duration
	if referenceTime.After(currentSunrise) && referenceTime.Before(nextSunrise) {
		timeFromSunrise = referenceTime.Sub(currentSunrise)
	} else if referenceTime.Before(currentSunrise) {
		// If before sunrise, we're in the previous day's cycle
		timeFromSunrise = 0
	} else {
		// If after next sunrise, we're in the next day's cycle
		timeFromSunrise = totalDuration
	}

	// Calculate which hora (1-24)
	horaNumber := int(timeFromSunrise/horaDuration) + 1
	if horaNumber > 24 {
		horaNumber = 24
	}
	if horaNumber < 1 {
		horaNumber = 1
	}

	// Calculate the ruling planet for this hora
	// The first hora of each day is ruled by the day's planetary lord
	// Then follows the traditional planetary sequence
	dayPlanetIndex := getPlanetIndex(VaraData[varaNumber].PlanetaryLord)
	horaPlanetIndex := (dayPlanetIndex + horaNumber - 1) % 7
	horaPlanet := HoraPlanets[horaPlanetIndex]

	span.SetAttributes(
		attribute.Float64("total_duration_hours", totalDuration.Hours()),
		attribute.Float64("hora_duration_hours", horaDuration.Hours()),
		attribute.Float64("time_from_sunrise_hours", timeFromSunrise.Hours()),
		attribute.Int("hora_number", horaNumber),
		attribute.Int("day_planet_index", dayPlanetIndex),
		attribute.Int("hora_planet_index", horaPlanetIndex),
		attribute.String("hora_planet", horaPlanet),
	)

	span.AddEvent("Hora calculated", trace.WithAttributes(
		attribute.Int("hora_number", horaNumber),
		attribute.String("hora_planet", horaPlanet),
	))

	return horaNumber, horaPlanet
}

// getPlanetIndex returns the index of a planet in the hora sequence
func getPlanetIndex(planet string) int {
	for i, p := range HoraPlanets {
		if p == planet {
			return i
		}
	}
	return 0 // Default to Sun if not found
}

// GetVaraFromGregorianDay is a convenience function to get Vara from Gregorian weekday
func (vc *VaraCalculator) GetVaraFromGregorianDay(ctx context.Context, gregorianDay time.Weekday, sunrise, nextSunrise time.Time, referenceTime time.Time) (*VaraInfo, error) {
	ctx, span := vc.observer.CreateSpan(ctx, "VaraCalculator.GetVaraFromGregorianDay")
	defer span.End()

	span.SetAttributes(
		attribute.Int("gregorian_day", int(gregorianDay)),
		attribute.String("sunrise", sunrise.Format("15:04:05")),
		attribute.String("next_sunrise", nextSunrise.Format("15:04:05")),
	)

	// Convert Gregorian weekday to Vara number
	varaNumber := int(gregorianDay) + 1
	if varaNumber > 7 {
		varaNumber = 1
	}

	return vc.calculateVaraFromSunrise(ctx, sunrise, nextSunrise, referenceTime)
}

// GetHoraForTime calculates the hora for a specific time within a Vara
func (vc *VaraCalculator) GetHoraForTime(ctx context.Context, specificTime time.Time, currentSunrise, nextSunrise time.Time, varaNumber int) (int, string, error) {
	ctx, span := vc.observer.CreateSpan(ctx, "VaraCalculator.GetHoraForTime")
	defer span.End()

	span.SetAttributes(
		attribute.String("specific_time", specificTime.Format("15:04:05")),
		attribute.String("current_sunrise", currentSunrise.Format("15:04:05")),
		attribute.String("next_sunrise", nextSunrise.Format("15:04:05")),
		attribute.Int("vara_number", varaNumber),
	)

	if varaNumber < 1 || varaNumber > 7 {
		return 0, "", fmt.Errorf("invalid vara number: %d, must be between 1 and 7", varaNumber)
	}

	horaNumber, horaPlanet := vc.calculateCurrentHora(ctx, currentSunrise, nextSunrise, specificTime, varaNumber)

	span.SetAttributes(
		attribute.Int("hora_number", horaNumber),
		attribute.String("hora_planet", horaPlanet),
	)

	return horaNumber, horaPlanet, nil
}

// IsAuspiciousVara returns true if the Vara is generally considered auspicious
func IsAuspiciousVara(vara *VaraInfo) bool {
	return vara.IsAuspicious
}

// GetVaraRecommendations returns recommendations based on the Vara
func GetVaraRecommendations(vara *VaraInfo) string {
	switch vara.Name {
	case "Ravivara": // Sunday
		return "Good for spiritual practices, government work, and leadership activities. Avoid starting new ventures."
	case "Somavara": // Monday
		return "Excellent for new beginnings, travel, and emotional healing. Good for all auspicious activities."
	case "Mangalavara": // Tuesday
		return "Avoid important activities. Not favorable for marriages, new ventures, or peaceful activities."
	case "Budhavara": // Wednesday
		return "Good for education, communication, business, and intellectual pursuits."
	case "Guruvara": // Thursday
		return "Most auspicious day. Excellent for all important activities, ceremonies, and new beginnings."
	case "Shukravara": // Friday
		return "Good for artistic pursuits, relationships, luxury items, and social activities."
	case "Shanivara": // Saturday
		return "Avoid important activities. Good for discipline, hard work, and dealing with obstacles."
	default:
		return "General vara with moderate influence."
	}
}

// GetHoraPlanetRecommendations returns recommendations based on the current hora planet
func GetHoraPlanetRecommendations(planet string) string {
	switch planet {
	case "Sun":
		return "Good for government work, leadership, and spiritual practices."
	case "Moon":
		return "Favorable for emotional matters, travel, and water-related activities."
	case "Mars":
		return "Good for physical activities, sports, and dealing with conflicts. Avoid peace negotiations."
	case "Mercury":
		return "Excellent for communication, education, business, and intellectual work."
	case "Jupiter":
		return "Most auspicious. Good for all activities, especially religious and educational."
	case "Venus":
		return "Good for artistic work, relationships, luxury, and social activities."
	case "Saturn":
		return "Good for discipline, hard work, and routine tasks. Avoid festivities."
	default:
		return "General planetary influence."
	}
}

// ValidateVaraCalculation validates a Vara calculation result
func ValidateVaraCalculation(vara *VaraInfo) error {
	if vara == nil {
		return fmt.Errorf("vara cannot be nil")
	}

	if vara.Number < 1 || vara.Number > 7 {
		return fmt.Errorf("invalid vara number: %d, must be between 1 and 7", vara.Number)
	}

	if vara.CurrentHora < 1 || vara.CurrentHora > 24 {
		return fmt.Errorf("invalid hora number: %d, must be between 1 and 24", vara.CurrentHora)
	}

	if vara.Duration <= 0 || vara.Duration > 30 {
		return fmt.Errorf("invalid vara duration: %f hours, must be positive and reasonable", vara.Duration)
	}

	if vara.EndTime.Before(vara.StartTime) {
		return fmt.Errorf("vara end time cannot be before start time")
	}

	if vara.Name == "" {
		return fmt.Errorf("vara name cannot be empty")
	}

	if vara.PlanetaryLord == "" {
		return fmt.Errorf("planetary lord cannot be empty")
	}

	if vara.HoraPlanet == "" {
		return fmt.Errorf("hora planet cannot be empty")
	}

	// Validate hora planet is one of the expected values
	validPlanet := false
	for _, planet := range HoraPlanets {
		if vara.HoraPlanet == planet {
			validPlanet = true
			break
		}
	}
	if !validPlanet {
		return fmt.Errorf("invalid hora planet: %s", vara.HoraPlanet)
	}

	return nil
}