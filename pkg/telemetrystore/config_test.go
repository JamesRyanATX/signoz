package telemetrystore

import (
	"context"
	"testing"
	"time"

	"github.com/SigNoz/signoz/pkg/config"
	"github.com/SigNoz/signoz/pkg/config/envprovider"
	"github.com/SigNoz/signoz/pkg/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWithEnvProvider(t *testing.T) {
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_DSN", "tcp://localhost:9000")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_MAX__IDLE__CONNS", "60")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_MAX__OPEN__CONNS", "150")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_DIAL__TIMEOUT", "5s")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_DEBUG", "true")

	conf, err := config.New(
		context.Background(),
		config.ResolverConfig{
			Uris: []string{"env:"},
			ProviderFactories: []config.ProviderFactory{
				envprovider.NewFactory(),
			},
		},
		[]factory.ConfigFactory{
			NewConfigFactory(),
		},
	)
	require.NoError(t, err)

	actual := Config{}
	err = conf.Unmarshal("telemetrystore", &actual)
	require.NoError(t, err)

	assert.NoError(t, actual.Validate())

	expected := NewConfigFactory().New().(Config)
	expected.Provider = "clickhouse"
	expected.Connection.MaxOpenConns = 150
	expected.Connection.MaxIdleConns = 60
	expected.Connection.DialTimeout = 5 * time.Second
	expected.Clickhouse.DSN = "tcp://localhost:9000"

	assert.Equal(t, expected, actual)
}

func TestNewWithEnvProviderWithQuerySettings(t *testing.T) {
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_SETTINGS_MAX__EXECUTION__TIME", "10")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_SETTINGS_MAX__EXECUTION__TIME__LEAF", "10")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_SETTINGS_TIMEOUT__BEFORE__CHECKING__EXECUTION__SPEED", "10")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_SETTINGS_MAX__BYTES__TO__READ", "1000000")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_SETTINGS_MAX__RESULT__ROWS", "10000")

	conf, err := config.New(
		context.Background(),
		config.ResolverConfig{
			Uris: []string{"env:"},
			ProviderFactories: []config.ProviderFactory{
				envprovider.NewFactory(),
			},
		},
		[]factory.ConfigFactory{
			NewConfigFactory(),
		},
	)
	require.NoError(t, err)

	actual := Config{}
	err = conf.Unmarshal("telemetrystore", &actual)

	require.NoError(t, err)

	expected := Config{
		Clickhouse: ClickhouseConfig{
			QuerySettings: QuerySettings{
				MaxExecutionTime:                    10,
				MaxExecutionTimeLeaf:                10,
				TimeoutBeforeCheckingExecutionSpeed: 10,
				MaxBytesToRead:                      1000000,
				MaxResultRows:                       10000,
			},
		},
	}

	assert.Equal(t, expected.Clickhouse.QuerySettings, actual.Clickhouse.QuerySettings)
}

func TestDatabaseNamesWithEnvProvider(t *testing.T) {
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_TRACE__DATABASE", "custom_traces")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_METRICS__DATABASE", "custom_metrics")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_LOGS__DATABASE", "custom_logs")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_METER__DATABASE", "custom_meter")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_METADATA__DATABASE", "custom_metadata")
	t.Setenv("SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_ANALYTICS__DATABASE", "custom_analytics")

	conf, err := config.New(
		context.Background(),
		config.ResolverConfig{
			Uris: []string{"env:"},
			ProviderFactories: []config.ProviderFactory{
				envprovider.NewFactory(),
			},
		},
		[]factory.ConfigFactory{
			NewConfigFactory(),
		},
	)
	require.NoError(t, err)

	actual := Config{}
	err = conf.Unmarshal("telemetrystore", &actual)
	require.NoError(t, err)

	assert.Equal(t, "custom_traces", actual.Clickhouse.TraceDatabase)
	assert.Equal(t, "custom_metrics", actual.Clickhouse.MetricsDatabase)
	assert.Equal(t, "custom_logs", actual.Clickhouse.LogsDatabase)
	assert.Equal(t, "custom_meter", actual.Clickhouse.MeterDatabase)
	assert.Equal(t, "custom_metadata", actual.Clickhouse.MetadataDatabase)
	assert.Equal(t, "custom_analytics", actual.Clickhouse.AnalyticsDatabase)
}

