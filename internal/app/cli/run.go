package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"git.lothric.net/examples/go/gogin/internal/app/api"
	"git.lothric.net/examples/go/gogin/internal/app/components"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
	"git.lothric.net/examples/go/gogin/internal/pkg/metrics"
	"git.lothric.net/examples/go/gogin/internal/pkg/status"
)

// appConfig defines the global application configuration and
// contains all parameters that could be specified during the
// application startup.
type appConfig struct {
	NodeName    string
	ServiceName string

	Http    httpConfig
	Log     logger.Config
	Status  status.Config
	Metrics metrics.Config
}

// httpConfig defines HTTP API server configuration
type httpConfig struct {

	// HttpPort is port that is used to host HTTP API server
	HttpPort uint16

	// Gin mode:
	//  - test
	//  - debug
	//  - release
	GinMode string
}

// runApp bootstraps and runs the application
func runApp(config appConfig) error {

	cleanLog, err := logger.New(config.Log)
	if err != nil {
		fmt.Printf("Failed to create the application logger: %s", err.Error())
		return err
	}

	log := cleanLog.WithFields(logger.Fields{
		logger.FieldNode:     config.NodeName,
		logger.FieldService:  config.ServiceName,
		logger.FieldPackage:  "cli",
		logger.FieldFunction: "runApp",
	})

	gin.SetMode(config.Http.GinMode)

	// --------------
	// Health status server
	statusServer, err := status.NewStatusServer()
	if err != nil {
		log.Error(err, "Failed to create the status server.")
		return err
	}

	go func(cfg status.Config) {
		statusServer.Serve(cfg)
	}(config.Status)

	// --------------
	// Prometheus metrics server
	metricsServer, err := metrics.NewPrometheusServer(config.Metrics)
	if err != nil {
		log.Error(err, "Failed to create the metrics server.")
		return err
	}

	go func() {
		metricsServer.Serve()
	}()

	// --------------
	// HTTP API server
	ctx := context.Background()

	componentFactory, err := components.NewComponentFactory(log)
	if err != nil {
		log.Error(err, "Failed to create component factory")
		return err
	}

	apiBuilder, err := api.NewApiBuilder(log, componentFactory)
	if err != nil {
		log.Error(err, "Failed to create HTTP API server.")
		return err
	}

	router, err := apiBuilder.BuildApi(ctx)
	if err != nil {
		log.Error(err, "Failed to Build API router.")
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Http.HttpPort),
		Handler: router,
	}

	go func() {
		// Serve HTTP API endpoints
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err, "Failed to listen and serve HTTP API endpoint.")
		}
	}()

	log.Info("The service is up and running.")

	// -------------------------
	// Wait for the app termination

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Requested app termination. Waiting 3 seconds...")

	// Wait 3 seconds for server to terminate
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Status
	log.Info("Gracefully terminating status server...")
	statusServer.Stop()

	// Metrics
	log.Info("Gracefully terminating metrics server...")
	if err := metricsServer.Stop(ctx); err != nil {
		log.Error(err, "Failed to gracefully shutdown prometheus metrics server.")
	}

	// HTTP API
	log.Info("Gracefully terminating API server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error(err, "Failed to gracefully shutdown HTTP API server.")
	}

	// Wait
	<-ctx.Done()
	log.Info("Timeout of 3 seconds has ended. Exiting.")

	return nil
}
