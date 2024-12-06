package scanner

import (
	"log"

	"github.com/aswin-kevin/nuclei-grpc/pkg/engine"
	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"github.com/projectdiscovery/nuclei/v3/pkg/testutils"
)

func EventToScanResult(event *output.ResultEvent) *pb.ScanResult {

	var info *pb.ScanResultInfo
	if event.Info.Classification != nil {

		info = &pb.ScanResultInfo{
			Name:        event.Info.Name,
			Description: event.Info.Description,
			Severity:    event.Info.SeverityHolder.Severity.String(),
			Remediation: event.Info.Remediation,
			Tags:        ToSliceSafe(event.Info.Tags),
			References:  make([]string, 0),
			Classification: &pb.ScanResultClassification{
				Cves:       ToSliceSafe(event.Info.Classification.CVEID),
				Cwes:       ToSliceSafe(event.Info.Classification.CWEID),
				Cpe:        event.Info.Classification.CPE,
				CvssVector: event.Info.Classification.CVSSMetrics,
				CvssScore:  event.Info.Classification.CVSSScore,
			},
		}

		if event.Info.Reference != nil {
			info.References = event.Info.Reference.ToSlice()
		}

	}

	var interaction *pb.Interaction
	if event.Interaction != nil {
		interaction = &pb.Interaction{
			Protocol:      event.Interaction.Protocol,
			UniqueId:      event.Interaction.UniqueID,
			FullId:        event.Interaction.FullId,
			Qtype:         event.Interaction.QType,
			RawRequest:    []byte(event.Interaction.RawRequest),
			RawResponse:   []byte(event.Interaction.RawResponse),
			SmtpFrom:      event.Interaction.SMTPFrom,
			RemoteAddress: event.Interaction.RemoteAddress,
			Timestamp:     event.Interaction.Timestamp.String(),
		}
	}

	return &pb.ScanResult{
		TemplateId:       event.TemplateID,
		Template:         event.Template,
		Info:             info,
		MatcherName:      event.MatcherName,
		ExtractorName:    event.ExtractorName,
		Type:             event.Type,
		Host:             event.Host,
		Path:             event.Path,
		Matched:          event.Matched,
		ExtractedResults: event.ExtractedResults,
		Request:          []byte(event.Request),
		Response:         []byte(event.Response),
		Ip:               event.IP,
		Timestamp:        event.Timestamp.String(),
		CurlCommand:      event.CURLCommand,
		MatcherStatus:    event.MatcherStatus,
		Interaction:      interaction,
	}
}

func ToSliceSafe(i interface{}) []string {
	if i == nil {
		return make([]string, 0)
	}
	return i.(interface{ ToSlice() []string }).ToSlice()
}

func Scan(in *pb.ScanRequest, stream pb.NucleiApi_ScanServer) error {

	// mockProgress := &testutils.MockProgressClient{}
	// reportingClient, _ := reporting.New(&reporting.Options{}, "")
	// defer reportingClient.Close()

	outputWriter := testutils.NewMockOutputWriter(false)
	outputWriter.WriteCallback = func(event *output.ResultEvent) {
		log.Printf("Got Result: %v\n", event.TemplateID)
		result := EventToScanResult(event)
		err := stream.Send(result)
		if err != nil {
			log.Printf("Error sending %v result to client: %v", event.TemplateID, err)
		}
	}

	// interactOpts := interactsh.DefaultOptions(outputWriter, reportingClient, mockProgress)

	// nucleiInteractOpts := nuclei.InteractshOpts{}
	// nucleiInteractOpts.Output = outputWriter

	// Load targets and execute nuclei engine
	err := engine.GlobalNucleiEngine.ExecuteNucleiWithOpts(
		in.Targets,
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Tags: in.Tags,
		}),
	)

	log.Println("Targets loaded")

	// GlobalNucleiEngine.GlobalResultCallback(func(event *output.ResultEvent) {
	// 	log.Printf("\n\nGot Result: %v\n\n", event.TemplateID)
	// 	result := eventToScanResult(event)
	// 	err := stream.Send(result)
	// 	if err != nil {
	// 		log.Printf("Error sending %v result to client: %v", event.TemplateID, err)
	// 	}
	// })

	if err != nil {
		log.Println("Error executing nuclei engine: ", err)
		return nil
	}

	return nil
}
