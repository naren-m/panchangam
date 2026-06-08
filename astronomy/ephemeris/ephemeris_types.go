package ephemeris

import (
	"context"
	"time"
)

// JulianDay represents a Julian day number.
type JulianDay float64

// PlanetaryPositions holds the positions of all planets.
type PlanetaryPositions struct {
	JulianDay JulianDay `json:"julian_day"`
	Sun       Position  `json:"sun"`
	Moon      Position  `json:"moon"`
	Mercury   Position  `json:"mercury"`
	Venus     Position  `json:"venus"`
	Mars      Position  `json:"mars"`
	Jupiter   Position  `json:"jupiter"`
	Saturn    Position  `json:"saturn"`
	Uranus    Position  `json:"uranus"`
	Neptune   Position  `json:"neptune"`
	Pluto     Position  `json:"pluto"`
}

// Position represents a celestial body's position.
type Position struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Distance  float64 `json:"distance"`
	Speed     float64 `json:"speed"`
}

// SolarPosition represents the Sun's position.
type SolarPosition struct {
	JulianDay         JulianDay `json:"julian_day"`
	Longitude         float64   `json:"longitude"`
	RightAscension    float64   `json:"right_ascension"`
	Declination       float64   `json:"declination"`
	Distance          float64   `json:"distance"`
	EquationOfTime    float64   `json:"equation_of_time"`
	MeanAnomaly       float64   `json:"mean_anomaly"`
	TrueAnomaly       float64   `json:"true_anomaly"`
	EccentricAnomaly  float64   `json:"eccentric_anomaly"`
	MeanLongitude     float64   `json:"mean_longitude"`
	ApparentLongitude float64   `json:"apparent_longitude"`
}

// LunarPosition represents the Moon's position.
type LunarPosition struct {
	JulianDay          JulianDay `json:"julian_day"`
	Longitude          float64   `json:"longitude"`
	Latitude           float64   `json:"latitude"`
	RightAscension     float64   `json:"right_ascension"`
	Declination        float64   `json:"declination"`
	Distance           float64   `json:"distance"`
	Phase              float64   `json:"phase"`
	PhaseAngle         float64   `json:"phase_angle"`
	Illumination       float64   `json:"illumination"`
	AngularDiameter    float64   `json:"angular_diameter"`
	MeanAnomaly        float64   `json:"mean_anomaly"`
	TrueAnomaly        float64   `json:"true_anomaly"`
	ArgumentOfLatitude float64   `json:"argument_of_latitude"`
	MeanLongitude      float64   `json:"mean_longitude"`
	TrueLongitude      float64   `json:"true_longitude"`
}

// HealthStatus represents the health status of an ephemeris provider.
type HealthStatus struct {
	Available    bool          `json:"available"`
	LastCheck    time.Time     `json:"last_check"`
	DataStartJD  float64       `json:"data_start_jd"`
	DataEndJD    float64       `json:"data_end_jd"`
	ResponseTime time.Duration `json:"response_time"`
	ErrorMessage string        `json:"error_message,omitempty"`
	Version      string        `json:"version,omitempty"`
	Source       string        `json:"source,omitempty"`
}

// EphemerisProvider defines the interface for ephemeris data providers.
type EphemerisProvider interface {
	GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error)
	GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error)
	GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error)
	IsAvailable(ctx context.Context) bool
	GetDataRange() (startJD, endJD JulianDay)
	GetHealthStatus(ctx context.Context) (*HealthStatus, error)
	GetProviderName() string
	GetVersion() string
	Close() error
}
