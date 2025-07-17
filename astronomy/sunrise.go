package astronomy

import (
	"math"
	"time"
)

const (
	// Solar depression angle for sunrise/sunset (geometric horizon)
	SolarDepressionAngle = 0.833

	// Degrees to radians conversion
	DegToRad = math.Pi / 180
	RadToDeg = 180 / math.Pi
)

// Location represents a geographic location with latitude and longitude
type Location struct {
	Latitude  float64
	Longitude float64
}

// SunTimes holds sunrise and sunset times
type SunTimes struct {
	Sunrise time.Time
	Sunset  time.Time
}

// CalculateSunTimes calculates sunrise and sunset times for a given location and date
func CalculateSunTimes(loc Location, date time.Time) (*SunTimes, error) {
	year, month, day := date.Date()
	
	// Convert to Julian day number
	jd := julianDayNumber(year, int(month), day)
	
	// Calculate centuries since J2000.0
	n := jd - 2451545.0
	
	// Mean longitude of the Sun
	L := math.Mod(280.460+0.9856474*n, 360.0)
	
	// Mean anomaly of the Sun
	g := math.Mod(357.528+0.9856003*n, 360.0) * DegToRad
	
	// Ecliptic longitude of the Sun
	lambda := L + 1.915*math.Sin(g) + 0.020*math.Sin(2*g)
	
	// Obliquity of the ecliptic
	epsilon := 23.439 - 0.0000004*n
	
	// Right ascension
	alpha := math.Atan2(math.Cos(epsilon*DegToRad)*math.Sin(lambda*DegToRad), math.Cos(lambda*DegToRad)) * RadToDeg
	
	// Declination
	delta := math.Asin(math.Sin(epsilon*DegToRad)*math.Sin(lambda*DegToRad)) * RadToDeg
	
	// Equation of time (in minutes)
	EqT := 4 * (L - alpha)
	
	// Hour angle for sunrise/sunset
	latRad := loc.Latitude * DegToRad
	deltaRad := delta * DegToRad
	
	// Calculate hour angle
	cosH := (math.Cos(90.833*DegToRad) - math.Sin(latRad)*math.Sin(deltaRad)) / (math.Cos(latRad) * math.Cos(deltaRad))
	
	// Check for polar day or polar night
	if cosH > 1 {
		// Polar night - sun never rises
		return &SunTimes{
			Sunrise: time.Date(year, month, day, 12, 0, 0, 0, date.Location()),
			Sunset:  time.Date(year, month, day, 12, 0, 0, 0, date.Location()),
		}, nil
	} else if cosH < -1 {
		// Polar day - sun never sets
		return &SunTimes{
			Sunrise: time.Date(year, month, day, 0, 0, 0, 0, date.Location()),
			Sunset:  time.Date(year, month, day, 23, 59, 59, 0, date.Location()),
		}, nil
	}
	
	// Hour angle in degrees
	H := math.Acos(cosH) * RadToDeg
	
	// Solar noon (in decimal hours UTC)
	solarNoon := 12.0 - loc.Longitude/15.0 - EqT/60.0
	
	// Sunrise and sunset times (in decimal hours UTC)
	sunriseDecimal := solarNoon - H/15.0
	sunsetDecimal := solarNoon + H/15.0
	
	// Convert to time
	sunriseTime := decimalHoursToTime(sunriseDecimal, year, month, day, time.UTC)
	sunsetTime := decimalHoursToTime(sunsetDecimal, year, month, day, time.UTC)
	
	return &SunTimes{
		Sunrise: sunriseTime,
		Sunset:  sunsetTime,
	}, nil
}

// julianDate calculates Julian date for noon of the given date
func julianDate(date time.Time) float64 {
	year := date.Year()
	month := int(date.Month())
	day := date.Day()
	
	if month <= 2 {
		year--
		month += 12
	}
	
	a := year / 100
	b := 2 - a + a/4
	
	jd := math.Floor(365.25*(float64(year)+4716)) + 
		math.Floor(30.6001*(float64(month)+1)) + 
		float64(day) + float64(b) - 1524.5
	
	return jd
}

