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

	var nucleiScanStrategy = "template-spray"
	var templateIdsToScan = make([]string, 0)

	var nucleiConcurrencyConfig = nuclei.Concurrency{
		TemplateConcurrency:           10, // number of templates to run concurrently (per host in host-spray mode)
		HostConcurrency:               5,  // number of hosts to scan concurrently (per template in template-spray mode)
		HeadlessHostConcurrency:       3,  // number of hosts to scan concurrently for headless templates (per template in template-spray mode)
		HeadlessTemplateConcurrency:   2,  // number of templates to run concurrently for headless templates (per host in host-spray mode)
		JavascriptTemplateConcurrency: 4,  // number of templates to run concurrently for javascript templates (per host in host-spray mode)
		TemplatePayloadConcurrency:    25, // max concurrent payloads to run for a template
		ProbeConcurrency:              50, // max concurrent HTTP probes to run
	}

	if in.ScanStrategy != "" {
		nucleiScanStrategy = in.ScanStrategy
	}

	if in.ScanConcurrencyConfig.TemplateConcurrency > 0 {
		nucleiConcurrencyConfig.TemplateConcurrency = int(in.ScanConcurrencyConfig.TemplateConcurrency)
	}

	if in.ScanConcurrencyConfig.HostConcurrency > 0 {
		nucleiConcurrencyConfig.HostConcurrency = int(in.ScanConcurrencyConfig.HostConcurrency)
	}

	if in.ScanConcurrencyConfig.HeadlessHostConcurrency > 0 {
		nucleiConcurrencyConfig.HeadlessHostConcurrency = int(in.ScanConcurrencyConfig.HeadlessHostConcurrency)
	}

	if in.ScanConcurrencyConfig.HeadlessTemplateConcurrency > 0 {
		nucleiConcurrencyConfig.HeadlessTemplateConcurrency = int(in.ScanConcurrencyConfig.HeadlessTemplateConcurrency)
	}

	if in.ScanConcurrencyConfig.JavascriptTemplateConcurrency > 0 {
		nucleiConcurrencyConfig.JavascriptTemplateConcurrency = int(in.ScanConcurrencyConfig.JavascriptTemplateConcurrency)
	}

	if in.ScanConcurrencyConfig.TemplatePayloadConcurrency > 0 {
		nucleiConcurrencyConfig.TemplatePayloadConcurrency = int(in.ScanConcurrencyConfig.TemplatePayloadConcurrency)
	}

	if in.ScanConcurrencyConfig.ProbeConcurrency > 0 {
		nucleiConcurrencyConfig.ProbeConcurrency = int(in.ScanConcurrencyConfig.ProbeConcurrency)
	}

	// If templates ids are provided, use them
	if len(in.TemplateIds) > 0 {
		templateIdsToScan = in.TemplateIds
	}

	if len(in.Templates) > 0 {
		userGivenTemplateIds := utils.GetTemplateIdsFromTemplateData(in.Templates)
		templateIdsToScan = append(templateIdsToScan, userGivenTemplateIds...)
	}

	scanLogger.Info().Msg("Templates to scan : " + strconv.Itoa(len(templateIdsToScan)))

	// Create nuclei engine with options
	ne, err := nuclei.NewNucleiEngineCtx(
		ctx,
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Severity: in.Severity,
			Tags:     in.Tags,
			Authors:  in.Authors,
			IDs:      templateIdsToScan,
		}), // Run with custom template filters
		nuclei.WithConcurrency(nucleiConcurrencyConfig), // Set concurrency
		nuclei.WithScanStrategy(nucleiScanStrategy),     // Set scan strategy
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
