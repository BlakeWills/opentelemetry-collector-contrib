package rabbitmq

import (
	"fmt"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config represents the exporter config settings within the collector's config.yaml
type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`
	VirtualHost                    string   `mapstructure:"virtual_host"`
	ConnectionString               string   `mapstructure:"connection_string"`
	Exchange                       Exchange `mapstructure:"exchange"`
}

// Exchange defines the configuration for the RabbitMQ exchange that we will publish too.
type Exchange struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
}

// Validate checks if the rabbitmq exporter configuration is valid
func (cfg *Config) Validate() error {

	if cfg.ConnectionString == "" {
		return fmt.Errorf("connection_string must not be empty")
	}

	return nil
}
