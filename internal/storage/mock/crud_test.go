package mock

import (
	"sync"
	"testing"

	"github.com/imirjar/metrx/internal/models"
)

func TestMemStorage_AddGauge(t *testing.T) {
	storage := MemStorage{
		mutex:   sync.Mutex{},
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}

	type want struct {
		gauges   map[string]float64
		counters map[string]int64
	}

	type counter struct {
		n string
		t string
		v float64
	}

	tests := []struct {
		name   string
		metric counter
		want   want
	}{
		{
			metric: counter{
				n: "g",
				t: "gauge",
				v: 123.1,
			},
			want: want{
				gauges:   map[string]float64{},
				counters: map[string]int64{},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := storage.AddGauge(
				models.Metrics{
					ID:    tt.metric.n,
					MType: tt.metric.n,
					Value: &tt.metric.v,
				}); err != nil {
				t.Errorf("MemStorage.AddGauge() error = %v", err)
			}
		})
	}
}
