package api

import (
	"context"
	"fmt"
	"time"
)

// Version represents API version information
type Version struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Patch int    `json:"patch"`
	Pre   string `json:"pre,omitempty"`
}

// String returns the version as a string
func (v Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		version += "-" + v.Pre
	}
	return version
}

// Location represents a geographical location for calculations
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone,omitempty"`
	Name      string  `json:"name,omitempty"`
}

// CalculationMethod represents different calculation approaches
type CalculationMethod string

const (
	MethodDrik  CalculationMethod = "drik"
	MethodVakya CalculationMethod = "vakya"
	MethodAuto  CalculationMethod = "auto"
)

// Region represents different regional variations
type Region string

const (
	RegionNorthIndia Region = "north_india"
	RegionSouthIndia Region = "south_india"
	RegionTamilNadu  Region = "tamil_nadu"
	RegionKerala     Region = "kerala"
	RegionBengal     Region = "bengal"
	RegionGujarat    Region = "gujarat"
	RegionMaha       Region = "maharashtra"
	RegionGlobal     Region = "global"
)

// CalendarSystem represents different calendar systems
type CalendarSystem string

const (
	CalendarPurnimanta CalendarSystem = "purnimanta" // North Indian (month ends on full moon)
	CalendarAmanta     CalendarSystem = "amanta"     // South Indian (month ends on new moon)
	CalendarLunar      CalendarSystem = "lunar"      // Pure lunar calendar
	CalendarSolar      CalendarSystem = "solar"      // Solar calendar
)

// PanchangamRequest represents a request for Panchangam data
type PanchangamRequest struct {
	Date              time.Time         `json:"date"`
	Location          Location          `json:"location"`
	CalculationMethod CalculationMethod `json:"calculation_method,omitempty"`
	Region            Region            `json:"region,omitempty"`
	CalendarSystem    CalendarSystem    `json:"calendar_system,omitempty"`
	Locale            string            `json:"locale,omitempty"`
	IncludeMuhurtas   bool              `json:"include_muhurtas,omitempty"`
	IncludeEvents     bool              `json:"include_events,omitempty"`
	IncludeRahukalam  bool              `json:"include_rahukalam,omitempty"`
}

// Tithi represents lunar day information
type Tithi struct {
	Number     int       `json:"number"`
	Name       string    `json:"name"`
	NameLocal  string    `json:"name_local,omitempty"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Percentage float64   `json:"percentage"`
	Lord       string    `json:"lord,omitempty"`
	Quality    string    `json:"quality,omitempty"`
	IsRunning  bool      `json:"is_running"`
}

// Nakshatra represents lunar mansion information
type Nakshatra struct {
	Number     int       `json:"number"`
	Name       string    `json:"name"`
	NameLocal  string    `json:"name_local,omitempty"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Percentage float64   `json:"percentage"`
	Pada       int       `json:"pada"`
	Lord       string    `json:"lord,omitempty"`
	Deity      string    `json:"deity,omitempty"`
	Symbol     string    `json:"symbol,omitempty"`
	Quality    string    `json:"quality,omitempty"`
	IsRunning  bool      `json:"is_running"`
}

// Yoga represents the combined solar-lunar position
type Yoga struct {
	Number     int       `json:"number"`
	Name       string    `json:"name"`
	NameLocal  string    `json:"name_local,omitempty"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Percentage float64   `json:"percentage"`
	Quality    string    `json:"quality,omitempty"`
	IsRunning  bool      `json:"is_running"`
}

// Karana represents half-tithi information
type Karana struct {
	Number     int       `json:"number"`
	Name       string    `json:"name"`
	NameLocal  string    `json:"name_local,omitempty"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Percentage float64   `json:"percentage"`
	Type       string    `json:"type"` // "movable" or "fixed"
	Quality    string    `json:"quality,omitempty"`
	IsRunning  bool      `json:"is_running"`
}

// Vara represents weekday information
type Vara struct {
	Number       int    `json:"number"`
	Name         string `json:"name"`
	NameLocal    string `json:"name_local,omitempty"`
	Lord         string `json:"lord,omitempty"`
	Color        string `json:"color,omitempty"`
	Gemstone     string `json:"gemstone,omitempty"`
	Significance string `json:"significance,omitempty"`
}

// SunMoonTimes represents sunrise, sunset, and moon timing information
type SunMoonTimes struct {
	Sunrise          time.Time `json:"sunrise"`
	Sunset           time.Time `json:"sunset"`
	SolarNoon        time.Time `json:"solar_noon"`
	DayLength        Duration  `json:"day_length"`
	Moonrise         time.Time `json:"moonrise,omitempty"`
	Moonset          time.Time `json:"moonset,omitempty"`
	MoonPhase        float64   `json:"moon_phase"`        // 0.0 = New Moon, 0.5 = Full Moon
	MoonIllumination float64   `json:"moon_illumination"` // Percentage
}

// Duration represents a time duration with human-readable format
type Duration struct {
	time.Duration
}

