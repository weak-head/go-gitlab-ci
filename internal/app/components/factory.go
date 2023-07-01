package components

import (
	"git.lothric.net/examples/go/gogin/internal/app/api/v1/handlers"
	"git.lothric.net/examples/go/gogin/internal/app/logic"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
	"git.lothric.net/examples/go/gogin/internal/pkg/metrics"
)

// componentFactory
type componentFactory struct {
	log logger.Log
}

// NewComponentFactory
func NewComponentFactory(log logger.Log) (*componentFactory, error) {
	return &componentFactory{
		log: log,
	}, nil
}

// CreateApiMetricsReporter
func (f *componentFactory) CreateApiMetricsReporter() (handlers.ApiMetricsReporter, error) {
	return metrics.NewReporter()
}

// CreateGistsLogic
func (f *componentFactory) CreateGistsLogic() (handlers.GistsLogic, error) {

	reporter, err := metrics.NewReporter()
	if err != nil {
		f.log.Error(err, "Failed to create Metrics Reporter")
		return nil, err
	}

	return logic.NewGistsLogic(f.log, reporter)
}
