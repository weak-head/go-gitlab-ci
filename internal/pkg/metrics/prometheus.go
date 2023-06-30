package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
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
)

// Config
type Config struct {
	Addr string
	Path string
}

// prometheusServer
type prometheusServer struct {
	server   *http.Server
	registry *prometheus.Registry
	conf     Config
}

// NewPrometheusServer
func NewPrometheusServer(conf Config) (*prometheusServer, error) {
	p := &prometheusServer{
		registry: prometheus.NewRegistry(),
		conf:     conf,
	}

	for _, c := range []prometheus.Collector{
		requestDuration,
		requestDurationsHistogram,
		requestsTotal,
		requestsFailures,
		collectors.NewBuildInfoCollector(),
	} {
		if err := p.registry.Register(c); err != nil {
			return nil, err
		}
	}

	p.server = &http.Server{
		Addr: p.conf.Addr,
		Handler: promhttp.HandlerFor(
			p.registry,
			promhttp.HandlerOpts{EnableOpenMetrics: true},
		),
	}

	return p, nil
}

// Serve
func (p *prometheusServer) Serve() error {
	if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop
func (p *prometheusServer) Stop(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

// reporter
type Reporter struct {
}

// NewReporter
func NewReporter() (*Reporter, error) {
	return &Reporter{}, nil
}

// ApiRequestProcessed
func (r *Reporter) ApiRequestProcessed(operation string, milliseconds float64) {
	requestDuration.WithLabelValues(operation).Observe(milliseconds)
	requestDurationsHistogram.WithLabelValues(operation).Observe(milliseconds)
	requestsTotal.WithLabelValues(operation).Inc()
}

// ApiRequestFailed
func (r *Reporter) ApiRequestFailed(operation string, failure string) {
	requestsTotal.WithLabelValues(operation).Inc()
	requestsFailures.WithLabelValues(operation, failure).Inc()
}
