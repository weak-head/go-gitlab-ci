package logic

import "git.lothric.net/examples/go/gogin/internal/pkg/logger"

// MetricsReporter defines the metrics reporter abstraction that is used
// to collect and report usage metrics.
type MetricsReporter interface {
}

type GistsLogic struct {
}

func NewGistsLogic(
	log logger.Log,
	reporter MetricsReporter,
) (*GistsLogic, error) {

	gl := &GistsLogic{}
	return gl, nil
}
