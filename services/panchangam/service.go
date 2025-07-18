// server.go
package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = log.Logger()

type PanchangamServer struct {
	observer observability.ObserverInterface
	ppb.UnimplementedPanchangamServer
}

func NewPanchangamServer() *PanchangamServer {
	return &PanchangamServer{
		observer: observability.Observer(),
	}
}

func (s *PanchangamServer) Get(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.GetPanchangamResponse, error) {
	ctx, span := s.observer.CreateSpan(ctx, "Get")
	defer span.End()
	// Create a child span for the service-level operation.
	logger.InfoContext(ctx, "Received request", "date", req.Date, "lat", req.Latitude, "lon", req.Longitude)
	
	// Validate request
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
	}
	
	d, err := s.fetchPanchangamData(ctx, req)
	if err != nil {
		return nil, err
	}
	
	response := &ppb.GetPanchangamResponse{
		PanchangamData: d,
	}
	time.Sleep(100 * time.Millisecond)
	logger.InfoContext(ctx, "Prepared response")

	return response, nil
}

func (s *PanchangamServer) fetchPanchangamData(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.PanchangamData, error) {
	ctx, span := s.observer.CreateSpan(ctx, "fetchPanchangamData")
	defer span.End()

	logger.InfoContext(ctx, "fetching panchangam data")
	// Simulate a delay in fetching data.
	time.Sleep(29 * time.Millisecond)

	// Parse the date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid date format: %v", err))
	}

	// Set timezone
	loc := time.Local // default to local timezone
	if req.Timezone != "" {
		parsedLoc, err := time.LoadLocation(req.Timezone)
		if err != nil {
			logger.WarnContext(ctx, "failed to load timezone, using local", "timezone", req.Timezone, "error", err)
		} else {
			loc = parsedLoc
		}
	}
	
	// Adjust date to the requested timezone
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)

	// Calculate sunrise and sunset
	location := astronomy.Location{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	
	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, date)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to calculate sun times: %v", err))
	}

	// Randomly return some error. This is just for testing.
	if rand.Intn(10)%2 == 0 {
		err := status.Error(codes.Internal, "failed to fetch panchangam data")
		logger.ErrorContext(ctx, "failed to fetch panchangam data", "error", err)
		return nil, err
	}
	
	return &ppb.PanchangamData{
		Date:        req.Date,
		Tithi:       "Some Tithi",
		Nakshatra:   "Some Nakshatra",
		Yoga:        "Some Yoga",
		Karana:      "Some Karana",
		SunriseTime: sunTimes.Sunrise.Format("15:04:05"),
		SunsetTime:  sunTimes.Sunset.Format("15:04:05"),
		Events: []*ppb.PanchangamEvent{
			{Name: "Some Event 1", Time: "08:00:00"},
			{Name: "Some Event 2", Time: "12:00:00"},
		},
	}, nil
}
