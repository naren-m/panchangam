// server.go
package panchangam

import (
    "context"

	ppb "github.com/naren-m/panchangam/proto/panchangam"
)

type PanchangamServer struct {
    ppb.UnimplementedPanchangamServer
}

func (s *PanchangamServer) Get(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.GetPanchangamResponse, error) {
    // Create a span for tracing the Get operation
    // Get the date from the request
    requestedDate := req.Date

    // Here you would implement your logic to fetch Panchangam data for the requested date
    // For this example, we'll just return some dummy data
    panchangamData := &ppb.PanchangamData{
        Date:        requestedDate,
        Tithi:       "Some Tithi", // Fixed placeholder for Tithi
        Nakshatra:   "Some Nakshatra",
        Yoga:        "Some Yoga",
        Karana:      "Some Karana", // Fixed placeholder for Karana
        SunriseTime: "06:00:00", // Example sunrise time
        SunsetTime:  "18:00:00", // Example sunset time
        Events: []*ppb.PanchangamEvent{
            {Name: "Some Event 1", Time: "08:00:00"},
            {Name: "Some Event 2", Time: "12:00:00"},
        },
    }

    // Create and return the response
    response := &ppb.GetPanchangamResponse{
        PanchangamData: panchangamData,
    }
    return response, nil
}
