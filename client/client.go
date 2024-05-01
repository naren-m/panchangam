package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	ppb "github.com/naren-m/panchangam/proto/panchangam"

)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a client instance
	client := ppb.NewPanchangamClient(conn)

	// Create a request
	request := &ppb.GetPanchangamRequest{
		Date: "2024-04-30", // Example date
	}

	// Call the RPC method
	response, err := client.Get(context.Background(), request)
	if err != nil {
		log.Fatalf("Error calling Get: %v", err)
	}

	// Process the response
	panchangamData := response.GetPanchangamData()
	fmt.Println("Panchangam Data:")
	fmt.Printf("Date: %s\n", panchangamData.GetDate())
	fmt.Printf("Date: %s\n", panchangamData.GetTithi())
	fmt.Printf("Date: %s\n", panchangamData.GetYoga())
	fmt.Printf("Date: %s\n", panchangamData.GetNakshatra())
	// Print other fields as needed
}
