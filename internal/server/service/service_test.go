package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imirjar/metrx/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestViewMetric(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockStorager(ctrl)

	var value float64 = 110
	gaugeMetric := models.Metrics{
		ID:    "TestGaugeA",
		MType: "gauge",
	}

	var delta int64 = 100
	counterMetric := models.Metrics{
		ID:    "TestCounterA",
		MType: "counter",
	}

	m.EXPECT().ReadMetric(context.Background(), gaugeMetric).
		Return(models.Metrics{
			ID:    "TestGaugeA",
			MType: "gauge",
			Value: &value,
		}, nil)
	m.EXPECT().ReadMetric(context.Background(), counterMetric).
		Return(models.Metrics{
			ID:    "TestCounterA",
			MType: "counter",
			Delta: &delta,
		}, nil)

	testService := ServerService{
		MemStorager: m,
	}

	tests := []struct {
		name string
		send models.Metrics
		want string
	}{
		{
			name: "ok gauge",
			send: models.Metrics{
				ID:    "TestGaugeA",
				MType: "gauge",
			},
			want: "110",
		},

		{
			name: "ok counter",
			send: models.Metrics{
				ID:    "TestCounterA",
				MType: "counter",
			},
			want: "100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := testService.ViewMetric(context.Background(), tt.send)
			if err != nil {
				t.Error(err)
			}

			val, err := result.GetVal()
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.want, val)
		})
	}
}

func TestUpdateMetric(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockStorager(ctrl)

	var value float64 = 110
	gaugeMetric := models.Metrics{
		ID:    "TestGaugeA",
		MType: "gauge",
		Value: &value,
	}

	var delta int64 = 100
	counterMetric := models.Metrics{
		ID:    "TestCounterA",
		MType: "counter",
		Delta: &delta,
	}

	m.EXPECT().AddMetric(context.Background(), gaugeMetric).Return(nil)
	m.EXPECT().ReadMetric(context.Background(), gaugeMetric).Return(gaugeMetric, nil)
	m.EXPECT().AddMetric(context.Background(), counterMetric).Return(nil)
	m.EXPECT().ReadMetric(context.Background(), counterMetric).Return(counterMetric, nil)

	testService := ServerService{
		MemStorager: m,
	}

	tests := []struct {
		name string
		send models.Metrics
		want error
	}{
		{
			name: "ok gauge",
			send: models.Metrics{
				ID:    "TestGaugeA",
				MType: "gauge",
				Value: &value,
			},
			want: nil,
		},

		{
			name: "ok counter",
			send: models.Metrics{
				ID:    "TestCounterA",
				MType: "counter",
				Delta: &delta,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := testService.UpdateMetric(context.Background(), tt.send)
			if err != nil {
				t.Error(err)
			}

			_, err = result.GetVal()
			if err != nil {
				t.Error(err)
			}

			// assert.Equal(t, tt.want, val)
		})
	}
}

func TestUpdateMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockStorager(ctrl)

	var value float64 = 110
	gaugeMetrics := []models.Metrics{
		{
			ID:    "TestGaugeA",
			MType: "gauge",
			Value: &value,
		},
	}

	var delta int64 = 100
	counterMetrics := []models.Metrics{
		{
			ID:    "TestCounterA",
			MType: "counter",
			Delta: &delta,
		},
	}

	m.EXPECT().AddMetrics(context.Background(), gaugeMetrics).Return(nil)
	m.EXPECT().AddMetrics(context.Background(), counterMetrics).Return(nil)

	testService := ServerService{
		MemStorager: m,
	}
	tests := []struct {
		name string
		send []models.Metrics
		want error
	}{
		{
			name: "ok gauge",
			send: []models.Metrics{
				{
					ID:    "TestGaugeA",
					MType: "gauge",
					Value: &value,
				},
			},

			want: nil,
		},

		{
			name: "ok counter",
			send: []models.Metrics{
				{
					ID:    "TestCounterA",
					MType: "counter",
					Delta: &delta,
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testService.UpdateMetrics(context.Background(), tt.send)

			assert.Equal(t, tt.want, result)
		})
	}
}
