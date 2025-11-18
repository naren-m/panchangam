package implementations

import (
	"fmt"
	"time"
)

// NaazhikaiConverter handles Tamil traditional time unit conversions
// Naazhikai (நாழிகை) is a traditional Tamil time unit
// 1 Naazhikai = 24 minutes
// 60 Naazhikai = 1 day (24 hours)
type NaazhikaiConverter struct{}

// TimeInNaazhikai represents time in Naazhikai units
type TimeInNaazhikai struct {
	Naazhikai int     // Whole naazhikai units (0-59)
	Vinaazhikai int   // Sub-units (60 vinaazhikai = 1 naazhikai, each ~24 seconds)
	Percentage float64 // Percentage of current naazhikai elapsed
}

// NewNaazhikaiConverter creates a new converter
func NewNaazhikaiConverter() *NaazhikaiConverter {
	return &NaazhikaiConverter{}
}

// ToNaazhikai converts a time.Time to Naazhikai units from sunrise
func (n *NaazhikaiConverter) ToNaazhikai(currentTime, sunrise time.Time) TimeInNaazhikai {
	// Calculate time elapsed since sunrise
	elapsed := currentTime.Sub(sunrise)

	// Convert to total minutes
	totalMinutes := elapsed.Minutes()

	// Calculate Naazhikai (1 Naazhikai = 24 minutes)
	totalNaazhikai := totalMinutes / 24.0
	naazhikai := int(totalNaazhikai)

	// Calculate Vinaazhikai (60 Vinaazhikai = 1 Naazhikai, so 1 Vinaazhikai = 24 seconds)
	fractionalPart := totalNaazhikai - float64(naazhikai)
	vinaazhikai := int(fractionalPart * 60)

	// Percentage of current naazhikai
	percentage := (fractionalPart * 60 - float64(vinaazhikai)) * 100

	return TimeInNaazhikai{
		Naazhikai:   naazhikai,
		Vinaazhikai: vinaazhikai,
		Percentage:  percentage,
	}
}

// FromNaazhikai converts Naazhikai units to duration from sunrise
func (n *NaazhikaiConverter) FromNaazhikai(naazhikai, vinaazhikai int) time.Duration {
	// Convert Naazhikai to minutes (1 Naazhikai = 24 minutes)
	totalMinutes := float64(naazhikai) * 24.0

	// Add Vinaazhikai (1 Vinaazhikai = 24 seconds = 0.4 minutes)
	totalMinutes += float64(vinaazhikai) * 0.4

	// Convert to duration
	return time.Duration(totalMinutes * float64(time.Minute))
}

// GetNaazhikaiTime returns the time in Naazhikai from sunrise
func (n *NaazhikaiConverter) GetNaazhikaiTime(currentTime, sunrise time.Time) string {
	tn := n.ToNaazhikai(currentTime, sunrise)
	return fmt.Sprintf("%d Naazhikai %d Vinaazhikai", tn.Naazhikai, tn.Vinaazhikai)
}

// GetNaazhikaiDetails returns detailed Naazhikai information
func (n *NaazhikaiConverter) GetNaazhikaiDetails(currentTime, sunrise, sunset time.Time) map[string]interface{} {
	tn := n.ToNaazhikai(currentTime, sunrise)

	// Calculate day length in Naazhikai
	dayLength := sunset.Sub(sunrise)
	dayNaazhikai := dayLength.Minutes() / 24.0

	return map[string]interface{}{
		"naazhikai":          tn.Naazhikai,
		"vinaazhikai":        tn.Vinaazhikai,
		"percentage":         tn.Percentage,
		"total_day_naazhikai": dayNaazhikai,
		"tamil_name":         "நாழிகை",
		"description":        "Traditional Tamil time unit (1 Naazhikai = 24 minutes)",
		"system":             "Tamil time measurement",
	}
}

