package server

import (
	"fmt"

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

func (s ServerService) Update(metric models.Metrics) error {

	switch metric.MType {
	case "gauge":
		err := s.MemStorager.AddGauge(metric)
		if err != nil {
			return err
		}
		return nil

	case "counter":
		err := s.MemStorager.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			return err
		}
		return nil
	default:
		return errServiceError
	}
}

func (s ServerService) View(metric models.Metrics) (models.Metrics, error) {

	rMetric, ok := s.MemStorager.ReadOne(metric)
	if !ok {
		return metric, errServiceError
	}
	return rMetric, nil

}
