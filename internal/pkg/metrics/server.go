package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config
type Config struct {
	// Addr
	Addr string

	// Path
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
