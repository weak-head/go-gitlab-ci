package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// TODO: comments
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of processed API requests.",
		},
		[]string{"operation"},
	)

	// TODO: comments
	requestsFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_errors_total",
			Help: "Total number of API request errors.",
		},
		[]string{"operation", "failure"},
	)

	// TODO: comments
	requestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "request_durations_seconds",
			Help:       "API request processing duration distributions.",
			Objectives: map[float64]float64{},
		},
		[]string{"operation"},
	)

	// TODO: comments
	requestDurationsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_durations_histogram_seconds",
			Help: "API request processing duration distributions.",
			// Buckets: prometheus.LinearBuckets(normMean-5*normDomain, .5*normDomain, 20),
		},
		[]string{"operation"},
	)
)
