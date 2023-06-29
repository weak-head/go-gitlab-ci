package status

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Config
type Config struct {
	RpcAddr string
}

// statusServer
type statusServer struct {
	grpcServer *grpc.Server
}

// NewStatusServer
func NewStatusServer() (*statusServer, error) {
	grpcServer := grpc.NewServer()

	health_service := health.NewServer()
	health_service.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, health_service)

	return &statusServer{
		grpcServer: grpcServer,
	}, nil
}

// Serve
func (s *statusServer) Serve(config Config) error {
	ln, err := net.Listen("tcp", config.RpcAddr)
	if err != nil {
		return err
	}

	if err := s.grpcServer.Serve(ln); err != nil {
		return err
	}

	return nil
}

// Stop
func (s *statusServer) Stop() {
	s.grpcServer.Stop()
}
