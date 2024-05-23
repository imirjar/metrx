package service

import (
	"context"
	"fmt"
	"log"

	"github.com/imirjar/metrx/internal/models"
)

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
