// server.go
package panchangam

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
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
	span := observability.SpanFromContext(ctx)
	logger.Info("Received request", "date", req.Date)
	d, _ := s.fetchPanchangamData(ctx, req.Date)
	response := &ppb.GetPanchangamResponse{
		PanchangamData: d,
	}
	time.Sleep(1 * time.Second)
	span.AddEvent("prepared")
	logger.InfoContext(ctx, "Prepared response")

	return response, nil
}

func (s *PanchangamServer) fetchPanchangamData(ctx context.Context, date string) (*ppb.PanchangamData, error) {
	span := observability.SpanFromContext(ctx)
	// ctx, span := tracer.Start(ctx, "prepareOrderItemsAndShippingQuoteFromCart")
	// defer span.End()

	time.Sleep(2 * time.Second)

	span.AddEvent("fetching panchangam data")
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
