package status

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	gh "google.golang.org/grpc/health/grpc_health_v1"
)

// Config defines the configuration for status server.
type Config struct {

	// RpcAddr is the endpoint the status server is listening.
	RpcAddr string
}

// statusServer is grpc-based status server that exposes
// service status via grpc on a particular port and is expected
// to be called by `grpc_health_probe`.
type statusServer struct {
	server *grpc.Server
}

// NewStatusServer creates a new health check status server
// that integrates with `grpc_health_probe` and responds on
// readiness and liveliness health probe checks.
func NewStatusServer() (*statusServer, error) {

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", gh.HealthCheckResponse_SERVING)

	grpcServer := grpc.NewServer()
	gh.RegisterHealthServer(grpcServer, healthServer)

	return &statusServer{
		server: grpcServer,
	}, nil
}

// Serve bootstrap the status server.
func (s *statusServer) Serve(config Config) error {
	ln, err := net.Listen("tcp", config.RpcAddr)
	if err != nil {
		return err
	}

	if err := s.server.Serve(ln); err != nil {
		return err
	}

	return nil
}

// Stop terminates the status server.
func (s *statusServer) Stop() {
	s.server.Stop()
}
