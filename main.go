package main

import (
	"context"
	"encoding/json"
	"fmt"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

func main() {
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

	// Set the templates directory to the user's home directory
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	panic(err)
	// }
	// templatesDir := filepath.Join(homeDir, "nuclei-templates")

	// Load targets and optionally probe non-http/https targets
	ne.LoadTargets([]string{"https://securin.io"}, false)

	fmt.Println("Targets loaded")

	// Execute the engine with JSON output callback
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		// Print the JSON output
		// fmt.Println("got results : ", event.Host, event.TemplateID, event.Type, event.Info)

		eventData, err := json.Marshal(event)
		if err != nil {
			fmt.Println("Error marshalling event:", err)
			return
		}
		fmt.Println(string(eventData))

	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution completed")
}
