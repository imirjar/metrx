package service

import (
	"testing"

	"github.com/imirjar/metrx/internal/models"
)

func TestCounter(t *testing.T) {
	service := New()
	tests := []struct { // добавляем слайс тестов
		name  string
		value models.Counter
	}{
		{
			name: "Create",
			value: models.Counter{
				Name:  "metric",
				Value: 10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := service.Counter(test.value.Name, test.value.Value); err != nil {
				t.Errorf("Sum() = %d", err)
			}
		})
	}
}

func TestGauge(t *testing.T) {
	service := New()
	tests := []struct { // добавляем слайс тестов
		name  string
		value models.Gauge
	}{
		{
			name: "Create",
			value: models.Gauge{
				Name:  "metric",
				Value: 10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := service.Gauge(test.value.Name, test.value.Value); err != nil {
				t.Errorf("Sum() = %d", err)
			}
		})
	}
}
