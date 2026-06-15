package astronomy

import "fmt"

// GetTithiTypeDescription returns a description of the Tithi type
func GetTithiTypeDescription(tithiType TithiType) string {
	switch tithiType {
	case TithiTypeNanda:
		return "Joyful, good for celebrations and new beginnings"
	case TithiTypeBhadra:
		return "Auspicious, good for all activities"
	case TithiTypeJaya:
		return "Victorious, good for achieving success"
	case TithiTypeRikta:
		return "Empty, avoid starting new ventures"
	case TithiTypePurna:
		return "Complete, excellent for completion of tasks"
	default:
		return "Unknown Tithi type"
	}
}

// ValidateTithiCalculation validates a Tithi calculation result
func ValidateTithiCalculation(tithi *TithiInfo) error {
	if tithi == nil {
		return fmt.Errorf("tithi cannot be nil")
	}

	if tithi.Number < 1 || tithi.Number > 30 {
		return fmt.Errorf("invalid tithi number: %d, must be between 1 and 30", tithi.Number)
	}

	if tithi.MoonSunDiff < 0 || tithi.MoonSunDiff >= 360 {
		return fmt.Errorf("invalid moon-sun difference: %f, must be between 0 and 360 degrees", tithi.MoonSunDiff)
	}

	if tithi.Duration <= 0 || tithi.Duration > 48 {
		return fmt.Errorf("invalid tithi duration: %f hours, must be positive and reasonable", tithi.Duration)
	}

	if tithi.EndTime.Before(tithi.StartTime) {
		return fmt.Errorf("tithi end time cannot be before start time")
	}

	if tithi.PakshaDay != 0 && (tithi.PakshaDay < 1 || tithi.PakshaDay > 15) {
		return fmt.Errorf("invalid paksha day: %d, must be between 1 and 15", tithi.PakshaDay)
	}

	if tithi.Paksha != "" && tithi.Paksha != "Shukla" && tithi.Paksha != "Krishna" {
		return fmt.Errorf("invalid paksha: %s, must be Shukla or Krishna", tithi.Paksha)
	}

	if tithi.CalendarSystem != "" {
		calendarSystem := normalizeTithiCalendarSystem(tithi.CalendarSystem)
		if calendarSystem != "Purnimanta" && calendarSystem != "Amanta" {
			return fmt.Errorf("invalid calendar system: %s, must be Purnimanta or Amanta", tithi.CalendarSystem)
		}
	}

	if tithi.Name == "" || (tithi.CalendarSystem != "" && tithi.TraditionalName == "") {
		return fmt.Errorf("tithi names cannot be empty")
	}

	return nil
}
