package cli

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// Service
	serviceName = "gogin"

	// Common
	configFile = "config"
	nodeName   = "node.name"

	// HTTP & Gin
	httpPort = "http.port"
	ginMode  = "http.gin.mode"

	// Logger
	logLevel     = "log.level"
	logFormatter = "log.formatter"

	// Health check
	statusRpcAddr = "status.rpc.addr"

	// Metrics
	metricsPrometheusAddr = "metrics.prometheus.addr"
	metricsPrometheusPath = "metrics.prometheus.path"
)

// cli
type cli struct {
	cfg appConfig
}

// NewCli creates a new Cobra Command
func NewCli() (*cobra.Command, error) {
	cli := &cli{}
	cmd := &cobra.Command{
		Use:     serviceName,
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return cmd, nil
}

// setupFlags configures the acceptable command line arguments
func setupFlags(cmd *cobra.Command) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// Common
	cmd.Flags().String(configFile, "", "Path to config file.")
	cmd.Flags().String(nodeName, hostname, "Unique server ID.")

	// HTTP & Gin
	cmd.Flags().String(httpPort, "8080", "HTTP API port.")
	cmd.Flags().String(ginMode, "release", "Gin mode.")

	// Log
	cmd.Flags().String(logLevel, "info", "Log level.")
	cmd.Flags().String(logFormatter, "json", "Log formatter.")

	// Status
	cmd.Flags().String(statusRpcAddr, ":8400", "Rpc address of status server.")

	// Metrics
	cmd.Flags().String(metricsPrometheusAddr, ":8880", "HTTP address of prometheus metrics endpoint.")
	cmd.Flags().String(metricsPrometheusPath, "/metrics", "HTTP URL endpoint of prometheus metrics endpoint.")

	return viper.BindPFlags(cmd.Flags())
}

// setupConfig reads the config file (if any), binds the environment
// variables to the viper keys and creates the application configuration
func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {

	// Try to load config from file
	if err := c.loadConfigFile(cmd); err != nil {
		fmt.Printf("Failed to load configuration file: %s \n", err.Error())
		return err
	}

	// Try to load config from environment variables
	if err := c.bindEnv(); err != nil {
		fmt.Printf("Failed to bind environment variables: %s \n", err.Error())
		return err
	}

	// Create and initialize application configuration
	if err := c.createConfig(); err != nil {
		fmt.Printf("Failed to create application configuration: %s \n", err.Error())
		return err
	}

	return nil
}

// loadConfigFile loads configuration from file,
// that could be specified by '--config' cli option.
func (c *cli) loadConfigFile(cmd *cobra.Command) error {
	file, err := cmd.Flags().GetString(configFile)
	if err != nil {
		return err
	}

	// Load config from the specified file
	if file != "" {
		viper.SetConfigFile(file)
		if err = viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return err
			}
		}
		return nil
	}

	// Search for in-place config and load it, if possible.
	// Ignore failed to find configuration file.
	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		fmt.Printf("Fatal error with configuration file: %s \n", err.Error())
		return err
	}

	// We can enable watchdog for config, if required:
	// viper.WatchConfig()

	return nil
}

// bindEnv binds configuration keys to environment variables.
func (c *cli) bindEnv() error {

	// HTTP & Gin
	viper.BindEnv(httpPort, "HTTP_PORT")
	viper.BindEnv(ginMode, "HTTP_GIN_MODE")

	// Log
	viper.BindEnv(logLevel, "LOG_LEVEL")
	viper.BindEnv(logFormatter, "LOG_FORMATTER")

	// Status
	viper.BindEnv(statusRpcAddr, "STATUS_RPC_ADDR")

	// Metrics
	viper.BindEnv(metricsPrometheusAddr, "METRICS_PROMETHEUS_ADDR")
	viper.BindEnv(metricsPrometheusPath, "METRICS_PROMETHEUS_PATH")

	return nil
}

// createConfig initializes application configuration
func (c *cli) createConfig() error {
	config := &c.cfg

	// Generic
	config.NodeName = viper.GetString(nodeName)
	config.ServiceName = serviceName

	// HTTP & Gin
	httpConfig := &config.Http
	httpConfig.HttpPort = viper.GetUint16(httpPort)
	httpConfig.GinMode = viper.GetString(ginMode)

	// Log
	logConfig := &config.Log
	logConfig.Level = viper.GetString(logLevel)
	logConfig.Formatter = viper.GetString(logFormatter)

	// Status
	statusConfig := &config.Status
	statusConfig.RpcAddr = viper.GetString(statusRpcAddr)

	// Metrics
	metricsConfig := &config.Metrics
	metricsConfig.Addr = viper.GetString(metricsPrometheusAddr)
	metricsConfig.Path = viper.GetString(metricsPrometheusPath)

	return nil
}

// run bootstraps and runs the application
func (c *cli) run(cmd *cobra.Command, args []string) error {

	// The configuration is already fully loaded an initialized
	return runApp(c.cfg)
}
