package astronomy

import "time"

const (
	// Lunar constants
	LunarParallax        = 0.125        // degrees, average lunar parallax
	LunarSemidiameter    = 0.25         // degrees, average lunar semidiameter
	LunarDepressionAngle = 0.125 + 0.25 // parallax + semidiameter

	// Lunar orbital constants
	LunarMeanDistance = 384400.0     // km, mean distance to moon
	LunarSynodicMonth = 29.530588853 // days, synodic month
	LunarEccentricity = 0.0549       // orbital eccentricity

	// Lunar phase constants
	NewMoon      = 0.0
	FirstQuarter = 0.25
	FullMoon     = 0.5
	LastQuarter  = 0.75
)

// LunarTimes holds moonrise and moonset times
type LunarTimes struct {
	Moonrise  time.Time
	Moonset   time.Time
	IsVisible bool // whether moon is visible (not below horizon all day)
}

// LunarPosition represents the moon's position
type LunarPosition struct {
	RightAscension float64 // degrees
	Declination    float64 // degrees
	Distance       float64 // km
	Phase          float64 // 0.0 = new, 0.5 = full
	Illumination   float64 // percentage illuminated
}

// LunarPhase represents moon phase information
type LunarPhase struct {
	Phase        float64   // 0.0-1.0, where 0=new moon, 0.5=full moon
	Illumination float64   // percentage illuminated (0-100)
	Name         string    // phase name
	Age          float64   // days since new moon
	NextPhase    time.Time // time of next major phase
}
