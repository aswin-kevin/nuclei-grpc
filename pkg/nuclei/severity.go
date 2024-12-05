package nuclei

import (
	"errors"
	"strings"

	"github.com/projectdiscovery/nuclei/v3/pkg/model/types/severity"
)

var severityMappings = map[severity.Severity]string{
	severity.Info:     "info",
	severity.Low:      "low",
	severity.Medium:   "medium",
	severity.High:     "high",
	severity.Critical: "critical",
	severity.Unknown:  "unknown",
}

func normalizeSeverityValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func stringToSeverity(severityString string) (severity.Severity, error) {
	normalizedValue := normalizeSeverityValue(severityString)
	for key, currentValue := range severityMappings {
		if normalizedValue == currentValue {
			return key, nil
		}
	}
	return -1, errors.New("Invalid severity: " + severityString)
}
