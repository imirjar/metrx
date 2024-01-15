package service

import (
	"fmt"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/internal/storage"
)

func New() *Service {
	return &Service{
		Storage: storage.New(),
	}
}

type Service struct {
	Storage storage.Storager
}

func (s *Service) UpdateMetric(mType string, name string, value string) error {
	if name == "" {
		return errMetricNameIncorrect
	}
	switch mType {

	case "gauge":
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errConvertationError
		}
		gauge := models.Gauge{
			Name:  name,
			Value: value,
		}
		_, err = s.Storage.AddGauge(gauge)
		if err != nil {
			return errStorageError
		}
		return nil

	case "counter":
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errConvertationError
		}
		counter := models.Counter{
			Name:  name,
			Value: value,
		}
		_, err = s.Storage.AddCounter(counter)
		if err != nil {
			return errStorageError
		}
		return nil
	default:
		return errServiceError
	}

}

func (s *Service) ViewGaugeByName(name string) (*models.Gauge, error) {
	if name == "" {
		return nil, errMetricNameIncorrect
	}

	return s.Storage.ReadGauge(name), nil

}

func (s *Service) ViewCounterByName(name string) (*models.Counter, error) {
	if name == "" {
		return nil, errMetricNameIncorrect
	}

	return s.Storage.ReadCounter(name), nil

}

func (s *Service) FindMetricValue(mType string, name string) any {
	return nil
}

func (s *Service) MetricList() string {
	gauges := s.Storage.ReadAllGauge()
	counters := s.Storage.ReadAllCounter()

	gaugeForm := "<a>Gauge</a>"
	for _, g := range gauges {
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", g.Name, g.Value)
	}

	counterForm := "<a>Counter</a>"
	for _, c := range counters {
		counterForm += fmt.Sprintf("<li>%s:%d</li>", c.Name, c.Value)
	}

	form := fmt.Sprintf("<html><ul>%s</ul><ul>%s</ul></html>", gaugeForm, counterForm)
	return form
}
