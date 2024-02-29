package server

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
)

// get all metrics as html page
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

func (s ServerService) View(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	switch metric.MType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(ctx, metric.ID)
		log.Println("###GAUGE OUT--->", metric.ID, ":", value)
		log.Printf("###GAUGE OUT F--->%f", value)
		if !ok {
			return metric, errServiceError
		}
		metric.Value = &value
		return metric, nil
	case "counter":
		delta, ok := s.MemStorager.ReadCounter(ctx, metric.ID)
		// log.Println("###COUNTER OUT--->", metric.ID, ":", delta)
		if !ok {
			return metric, errServiceError
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errServiceError
	}
}

func (s ServerService) Update(ctx context.Context, metric models.Metrics) (models.Metrics, error) {

	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return metric, errServiceError
		}
		value, err := s.MemStorager.AddGauge(ctx, metric.ID, *metric.Value)
		if err != nil {
			return metric, err
		}
		metric.Value = &value
		return metric, nil

	case "counter":
		if metric.Delta == nil {
			return metric, errServiceError
		}
		delta, err := s.MemStorager.AddCounter(ctx, metric.ID, *metric.Delta)
		if err != nil {
			return metric, err
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errServiceError
	}
}

func (s ServerService) ViewPath(ctx context.Context, name, mType string) (string, error) {
	switch mType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(ctx, name)
		// log.Println("###GAUGE OUT--->", name, ":", value)
		if !ok {
			return "", errServiceError
		}
		return fmt.Sprint(value), nil
	case "counter":
		delta, ok := s.MemStorager.ReadCounter(ctx, name)
		// log.Println("###GAUGE OUT--->", name, ":", delta)
		if !ok {
			return "", errServiceError
		}
		return fmt.Sprint(delta), nil
	default:
		return "", errServiceError
	}
}

func (s ServerService) UpdatePath(ctx context.Context, name, mType, mValue string) (string, error) {

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return "", errServiceError
		}
		_, err = s.MemStorager.AddGauge(ctx, name, value)
		if err != nil {
			return "", err
		}
		return mValue, nil

	case "counter":
		delta, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return "", errServiceError
		}
		newDelta, err := s.MemStorager.AddCounter(ctx, name, delta)
		if err != nil {
			return "", err
		}

		return fmt.Sprint(newDelta), nil
	default:
		return "", errServiceError
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
