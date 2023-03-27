package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
)

type rabbitmqLogsExporter struct {
	connection   *amqp.Connection
	channel      AmpqChannel
	exchangeName string
	marshaler    *plog.JSONMarshaler
}

// AmpqChannel is an interface that allows us to mock the channel within tests
type AmpqChannel interface {
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Close() error
}

func (rabbitmqLogsExporter *rabbitmqLogsExporter) logsDataPusher(ctx context.Context, logData plog.Logs) error {
	body, err := rabbitmqLogsExporter.marshaler.MarshalLogs(logData)
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}

	err = rabbitmqLogsExporter.channel.PublishWithContext(ctx,
		rabbitmqLogsExporter.exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (rabbitmqLogsExporter *rabbitmqLogsExporter) Close() error {
	if err := rabbitmqLogsExporter.connection.Close(); err != nil {
		return err
	}

	if err := rabbitmqLogsExporter.channel.Close(); err != nil {
		return err
	}

	return nil
}

func newLogsExporter(config *Config, _ exporter.CreateSettings) (*rabbitmqLogsExporter, error) {
	amqpCfg := amqp.Config{
		Vhost:      config.VirtualHost,
		Properties: amqp.NewConnectionProperties(),
	}

	amqpCfg.Properties.SetClientConnectionName("otel-collector")

	conn, err := amqp.DialConfig(config.ConnectionString, amqpCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		config.Exchange.Name,
		config.Exchange.Type,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	exp := &rabbitmqLogsExporter{
		connection:   conn,
		channel:      channel,
		exchangeName: config.Exchange.Name,
		marshaler:    &plog.JSONMarshaler{},
	}

	return exp, nil
}
