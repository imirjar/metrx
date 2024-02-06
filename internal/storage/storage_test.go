package storage

import (
	"testing"

	"github.com/imirjar/metrx/config"
)

func TestMemStorage_AddGauge_ReadGauge(t *testing.T) {

	tests := []struct {
		name     string
		mName    string
		mValue   float64
		expected float64
	}{
		{
			name:     "All right",
			mName:    "gaugeMetric",
			mValue:   10.123,
			expected: 10.123,
		},
		{
			name:     "Zero value",
			mName:    "gaugeMetric",
			mValue:   0,
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := NewStorage(*config.NewServerConfig())
			memStorage.AddGauge(tt.mName, tt.mValue)

			gauge, ok := memStorage.ReadGauge(tt.mName)

			if !ok || gauge != tt.expected {
				t.Errorf("Value = %f, want %f", gauge, tt.expected)
			}
		})
	}
}

func TestMemStorage_AddCounter_ReadCounter(t *testing.T) {

	tests := []struct {
		name     string
		mName    string
		mValue   int64
		expected int64
	}{
		{
			name:     "All right",
			mName:    "gaugeMetric",
			mValue:   10,
			expected: 10,
		},
		{
			name:     "Zero value",
			mName:    "gaugeMetric",
			mValue:   0,
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := NewStorage(*config.NewServerConfig())
			memStorage.AddCounter(tt.mName, tt.mValue)

			gauge, ok := memStorage.ReadCounter(tt.mName)

			if !ok || gauge != tt.expected {
				t.Errorf("Value = %d, want %d", gauge, tt.expected)
			}
		})
	}
}
