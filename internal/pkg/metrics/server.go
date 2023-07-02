package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config is a prometheus server configuration.
type Config struct {
	// Addr is endpoint of the prometheus metrics server.
	// For example: ":8880"
	Addr string

	// Path is route that the metrics should be exposed on.
	// For example: "/metrics"
	Path string
}

// prometheusServer is Prometheus HTTP server for metrics collection.
type prometheusServer struct {
	server   *http.Server
	registry *prometheus.Registry
	conf     Config
}

// NewPrometheusServer creates a new instance of Prometheus HTTP server.
func NewPrometheusServer(conf Config) (*prometheusServer, error) {

	// Prometheus HTTP Server with metrics registry
	p := &prometheusServer{
		registry: prometheus.NewRegistry(),
		conf:     conf,
	}

	// Register all defined metrics
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

	// Report metrics on the specified route
	mux := http.NewServeMux()
	mux.Handle(
		p.conf.Path,
		promhttp.HandlerFor(
			p.registry,
			promhttp.HandlerOpts{EnableOpenMetrics: true}),
	)

	p.server = &http.Server{
		Addr:    p.conf.Addr,
		Handler: mux,
	}

	return p, nil
}

// Serve starts prometheus HTTP server.
func (p *prometheusServer) Serve() error {
	if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops prometheus HTTP server.
func (p *prometheusServer) Stop(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}
