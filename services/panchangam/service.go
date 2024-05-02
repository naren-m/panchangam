// server.go
package panchangam

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/logging"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PanchangamServer struct {
	span *logging.Span
	tracer trace.Tracer
    ppb.UnimplementedPanchangamServer
}

func NewPanchangamServer(span *logging.Span, tracer trace.Tracer) *PanchangamServer {
	return &PanchangamServer{
		span: span,
		tracer: tracer,
	}
}

func (s *PanchangamServer) Get(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.GetPanchangamResponse, error) {
    // Create a child span for the service-level operation.
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("app.date", req.GetDate()),
	)


    d, _ := s.fetchPanchangamData(ctx, req.Date)
    response := &ppb.GetPanchangamResponse{
        PanchangamData: d,
    }
	span.AddEvent("prepared")

    return response, nil
}

func (s *PanchangamServer) fetchPanchangamData(ctx context.Context, date string) (*ppb.PanchangamData, error) {
	ctx, span := s.tracer.Start(ctx, "prepareOrderItemsAndShippingQuoteFromCart")
	defer span.End()

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