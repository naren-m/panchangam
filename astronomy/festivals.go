package astronomy

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Festival represents a Hindu festival with its details
type Festival struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`        // "major", "minor", "regional"
	Significance string   `json:"significance"`
	Observances []string  `json:"observances"`
}

// FestivalCalendar contains festival detection logic
type FestivalCalendar struct {
	// Fixed date festivals (Gregorian calendar based)
	fixedFestivals map[string]Festival
	// Lunar festivals (based on Tithi)
	lunarFestivals map[int][]Festival
}

// NewFestivalCalendar creates a new festival calendar
func NewFestivalCalendar() *FestivalCalendar {
	fc := &FestivalCalendar{
		fixedFestivals: make(map[string]Festival),
		lunarFestivals: make(map[int][]Festival),
	}
	
	fc.initializeFixedFestivals()
	fc.initializeLunarFestivals()
	
	return fc
}

// initializeFixedFestivals sets up Gregorian calendar-based festivals
func (fc *FestivalCalendar) initializeFixedFestivals() {
	// Major fixed festivals
	fc.fixedFestivals["01-26"] = Festival{
		Name:        "Republic Day",
		Type:        "national",
		Significance: "India's Constitution came into effect",
		Observances: []string{"Flag hoisting", "Parades", "Cultural programs"},
	}
	
	fc.fixedFestivals["08-15"] = Festival{
		Name:        "Independence Day",
		Type:        "national",
		Significance: "India's independence from British rule",
		Observances: []string{"Flag hoisting", "Patriotic ceremonies"},
	}
	
	fc.fixedFestivals["10-02"] = Festival{
		Name:        "Gandhi Jayanti",
		Type:        "national",
		Significance: "Birthday of Mahatma Gandhi",
		Observances: []string{"Prayer meetings", "Spinning wheel ceremonies"},
	}
}

// initializeLunarFestivals sets up Tithi-based festivals
func (fc *FestivalCalendar) initializeLunarFestivals() {
	// Ekadashi (11th lunar day) - occurs twice per month
	ekadashi := Festival{
		Name:        "Ekadashi",
		Type:        "major",
		Significance: "Sacred to Lord Vishnu, fasting day",
		Observances: []string{"Fasting", "Prayer", "Meditation", "Charity"},
	}
	fc.lunarFestivals[11] = append(fc.lunarFestivals[11], ekadashi)
	
	// Amavasya (New Moon - 30th Tithi)
	amavasya := Festival{
		Name:        "Amavasya",
		Type:        "minor",
		Significance: "New moon day, ancestor worship",
		Observances: []string{"Ancestral prayers", "Charity", "Meditation"},
	}
	fc.lunarFestivals[30] = append(fc.lunarFestivals[30], amavasya)
	
	// Purnima (Full Moon - 15th Tithi)
	purnima := Festival{
		Name:        "Purnima",
		Type:        "minor",
		Significance: "Full moon day, auspicious for prayers",
		Observances: []string{"Prayers", "Meditation", "Charity", "Fasting"},
	}
	fc.lunarFestivals[15] = append(fc.lunarFestivals[15], purnima)
	
	// Chaturthi (4th Tithi) - Ganesh Chaturthi varies by month
	chaturthi := Festival{
		Name:        "Chaturthi",
		Type:        "minor",
		Significance: "Sacred to Lord Ganesha",
		Observances: []string{"Ganesha prayers", "Offerings", "Modak preparation"},
	}
	fc.lunarFestivals[4] = append(fc.lunarFestivals[4], chaturthi)
	
	// Navami (9th Tithi) - Sacred to Devi
	navami := Festival{
		Name:        "Navami",
		Type:        "minor",
		Significance: "Sacred to Divine Mother",
		Observances: []string{"Devi prayers", "Fasting", "Scripture reading"},
	}
	fc.lunarFestivals[9] = append(fc.lunarFestivals[9], navami)
}

// GetFestivalsForDate returns festivals for a specific date
func (fc *FestivalCalendar) GetFestivalsForDate(ctx context.Context, date time.Time, tithiNumber int) ([]Festival, error) {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "GetFestivalsForDate")
	defer span.End()
	
	span.SetAttributes(
		attribute.String("date", date.Format("2006-01-02")),
		attribute.Int("tithi_number", tithiNumber),
	)
	
	var festivals []Festival
	
	// Check fixed festivals (Gregorian calendar)
	monthDay := date.Format("01-02")
	if festival, exists := fc.fixedFestivals[monthDay]; exists {
		festival.Date = date
		festivals = append(festivals, festival)
		span.AddEvent("Fixed festival found", trace.WithAttributes(
			attribute.String("festival_name", festival.Name),
			attribute.String("festival_type", festival.Type),
		))
	}
	
	// Check lunar festivals (Tithi-based)
	if lunarFestivals, exists := fc.lunarFestivals[tithiNumber]; exists {
		for _, festival := range lunarFestivals {
			festival.Date = date
			
			// Add month-specific naming for certain festivals
			festival.Name = fc.getMonthSpecificName(festival.Name, date, tithiNumber)
			festivals = append(festivals, festival)
			
			span.AddEvent("Lunar festival found", trace.WithAttributes(
				attribute.String("festival_name", festival.Name),
				attribute.String("festival_type", festival.Type),
				attribute.Int("tithi", tithiNumber),
			))
		}
	}
	
	// Add seasonal festivals based on month
	seasonalFestivals := fc.getSeasonalFestivals(date)
	festivals = append(festivals, seasonalFestivals...)
	
	span.SetAttributes(
		attribute.Int("total_festivals", len(festivals)),
	)
	
	return festivals, nil
}

// getMonthSpecificName returns month-specific festival names
func (fc *FestivalCalendar) getMonthSpecificName(baseName string, date time.Time, tithiNumber int) string {
	month := date.Month()
	
	switch baseName {
	case "Ekadashi":
		// Different Ekadashi names based on month
		ekadashiNames := map[time.Month]string{
			time.January:   "Pausha Putrada Ekadashi",
			time.February:  "Magha Shattila Ekadashi", 
			time.March:     "Phalguna Vijaya Ekadashi",
			time.April:     "Chaitra Kamada Ekadashi",
			time.May:       "Vaishakha Mohini Ekadashi",
			time.June:      "Jyeshtha Nirjala Ekadashi",
			time.July:      "Ashadha Yogini Ekadashi",
			time.August:    "Shravana Kamika Ekadashi",
			time.September: "Bhadrapada Aja Ekadashi",
			time.October:   "Ashwin Indira Ekadashi",
			time.November:  "Kartik Rama Ekadashi",
			time.December:  "Margashirsha Mokshada Ekadashi",
		}
		if name, exists := ekadashiNames[month]; exists {
			return name
		}
		
	case "Purnima":
		// Different Purnima names based on month
		purnimaNames := map[time.Month]string{
			time.January:   "Pausha Purnima",
			time.February:  "Magha Purnima",
			time.March:     "Holi Purnima",
			time.April:     "Chaitra Purnima",
			time.May:       "Buddha Purnima",
			time.June:      "Vat Purnima",
			time.July:      "Guru Purnima",
			time.August:    "Raksha Bandhan",
			time.September: "Bhadrapada Purnima",
			time.October:   "Sharad Purnima",
			time.November:  "Kartik Purnima",
			time.December:  "Margashirsha Purnima",
		}
		if name, exists := purnimaNames[month]; exists {
			return name
		}
		
	case "Amavasya":
		// Different Amavasya names based on month
		amavasyas := map[time.Month]string{
			time.October:  "Diwali Amavasya",
			time.November: "Kartik Amavasya",
		}
		if name, exists := amavasyas[month]; exists {
			return name
		}
	}
	
	return baseName
}

// getSeasonalFestivals returns seasonal festivals for specific months
func (fc *FestivalCalendar) getSeasonalFestivals(date time.Time) []Festival {
	var festivals []Festival
	month := date.Month()
	
	switch month {
	case time.March:
		if date.Day() >= 20 && date.Day() <= 22 {
			festivals = append(festivals, Festival{
				Name:        "Spring Equinox",
				Date:        date,
				Type:        "seasonal",
				Significance: "Beginning of spring season",
				Observances: []string{"Nature worship", "Spring cleaning", "New plantings"},
			})
		}
		
	case time.June:
		if date.Day() >= 20 && date.Day() <= 22 {
			festivals = append(festivals, Festival{
				Name:        "Summer Solstice",
				Date:        date,
				Type:        "seasonal", 
				Significance: "Longest day of the year",
				Observances: []string{"Sun worship", "Early morning prayers"},
			})
		}
		
	case time.September:
		if date.Day() >= 22 && date.Day() <= 24 {
			festivals = append(festivals, Festival{
				Name:        "Autumn Equinox",
				Date:        date,
				Type:        "seasonal",
				Significance: "Beginning of autumn season",
				Observances: []string{"Harvest celebrations", "Ancestor worship"},
			})
		}
		
	case time.December:
		if date.Day() >= 20 && date.Day() <= 22 {
			festivals = append(festivals, Festival{
				Name:        "Winter Solstice",
				Date:        date,
				Type:        "seasonal",
				Significance: "Longest night of the year",
				Observances: []string{"Light festivals", "Fire rituals"},
			})
		}
	}
	
	return festivals
}

// GetUpcomingFestivals returns festivals in the next N days
func (fc *FestivalCalendar) GetUpcomingFestivals(ctx context.Context, startDate time.Time, days int) ([]Festival, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetUpcomingFestivals")
	defer span.End()
	
	span.SetAttributes(
		attribute.String("start_date", startDate.Format("2006-01-02")),
		attribute.Int("days", days),
	)
	
	var allFestivals []Festival
	
	for i := 0; i < days; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		
		// For this basic implementation, we'll use a simplified tithi calculation
		// In a full implementation, this would use the actual Tithi calculator
		dayOfMonth := currentDate.Day()
		approximateTithi := (dayOfMonth % 30) + 1
		
		festivals, err := fc.GetFestivalsForDate(ctx, currentDate, approximateTithi)
		if err != nil {
			span.RecordError(err)
			continue
		}
		
		allFestivals = append(allFestivals, festivals...)
	}
	
	span.SetAttributes(
		attribute.Int("total_upcoming_festivals", len(allFestivals)),
	)
	
	return allFestivals, nil
}

// Helper function to get festival names as strings (for API compatibility)
func GetFestivalNamesForDate(ctx context.Context, date time.Time, tithiNumber int) ([]string, error) {
	fc := NewFestivalCalendar()
	festivals, err := fc.GetFestivalsForDate(ctx, date, tithiNumber)
	if err != nil {
		return nil, err
	}
	
	var names []string
	for _, festival := range festivals {
		names = append(names, festival.Name)
	}
	
	return names, nil
}