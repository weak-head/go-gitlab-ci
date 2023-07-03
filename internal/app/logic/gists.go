package logic

import (
	"context"
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
	// NOTE: define me
}

// GistsLogic implements business rules for the Gists.
type GistsLogic struct {
	log      logger.Log
	reporter MetricsReporter
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

	return &GistsLogic{
		log:      log,
		reporter: reporter,
	}, nil
}

// GetGists returns the filtered list of Gists for the specified 'language'.
func (g *GistsLogic) GetGists(ctx context.Context, language string) error {
	log := logger.FromContext(g.log, ctx, "GetGists")
	log.Info("Handling GetGists")

	// NOTE: implement me

	return nil
}
