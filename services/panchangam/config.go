package panchangam

import "time"

// Config holds the configuration for the Panchangam service
type Config struct {
	// Cache settings
	CacheSize int
	CacheTTL  time.Duration

	// Calendar settings
	DefaultCalendarSystem string
	RegionCalendarSystems map[string]string
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		CacheSize:             1000,
		CacheTTL:              24 * time.Hour,
		DefaultCalendarSystem: "Purnimanta",
		RegionCalendarSystems: map[string]string{
			// South Indian states typically use Amanta system
			"Tamil Nadu":     "Amanta",
			"Kerala":         "Amanta",
			"Karnataka":      "Amanta",
			"Andhra Pradesh": "Amanta",
			"Telangana":      "Amanta",

			// North Indian states typically use Purnimanta system
			"Maharashtra":      "Purnimanta",
			"Gujarat":          "Purnimanta",
			"Rajasthan":        "Purnimanta",
			"Uttar Pradesh":    "Purnimanta",
			"Madhya Pradesh":   "Purnimanta",
			"Bihar":            "Purnimanta",
			"West Bengal":      "Purnimanta",
			"Odisha":           "Purnimanta",
			"Punjab":           "Purnimanta",
			"Haryana":          "Purnimanta",
			"Himachal Pradesh": "Purnimanta",
			"Uttarakhand":      "Purnimanta",
			"Delhi":            "Purnimanta",

			// Default for other regions
			"California": "Purnimanta",
			"New York":   "Purnimanta",
			"Texas":      "Purnimanta",
			"New Jersey": "Purnimanta",
		},
	}
}