// ConvertDurationToNaazhikai converts a time.Duration to Naazhikai units
func (n *NaazhikaiConverter) ConvertDurationToNaazhikai(duration time.Duration) TimeInNaazhikai {
	totalMinutes := duration.Minutes()
	totalNaazhikai := totalMinutes / 24.0
	naazhikai := int(totalNaazhikai)
	fractionalPart := totalNaazhikai - float64(naazhikai)
	vinaazhikai := int(fractionalPart * 60)
	percentage := (fractionalPart * 60 - float64(vinaazhikai)) * 100

	return TimeInNaazhikai{
		Naazhikai:   naazhikai,
		Vinaazhikai: vinaazhikai,
		Percentage:  percentage,
	}
}

// GetNaazhikaiPeriodName returns the traditional Tamil name for time of day
func (n *NaazhikaiConverter) GetNaazhikaiPeriodName(naazhikai int) string {
	// Traditional Tamil time periods based on Naazhikai
	switch {
	case naazhikai >= 0 && naazhikai < 5:
		return "காலை (Morning - Kaalai)"
	case naazhikai >= 5 && naazhikai < 15:
		return "முற்பகல் (Forenoon - Murpagal)"
	case naazhikai >= 15 && naazhikai < 20:
		return "மதியம் (Noon - Madhiyam)"
	case naazhikai >= 20 && naazhikai < 30:
		return "பிற்பகல் (Afternoon - Pirpagal)"
	case naazhikai >= 30 && naazhikai < 40:
		return "மாலை (Evening - Maalai)"
	case naazhikai >= 40 && naazhikai < 60:
		return "இரவு (Night - Iravu)"
	default:
		return "Unknown"
	}
}

// CalculateMuhurtaInNaazhikai calculates muhurta timing in Naazhikai units
func (n *NaazhikaiConverter) CalculateMuhurtaInNaazhikai(sunrise, sunset time.Time) []NaazhikaiMuhurta {
	var muhurtas []NaazhikaiMuhurta

	// Day is divided into specific Naazhikai periods for muhurtas
	dayLength := sunset.Sub(sunrise)
	naazhikaiPerMuhurta := 2 // Traditional 2 Naazhikai per muhurta

	totalDayNaazhikai := int(dayLength.Minutes() / 24.0)
	muhurtaCount := totalDayNaazhikai / naazhikaiPerMuhurta

	for i := 0; i < muhurtaCount; i++ {
		startNaazhikai := i * naazhikaiPerMuhurta
		endNaazhikai := startNaazhikai + naazhikaiPerMuhurta

		startTime := sunrise.Add(n.FromNaazhikai(startNaazhikai, 0))
		endTime := sunrise.Add(n.FromNaazhikai(endNaazhikai, 0))

		quality := n.getMuhurtaQuality(startNaazhikai)

		muhurtas = append(muhurtas, NaazhikaiMuhurta{
			MuhurtaNumber:   i + 1,
			StartNaazhikai:  startNaazhikai,
			EndNaazhikai:    endNaazhikai,
			StartTime:       startTime,
			EndTime:         endTime,
			Quality:         quality,
			TamilName:       n.getNaazhikaiMuhurtaName(i + 1),
		})
	}

	return muhurtas
}

// NaazhikaiMuhurta represents a muhurta in Naazhikai units
type NaazhikaiMuhurta struct {
	MuhurtaNumber  int
	StartNaazhikai int
	EndNaazhikai   int
	StartTime      time.Time
	EndTime        time.Time
	Quality        string
	TamilName      string
}

// getMuhurtaQuality determines the quality of a muhurta based on Naazhikai
func (n *NaazhikaiConverter) getMuhurtaQuality(startNaazhikai int) string {
	// Simplified quality assessment
	// In traditional Tamil astrology, certain Naazhikai periods are more auspicious
	auspiciousNaazhikai := []int{4, 6, 8, 12, 16, 20, 24, 28, 32}

	for _, auspicious := range auspiciousNaazhikai {
		if startNaazhikai == auspicious {
			return "auspicious"
		}
	}

	inauspiciousNaazhikai := []int{10, 14, 18, 22, 26, 30}
	for _, inauspicious := range inauspiciousNaazhikai {
		if startNaazhikai == inauspicious {
			return "inauspicious"
		}
	}

	return "neutral"
}

