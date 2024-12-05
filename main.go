package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aswin-kevin/nuclei-grpc/pkg/server"
	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const listenAddress = "localhost:8555"

func main1() {
	ctx := context.Background()

	// Create nuclei engine with options
	ne, err := nuclei.NewNucleiEngineCtx(
		ctx,
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Tags: []string{"tech"},
		}), // Run critical severity templates only
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Engine created")

	defer ne.Close()

	// Load targets and optionally probe non-http/https targets
	ne.LoadTargets([]string{"https://securin.io"}, false)

	fmt.Println("Targets loaded")

	// Execute the engine with JSON output callback
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		// Print the JSON output
		fmt.Println("got results : ", event.Host, event.TemplateID, event.Type, event.Info)
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution completed")
}

func main() {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Started nuclei-api server on:", listenAddress)

	s := grpc.NewServer()
	pb.RegisterNucleiApiServer(s, &server.Server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
