package service

import (
	"context"
	"log"

	"github.com/imirjar/metrx/internal/models"
)

// Get Metric from store
func (s ServerService) ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	// log.Println("SERVICE ViewPath")
	return s.MemStorager.ReadMetric(ctx, metric)
}

func (s ServerService) ViewMetrics(ctx context.Context) (map[string][]models.Metrics, error) {
	// log.Println("SERVICE ViewPath")
	ms := make(map[string][]models.Metrics)
	var err error
	ms["gauges"], err = s.MemStorager.ReadMetrics(ctx, "gauge")
	if err != nil {
		return ms, err
	}

	ms["counters"], err = s.MemStorager.ReadMetrics(ctx, "counter")
	if err != nil {
		return ms, err
	}

	return ms, nil
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
func New() *ServerService {
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
