package engine

import (
	"context"
	"errors"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
)

var GlobalNucleiEngine *nuclei.ThreadSafeNucleiEngine

func InitializeNucleiEngine() error {
	ctx := context.Background()

	// Create nuclei engine with options
	ne, err := nuclei.NewThreadSafeNucleiEngineCtx(ctx)

	if err != nil {
		return errors.New("Got error while creating nuclei engine :" + err.Error())
	}

	GlobalNucleiEngine = ne
	return nil
}
