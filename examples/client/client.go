package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	scanRequest = &pb.ScanRequest{
		Targets: []string{"https://hotstar.com"},
		Tags:    []string{"dns"},
	}
	listenAddress = "localhost:8555"
)

func main() {
	if len(scanRequest.Targets) == 0 {
		log.Fatal("Target is required")
		return
	}

	log.Println("Requesting scan for targets: ", scanRequest.Targets)

	conn, err := grpc.Dial(listenAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNucleiApiClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*60)
	defer cancel()
	stream, err := c.Scan(ctx, scanRequest)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.Scan failed: %v", err)
		}
		log.Println("Result :", result)
	}
}
