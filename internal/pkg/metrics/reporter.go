package metrics

import (
	"errors"

	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

var (
	// ErrNoLoggerProvided happens when logger is not provided.
	ErrNoLoggerProvided = errors.New("no logger provided")
)

// Reporter
type Reporter struct {
	log logger.Log
}

// NewReporter
func NewReporter(log logger.Log) (*Reporter, error) {
	if log == nil {
		return nil, ErrNoLoggerProvided
	}

	return &Reporter{
		log: log,
	}, nil
}

// ApiRequestProcessed
func (r *Reporter) ApiRequestProcessed(operation string, milliseconds float64) {
	log := r.log.WithField(logger.FieldFunction, "ApiRequestProcessed")
	log.Info("Api request has been processed")

	requestDuration.WithLabelValues(operation).Observe(milliseconds)
	requestDurationsHistogram.WithLabelValues(operation).Observe(milliseconds)
	requestsTotal.WithLabelValues(operation).Inc()
}

// ApiRequestFailed
func (r *Reporter) ApiRequestFailed(operation string, failure string) {
	log := r.log.WithField(logger.FieldFunction, "ApiRequestProcessed")
	log.Info("Api request processing has failed")

	requestsTotal.WithLabelValues(operation).Inc()
	requestsFailures.WithLabelValues(operation, failure).Inc()
}
