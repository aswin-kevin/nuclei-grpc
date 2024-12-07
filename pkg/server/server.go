package server

import (
	"github.com/aswin-kevin/nuclei-grpc/pkg/logger"
	"github.com/aswin-kevin/nuclei-grpc/pkg/scanner"
	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"
	"github.com/aswin-kevin/nuclei-grpc/pkg/utils"
)

type Server struct {
	pb.UnimplementedNucleiApiServer
}

func (s *Server) Scan(in *pb.ScanRequest, stream pb.NucleiApi_ScanServer) error {
	scanId, _ := utils.GenerateUUID()

	logger.GlobalLogger.Info().Msg("Received a request to scan : " + scanId)

	// creating a sub logger for the scan
	scanLogger := logger.GlobalLogger.With().Str("SCAN ID", scanId).Logger()

	if len(in.Targets) == 0 {
		scanLogger.Error().Msg("No targets provided")
		scanLogger.Error().Msg("Closing the stream")
		return nil
	}

	scanLogger.Info().Msg("Starting scan")

	scanner.Scan(in, stream, &scanLogger)
	return nil
}
