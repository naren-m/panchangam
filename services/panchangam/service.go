// server.go
package panchangam

import (
	"context"
	"time"

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
	// Create a child span for the service-level operation.
	logger.Info("Received request", "date", req.Date)
	d, _ := s.fetchPanchangamData(ctx, req.Date)
	response := &ppb.GetPanchangamResponse{
		PanchangamData: d,
	}
	time.Sleep(100 * time.Millisecond)
	logger.InfoContext(ctx, "Prepared response")

	return response, nil
}

func (s *PanchangamServer) fetchPanchangamData(ctx context.Context, date string) (*ppb.PanchangamData, error) {
	// span := observability.SpanFromContext(ctx)
	// // ctx, span := tracer.Start(ctx, "prepareOrderItemsAndShippingQuoteFromCart")
	// // defer span.End()
	// span := observability.SpanFromContext(ctx)

	logger.InfoContext(ctx, "fetching panchangam data")
	// Simulate a delay in fetching data.
	time.Sleep(29 * time.Millisecond)

	// Randomly return some error. This is just for testing.
	if rand.Intn(10)%2 == 0 {
		err := status.Error(codes.Internal, "failed to fetch panchangam data")
		logger.ErrorContext(ctx, "failed to fetch panchangam data", "error", err)
		return nil, err
	}
	return &ppb.PanchangamData{
		Date:        date,
		Tithi:       "Some Tithi",
		Nakshatra:   "Some Nakshatra",
		Yoga:        "Some Yoga",
		Karana:      "Some Karana",
		SunriseTime: "06:00:00",
		SunsetTime:  "18:00:00",
		Events: []*ppb.PanchangamEvent{
			{Name: "Some Event 1", Time: "08:00:00"},
			{Name: "Some Event 2", Time: "12:00:00"},
		},
	}, nil
}
