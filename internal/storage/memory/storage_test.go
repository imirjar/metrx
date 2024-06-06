package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/imirjar/metrx/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	store := Storage{
		MemStorage: MemStorage{
			mutex:   sync.Mutex{},
			Gauge:   make(map[string]float64),
			Counter: make(map[string]int64),
		},
	}

	var value float64 = 100
	var delta int64 = 100

	type output struct {
		metric models.Metrics
		err    error
	}

	tests := []struct {
		name  string
		input models.Metrics
		output
	}{
		{
			name: "ok",
			input: models.Metrics{
				ID:    "TestGaugeA",
				MType: "gauge",
				Value: &value,
			},
			output: output{
				metric: models.Metrics{
					ID:    "TestGaugeA",
					MType: "gauge",
					Value: &value,
				},
				err: nil,
			},
		},
		{
			name: "ok2",
			input: models.Metrics{
				ID:    "TestCounterA",
				MType: "counter",
				Delta: &delta,
			},
			output: output{
				metric: models.Metrics{
					ID:    "TestCounterA",
					MType: "counter",
					Delta: &delta,
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := store.AddMetric(context.Background(), tt.input); err != nil {
				t.Error(err)
			}
			m, err := store.ReadMetric(context.Background(), tt.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, m, tt.output.metric)
		})
	}

}

func TestMetrics(t *testing.T) {
	store := Storage{
		MemStorage: MemStorage{
			mutex:   sync.Mutex{},
			Gauge:   make(map[string]float64),
			Counter: make(map[string]int64),
		},
	}

	var value float64 = 100
	var delta int64 = 100

	tests := []struct {
		name   string
		input  []models.Metrics
		output error
	}{
		{
			name: "ok",
			input: []models.Metrics{
				models.Metrics{
					ID:    "TestGaugeA",
					MType: "gauge",
					Value: &value,
				},
			},
			output: nil,
		},
		{
			name: "ok2",
			input: []models.Metrics{
				models.Metrics{
					ID:    "TestCounterA",
					MType: "counter",
					Delta: &delta,
				},
			},
			output: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := store.AddMetrics(context.Background(), tt.input); err != nil {
				t.Error(err)
			}
			_, err := store.ReadMetrics(context.Background(), "gauge")
			if err != nil {
				t.Error(err)
			}
		})
	}
}
