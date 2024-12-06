package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aswin-kevin/nuclei-grpc/pkg/engine"
	"github.com/aswin-kevin/nuclei-grpc/pkg/scanner"
	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

type Server struct {
	pb.UnimplementedNucleiApiServer
}

func (s *Server) Scan(in *pb.ScanRequest, stream pb.NucleiApi_ScanServer) error {
	log.Println("Got a request to scan: ", in.Targets)

	engine.GlobalNucleiEngine.GlobalResultCallback(func(event *output.ResultEvent) {
		log.Printf("\n\nGot Result: %v\n\n", event.TemplateID)

		data, _ := json.Marshal(event)
		fmt.Println(string(data))

		result := scanner.EventToScanResult(event)
		err := stream.Send(result)
		if err != nil {
			log.Printf("Error sending %v result to client: %v", event.TemplateID, err)
		}
	})

	scanner.Scan(in, stream)
	return nil
}