// solarPosition calculates equation of time and solar declination
func solarPosition(jd float64) (float64, float64) {
	n := jd - 2451545.0
	L := math.Mod(280.460+0.9856474*n, 360)
	g := math.Mod(357.528+0.9856003*n, 360) * DegToRad
	
	// Solar ecliptic longitude
	lambda := L + 1.915*math.Sin(g) + 0.020*math.Sin(2*g)
	
	// Right ascension
	ra := math.Atan2(math.Cos(23.44*DegToRad)*math.Sin(lambda*DegToRad), math.Cos(lambda*DegToRad)) * RadToDeg
	ra = math.Mod(ra+360, 360)
	
	// Equation of time (in minutes)
	eqTime := 4 * (L - ra)
	
	// Solar declination
	decl := math.Asin(math.Sin(23.44*DegToRad) * math.Sin(lambda*DegToRad))
	
	return eqTime, decl
}

// calculateRiseSet calculates sunrise and sunset times in minutes from midnight
func calculateRiseSet(latitude, longitude, jd, eqTime, decl float64) (float64, float64) {
	latRad := latitude * DegToRad
	
	// Hour angle with solar depression angle
	cosH := (math.Cos(SolarDepressionAngle*DegToRad) - math.Sin(latRad)*math.Sin(decl)) / 
		(math.Cos(latRad) * math.Cos(decl))
	
	// Check for polar day or polar night
	if cosH > 1 {
		// Polar night - sun never rises
		return 0, 0
	} else if cosH < -1 {
		// Polar day - sun never sets
		return 0, 24 * 60
	}
	
	H := math.Acos(cosH) * RadToDeg
	
	// Time corrections (equation of time and longitude correction)
	timeCorrection := eqTime + longitude*4
	
	// Sunrise and sunset times (minutes from midnight UTC)
	sunrise := 720 - 4*H - timeCorrection
	sunset := 720 + 4*H - timeCorrection
	
	// Ensure times are within 0-1440 minutes (24 hours)
	sunrise = math.Mod(sunrise+1440, 1440)
	sunset = math.Mod(sunset+1440, 1440)
	
	return sunrise, sunset
}

// GetSunriseTime returns just the sunrise time for a location and date
func GetSunriseTime(loc Location, date time.Time) (time.Time, error) {
	sunTimes, err := CalculateSunTimes(loc, date)
	if err != nil {
		return time.Time{}, err
	}
	return sunTimes.Sunrise, nil
}

// GetSunsetTime returns just the sunset time for a location and date
func GetSunsetTime(loc Location, date time.Time) (time.Time, error) {
	sunTimes, err := CalculateSunTimes(loc, date)
	if err != nil {
		return time.Time{}, err
	}
	return sunTimes.Sunset, nil
}

// julianDayNumber calculates Julian day number for the given date
func julianDayNumber(year, month, day int) float64 {
	if month <= 2 {
		year--
		month += 12
	}
	
	a := year / 100
	b := 2 - a + a/4
	
	jd := math.Floor(365.25*(float64(year)+4716)) + 
		math.Floor(30.6001*(float64(month)+1)) + 
		float64(day) + float64(b) - 1524.5
	
	return jd
}

// decimalHoursToTime converts decimal hours to time.Time
func decimalHoursToTime(decimalHours float64, year int, month time.Month, day int, loc *time.Location) time.Time {
	// Ensure the decimal hours is within 0-24 range
	decimalHours = math.Mod(decimalHours+24, 24)
	
	hours := int(decimalHours)
	minutes := int((decimalHours - float64(hours)) * 60)
	seconds := int(((decimalHours - float64(hours)) * 60 - float64(minutes)) * 60)
	
	return time.Date(year, month, day, hours, minutes, seconds, 0, loc)
}