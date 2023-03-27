package rabbitmq

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "rabbitmq"
	// The stability level of the exporter.
	stability = component.StabilityLevelAlpha
)

func createDefaultConfig() component.Config {
	return &Config{
		TimeoutSettings:  exporterhelper.NewDefaultTimeoutSettings(),
		RetrySettings:    exporterhelper.NewDefaultRetrySettings(),
		QueueSettings:    exporterhelper.NewDefaultQueueSettings(),
		VirtualHost:      "/",
		ConnectionString: "",
		Exchange: Exchange{
			Name: "otel-logs",
			Type: "direct",
		},
	}
}

func createLogsExporter(ctx context.Context, params exporter.CreateSettings, baseCfg component.Config) (exporter.Logs, error) {
	cfg := baseCfg.(*Config)
	exp, err := newLogsExporter(cfg, params)

	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	return exporterhelper.NewLogsExporter(
		ctx,
		params,
		baseCfg,
		exp.logsDataPusher,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(cfg.TimeoutSettings),
		exporterhelper.WithRetry(cfg.RetrySettings),
		exporterhelper.WithQueue(cfg.QueueSettings),
	)
}

// NewFactory creates a factory for RabbitMq exporter
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, stability))
}
