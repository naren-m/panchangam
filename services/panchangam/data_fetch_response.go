package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	ppb "github.com/naren-m/panchangam/proto"
)

func buildPanchangamData(
	ctx context.Context,
	requestTime *panchangamRequestTime,
	sunTimes *astronomy.SunTimes,
	elements *panchangamElements,
	optionalData optionalPanchangamData,
) *ppb.PanchangamData {
	localSunrise := sunTimes.Sunrise.In(requestTime.timezoneLocation)
	localSunset := sunTimes.Sunset.In(requestTime.timezoneLocation)
	events := buildPanchangamEvents(requestTime.timezoneLocation, localSunrise, localSunset, elements, optionalData)

	logger.DebugContext(ctx, "Sun times converted to local timezone",
		"utc_sunrise", sunTimes.Sunrise.Format("15:04:05 MST"),
		"utc_sunset", sunTimes.Sunset.Format("15:04:05 MST"),
		"local_sunrise", localSunrise.Format("15:04:05 MST"),
		"local_sunset", localSunset.Format("15:04:05 MST"),
		"timezone", requestTime.timezoneLocation.String())

	tithiDisplay := fmt.Sprintf("%s - %s Paksha Day %d (%s)",
		elements.tithi.TraditionalName,
		elements.tithi.Paksha,
		elements.tithi.PakshaDay,
		elements.calendarSystem)

	return &ppb.PanchangamData{
		Date:           requestTime.responseDate,
		Tithi:          tithiDisplay,
		Nakshatra:      fmt.Sprintf("%s (%d)", elements.nakshatra.Name, elements.nakshatra.Number),
		Yoga:           fmt.Sprintf("%s (%d)", elements.yoga.Name, elements.yoga.Number),
		Karana:         fmt.Sprintf("%s (%d)", elements.karana.Name, elements.karana.Number),
		SunriseTime:    localSunrise.Format("15:04:05"),
		SunsetTime:     localSunset.Format("15:04:05"),
		Events:         events,
		Timezone:       requestTime.timezoneLocation.String(),
		TimezoneOffset: requestTime.timezoneInfo.Formatted,
		IsDst:          requestTime.timezoneInfo.IsDST,
	}
}

func buildPanchangamEvents(
	loc *time.Location,
	localSunrise time.Time,
	localSunset time.Time,
	elements *panchangamElements,
	optionalData optionalPanchangamData,
) []*ppb.PanchangamEvent {
	tithiStart := elements.tithi.StartTime.In(loc)
	tithiEnd := elements.tithi.EndTime.In(loc)
	if !tithiEnd.After(tithiStart) {
		tithiEnd = tithiEnd.Add(24 * time.Hour)
	}

	events := []*ppb.PanchangamEvent{
		{
			Name:      "Sunrise",
			Time:      localSunrise.Format("15:04:05"),
			EventType: "SUNRISE",
		},
		{
			Name:      "Sunset",
			Time:      localSunset.Format("15:04:05"),
			EventType: "SUNSET",
		},
		{
			Name:      fmt.Sprintf("Tithi: %s (%s Paksha)", elements.tithi.TraditionalName, elements.tithi.Paksha),
			Time:      tithiStart.Format("15:04:05"),
			EventType: "TITHI",
		},
		{
			Name:      fmt.Sprintf("Tithi starts: %s", elements.tithi.TraditionalName),
			Time:      tithiStart.Format("15:04:05"),
			EventType: "TITHI_START",
		},
		{
			Name:      fmt.Sprintf("Tithi ends: %s", elements.tithi.TraditionalName),
			Time:      tithiEnd.Format("15:04:05"),
			EventType: "TITHI_END",
		},
		{
			Name:      fmt.Sprintf("Nakshatra: %s", elements.nakshatra.Name),
			Time:      "00:00:00",
			EventType: "NAKSHATRA",
		},
		{
			Name:      fmt.Sprintf("Raasi: %s", elements.raasi),
			Time:      "00:00:00",
			EventType: "RAASI",
		},
		{
			Name:      fmt.Sprintf("Yoga: %s", elements.yoga.Name),
			Time:      "00:00:00",
			EventType: "YOGA",
		},
		{
			Name:      fmt.Sprintf("Karana: %s", elements.karana.Name),
			Time:      "00:00:00",
			EventType: "KARANA",
		},
		{
			Name:      fmt.Sprintf("Vara: %s", elements.vara.Name),
			Time:      localSunrise.Format("15:04:05"),
			EventType: "VARA",
		},
	}

	if optionalData.lunarTimes != nil && optionalData.lunarTimes.IsVisible {
		events = append(events,
			&ppb.PanchangamEvent{
				Name:      "Moonrise",
				Time:      optionalData.lunarTimes.Moonrise.Format("15:04:05"),
				EventType: "MOONRISE",
			},
			&ppb.PanchangamEvent{
				Name:      "Moonset",
				Time:      optionalData.lunarTimes.Moonset.Format("15:04:05"),
				EventType: "MOONSET",
			},
		)
	}

	if optionalData.lunarPhase != nil {
		events = append(events, &ppb.PanchangamEvent{
			Name:      fmt.Sprintf("Moon Phase: %s (%.1f%% illuminated)", optionalData.lunarPhase.Name, optionalData.lunarPhase.Illumination),
			Time:      "00:00:00",
			EventType: "MOON_PHASE",
		})
	}

	if optionalData.traditionalPeriods != nil {
		events = append(events,
			&ppb.PanchangamEvent{
				Name:      "Rahu Kalam",
				Time:      optionalData.traditionalPeriods.RahuKalam.Start.Format("15:04:05"),
				EventType: "RAHU_KALAM",
			},
			&ppb.PanchangamEvent{
				Name:      "Yamagandam",
				Time:      optionalData.traditionalPeriods.Yamagandam.Start.Format("15:04:05"),
				EventType: "YAMAGANDAM",
			},
			&ppb.PanchangamEvent{
				Name:      "Gulika Kalam",
				Time:      optionalData.traditionalPeriods.GulikaKalam.Start.Format("15:04:05"),
				EventType: "GULIKA_KALAM",
			},
		)

		if optionalData.traditionalPeriods.AbhijitMuhurta.Auspicious {
			events = append(events, &ppb.PanchangamEvent{
				Name:      "Abhijit Muhurta",
				Time:      optionalData.traditionalPeriods.AbhijitMuhurta.Start.Format("15:04:05"),
				EventType: "ABHIJIT_MUHURTA",
			})
		}
	}

	for _, festival := range optionalData.festivals {
		events = append(events, &ppb.PanchangamEvent{
			Name:      fmt.Sprintf("Festival: %s", festival),
			Time:      "00:00:00",
			EventType: "FESTIVAL",
		})
	}

	return events
}
