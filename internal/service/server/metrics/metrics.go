package metrics

import (
	"fmt"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/storage/mock"
)

func NewMetricsService(cfg config.ServerConfig) *MetricsService {
	store := mock.NewMockStorage(&cfg.StorageConfig)
	metricsService := MetricsService{
		MemStorager: store,
		cfg:         &cfg.ServiceConfig,
	}
	// run dump auto-exporter

	return &metricsService
}

type MetricsService struct {
	MemStorager Storager
	cfg         *config.ServiceConfig
}

type Storager interface {
	AddGauge(mName string, mValue float64)
	AddCounter(mName string, mValue int64)
	ReadAllGauge() map[string]float64
	ReadAllCounter() map[string]int64
	ReadGauge(mName string) (float64, bool)
	ReadCounter(mName string) (int64, bool)
}

func (s MetricsService) UpdateGauge(mName string, mValue float64) error {
	if mName == "" {
		return errMetricNameIncorrect
	}

	s.MemStorager.AddGauge(mName, mValue)
	// if s.cfg.Interval == 0 {
	// 	s.Backup()
	// }
	return nil
}

// оupdate counter
func (s MetricsService) UpdateCounter(mName string, mValue int64) error {
	if mName == "" {
		return errMetricNameIncorrect
	}

	// if counter exists -> counter += new value
	curV, ok := s.MemStorager.ReadCounter(mName)
	if ok {
		s.MemStorager.AddCounter(mName, curV+mValue)
	} else {
		s.MemStorager.AddCounter(mName, mValue)
	}
	// if s.cfg.Interval == 0 {
	// 	s.Backup()
	// }
	//no error
	return nil
}

// get gauge metric
func (s MetricsService) ViewGaugeByName(mName string) (float64, error) {
	if mName == "" {
		return 0, errMetricNameIncorrect
	}

	gauge, ok := s.MemStorager.ReadGauge(mName)
	if !ok {
		return gauge, errServiceError
	}

	return gauge, nil
}

// get counter metric
func (s MetricsService) ViewCounterByName(mName string) (int64, error) {
	if mName == "" {
		return 0, errMetricNameIncorrect
	}

	counter, ok := s.MemStorager.ReadCounter(mName)
	if !ok {
		return counter, errServiceError
	}

	return counter, nil
}

// get all metrics as html page
func (s MetricsService) MetricPage() string {
	gauges := s.MemStorager.ReadAllGauge()
	counters := s.MemStorager.ReadAllCounter()

	gaugeForm := "<a>Gauge</a>"
	for i, g := range gauges {
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", i, g)
	}

	counterForm := "<a>Counter</a>"
	for i, c := range counters {
		counterForm += fmt.Sprintf("<li>%s:%d</li>", i, c)
	}

	form := fmt.Sprintf("<html><ul>%s</ul><ul>%s</ul></html>", gaugeForm, counterForm)
	// if s.cfg.Interval == 0 {
	// 	s.Backup()
	// }
	return form
}
