package components

import (
	"errors"

	"git.lothric.net/examples/go/gogin/internal/app/api/v1/handlers"
	"git.lothric.net/examples/go/gogin/internal/app/logic"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
	"git.lothric.net/examples/go/gogin/internal/pkg/metrics"
)

var (
	// ErrNoLoggerProvided happens when logger is not provided.
	ErrNoLoggerProvided = errors.New("no logger provided")
)

// componentFactory is a factory that creates components that are required
// for construction of API Handlers.
type componentFactory struct {
	log logger.Log
}

// NewComponentFactory creates a new instance of the component factory.
func NewComponentFactory(
	log logger.Log,
) (*componentFactory, error) {
	if log == nil {
		return nil, ErrNoLoggerProvided
	}

	return &componentFactory{
		log: log,
	}, nil
}

// CreateApiMetricsReporter create a new API Metrics Reporter.
func (f *componentFactory) CreateApiMetricsReporter() (handlers.ApiMetricsReporter, error) {
	log := f.log.WithField(logger.FieldFunction, "CreateApiMetricsReporter")
	log.Info("Creating Api Metrics Reporter")

	return metrics.NewReporter(log)
}

// CreateGistsLogic creates a business logic for Gists.
func (f *componentFactory) CreateGistsLogic() (handlers.GistsLogic, error) {
	log := f.log.WithField(logger.FieldFunction, "CreateGistsLogic")
	log.Info("Creating Gists logic")

	reporter, err := metrics.NewReporter(log)
	if err != nil {
		log.Error(err, "Failed to create business logic Metrics Reporter")
		return nil, err
	}

	return logic.NewGistsLogic(log, reporter)
}
