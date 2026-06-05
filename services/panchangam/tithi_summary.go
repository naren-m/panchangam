package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultTithiSummaryTimezone = "UTC"
	defaultTithiSummaryMethod   = "Drik"
	defaultTithiSummaryLocale   = "en"
	tithiBoundarySearchStep     = time.Hour
	tithiBoundarySearchWindow   = 72 * time.Hour
	tithiBoundaryPrecision      = time.Second
)

func (s *PanchangamServer) GetTithiSummary(ctx context.Context, req *ppb.GetTithiSummaryRequest) (*ppb.GetTithiSummaryResponse, error) {
	ctx, span := s.observer.CreateSpan(ctx, "PanchangamServer.GetTithiSummary")
	defer span.End()

	if req == nil {
		err := status.Error(codes.InvalidArgument, "request cannot be nil")
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("request.at", req.At),
		attribute.Float64("request.latitude", req.Latitude),
		attribute.Float64("request.longitude", req.Longitude),
		attribute.String("request.timezone", req.Timezone),
		attribute.String("request.region", req.Region),
		attribute.String("request.method", req.CalculationMethod),
		attribute.String("request.locale", req.Locale),
		attribute.String("request.calendar_system", req.CalendarSystem),
	)

	input, err := validateTithiSummaryRequest(req)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("calculation.timezone", input.timezone),
		attribute.String("calculation.region", input.region),
		attribute.String("calculation.method", input.method),
		attribute.String("calculation.locale", input.locale),
		attribute.String("calculation.calendar_system", input.calendarSystem),
	)

	tithi, err := s.getTithiAt(ctx, input.at, input.calendarSystem)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate tithi: %v", err))
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "tithi_summary_current_tithi",
			Component: "panchangam_service",
		})
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	startTime, err := s.findTithiBoundary(ctx, input.at, tithi.Number, input.calendarSystem, -1)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate tithi start time: %v", err))
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	endTime, err := s.findTithiBoundary(ctx, input.at, tithi.Number, input.calendarSystem, 1)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate tithi end time: %v", err))
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	localAt := input.at.In(input.location)
	nakshatra, yoga, karana, vara, err := s.getPanchaAngaSummary(ctx, localAt, input.latitude, input.longitude)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	day, err := getTithiSummaryDay(ctx, localAt, input.latitude, input.longitude, input.location)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate day details: %v", err))
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	resp := &ppb.GetTithiSummaryResponse{
		Date: localAt.Format("2006-01-02"),
		Tithi: &ppb.TithiSummary{
			Name:            tithi.Name,
			TraditionalName: tithi.TraditionalName,
			Number:          int32(tithi.Number),
			Paksha:          tithi.Paksha,
			PakshaDay:       int32(tithi.PakshaDay),
			Type:            string(tithi.Type),
			StartTime:       startTime.UTC().Format(time.RFC3339),
			EndTime:         endTime.UTC().Format(time.RFC3339),
		},
		PanchaAnga: &ppb.PanchaAngaSummary{
			Nakshatra: nakshatra.Name,
			Yoga:      yoga.Name,
			Karana:    karana.Name,
			Vara:      vara.GregorianDay,
		},
		Calculation: &ppb.TithiCalculationSummary{
			Timezone:       input.timezone,
			Region:         input.region,
			CalendarSystem: input.calendarSystem,
			Method:         input.method,
			Locale:         input.locale,
		},
		Day:           day,
		GeneratedAt:   input.at.UTC().Format(time.RFC3339),
		NextRefreshAt: endTime.UTC().Format(time.RFC3339),
	}

	span.AddEvent("Tithi summary response prepared", traceAttributes(
		"response.date", resp.Date,
		"response.tithi", resp.Tithi.Name,
		"response.sunrise", resp.Day.SunriseTime,
		"response.sunset", resp.Day.SunsetTime,
		"response.next_refresh_at", resp.NextRefreshAt,
	)...)

	return resp, nil
}

type tithiSummaryInput struct {
	at             time.Time
	latitude       float64
	longitude      float64
	timezone       string
	location       *time.Location
	region         string
	method         string
	locale         string
	calendarSystem string
}

func validateTithiSummaryRequest(req *ppb.GetTithiSummaryRequest) (*tithiSummaryInput, error) {
	if req.At == "" {
		return nil, status.Error(codes.InvalidArgument, "at parameter is required")
	}

	at, err := time.Parse(time.RFC3339, req.At)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid at timestamp: %v", err))
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
	}

	timezoneName := req.Timezone
	if timezoneName == "" {
		timezoneName = defaultTithiSummaryTimezone
	}

	location, err := NewTimezoneParser().ParseTimezone(timezoneName)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid timezone: %v", err))
	}

	method := req.CalculationMethod
	if method == "" {
		method = defaultTithiSummaryMethod
	}

	locale := req.Locale
	if locale == "" {
		locale = defaultTithiSummaryLocale
	}

	calendarSystem := req.CalendarSystem
	if calendarSystem == "" {
		calendarSystem = getCalendarSystemForRegion(req.Region)
	}
	if calendarSystem != "Purnimanta" && calendarSystem != "Amanta" {
		return nil, status.Error(codes.InvalidArgument, "calendar_system must be Purnimanta or Amanta")
	}

	return &tithiSummaryInput{
		at:             at,
		latitude:       req.Latitude,
		longitude:      req.Longitude,
		timezone:       timezoneName,
		location:       location,
		region:         req.Region,
		method:         method,
		locale:         locale,
		calendarSystem: calendarSystem,
	}, nil
}

