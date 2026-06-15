package ephemeris

import (
	"math"
	"time"
)

// TimeToJulianDay converts a time.Time to Julian day number.
func TimeToJulianDay(t time.Time) JulianDay {
	utc := t.UTC()
	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()

	if month <= 2 {
		year--
		month += 12
	}

	a := year / 100
	b := 2 - a + a/4

	jd := math.Floor(365.25*float64(year+4716)) +
		math.Floor(30.6001*float64(month+1)) +
		float64(day) + float64(b) - 1524.5

	hour := float64(utc.Hour())
	minute := float64(utc.Minute())
	second := float64(utc.Second())

	jd += hour/24.0 + minute/1440.0 + second/86400.0

	return JulianDay(jd)
}

// JulianDayToTime converts a Julian day number to time.Time.
func JulianDayToTime(jd JulianDay) time.Time {
	z := math.Floor(float64(jd) + 0.5)
	f := float64(jd) + 0.5 - z

	var a float64
	if z < 2299161 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}

	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	day := int(b - d - math.Floor(30.6001*e) + f)
	var month int
	if e < 14 {
		month = int(e - 1)
	} else {
		month = int(e - 13)
	}

	var year int
	if month > 2 {
		year = int(c - 4716)
	} else {
		year = int(c - 4715)
	}

	dayFraction := f
	hours := dayFraction * 24
	hour := int(hours)
	minutes := (hours - float64(hour)) * 60
	minute := int(minutes)
	seconds := (minutes - float64(minute)) * 60
	second := int(seconds)
	nanosecond := int((seconds - float64(second)) * 1e9)

	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, time.UTC)
}
