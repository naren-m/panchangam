package panchangam

import (
	"context"
	"fmt"

	"github.com/naren-m/panchangam/astronomy"
	ppb "github.com/naren-m/panchangam/proto"
)

func (s *PanchangamServer) fetchPanchangamData(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.PanchangamData, error) {
	ctx, span := s.observer.CreateSpan(ctx, "fetchPanchangamData")
	defer span.End()

	logger.InfoContext(ctx, "Starting panchangam data fetch",
		"operation", "fetchPanchangamData",
		"date", req.Date,
		"location", fmt.Sprintf("%.4f,%.4f", req.Latitude, req.Longitude),
	)

	span.SetAttributes(
		traceAttribute("operation", "fetchPanchangamData"),
		traceAttribute("date", req.Date),
		traceAttribute("location", fmt.Sprintf("%.4f,%.4f", req.Latitude, req.Longitude)),
	)

	requestTime, err := parseRequestTime(ctx, req, span)
	if err != nil {
		return nil, err
	}

	location := requestLocation(req)
	logger.DebugContext(ctx, "Starting astronomical calculations",
		"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude))

	sunTimes, err := calculateSunTimes(ctx, span, requestTime.date, location)
	if err != nil {
		return nil, err
	}

	calculationTime := requestTime.calculationTime
	if !requestTime.hasClockTime {
		calculationTime = sunTimes.Sunrise.In(requestTime.timezoneLocation)
	}

	elements, err := s.calculateRequiredElements(ctx, span, req.Date, calculationTime, requestTime.date, location, req.Region)
	if err != nil {
		return nil, err
	}

	logger.InfoContext(ctx, "Building panchangam data response with real calculations",
		"operation", "buildResponse",
		"date", req.Date,
		"tithi", elements.tithi.Name,
		"nakshatra", elements.nakshatra.Name,
		"raasi", elements.raasi,
		"yoga", elements.yoga.Name,
		"karana", elements.karana.Name,
		"vara", elements.vara.Name,
		"sunrise", sunTimes.Sunrise.Format("15:04:05"),
		"sunset", sunTimes.Sunset.Format("15:04:05"))

	optionalData := calculateOptionalPanchangamData(ctx, requestTime.date, location, elements.tithi.Number)
	data := buildPanchangamData(ctx, requestTime, sunTimes, elements, optionalData)

	logger.InfoContext(ctx, "Panchangam data fetched successfully",
		"operation", "fetchPanchangamData",
		"date", data.Date,
		"tithi", data.Tithi,
		"nakshatra", data.Nakshatra,
		"events_count", len(data.Events))

	span.AddEvent("Data fetch completed", traceAttributes(
		"success", "true",
		"events_count", fmt.Sprintf("%d", len(data.Events)),
	)...)

	return data, nil
}

func requestLocation(req *ppb.GetPanchangamRequest) astronomy.Location {
	return astronomy.Location{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
}
