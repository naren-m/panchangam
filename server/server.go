package main

import (
	"context"
	"fmt"
	"net"

	ppb "github.com/naren-m/panchangam/proto/panchangam"

	"google.golang.org/grpc"
)

type PanchangamServer struct{
	ppb.UnimplementedPanchangamServer
}

func (s *PanchangamServer) Get(ctx context.Context, req *ppb.GetPanchangamRequest) (*ppb.GetPanchangamResponse, error) {
	// Get the date from the request
	requestedDate := req.Date

	// Here you would implement your logic to fetch Panchangam data for the requested date
	// For this example, we'll just return some dummy data
	panchangamData := &ppb.PanchangamData{
		Date:        requestedDate,
		Tithi:       "Some Tithi",
		Nakshatra:   "Some Nakshatra",
		Yoga:        "Some Yoga",
		Karana:      "Some Karana",
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

func main() {
	// Create a listener on TCP port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the Panchangam service with the server
	ppb.RegisterPanchangamServer(grpcServer, &PanchangamServer{})

	fmt.Println("Starting server on port :50051")
	// Start serving requests
	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Println("Failed to serve:", err)
		return
	}
}
