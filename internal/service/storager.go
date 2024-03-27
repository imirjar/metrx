package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
)

func (s ServerService) MetricPage(ctx context.Context) (string, error) {
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

func (s ServerService) ViewPath(ctx context.Context, name, mType string) (string, error) {
	switch mType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(ctx, name)
		if !ok {
			return "", errGaugeDoesNotMatched
		}
		return fmt.Sprint(value), nil
	case "counter":
		delta, ok := s.MemStorager.ReadCounter(ctx, name)
		if !ok {
			return "", errCounterDoesNotMatched
		}
		return fmt.Sprint(delta), nil
	default:
		return "", errMetricTypeError
	}
}

func (s ServerService) UpdatePath(ctx context.Context, name, mType, mValue string) (string, error) {
	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return "", errConvertationError
		}
		newValue, err := s.MemStorager.AddGauge(ctx, name, value)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(newValue), nil

	case "counter":
		delta, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return "", errConvertationError
		}

		newDelta, err := s.MemStorager.AddCounter(ctx, name, delta)
		if err != nil {
			return "", err
		}

		return fmt.Sprint(newDelta), nil
	default:
		return "", errMetricTypeError
	}
}

func (s ServerService) BatchUpdate(ctx context.Context, metrics []models.Metrics) error {
	var (
		gauges   = map[string]float64{}
		counters = map[string]int64{}
	)

	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			gauges[metric.ID] = *metric.Value
		case "counter":
			counters[metric.ID] = *metric.Delta
			// if _, ok := counters[metric.ID]; ok {
			// 	counters[metric.ID] = counters[metric.ID] + *metric.Delta
			// } else {
			// 	counters[metric.ID] = *metric.Delta
			// }

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
