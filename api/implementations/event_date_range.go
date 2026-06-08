package implementations

import (
	"fmt"
	"time"
)

func validateEventDateRange(start, end time.Time) error {
	if !end.Before(start) {
		return nil
	}
	return fmt.Errorf("end date %s is before start date %s", end.Format("2006-01-02"), start.Format("2006-01-02"))
}