func TestDatabaseNamesDefaults(t *testing.T) {
	conf, err := config.New(
		context.Background(),
		config.ResolverConfig{
			Uris: []string{"env:"},
			ProviderFactories: []config.ProviderFactory{
				envprovider.NewFactory(),
			},
		},
		[]factory.ConfigFactory{
			NewConfigFactory(),
		},
	)
	require.NoError(t, err)

	actual := Config{}
	err = conf.Unmarshal("telemetrystore", &actual)
	require.NoError(t, err)

	assert.Equal(t, DefaultTraceDatabase, actual.Clickhouse.TraceDatabase)
	assert.Equal(t, DefaultMetricsDatabase, actual.Clickhouse.MetricsDatabase)
	assert.Equal(t, DefaultLogsDatabase, actual.Clickhouse.LogsDatabase)
	assert.Equal(t, DefaultMeterDatabase, actual.Clickhouse.MeterDatabase)
	assert.Equal(t, DefaultMetadataDatabase, actual.Clickhouse.MetadataDatabase)
	assert.Equal(t, DefaultAnalyticsDatabase, actual.Clickhouse.AnalyticsDatabase)
}

func TestIndividualDatabaseNames(t *testing.T) {
	tests := []struct {
		name         string
		envVar       string
		envValue     string
		expectedFunc func(Config) string
		defaultValue string
	}{
		{
			name:         "TraceDatabase",
			envVar:       "SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_TRACE__DATABASE",
			envValue:     "test_traces",
			expectedFunc: func(c Config) string { return c.Clickhouse.TraceDatabase },
			defaultValue: DefaultTraceDatabase,
		},
		{
			name:         "MetricsDatabase",
			envVar:       "SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_METRICS__DATABASE",
			envValue:     "test_metrics",
			expectedFunc: func(c Config) string { return c.Clickhouse.MetricsDatabase },
			defaultValue: DefaultMetricsDatabase,
		},
		{
			name:         "LogsDatabase",
			envVar:       "SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_LOGS__DATABASE",
			envValue:     "test_logs",
			expectedFunc: func(c Config) string { return c.Clickhouse.LogsDatabase },
			defaultValue: DefaultLogsDatabase,
		},
		{
			name:         "MeterDatabase",
			envVar:       "SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_METER__DATABASE",
			envValue:     "test_meter",
			expectedFunc: func(c Config) string { return c.Clickhouse.MeterDatabase },
			defaultValue: DefaultMeterDatabase,
		},
		{
			name:         "MetadataDatabase",
			envVar:       "SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_METADATA__DATABASE",
			envValue:     "test_metadata",
			expectedFunc: func(c Config) string { return c.Clickhouse.MetadataDatabase },
			defaultValue: DefaultMetadataDatabase,
		},
		{
			name:         "AnalyticsDatabase",
			envVar:       "SIGNOZ_TELEMETRYSTORE_CLICKHOUSE_ANALYTICS__DATABASE",
			envValue:     "test_analytics",
			expectedFunc: func(c Config) string { return c.Clickhouse.AnalyticsDatabase },
			defaultValue: DefaultAnalyticsDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.envVar, tt.envValue)

			conf, err := config.New(
				context.Background(),
				config.ResolverConfig{
					Uris: []string{"env:"},
					ProviderFactories: []config.ProviderFactory{
						envprovider.NewFactory(),
					},
				},
				[]factory.ConfigFactory{
					NewConfigFactory(),
				},
			)
			require.NoError(t, err)

			actual := Config{}
			err = conf.Unmarshal("telemetrystore", &actual)
			require.NoError(t, err)

			assert.Equal(t, tt.envValue, tt.expectedFunc(actual))
		})

		t.Run(tt.name+"_Default", func(t *testing.T) {
			conf, err := config.New(
				context.Background(),
				config.ResolverConfig{
					Uris: []string{"env:"},
					ProviderFactories: []config.ProviderFactory{
						envprovider.NewFactory(),
					},
				},
				[]factory.ConfigFactory{
					NewConfigFactory(),
				},
			)
			require.NoError(t, err)

			actual := Config{}
			err = conf.Unmarshal("telemetrystore", &actual)
			require.NoError(t, err)

			assert.Equal(t, tt.defaultValue, tt.expectedFunc(actual))
		})
	}
}
