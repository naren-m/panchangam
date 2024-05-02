// server.go
package panchangam

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/logging"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
)

type PanchangamServer struct {
	span *logging.Span
    ppb.UnimplementedPanchangamServer
}

func NewPanchangamServer(span *logging.Span) *PanchangamServer {
	return &PanchangamServer{
		span: span,
	}
}

func (s *PanchangamServer) Get(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.GetPanchangamResponse, error) {
    // Create a child span for the service-level operation.
    span := logging.NewSpan(ctx, "PanchangamServer.Get", logrus.DebugLevel)
    defer span.End()

	var myKey = attribute.Key("PanchangamServer")
	span.Span.SetAttributes(myKey.String("Get"))

	span.Trace("Fetching panchangam data", logrus.Fields{"date": req.Date, })

    d, _ := s.fetchPanchangamData(span.Ctx, req.Date)
    response := &ppb.GetPanchangamResponse{
        PanchangamData: d,
    }

    span1 := logging.NewSpan(s.span.Ctx, "Second PanchangamServer.Get", logrus.DebugLevel)
    defer span1.End()

    return response, nil
}

func (s *PanchangamServer) fetchPanchangamData(ctx context.Context, date string) (*ppb.PanchangamData, error) {
    span := logging.NewSpan(ctx, "fetchPanchangamData", logrus.DebugLevel)
    defer span.End()

	span.Trace("in the database", logrus.Fields{"date": date, })
	time.Sleep(2 * time.Second)

	span.Trace("in the database, after sleep", logrus.Fields{"date": date, })
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