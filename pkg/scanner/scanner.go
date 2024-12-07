package scanner

import (
	"context"
	"strconv"

	pb "github.com/aswin-kevin/nuclei-grpc/pkg/service"
	"github.com/aswin-kevin/nuclei-grpc/pkg/utils"
	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"github.com/rs/zerolog"
)

func eventToScanResult(event *output.ResultEvent) *pb.ScanResult {

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

func Scan(in *pb.ScanRequest, stream pb.NucleiApi_ScanServer, scanLogger *zerolog.Logger) error {
	utils.IncreaseNucleiInstanceCount()

	ctx := context.Background()

	// Create nuclei engine with options
	ne, err := nuclei.NewNucleiEngineCtx(
		ctx,
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Tags:    in.Tags,
			Authors: in.Authors,
			IDs:     in.Templates,
		}), // Run with custom template filters

	)

	if err != nil {
		scanLogger.Error().Err(err).Msg("Got error while creating nuclei engine")
		return nil
	}

	defer func() {
		utils.DecreaseNucleiInstanceCount()
		if utils.GetNucleiInstanceCount() == 0 {
			ne.Close()
			scanLogger.Info().Msg("All nuclei instances are closed")
		} else {
			scanLogger.Info().Msg("Nuclei instance is not closed due to other active scans " + strconv.Itoa(utils.GetNucleiInstanceCount()))
		}
	}()

	scanLogger.Info().Msg("New nuclei engine instance created")

	// defer ne.Close()

	// Load targets and optionally probe non-http/https targets
	ne.LoadTargets(in.Targets, false)

	scanLogger.Info().Msg("Targets are loaded into nuclei engine")

	// Execute the engine with JSON output callback
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {

		scanLogger.Info().Msg("FOUND : " + event.TemplateID)

		result := eventToScanResult(event)
		err := stream.Send(result)
		if err != nil {
			scanLogger.Error().Err(err).Msg("Error sending result to client :" + event.TemplateID)
		}
	})

	if err != nil {
		scanLogger.Error().Err(err).Msg("Error executing nuclei engine")
		return nil
	}

	scanLogger.Info().Msg("Nuclei engine scan completed")
	return nil
}
