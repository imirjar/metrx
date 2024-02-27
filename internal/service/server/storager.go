package server

import (
	"fmt"
	"log"

	"github.com/imirjar/metrx/internal/models"
)

// get all metrics as html page
func (s ServerService) MetricPage() (string, error) {
	gauges, err := s.MemStorager.ReadAllGauges()
	if err != nil {
		return "", err
	}
	counters, err := s.MemStorager.ReadAllCounters()
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
	if s.cfg.Interval == 0 {
		s.Backup()
	}
	return form, nil
}

func (s ServerService) Update(metric models.Metrics) (models.Metrics, error) {

	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return metric, errServiceError
		}
		value, err := s.MemStorager.AddGauge(metric.ID, *metric.Value)
		if err != nil {
			return metric, err
		}
		metric.Value = &value
		return metric, nil

	case "counter":
		if metric.Delta == nil {
			return metric, errServiceError
		}
		delta, err := s.MemStorager.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			return metric, err
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errServiceError
	}
}

func (s ServerService) UpdateGauge(metric models.Metrics) (models.Metrics, error) {

	switch metric.MType {
	case "gauge":
		value, err := s.MemStorager.AddGauge(metric.ID, *metric.Value)
		if err != nil {
			return metric, err
		}
		metric.Value = &value
		return metric, nil

	case "counter":
		delta, err := s.MemStorager.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			return metric, err
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errServiceError
	}
}

func (s ServerService) BatchUpdate(metrics []models.Metrics) error {
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
		}
	}

	err := s.MemStorager.AddGauges(gauges)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	err = s.MemStorager.AddCounters(counters)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func (s ServerService) View(metric models.Metrics) (models.Metrics, error) {
	switch metric.MType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(metric)
		if !ok {
			return metric, errServiceError
		}
		metric.Value = &value
		return metric, nil
	case "counter":
		delta, ok := s.MemStorager.ReadCounter(metric)
		if !ok {
			return metric, errServiceError
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errServiceError
	}

}
