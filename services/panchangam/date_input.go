package panchangam

import (
	"strings"
	"time"
)

func parsePanchangamDateInput(value string) (time.Time, time.Time, bool, error) {
	value = strings.TrimSpace(value)

	if parsedDate, err := time.Parse("2006-01-02", value); err == nil {
		return parsedDate, time.Time{}, false, nil
	}

	parsedDateTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, time.Time{}, false, err
	}
	return time.Time{}, parsedDateTime, true, nil
}
