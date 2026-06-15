package panchangam

// calendarSystemByRegion maps regions to their traditional calendar system.
var calendarSystemByRegion = map[string]string{
	"Tamil Nadu":     "Amanta",
	"Kerala":         "Amanta",
	"Gujarat":        "Amanta",
	"Karnataka":      "Amanta",
	"Andhra Pradesh": "Purnimanta",
	"Telangana":      "Purnimanta",
	"Maharashtra":    "Purnimanta",
	"Uttar Pradesh":  "Purnimanta",
	"Bihar":          "Purnimanta",
	"West Bengal":    "Purnimanta",
	"Rajasthan":      "Purnimanta",
	"Madhya Pradesh": "Purnimanta",
	"Punjab":         "Purnimanta",
	"Odisha":         "Purnimanta",
	"Hyderabad":      "Purnimanta",
	"Chennai":        "Amanta",
	"Bangalore":      "Amanta",
	"Mumbai":         "Purnimanta",
	"Delhi":          "Purnimanta",
	"New York":       "Purnimanta",
	"Texas":          "Purnimanta",
	"New Jersey":     "Purnimanta",
	"California":     "Purnimanta",
}

func getCalendarSystemForRegion(region string) string {
	if system, exists := calendarSystemByRegion[region]; exists {
		return system
	}
	return "Purnimanta"
}

func moonRaasiName(longitude float64) string {
	names := []string{
		"Mesha", "Vrishabha", "Mithuna", "Karka", "Simha", "Kanya",
		"Tula", "Vrischika", "Dhanu", "Makara", "Kumbha", "Meena",
	}
	for longitude < 0 {
		longitude += 360
	}
	for longitude >= 360 {
		longitude -= 360
	}
	return names[int(longitude/30)%len(names)]
}
