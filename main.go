package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/aswin-kevin/nuclei-grpc/pkg/logger"
	"github.com/aswin-kevin/nuclei-grpc/pkg/server"
	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"
	"github.com/aswin-kevin/nuclei-grpc/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"os"

	"github.com/spf13/cobra"
)

var (
	address         string
	port            string
	healthCheckPort string
)

const (
	defaultAddress         = "localhost"
	defaultPort            = "8555"
	defaultHealthCheckPort = "8556"
)

var rootCmd = &cobra.Command{
	Use:   "nuclei-grpc",
	Short: "Nuclei gRPC server",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		startHealthCheckServer()
		startServer()
	},
}

func init() {
	logger.InitializeGlobalLogger()
	utils.LoadAllNucleiTemplatesMetadata()
	startCmd.Flags().StringVarP(&address, "address", "a", defaultAddress, "Address to listen on")
	startCmd.Flags().StringVarP(&port, "port", "p", defaultPort, "Port to listen on")
	startCmd.Flags().StringVarP(&healthCheckPort, "hport", "r", defaultHealthCheckPort, "REST API Port to listen on - health check")
	rootCmd.AddCommand(startCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}

func startServer() {
	listenAddress := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		logger.GlobalLogger.Fatal().Err(err).Msg("Failed to listen -> Closing server")
	}
	logger.GlobalLogger.Info().Msg("Started nuclei-api server on: " + listenAddress)

	s := grpc.NewServer()
	pb.RegisterNucleiApiServer(s, &server.Server{})

	// it gives server metadata to client
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		logger.GlobalLogger.Fatal().Err(err).Msg("Failed to serve -> Closing server")
	}
}

func startHealthCheckServer() {
	restAPIListenAddress := fmt.Sprintf("%s:%s", address, defaultHealthCheckPort)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("I AM HEALTHY"))
	})

	go func() {
		logger.GlobalLogger.Info().Msg("Started REST API health check server on : " + restAPIListenAddress)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", defaultHealthCheckPort), nil); err != nil {
			logger.GlobalLogger.Fatal().Err(err).Msg("Failed to start REST API health check server")
		}
	}()
}
