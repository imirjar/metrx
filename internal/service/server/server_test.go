package server

import (
	"context"
	"math/rand"
	"testing"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/models"
)

func TestServerGauge(t *testing.T) {
	cfg := config.Testcfg
	server := NewServerService(cfg)
	var value float64

	tests := []struct { // добавляем слайс тестов
		name        string
		metric      models.Metrics
		expectedErr error
	}{
		{
			name: "gauge",
			metric: models.Metrics{
				ID:    "gauge",
				MType: "gauge",
				Value: &value,
			},
			expectedErr: nil,
		},
		{
			name: "counter",
			metric: models.Metrics{
				ID:    "counter",
				MType: "counter",
				Value: &value,
			},
			expectedErr: nil,
		},
	}
	for _, test := range tests {
		value = rand.Float64()
		t.Run(test.name, func(t *testing.T) {
			update, err := server.Update(context.Background(), test.metric)
			if err != nil {
				return
			}

			view, err := server.View(context.Background(), test.metric)
			if err != nil {
				return
			}
			if *update.Value != *view.Value {
				t.Errorf(`Value: %f Expected: %f \n`, *update.Value, *view.Value)
			}
		})
	}
}
