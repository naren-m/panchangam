package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type panchangamElements struct {
	tithi          *astronomy.TithiInfo
	nakshatra      *astronomy.NakshatraInfo
	yoga           *astronomy.YogaInfo
	karana         *astronomy.KaranaInfo
	vara           *astronomy.VaraInfo
	raasi          string
	calendarSystem string
}

func (s *PanchangamServer) calculateRequiredElements(
	ctx context.Context,
	span trace.Span,
	requestDate string,
	calculationTime time.Time,
	date time.Time,
	location astronomy.Location,
	region string,
) (*panchangamElements, error) {
	logger.InfoContext(ctx, "Calculating Panchangam elements",
		"operation", "calculatePanchangamElements",
		"date", requestDate)

	calendarSystem := getCalendarSystemForRegion(region)

	tithi, err := s.tithiCalc.GetTithiForTimeWithCalendarSystem(ctx, calculationTime, calendarSystem)
	if err != nil {
		return nil, recordRequiredElementError(ctx, span, "tithi", "tithi_calculation", "failed to calculate tithi", err)
	}

	nakshatra, err := s.nakshatraCalc.GetNakshatraForTime(ctx, calculationTime)
	if err != nil {
		return nil, recordRequiredElementError(ctx, span, "nakshatra", "nakshatra_calculation", "failed to calculate nakshatra", err)
	}

	yoga, err := s.yogaCalc.GetYogaForTime(ctx, calculationTime)
	if err != nil {
		return nil, recordRequiredElementError(ctx, span, "yoga", "yoga_calculation", "failed to calculate yoga", err)
	}

	karana, err := s.karanaCalc.GetKaranaForTime(ctx, calculationTime)
	if err != nil {
		return nil, recordRequiredElementError(ctx, span, "karana", "karana_calculation", "failed to calculate karana", err)
	}

	vara, err := s.varaCalc.GetVaraForDate(ctx, date, location)
	if err != nil {
		return nil, recordRequiredElementError(ctx, span, "vara", "vara_calculation", "failed to calculate vara", err)
	}

	return &panchangamElements{
		tithi:          tithi,
		nakshatra:      nakshatra,
		yoga:           yoga,
		karana:         karana,
		vara:           vara,
		raasi:          moonRaasiName(nakshatra.MoonLongitude),
		calendarSystem: calendarSystem,
	}, nil
}

func recordRequiredElementError(ctx context.Context, span trace.Span, label, operation, message string, err error) error {
	grpcErr := status.Error(codes.Internal, fmt.Sprintf("%s: %v", message, err))
	observability.RecordError(ctx, grpcErr, observability.ErrorContext{
		Severity:  observability.SeverityHigh,
		Category:  observability.CategoryCalculation,
		Operation: operation,
		Component: "panchangam_service",
	})
	logger.ErrorContext(ctx, fmt.Sprintf("Failed to calculate %s", label), "error", grpcErr)
	span.RecordError(grpcErr)
	return grpcErr
}
