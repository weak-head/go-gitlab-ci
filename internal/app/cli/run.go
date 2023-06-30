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
	"github.com/spf13/viper"

	"git.lothric.net/examples/go/gogin/internal/app/api"
	"git.lothric.net/examples/go/gogin/internal/app/components"
	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
	"git.lothric.net/examples/go/gogin/internal/pkg/metrics"
	"git.lothric.net/examples/go/gogin/internal/pkg/status"
)

// cfg is application configuration
type cfg struct {
	HttpPort string
	GinMode  string

	LogConfig     logger.Config
	StatusConfig  status.Config
	MetricsConfig metrics.Config
}

// runApp bootstraps and runs the application
func runApp(config cfg) error {

	var log logger.Log
	log, err := logger.New(config.LogConfig)
	if err != nil {
		fmt.Printf("Failed to create the application logger: %s", err.Error())
		return err
	}

	appLog := log.WithFields(logger.Fields{
		logger.FieldNode:    viper.GetString(nodeName),
		logger.FieldService: "gogin",
	})

	log = appLog.WithFields(logger.Fields{
		logger.FieldPackage:  "main",
		logger.FieldFunction: "cli.run",
	})

	if config.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// --------------
	// Health status server
	statusServer, err := status.NewStatusServer()
	if err != nil {
		log.Error(err, "Failed to create the status server.")
		return err
	}

	go func(cfg status.Config) {
		statusServer.Serve(cfg)
	}(config.StatusConfig)

	// --------------
	// Prometheus metrics server
	metricsServer, err := metrics.NewPrometheusServer(config.MetricsConfig)
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

	componentFactory, err := components.NewComponentFactory()
	if err != nil {
		log.Error(err, "Failed to create component factory")
		return err
	}

	router, err := api.BuildApiEngine(ctx, log, componentFactory)
	if err != nil {
		log.Error(err, "Failed to create HTTP API server.")
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", viper.GetString(httpPort)),
		Handler: router,
	}

	go func() {
		// Serve HTTP API endpoints
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err, "Failed to listen and serve HTTP API endpoint.")
		}
	}()

	log.Info("The gogin micro-service is up and running.")

	// -------------------------
	// Wait for the app termination

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Requested app termination. Waiting 5 seconds...")

	// Wait 5 seconds for server to terminate
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Status
	statusServer.Stop()

	// Metrics
	if err := metricsServer.Stop(ctx); err != nil {
		log.Error(err, "Failed to gracefully shutdown prometheus metrics server.")
	}

	// HTTP API
	if err := srv.Shutdown(ctx); err != nil {
		log.Error(err, "Failed to gracefully shutdown HTTP API server.")
	}

	// Wait
	<-ctx.Done()
	log.Info("Timeout of 5 seconds has ended. Gogin exiting.")

	return nil
}
