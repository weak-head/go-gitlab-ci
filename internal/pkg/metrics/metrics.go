package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (

	// requestsTotal is a total number of processed API requests.
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of processed API requests.",
		},
		[]string{"operation"},
	)

	// requestsFailures is a total number of API request errors.
	requestsFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_errors_total",
			Help: "Total number of API request errors.",
		},
		[]string{"operation", "failure"},
	)

	// requestDuration is API request processing duration distributions.
	requestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "request_durations_seconds",
			Help:       "API request processing duration distributions.",
			Objectives: map[float64]float64{},
		},
		[]string{"operation"},
	)

	// requestDurationsHistogram is API request processing duration distributions.
	requestDurationsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_durations_histogram_seconds",
			Help: "API request processing duration distributions.",
			// Start at 10 milliseconds, add 20 buckets, 10 milliseconds each
			Buckets: prometheus.LinearBuckets(0.01, 0.01, 20),
		},
		[]string{"operation"},
	)
)