// MarshalJSON implements JSON marshaling for Duration
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// Event represents special events, festivals, or occurrences
type Event struct {
	Name         string                 `json:"name"`
	NameLocal    string                 `json:"name_local,omitempty"`
	Type         EventType              `json:"type"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time,omitempty"`
	Significance string                 `json:"significance,omitempty"`
	Region       Region                 `json:"region,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// EventType represents different types of events
type EventType string

const (
	EventTypeFestival    EventType = "festival"
	EventTypeEkadashi    EventType = "ekadashi"
	EventTypeAmavasya    EventType = "amavasya"
	EventTypePurnima     EventType = "purnima"
	EventTypeVrat        EventType = "vrat"
	EventTypeRahukalam   EventType = "rahukalam"
	EventTypeYamagandam  EventType = "yamagandam"
	EventTypeGulikakalam EventType = "gulikakalam"
	EventTypeAbhijit     EventType = "abhijit"
	EventTypeBrahma      EventType = "brahma_muhurta"
	EventTypeSankashti   EventType = "sankashti"
	EventTypeAshtami     EventType = "ashtami"
	EventTypeNavami      EventType = "navami"
	EventTypeSolar       EventType = "solar_event"
	EventTypeLunar       EventType = "lunar_event"
	EventTypePlanetary   EventType = "planetary_event"
)

// Muhurta represents auspicious time periods
type Muhurta struct {
	Name         string                 `json:"name"`
	NameLocal    string                 `json:"name_local,omitempty"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Quality      MuhurtaQuality         `json:"quality"`
	Purpose      []string               `json:"purpose,omitempty"` // What activities it's good for
	Avoid        []string               `json:"avoid,omitempty"`   // What to avoid
	Significance string                 `json:"significance,omitempty"`
	Region       Region                 `json:"region,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// MuhurtaQuality represents the quality of a muhurta
type MuhurtaQuality string

const (
	QualityAuspicious   MuhurtaQuality = "auspicious"
	QualityInauspicious MuhurtaQuality = "inauspicious"
	QualityNeutral      MuhurtaQuality = "neutral"
	QualityHighly       MuhurtaQuality = "highly_auspicious"
	QualityMildly       MuhurtaQuality = "mildly_auspicious"
	QualityAvoid        MuhurtaQuality = "avoid"
)

// PanchangamData represents complete Panchangam information for a day
type PanchangamData struct {
	Date           time.Time      `json:"date"`
	Location       Location       `json:"location"`
	Region         Region         `json:"region"`
	CalendarSystem CalendarSystem `json:"calendar_system"`

	// Five main elements
	Tithi     Tithi     `json:"tithi"`
	Nakshatra Nakshatra `json:"nakshatra"`
	Yoga      Yoga      `json:"yoga"`
	Karana    Karana    `json:"karana"`
	Vara      Vara      `json:"vara"`

	// Astronomical data
	SunMoonTimes SunMoonTimes `json:"sun_moon_times"`

	// Events and muhurtas
	Events   []Event   `json:"events,omitempty"`
	Muhurtas []Muhurta `json:"muhurtas,omitempty"`

	// Metadata
	CalculationMethod CalculationMethod `json:"calculation_method"`
	Ayanamsa          float64           `json:"ayanamsa,omitempty"`
	JulianDay         float64           `json:"julian_day,omitempty"`
	Locale            string            `json:"locale,omitempty"`
	Version           Version           `json:"version"`
	GeneratedAt       time.Time         `json:"generated_at"`
}

// PanchangamAPI defines the core interface for Panchangam calculations
type PanchangamAPI interface {
	// GetPanchangam returns Panchangam data for a specific request
	GetPanchangam(ctx context.Context, req PanchangamRequest) (*PanchangamData, error)

	// GetDateRange returns Panchangam data for a range of dates
	GetDateRange(ctx context.Context, start, end time.Time, location Location, options ...RequestOption) ([]*PanchangamData, error)

	// GetVersion returns the API version
	GetVersion() Version

	// GetSupportedRegions returns all supported regions
	GetSupportedRegions() []Region

	// GetSupportedMethods returns all supported calculation methods
	GetSupportedMethods() []CalculationMethod

	// GetSupportedCalendars returns all supported calendar systems
	GetSupportedCalendars() []CalendarSystem
}

// RequestOption allows for functional options pattern
type RequestOption func(*PanchangamRequest)

// WithCalculationMethod sets the calculation method
func WithCalculationMethod(method CalculationMethod) RequestOption {
	return func(req *PanchangamRequest) {
		req.CalculationMethod = method
	}
}

// WithRegion sets the regional variation
func WithRegion(region Region) RequestOption {
	return func(req *PanchangamRequest) {
		req.Region = region
	}
}

// WithCalendarSystem sets the calendar system
func WithCalendarSystem(system CalendarSystem) RequestOption {
	return func(req *PanchangamRequest) {
		req.CalendarSystem = system
	}
}

// WithLocale sets the locale for localized names
func WithLocale(locale string) RequestOption {
	return func(req *PanchangamRequest) {
		req.Locale = locale
	}
}

// WithMuhurtas includes muhurta calculations
func WithMuhurtas(include bool) RequestOption {
	return func(req *PanchangamRequest) {
		req.IncludeMuhurtas = include
	}
}

// WithEvents includes event calculations
func WithEvents(include bool) RequestOption {
	return func(req *PanchangamRequest) {
		req.IncludeEvents = include
	}
}

// WithRahukalam includes Rahukalam calculations
func WithRahukalam(include bool) RequestOption {
	return func(req *PanchangamRequest) {
		req.IncludeRahukalam = include
	}
}
