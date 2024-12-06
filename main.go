package main

import (
	"log"
	"net"

	"github.com/aswin-kevin/nuclei-grpc/pkg/scanner"
	"github.com/aswin-kevin/nuclei-grpc/pkg/server"
	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const listenAddress = "localhost:8555"

func init() {
	engineErr := scanner.InitializeNucleiEngine()
	if engineErr != nil {
		log.Fatalf(engineErr.Error())
	}
	log.Println("Initialized nuclei engine")
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
		scanner.GlobalNucleiEngine.Close()
		log.Fatalf("failed to serve: %v", err)
	}
}
