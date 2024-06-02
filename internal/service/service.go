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
	gauges, err := s.MemStorager.ReadAllGauges(ctx)
	if err != nil {
		return "", err
	}
	counters, err := s.MemStorager.ReadAllCounters(ctx)
	if err != nil {
		return "", err
	}

	gaugeForm := "<a>Gauge</a>"
	for k, g := range gauges {
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", k, g)
	}

	counterForm := "<a>Counter</a>"
	for k, c := range counters {
		counterForm += fmt.Sprintf("<li>%s:%d</li>", k, c)
	}

	form := fmt.Sprintf("<html><ul>%s</ul><ul>%s</ul></html>", gaugeForm, counterForm)

	return form, nil
}

// Get Metric from store
func (s ServerService) ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	// log.Println("SERVICE ViewPath")
	switch metric.MType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(ctx, metric.ID)
		if !ok {
			return metric, errGaugeDoesNotMatched
		}
		metric.Value = &value
		return metric, nil
	case "counter":
		delta, ok := s.MemStorager.ReadCounter(ctx, metric.ID)
		if !ok {
			return metric, errCounterDoesNotMatched
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errMetricTypeError
	}
}

// Set or update Metric in store
func (s ServerService) UpdateMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	log.Println("SERVICE UpdatePath")
	switch metric.MType {
	case "gauge":
		newValue, err := s.MemStorager.AddGauge(ctx, metric.ID, *metric.Value)
		if err != nil {
			return metric, err
		}
		metric.Value = &newValue
		return metric, nil

	case "counter":
		newDelta, err := s.MemStorager.AddCounter(ctx, metric.ID, *metric.Delta)
		if err != nil {
			return metric, err
		}
		metric.Delta = &newDelta
		return metric, nil
	default:
		return metric, errMetricTypeError
	}
}

// Set or update list of Metrics in store
func (s ServerService) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	log.Println("SERVICE BatchUpdate")
	var (
		gauges   = map[string]float64{}
		counters = map[string]int64{}
	)

	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			gauges[metric.ID] = *metric.Value
		case "counter":
			if _, ok := counters[metric.ID]; ok {
				counters[metric.ID] = counters[metric.ID] + *metric.Delta
			} else {
				counters[metric.ID] = *metric.Delta
			}

		}
	}

	err := s.MemStorager.AddGauges(ctx, gauges)
	if err != nil {
		return err
	}

	err = s.MemStorager.AddCounters(ctx, counters)
	if err != nil {
		return err
	}

	return nil
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
	AddGauges(ctx context.Context, gauges map[string]float64) error
	AddCounters(ctx context.Context, counters map[string]int64) error
	AddGauge(ctx context.Context, name string, value float64) (float64, error)
	AddCounter(ctx context.Context, name string, delta int64) (int64, error)
	ReadGauge(ctx context.Context, name string) (float64, bool)
	ReadCounter(ctx context.Context, name string) (int64, bool)
	ReadAllGauges(ctx context.Context) (map[string]float64, error)
	ReadAllCounters(ctx context.Context) (map[string]int64, error)
}
