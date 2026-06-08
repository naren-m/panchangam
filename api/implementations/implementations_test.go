package implementations

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

var testLocation = api.Location{
	Latitude:  19.0760,
	Longitude: 72.8777,
	Timezone:  "Asia/Kolkata",
	Name:      "Mumbai",
}

var testDate = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
