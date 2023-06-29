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
	cfg cfg
}

// NewCli
func NewCli() (*cobra.Command, error) {
	cli := &cli{}
	cmd := &cobra.Command{
		Use:     "gogin",
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
	if err := c.loadConfigFile(cmd); err != nil {
		fmt.Printf("Failed to load configuration file: %s \n", err.Error())
		return err
	}

	if err := c.bindEnv(); err != nil {
		fmt.Printf("Failed to bind environment variables: %s \n", err.Error())
		return err
	}

	if err := c.createConfig(); err != nil {
		fmt.Printf("Failed to create application configuration: %s \n", err.Error())
		return err
	}

	return nil
}

// loadConfigFile
func (c *cli) loadConfigFile(cmd *cobra.Command) error {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err = viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return err
			}
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		// Ignore failed to find configuration file
		if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			fmt.Printf("Fatal error with configuration file: %s \n", err.Error())
			return err
		}
	}
	// viper.WatchConfig()

	return nil
}

// bindEnv
func (c *cli) bindEnv() error {

	// HTTP & Gin
	viper.BindEnv(httpPort, "HTTP_PORT")
	viper.BindEnv(ginMode, "GIN_MODE")

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

func (c *cli) createConfig() error {
	config := &c.cfg

	// HTTP & Gin
	config.HttpPort = viper.GetString(httpPort)
	config.GinMode = viper.GetString(ginMode)

	// Log
	logConfig := &config.LogConfig
	logConfig.Level = viper.GetString(logLevel)
	logConfig.Formatter = viper.GetString(logFormatter)

	// Status
	statusConfig := &config.StatusConfig
	statusConfig.RpcAddr = viper.GetString(statusRpcAddr)

	// Metrics
	metricsConfig := &config.MetricsConfig
	metricsConfig.Addr = viper.GetString(metricsPrometheusAddr)
	metricsConfig.Path = viper.GetString(metricsPrometheusPath)

	return nil
}

// run bootstraps and runs the application
func (c *cli) run(cmd *cobra.Command, args []string) error {
	return runApp(c.cfg)
}
