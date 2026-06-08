//go:build integration
// +build integration

package gateway

import "strconv"

func buildQueryString(date string, lat, lng float64, tz string) string {
	return "date=" + date + "&lat=" + floatToString(lat) + "&lng=" + floatToString(lng) + "&tz=" + tz
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func compareData(data1, data2 map[string]interface{}) bool {
	fields := []string{"date", "tithi", "nakshatra", "sunrise_time", "sunset_time"}

	for _, field := range fields {
		val1, ok1 := data1[field]
		val2, ok2 := data2[field]

		if ok1 != ok2 || val1 != val2 {
			return false
		}
	}

	return true
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
