package server

import (
	"fmt"

	"github.com/imirjar/metrx/internal/models"
)

// get all metrics as html page
func (s ServerService) MetricPage() (string, error) {
	gauges, err := s.MemStorager.ReadAll("gauge")
	if err != nil {
		return "", err
	}
	counters, err := s.MemStorager.ReadAll("counter")
	if err != nil {
		return "", err
	}

	gaugeForm := "<a>Gauge</a>"
	for i, g := range gauges {
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", i, g)
	}

	counterForm := "<a>Counter</a>"
	for i, c := range counters {
		counterForm += fmt.Sprintf("<li>%s:%d</li>", i, c)
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
		_, exists := s.MemStorager.ReadOne(metric)
		if exists {
			err := s.MemStorager.Update(metric)
			if err != nil {
				return err
			}
			return nil
		} else {
			err := s.MemStorager.Create(metric)
			if err != nil {
				return err
			}
			return nil
		}

	case "counter":
		nMetric, exists := s.MemStorager.ReadOne(metric)
		if exists {
			value := *metric.Delta + *nMetric.Delta
			metric.Delta = &value
			err := s.MemStorager.Update(metric)
			if err != nil {
				return err
			}
			return nil
		} else {
			err := s.MemStorager.Create(metric)
			if err != nil {
				return err
			}
			return nil
		}
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
