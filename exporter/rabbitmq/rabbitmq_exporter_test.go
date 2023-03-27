package rabbitmq

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
	"testing"
)

type MockMessage struct {
	MessageBody  []byte
	ExchangeName string
}

type MockAmpqChannel struct {
	PublishedMessages []MockMessage
}

func (mockAmpqChannel *MockAmpqChannel) PublishWithContext(_ context.Context, exchange, _ string, _, _ bool, msg amqp.Publishing) error {
	mockMsg := &MockMessage{
		MessageBody:  msg.Body,
		ExchangeName: exchange,
	}

	mockAmpqChannel.PublishedMessages = append(mockAmpqChannel.PublishedMessages, *mockMsg)
	return nil
}

func (_ *MockAmpqChannel) Close() error {
	return nil
}

func getExpectedJsonMessage(t *testing.T, marshaler plog.JSONMarshaler, logs plog.Logs) string {
	d, err := marshaler.MarshalLogs(logs)
	if err != nil {
		t.Errorf("failed to marshal logs: %s", err)
	}
	return string(d)
}

func TestLogsDataPublisher(t *testing.T) {
	lr := testdata.GenerateLogsOneLogRecord()
	marshaler := &plog.JSONMarshaler{}

	ch := &MockAmpqChannel{}

	exp := &rabbitmqLogsExporter{
		connection:   nil,
		channel:      ch,
		exchangeName: "otel-test-logs-exchange",
		marshaler:    marshaler,
	}

	err := exp.logsDataPusher(context.Background(), lr)
	require.NoError(t, err)

	assert.Len(t, ch.PublishedMessages, 1)
	pubMsg := ch.PublishedMessages[0]
	assert.Equal(t, getExpectedJsonMessage(t, *marshaler, lr), string(pubMsg.MessageBody))
}

func TestLogsDataPublisher_ExchangeName(t *testing.T) {
	lr := testdata.GenerateLogsOneLogRecord()
	marshaler := &plog.JSONMarshaler{}

	ch := &MockAmpqChannel{}

	exp := &rabbitmqLogsExporter{
		connection:   nil,
		channel:      ch,
		exchangeName: "otel-test-logs-exchange",
		marshaler:    marshaler,
	}

	err := exp.logsDataPusher(context.Background(), lr)
	require.NoError(t, err)

	assert.Len(t, ch.PublishedMessages, 1)
	pubMsg := ch.PublishedMessages[0]
	assert.Equal(t, exp.exchangeName, pubMsg.ExchangeName)
}
