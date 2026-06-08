package astronomy

import "time"

// TithiType represents the categorization of Tithi
type TithiType string

const (
	TithiTypeNanda  TithiType = "Nanda"  // 1, 6, 11 (Joyful)
	TithiTypeBhadra TithiType = "Bhadra" // 2, 7, 12 (Auspicious)
	TithiTypeJaya   TithiType = "Jaya"   // 3, 8, 13 (Victorious)
	TithiTypeRikta  TithiType = "Rikta"  // 4, 9, 14 (Empty)
	TithiTypePurna  TithiType = "Purna"  // 5, 10, 15 (Full/Complete)
)

// TithiInfo represents a Tithi with its properties
type TithiInfo struct {
	Number          int       `json:"number"`           // 1-30 (Purnimanta) or adjusted (Amanta)
	Name            string    `json:"name"`             // Traditional Sanskrit name of the Tithi
	Type            TithiType `json:"type"`             // Category (Nanda, Bhadra, Jaya, Rikta, Purna)
	StartTime       time.Time `json:"start_time"`       // When this Tithi begins
	EndTime         time.Time `json:"end_time"`         // When this Tithi ends
	Duration        float64   `json:"duration"`         // Duration in hours
	IsShukla        bool      `json:"is_shukla"`        // true for Shukla Paksha, false for Krishna Paksha
	Paksha          string    `json:"paksha"`           // "Shukla" or "Krishna"
	PakshaDay       int       `json:"paksha_day"`       // 1-15 within the paksha
	TraditionalName string    `json:"traditional_name"` // Traditional Sanskrit name (Dvithiya, Thuthiya, etc.)
	MoonSunDiff     float64   `json:"moon_sun_diff"`    // Moon longitude - Sun longitude in degrees
	CalendarSystem  string    `json:"calendar_system"`  // "Purnimanta" or "Amanta"
}
