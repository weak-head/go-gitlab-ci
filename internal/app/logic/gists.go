package logic

import (
	"errors"

	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

var (
	// ErrNoLoggerProvided happens when logger is not provided.
	ErrNoLoggerProvided = errors.New("no logger provided")

	// ErrNoReporterProvided happens when metrics reporter is not provided.
	ErrNoReporterProvided = errors.New("no metrics reporter provided")
)

// MetricsReporter defines a metrics reporter that is used
// to collect and report usage metrics.
type MetricsReporter interface {
	// TODO: define me
}

// GistsLogic implements business rules for the Gists.
type GistsLogic struct {
	// TODO: define me
}

// NewGistsLogic creates a new instance of GistsLogic that
// defines business rules and logic to handle them.
func NewGistsLogic(
	log logger.Log,
	reporter MetricsReporter,
) (*GistsLogic, error) {

	if log == nil {
		return nil, ErrNoLoggerProvided
	}

	if reporter == nil {
		return nil, ErrNoReporterProvided
	}

	return &GistsLogic{}, nil
}