// getNaazhikaiMuhurtaName returns traditional Tamil muhurta names
func (n *NaazhikaiConverter) getNaazhikaiMuhurtaName(muhurtaNumber int) string {
	// Traditional Tamil muhurta names (simplified)
	names := map[int]string{
		1:  "ரோதனம் (Rodhanam)",
		2:  "விஜயம் (Vijayam)",
		3:  "அமிர்தம் (Amirtham)",
		4:  "முரஞ்ஜம் (Muranjam)",
		5:  "அபிஜித் (Abhijit)",
		6:  "சாயாஹ்னம் (Saayaahnam)",
		7:  "அமல (Amala)",
		8:  "சர்வார்த்த சித்தி (Sarvaartha Siddhi)",
		9:  "சுகந்தம் (Sugandham)",
		10: "அஷ்ட நாகம் (Ashta Naagam)",
		11: "வசுமத (Vasumadha)",
		12: "புண்ய (Punya)",
	}

	if name, exists := names[muhurtaNumber]; exists {
		return name
	}

	return fmt.Sprintf("முகூர்த்தம் %d", muhurtaNumber)
}

// FormatNaazhikaiTime formats Naazhikai time in Tamil numerals (optional)
func (n *NaazhikaiConverter) FormatNaazhikaiTime(naazhikai, vinaazhikai int, useTamilNumerals bool) string {
	if !useTamilNumerals {
		return fmt.Sprintf("%d நாழிகை %d விநாழிகை", naazhikai, vinaazhikai)
	}

	// Tamil numerals conversion (simplified)
	tamilNaazhikai := n.toTamilNumerals(naazhikai)
	tamilVinaazhikai := n.toTamilNumerals(vinaazhikai)

	return fmt.Sprintf("%s நாழிகை %s விநாழிகை", tamilNaazhikai, tamilVinaazhikai)
}

// toTamilNumerals converts Arabic numerals to Tamil numerals
func (n *NaazhikaiConverter) toTamilNumerals(num int) string {
	// Tamil numerals: ௧ ௨ ௩ ௪ ௫ ௬ ௭ ௮ ௯ ௰
	tamilDigits := []rune{'௦', '௧', '௨', '௩', '௪', '௫', '௬', '௭', '௮', '௯'}

	if num == 0 {
		return string(tamilDigits[0])
	}

	var result []rune
	for num > 0 {
		digit := num % 10
		result = append([]rune{tamilDigits[digit]}, result...)
		num /= 10
	}

	return string(result)
}

// GetSunriseInNaazhikai calculates sunrise time in Naazhikai from midnight
func (n *NaazhikaiConverter) GetSunriseInNaazhikai(sunrise time.Time) TimeInNaazhikai {
	midnight := time.Date(sunrise.Year(), sunrise.Month(), sunrise.Day(), 0, 0, 0, 0, sunrise.Location())
	return n.ToNaazhikai(sunrise, midnight)
}

// CompareWithModernTime provides a comparison between Naazhikai and modern time
func (n *NaazhikaiConverter) CompareWithModernTime(currentTime, sunrise time.Time) map[string]interface{} {
	tn := n.ToNaazhikai(currentTime, sunrise)

	return map[string]interface{}{
		"modern_time": currentTime.Format("15:04:05"),
		"naazhikai_time": fmt.Sprintf("%d நாழிகை %d விநாழிகை", tn.Naazhikai, tn.Vinaazhikai),
		"tamil_period": n.GetNaazhikaiPeriodName(tn.Naazhikai),
		"naazhikai_value": tn.Naazhikai,
		"vinaazhikai_value": tn.Vinaazhikai,
		"minutes_from_sunrise": currentTime.Sub(sunrise).Minutes(),
		"conversion_factor": "1 Naazhikai = 24 minutes",
	}
}
