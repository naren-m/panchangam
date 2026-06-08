package implementations

// VakyaConstants holds traditional Vakya calculation constants
type VakyaConstants struct {
	// Traditional mean motion values (degrees per day)
	SunMeanMotion  float64 // ~0.985647 degrees per day
	MoonMeanMotion float64 // ~13.176358 degrees per day

	// Traditional epoch (reference point)
	EpochJD float64 // Kaliyuga start or other traditional epoch

	// Ayanamsa value for traditional calculations
	TraditionalAyanamsa float64

	// Correction factors for different celestial bodies
	SunCorrection  float64
	MoonCorrection float64
}

// DrikGanitaConfig holds modern calculation configuration
type DrikGanitaConfig struct {
	// Modern ephemeris to use (Swiss, JPL, etc.)
	EphemerisType string

	// Ayanamsa system to use
	AyanamsaSystem string

	// Precision level
	PrecisionLevel string

	// Use atmospheric refraction corrections
	UseAtmosphericCorrection bool

	// Use delta-T corrections
	UseDeltaTCorrection bool
}
