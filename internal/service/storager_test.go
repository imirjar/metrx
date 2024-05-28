package service

import (
	"context"
	"testing"

	"github.com/imirjar/metrx/internal/models"
)

var (
	i int64   = 100
	f float64 = 123.998
)

func TestServerService_PathHandler(t *testing.T) {
	service := NewServerService()

	tests := []struct {
		name    string
		metric  models.Metrics
		want    models.Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "OK Gauge",
			metric: models.Metrics{
				ID:    "gaugeMetric",
				MType: "gauge",
				Value: &f,
			},
			want: models.Metrics{
				ID:    "gaugeMetric",
				MType: "gauge",
				Value: &f,
			},
			wantErr: false,
		},
		{
			name: "OK Counter",
			metric: models.Metrics{
				ID:    "counterMetric",
				MType: "counter",
				Delta: &i,
			},
			want: models.Metrics{
				ID:    "gaugeMetric",
				MType: "gauge",
				Value: &f,
			},
			wantErr: false,
		},
		{
			name: "OK Add Counter",
			metric: models.Metrics{
				ID:    "counterMetric",
				MType: "counter",
				Delta: &i,
			},
			want: models.Metrics{
				ID:    "gaugeMetric",
				MType: "gauge",
				Value: &f,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.UpdateMetric(context.Background(), tt.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerService.ViewPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("ServerService.ViewPath() = %v, want %v", result, tt.want)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ViewMetric(context.Background(), tt.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerService.ViewPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("ServerService.ViewPath() = %v, want %v", result, tt.want)
			}
		})
	}
}
