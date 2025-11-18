package skyview

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// Observer represents the observer's location on Earth
type Observer struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"` // in meters
	Timezone  string  `json:"timezone"`
}

// CelestialBody represents a celestial object in the sky
type CelestialBody struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	SanskritName     string             `json:"sanskrit_name"`
	HindiName        string             `json:"hindi_name"`
	Type             string             `json:"type"` // sun, moon, planet
	EclipticCoords   EclipticCoordinates `json:"ecliptic_coords"`
	EquatorialCoords *EquatorialCoordinates `json:"equatorial_coords,omitempty"`
	HorizontalCoords *HorizontalCoordinates `json:"horizontal_coords,omitempty"`
	Magnitude        float64            `json:"magnitude"`
	Color            string             `json:"color"`
	IsVisible        bool               `json:"is_visible"`
	RiseTime         *time.Time         `json:"rise_time,omitempty"`
	SetTime          *time.Time         `json:"set_time,omitempty"`
	TransitTime      *time.Time         `json:"transit_time,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// EclipticCoordinates represents ecliptic coordinates
type EclipticCoordinates struct {
	Longitude float64 `json:"longitude"` // degrees
	Latitude  float64 `json:"latitude"`  // degrees
	Distance  float64 `json:"distance"`  // AU
}

// EquatorialCoordinates represents equatorial coordinates
type EquatorialCoordinates struct {
	RightAscension float64 `json:"right_ascension"` // degrees
	Declination    float64 `json:"declination"`     // degrees
	Distance       float64 `json:"distance"`        // AU
}

// HorizontalCoordinates represents horizontal (alt-az) coordinates
type HorizontalCoordinates struct {
	Azimuth  float64 `json:"azimuth"`  // degrees from north
	Altitude float64 `json:"altitude"` // degrees above horizon
	Distance float64 `json:"distance"` // AU
}

// SkyViewResponse represents the complete sky view data
type SkyViewResponse struct {
	Timestamp       time.Time       `json:"timestamp"`
	Observer        Observer        `json:"observer"`
	Bodies          []CelestialBody `json:"bodies"`
	VisibleBodies   []CelestialBody `json:"visible_bodies"`
	JulianDay       float64         `json:"julian_day"`
	LocalSiderealTime float64       `json:"local_sidereal_time"`
}

// SkyViewService provides sky visualization data
type SkyViewService struct {
	ephemerisProvider ephemeris.EphemerisProvider
}

// NewSkyViewService creates a new sky view service
func NewSkyViewService(provider ephemeris.EphemerisProvider) *SkyViewService {
	return &SkyViewService{
		ephemerisProvider: provider,
	}
}

// GetSkyView returns the complete sky view for a given time and observer
func (s *SkyViewService) GetSkyView(ctx context.Context, observer Observer, t time.Time) (*SkyViewResponse, error) {
	if s.ephemerisProvider == nil {
		return nil, fmt.Errorf("ephemeris provider is not initialized")
	}

	// Convert time to Julian day
	jd := ephemeris.JulianDay(DateToJulianDay(t))

	// Get planetary positions from ephemeris
	positions, err := s.ephemerisProvider.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		return nil, fmt.Errorf("failed to get planetary positions: %w", err)
	}

	// Get detailed Moon position for additional metadata
	moonPos, err := s.ephemerisProvider.GetMoonPosition(ctx, jd)
	if err != nil {
		return nil, fmt.Errorf("failed to get moon position: %w", err)
	}

	// Calculate Local Sidereal Time
	lst := GetLocalSiderealTime(observer.Longitude, t)

	// Create celestial bodies list
	bodies := []CelestialBody{}
	visibleBodies := []CelestialBody{}

	// Add Sun
	sun := s.createCelestialBody("sun", "Sun", "Surya", "सूर्य", "sun",
		positions.Sun, -26.7, "#ffee00", observer, lst)
	bodies = append(bodies, sun)
	if sun.IsVisible {
		visibleBodies = append(visibleBodies, sun)
	}

	// Add Moon
	moon := s.createCelestialBody("moon", "Moon", "Chandra", "चन्द्र", "moon",
		positions.Moon, -12.6, "#ffffff", observer, lst)
	moon.Metadata = map[string]interface{}{
		"phase":         moonPos.Phase,
		"illumination":  moonPos.Illumination,
		"phase_angle":   moonPos.PhaseAngle,
		"angular_diameter": moonPos.AngularDiameter,
	}
	bodies = append(bodies, moon)
	if moon.IsVisible {
		visibleBodies = append(visibleBodies, moon)
	}

	// Add planets
	planetData := []struct {
		id           string
		name         string
		sanskritName string
		hindiName    string
		position     ephemeris.Position
		magnitude    float64
		color        string
	}{
		{"mercury", "Mercury", "Budha", "बुध", positions.Mercury, -0.4, "#b8b8b8"},
		{"venus", "Venus", "Shukra", "शुक्र", positions.Venus, -4.4, "#ffc649"},
		{"mars", "Mars", "Mangala", "मङ्गल", positions.Mars, -1.5, "#ff6347"},
		{"jupiter", "Jupiter", "Guru", "गुरु", positions.Jupiter, -2.5, "#daa520"},
		{"saturn", "Saturn", "Shani", "शनि", positions.Saturn, 0.5, "#f4e0b8"},
		{"uranus", "Uranus", "Arun", "अरुण", positions.Uranus, 5.7, "#4fd0e0"},
		{"neptune", "Neptune", "Varun", "वरुण", positions.Neptune, 7.8, "#4169e1"},
	}

	for _, p := range planetData {
		body := s.createCelestialBody(p.id, p.name, p.sanskritName, p.hindiName, "planet",
			p.position, p.magnitude, p.color, observer, lst)
		bodies = append(bodies, body)
		if body.IsVisible {
			visibleBodies = append(visibleBodies, body)
		}
	}

	response := &SkyViewResponse{
		Timestamp:         t,
		Observer:          observer,
		Bodies:            bodies,
		VisibleBodies:     visibleBodies,
		JulianDay:         float64(jd),
		LocalSiderealTime: lst,
	}

	return response, nil
}

// createCelestialBody creates a celestial body with all calculated coordinates
func (s *SkyViewService) createCelestialBody(id, name, sanskritName, hindiName, bodyType string,
	pos ephemeris.Position, magnitude float64, color string, observer Observer, lst float64) CelestialBody {

	// Ecliptic coordinates (from ephemeris)
	ecliptic := EclipticCoordinates{
		Longitude: pos.Longitude,
		Latitude:  pos.Latitude,
		Distance:  pos.Distance,
	}

	// Convert to equatorial coordinates
	equatorial := eclipticToEquatorial(ecliptic)

	// Convert to horizontal coordinates
	horizontal := equatorialToHorizontal(equatorial, observer, lst)

	// Determine visibility (above horizon)
	isVisible := horizontal.Altitude > 0

	body := CelestialBody{
		ID:               id,
		Name:             name,
		SanskritName:     sanskritName,
		HindiName:        hindiName,
		Type:             bodyType,
		EclipticCoords:   ecliptic,
		EquatorialCoords: &equatorial,
		HorizontalCoords: &horizontal,
		Magnitude:        magnitude,
		Color:            color,
		IsVisible:        isVisible,
		Metadata:         make(map[string]interface{}),
	}

	return body
}

// Coordinate transformation functions

const (
	degToRad   = math.Pi / 180.0
	radToDeg   = 180.0 / math.Pi
	j2000Obliquity = 23.43929111 // Earth's obliquity at J2000 epoch in degrees
)

func eclipticToEquatorial(ecliptic EclipticCoordinates) EquatorialCoordinates {
	// For J2000 epoch (can be improved with precession calculation)
	lambda := ecliptic.Longitude * degToRad
	beta := ecliptic.Latitude * degToRad
	eps := j2000Obliquity * degToRad

	// Right Ascension
	ra := math.Atan2(
		math.Sin(lambda)*math.Cos(eps) - math.Tan(beta)*math.Sin(eps),
		math.Cos(lambda),
	)

	// Declination
	dec := math.Asin(
		math.Sin(beta)*math.Cos(eps) + math.Cos(beta)*math.Sin(eps)*math.Sin(lambda),
	)

	// Normalize RA to 0-360 degrees
	raDeg := ra * radToDeg
	if raDeg < 0 {
		raDeg += 360
	}

	return EquatorialCoordinates{
		RightAscension: raDeg,
		Declination:    dec * radToDeg,
		Distance:       ecliptic.Distance,
	}
}

func equatorialToHorizontal(equatorial EquatorialCoordinates, observer Observer, lst float64) HorizontalCoordinates {
	// Hour angle
	ha := (lst - equatorial.RightAscension) * degToRad

	dec := equatorial.Declination * degToRad
	lat := observer.Latitude * degToRad

	// Altitude
	alt := math.Asin(
		math.Sin(dec)*math.Sin(lat) + math.Cos(dec)*math.Cos(lat)*math.Cos(ha),
	)

	// Azimuth
	az := math.Atan2(
		-math.Cos(dec)*math.Sin(ha),
		math.Sin(dec)*math.Cos(lat) - math.Cos(dec)*math.Sin(lat)*math.Cos(ha),
	)

	// Convert to degrees and normalize azimuth to 0-360
	azDeg := az * radToDeg
	if azDeg < 0 {
		azDeg += 360
	}

	return HorizontalCoordinates{
		Azimuth:  azDeg,
		Altitude: alt * radToDeg,
		Distance: equatorial.Distance,
	}
}

// GetLocalSiderealTime calculates Local Sidereal Time in degrees
func GetLocalSiderealTime(longitude float64, t time.Time) float64 {
	jd := DateToJulianDay(t)
	T := (jd - 2451545.0) / 36525.0

	// Greenwich mean sidereal time at 0h UT
	gmst0 := 280.46061837 +
		360.98564736629*(jd-2451545.0) +
		0.000387933*T*T -
		T*T*T/38710000.0

	// Normalize to 0-360
	gmst0 = math.Mod(gmst0, 360)
	if gmst0 < 0 {
		gmst0 += 360
	}

	// Add hour angle
	ut := float64(t.Hour()) + float64(t.Minute())/60.0 + float64(t.Second())/3600.0

	lst := gmst0 + ut*15 + longitude

	// Normalize to 0-360
	lst = math.Mod(lst, 360)
	if lst < 0 {
		lst += 360
	}

	return lst
}

// DateToJulianDay converts a date to Julian Day number
func DateToJulianDay(date time.Time) float64 {
	year := date.Year()
	month := int(date.Month())
	day := date.Day()
	hour := date.Hour()
	minute := date.Minute()
	second := date.Second()

	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3

	jdn := float64(day) + float64((153*m+2)/5) +
		float64(365*y) + float64(y/4) -
		float64(y/100) + float64(y/400) - 32045.0

	// Add time of day
	jdn += (float64(hour)-12.0)/24.0 + float64(minute)/1440.0 + float64(second)/86400.0

	return jdn
}
