package astronomy

import "time"

func siderealLongitude(tropicalLongitude float64, date time.Time) float64 {
	return normalizeLongitude(tropicalLongitude - lahiriAyanamsha(date))
}

func lahiriAyanamsha(date time.Time) float64 {
	year := decimalYear(date)
	return 23.852778 + ((year - 2000.0) * 0.013969)
}

func decimalYear(date time.Time) float64 {
	start := time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
	end := time.Date(date.Year()+1, time.January, 1, 0, 0, 0, 0, date.Location())
	return float64(date.Year()) + date.Sub(start).Seconds()/end.Sub(start).Seconds()
}
