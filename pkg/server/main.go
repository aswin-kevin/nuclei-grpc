package server

import (
	"log"

	"github.com/aswin-kevin/nuclei-grpc/pkg/nuclei"
	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"
)

type Server struct {
	pb.UnimplementedNucleiApiServer
}

func (s *Server) Scan(in *pb.ScanRequest, stream pb.NucleiApi_ScanServer) error {
	log.Println("Got a request to scan: ", in.Targets)
	nuclei.Scan(in, stream)
	return nil
}
