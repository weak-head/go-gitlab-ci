package logic

import "git.lothric.net/examples/go/gogin/internal/pkg/logger"

// Reporter defines the metrics reporter abstraction that is used
// to collect and report usage metrics.
type Reporter interface {

	// ApiRequestProcessed should be called when an API request
	// has been successfully processed.
	ApiRequestProcessed(operation string, milliseconds float64)

	// ApiRequestFailed should be called when an API request
	// failed to be processed.
	ApiRequestFailed(operation string, failure string)
}

type gistsLogic struct {
}

func NewGistsLogic(
	log logger.Log,
	reporter Reporter,
) (*gistsLogic, error) {

	gl := &gistsLogic{}
	return gl, nil
}
