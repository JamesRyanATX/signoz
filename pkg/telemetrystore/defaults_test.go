package telemetrystore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultDatabaseNames(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "DefaultTraceDatabase",
			constant: DefaultTraceDatabase,
			expected: "signoz_traces",
		},
		{
			name:     "DefaultMetricsDatabase",
			constant: DefaultMetricsDatabase,
			expected: "signoz_metrics",
		},
		{
			name:     "DefaultLogsDatabase",
			constant: DefaultLogsDatabase,
			expected: "signoz_logs",
		},
		{
			name:     "DefaultMeterDatabase",
			constant: DefaultMeterDatabase,
			expected: "signoz_meter",
		},
		{
			name:     "DefaultMetadataDatabase",
			constant: DefaultMetadataDatabase,
			expected: "signoz_metadata",
		},
		{
			name:     "DefaultAnalyticsDatabase",
			constant: DefaultAnalyticsDatabase,
			expected: "signoz_analytics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

func TestDefaultDatabaseNamesAreUnique(t *testing.T) {
	databases := []string{
		DefaultTraceDatabase,
		DefaultMetricsDatabase,
		DefaultLogsDatabase,
		DefaultMeterDatabase,
		DefaultMetadataDatabase,
		DefaultAnalyticsDatabase,
	}

	seen := make(map[string]bool)
	for _, db := range databases {
		assert.False(t, seen[db], "Database name %s appears more than once", db)
		seen[db] = true
	}
}

func TestDefaultDatabaseNamesHaveSignozPrefix(t *testing.T) {
	databases := map[string]string{
		"DefaultTraceDatabase":     DefaultTraceDatabase,
		"DefaultMetricsDatabase":   DefaultMetricsDatabase,
		"DefaultLogsDatabase":      DefaultLogsDatabase,
		"DefaultMeterDatabase":     DefaultMeterDatabase,
		"DefaultMetadataDatabase":  DefaultMetadataDatabase,
		"DefaultAnalyticsDatabase": DefaultAnalyticsDatabase,
	}

	for name, db := range databases {
		assert.Contains(t, db, "signoz_", "%s should contain 'signoz_' prefix", name)
	}
}
