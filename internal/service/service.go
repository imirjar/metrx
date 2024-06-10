package service

import (
	"context"
	"fmt"
	"log"

	"github.com/imirjar/metrx/internal/models"
)

// Generate html page witch consist all of metrics
func (s ServerService) MetricPage(ctx context.Context) (string, error) {
	log.Println("SERVICE MetricPage")
	gauges, err := s.MemStorager.ReadMetrics(ctx, "gauge")
	if err != nil {
		return "", err
	}
	counters, err := s.MemStorager.ReadMetrics(ctx, "counter")
	if err != nil {
		return "", err
	}

	gaugeForm := "<a>Gauge</a>"
	for _, m := range gauges {
		// log.Print(m)
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", m.ID, *m.Value)
	}

	counterForm := "<a>Counter</a>"
	for _, m := range counters {
		// log.Print(m)
		counterForm += fmt.Sprintf("<li>%s:%d</li>", m.ID, *m.Delta)
	}

	form := fmt.Sprintf("<html><ul>%s</ul><ul>%s</ul></html>", gaugeForm, counterForm)

	return form, nil
}

// Get Metric from store
func (s ServerService) ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	// log.Println("SERVICE ViewPath")
	return s.MemStorager.ReadMetric(ctx, metric)
}

// Set or update Metric in store
func (s ServerService) UpdateMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	log.Println("SERVICE UpdatePath")
	if err := s.MemStorager.AddMetric(ctx, metric); err != nil {
		return metric, err
	}
	return s.MemStorager.ReadMetric(ctx, metric)
}

// Set or update list of Metrics in store
func (s ServerService) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	log.Println("SERVICE BatchUpdate")
	return s.MemStorager.AddMetrics(ctx, metrics)
}

// Create new Service layer
func NewServerService() *ServerService {
	return &ServerService{}
}

// Service layer whitch consist app main logic
type ServerService struct {
	MemStorager Storager
}

// Implement storage layer methods for collecting and reading metrics
type Storager interface {
	AddMetrics(ctx context.Context, metrics []models.Metrics) error
	AddMetric(ctx context.Context, metric models.Metrics) error
	ReadMetrics(ctx context.Context, mType string) ([]models.Metrics, error)
	ReadMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error)
}
