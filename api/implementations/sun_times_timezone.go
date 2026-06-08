package implementations

import (
	"fmt"
	"strings"
	"time"

	"github.com/naren-m/panchangam/astronomy"
)

func convertSunTimesToTimezone(sunTimes *astronomy.SunTimes, timezone string) error {
	timezone = strings.TrimSpace(timezone)
	if timezone == "" {
		return nil
	}

	tz, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone %q: %w", timezone, err)
	}

	sunTimes.Sunrise = sunTimes.Sunrise.In(tz)
	sunTimes.Sunset = sunTimes.Sunset.In(tz)
	return nil
}