func (s *PanchangamServer) getTithiAt(ctx context.Context, at time.Time, calendarSystem string) (*astronomy.TithiInfo, error) {
	jd := ephemeris.TimeToJulianDay(at)
	positions, err := s.ephemerisManager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		return nil, err
	}

	return s.tithiCalc.GetTithiFromLongitudesWithCalendarSystem(
		ctx,
		positions.Sun.Longitude,
		positions.Moon.Longitude,
		at,
		calendarSystem,
	)
}

func (s *PanchangamServer) findTithiBoundary(ctx context.Context, at time.Time, currentTithiNumber int, calendarSystem string, direction int) (time.Time, error) {
	if direction != -1 && direction != 1 {
		return time.Time{}, fmt.Errorf("direction must be -1 or 1")
	}

	currentSide := at.UTC()
	step := time.Duration(direction) * tithiBoundarySearchStep
	searchLimit := at.UTC().Add(time.Duration(direction) * tithiBoundarySearchWindow)

	for probe := currentSide.Add(step); withinBoundarySearch(probe, searchLimit, direction); probe = probe.Add(step) {
		probeTithi, err := s.getTithiAt(ctx, probe, calendarSystem)
		if err != nil {
			return time.Time{}, err
		}

		if probeTithi.Number != currentTithiNumber {
			if direction < 0 {
				return s.bisectTithiBoundary(ctx, probe, currentSide, currentTithiNumber, calendarSystem, false)
			}
			return s.bisectTithiBoundary(ctx, currentSide, probe, currentTithiNumber, calendarSystem, true)
		}

		currentSide = probe
	}

	return time.Time{}, fmt.Errorf("no tithi transition found within %s of %s", tithiBoundarySearchWindow, at.Format(time.RFC3339))
}

func withinBoundarySearch(probe, limit time.Time, direction int) bool {
	if direction < 0 {
		return !probe.Before(limit)
	}
	return !probe.After(limit)
}

func (s *PanchangamServer) bisectTithiBoundary(ctx context.Context, low, high time.Time, currentTithiNumber int, calendarSystem string, lowIsCurrent bool) (time.Time, error) {
	for high.Sub(low) > tithiBoundaryPrecision {
		mid := low.Add(high.Sub(low) / 2)
		midTithi, err := s.getTithiAt(ctx, mid, calendarSystem)
		if err != nil {
			return time.Time{}, err
		}

		if lowIsCurrent {
			if midTithi.Number == currentTithiNumber {
				low = mid
			} else {
				high = mid
			}
			continue
		}

		if midTithi.Number == currentTithiNumber {
			high = mid
		} else {
			low = mid
		}
	}

	return high.UTC().Truncate(time.Second), nil
}

func (s *PanchangamServer) getPanchaAngaSummary(ctx context.Context, at time.Time, latitude, longitude float64) (*astronomy.NakshatraInfo, *astronomy.YogaInfo, *astronomy.KaranaInfo, *astronomy.VaraInfo, error) {
	nakshatra, err := s.nakshatraCalc.GetNakshatraForDate(ctx, at)
	if err != nil {
		return nil, nil, nil, nil, status.Error(codes.Internal, fmt.Sprintf("failed to calculate nakshatra: %v", err))
	}

	yoga, err := s.yogaCalc.GetYogaForDate(ctx, at)
	if err != nil {
		return nil, nil, nil, nil, status.Error(codes.Internal, fmt.Sprintf("failed to calculate yoga: %v", err))
	}

	karana, err := s.karanaCalc.GetKaranaForDate(ctx, at)
	if err != nil {
		return nil, nil, nil, nil, status.Error(codes.Internal, fmt.Sprintf("failed to calculate karana: %v", err))
	}

	vara, err := s.varaCalc.GetVaraForDate(ctx, at, astronomy.Location{
		Latitude:  latitude,
		Longitude: longitude,
	})
	if err != nil {
		return nil, nil, nil, nil, status.Error(codes.Internal, fmt.Sprintf("failed to calculate vara: %v", err))
	}

	return nakshatra, yoga, karana, vara, nil
}

func getTithiSummaryDay(ctx context.Context, at time.Time, latitude, longitude float64, location *time.Location) (*ppb.DaySummary, error) {
	geo := astronomy.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, geo, at)
	if err != nil {
		return nil, err
	}

	traditionalPeriods, err := astronomy.CalculateTraditionalPeriodsWithContext(ctx, geo, at)
	if err != nil {
		return nil, err
	}
	if traditionalPeriods.AbhijitMuhurta == nil {
		return nil, fmt.Errorf("abhijit muhurta is unavailable")
	}

	return &ppb.DaySummary{
		SunriseTime: sunTimes.Sunrise.In(location).Format("15:04"),
		SunsetTime:  sunTimes.Sunset.In(location).Format("15:04"),
		AbhijitMuhurta: &ppb.TimeWindow{
			Name:       "Abhijit",
			StartTime:  traditionalPeriods.AbhijitMuhurta.Start.In(location).Format("15:04"),
			EndTime:    traditionalPeriods.AbhijitMuhurta.End.In(location).Format("15:04"),
			Auspicious: traditionalPeriods.AbhijitMuhurta.Auspicious,
		},
	}, nil
}
