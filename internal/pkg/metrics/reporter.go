package metrics

import (
	"errors"

	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

var (
	// ErrNoLoggerProvided happens when logger is not provided.
	ErrNoLoggerProvided = errors.New("no logger provided")
)

// reporter collects and reports application metrics and usage statistics.
type reporter struct {
	log logger.Log
}

// NewReporter creates a new instance of metrics reporter.
func NewReporter(log logger.Log) (*reporter, error) {
	if log == nil {
		return nil, ErrNoLoggerProvided
	}

	return &reporter{
		log: log,
	}, nil
}

// ApiRequestProcessed tracks a successfully processed API request.
func (r *reporter) ApiRequestProcessed(operation string, milliseconds float64) {
	log := r.log.WithField(logger.FieldFunction, "ApiRequestProcessed")
	log.Info("Api request has been processed")

	requestDuration.WithLabelValues(operation).Observe(milliseconds)
	requestDurationsHistogram.WithLabelValues(operation).Observe(milliseconds)
	requestsTotal.WithLabelValues(operation).Inc()
}

// ApiRequestFailed tracks an API request that has failed to be processed.
func (r *reporter) ApiRequestFailed(operation string, failure string) {
	log := r.log.WithField(logger.FieldFunction, "ApiRequestProcessed")
	log.Info("Api request processing has failed")

	requestsTotal.WithLabelValues(operation).Inc()
	requestsFailures.WithLabelValues(operation, failure).Inc()
}
