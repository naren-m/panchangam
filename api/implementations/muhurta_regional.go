package implementations

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
)

// getRegionalMuhurtas returns region-specific muhurtas
func (m *MuhurtaPlugin) getRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location, region api.Region, sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta

	switch region {
	case api.RegionTamilNadu:
		muhurtas = append(muhurtas, m.getTamilMuhurtas(sunTimes, dayLength)...)
	case api.RegionKerala:
		muhurtas = append(muhurtas, m.getKeralaMuhurtas(sunTimes, dayLength)...)
	case api.RegionBengal:
		muhurtas = append(muhurtas, m.getBengalMuhurtas(sunTimes, dayLength)...)
	}

	return muhurtas
}

// getTamilMuhurtas returns Tamil Nadu specific muhurtas
func (m *MuhurtaPlugin) getTamilMuhurtas(sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta

	// Naazhikai-based calculations (Tamil time units)
	// 1 Naazhikai = 24 minutes, 60 Naazhikai = 1 day
	naazhikaiDuration := 24 * time.Minute

	// Shubha Muhurta (auspicious period in Tamil tradition)
	// Typically 2 Naazhikai in the morning
	shubhaStart := sunTimes.Sunrise.Add(2 * naazhikaiDuration)
	shubhaEnd := shubhaStart.Add(2 * naazhikaiDuration)

	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Shubha Muhurta",
		NameLocal:    "ஶுப முஹூர்த்த",
		StartTime:    shubhaStart,
		EndTime:      shubhaEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"business", "education", "ceremonies"},
		Significance: "Tamil traditional auspicious period",
		Region:       api.RegionTamilNadu,
		Metadata: map[string]interface{}{
			"duration_naazhikai": 2,
			"calculation_system": "tamil_naazhikai",
		},
	})

	return muhurtas
}

// getKeralaMuhurtas returns Kerala specific muhurtas
func (m *MuhurtaPlugin) getKeralaMuhurtas(sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta

	// Kerala specific Malyalam calendar muhurtas
	// Uchcha Kalam (auspicious time)
	ucchaStart := sunTimes.Sunrise.Add(dayLength / 4)
	ucchaEnd := ucchaStart.Add(90 * time.Minute)

	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Uchcha Kalam",
		NameLocal:    "ഉച്ച കാലം",
		StartTime:    ucchaStart,
		EndTime:      ucchaEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"religious_ceremonies", "important_decisions"},
		Significance: "Kerala traditional auspicious period",
		Region:       api.RegionKerala,
		Metadata: map[string]interface{}{
			"calculation_system": "malayalam_calendar",
		},
	})

	return muhurtas
}

// getBengalMuhurtas returns Bengal specific muhurtas
func (m *MuhurtaPlugin) getBengalMuhurtas(sunTimes *astronomy.SunTimes, dayLength time.Duration) []api.Muhurta {
	var muhurtas []api.Muhurta

	// Labha Kaal (beneficial time in Bengali tradition)
	labhaStart := sunTimes.Sunrise.Add(dayLength / 3)
	labhaEnd := labhaStart.Add(72 * time.Minute)

	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Labha Kaal",
		NameLocal:    "লাভ কাল",
		StartTime:    labhaStart,
		EndTime:      labhaEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"business", "financial_transactions", "investments"},
		Significance: "Bengali traditional beneficial period for gains",
		Region:       api.RegionBengal,
		Metadata: map[string]interface{}{
			"calculation_system": "bengali_calendar",
		},
	})

	return muhurtas
}
