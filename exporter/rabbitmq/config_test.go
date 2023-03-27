package rabbitmq

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id           component.ID
		expected     component.Config
		errorMessage string
	}{
		{
			id:           component.NewIDWithName("rabbitmq", ""),
			errorMessage: "connection_string must not be empty",
		},
		{
			id: component.NewIDWithName("rabbitmq", "1"),
			expected: &Config{
				ConnectionString: "ampq://guest:guest@localhost",
				VirtualHost:      "/",
				Exchange: Exchange{
					Name: "otel-logs-exchange",
					Type: "direct",
				},
				TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
				QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
				RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())

			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			if tt.errorMessage != "" {
				assert.EqualError(t, component.ValidateConfig(cfg), tt.errorMessage)
				return
			}

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
