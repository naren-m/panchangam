package ephemeris

import "time"

// RetrogradeMotion indicates whether a planet is in retrograde motion.
type RetrogradeMotion string

const (
	// MotionDirect indicates normal forward motion.
	MotionDirect RetrogradeMotion = "direct"

	// MotionRetrograde indicates backward retrograde motion.
	MotionRetrograde RetrogradeMotion = "retrograde"

	// MotionStationary indicates the planet is at a stationary point.
	MotionStationary RetrogradeMotion = "stationary"
)

// PlanetaryStation represents a stationary point where planet changes direction.
type PlanetaryStation struct {
	Planet      string
	JulianDay   JulianDay
	Time        time.Time
	Longitude   float64
	StationType StationType
	Speed       float64
}

// StationType indicates the type of planetary station.
type StationType string

const (
	// StationRetrograde indicates planet is becoming retrograde.
	StationRetrograde StationType = "station_retrograde"

	// StationDirect indicates planet is becoming direct.
	StationDirect StationType = "station_direct"
)

// RetrogradePeriod represents a period of retrograde motion.
type RetrogradePeriod struct {
	Planet           string
	StartJD          JulianDay
	EndJD            JulianDay
	StartTime        time.Time
	EndTime          time.Time
	StartLongitude   float64
	EndLongitude     float64
	Duration         time.Duration
	MaxRetroDistance float64
}

// MotionAnalysis provides comprehensive analysis of planetary motion.
type MotionAnalysis struct {
	JulianDay      JulianDay
	Planet         string
	Motion         RetrogradeMotion
	Speed          float64
	Longitude      float64
	IsNearStation  bool
	NextStation    *PlanetaryStation
	CurrentPeriod  *RetrogradePeriod
	RecentStations []PlanetaryStation
}
