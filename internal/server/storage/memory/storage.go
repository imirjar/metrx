package memory

import (
	"context"
	"log"
	"time"

	"github.com/imirjar/metrx/internal/models"
)

// Create new storage layer
func NewMemoryStorage(filePath string) *Storage {

	storage := Storage{
		DumpPath:   filePath,
		AutoExport: false,
		MemStorage: MemStorage{

			Gauge:   map[string]float64{},
			Counter: map[string]int64{},
		},
	}

	return &storage
}

// "In memory" storage
type Storage struct {
	DumpPath   string
	AutoExport bool
	MemStorage
}

func (s *Storage) Configure(filePath string, autoImport bool, interval time.Duration) {
	if autoImport {
		s.MemStorage.Import(filePath)
	}

	if interval == 0 {
		s.AutoExport = true
	}

	if interval > 0 {
		go func() {
			defer s.MemStorage.Export(filePath)
			for {
				time.Sleep(interval)
				s.MemStorage.Export(filePath)
			}
		}()
	}
}

func (s *Storage) AddMetrics(ctx context.Context, metrics []models.Metrics) error {
	for _, m := range metrics {
		// log.Println("STORAGE", m)
		switch m.MType {
		case "gauge":
			s.MemStorage.Gauge[m.ID] = *m.Value
		case "counter":
			s.MemStorage.Counter[m.ID] = *m.Delta
		default:
			return errMetricTypeError
		}
	}

	if s.AutoExport {
		s.MemStorage.Export(s.DumpPath)
	}
	return nil
}

func (s *Storage) AddMetric(ctx context.Context, metric models.Metrics) error {

	switch metric.MType {
	case "gauge":
		s.MemStorage.Gauge[metric.ID] = *metric.Value
		return nil
	case "counter":
		log.Println("####1", s.MemStorage.Counter[metric.ID])
		s.MemStorage.Counter[metric.ID] = s.MemStorage.Counter[metric.ID] + *metric.Delta
		log.Println("####2", s.MemStorage.Counter[metric.ID])
		return nil
	default:
		return errMetricTypeError
	}
}

func (s *Storage) ReadMetrics(ctx context.Context, mType string) ([]models.Metrics, error) {

	var metrics = make([]models.Metrics, 0)

	switch mType {
	case "gauge":
		for n, v := range s.MemStorage.Gauge {
			v := v // otherwise all of metrics have the same value
			m := models.Metrics{
				ID:    n,
				MType: mType,
				Value: &v,
			}
			metrics = append(metrics, m)
		}
		// log.Print(metrics)
		return metrics, nil //s.MemStorage.Gauge[m.ID]
	case "counter":
		for n, d := range s.MemStorage.Counter {
			d := d // otherwise all of metrics have the same value
			m := models.Metrics{
				ID:    n,
				MType: mType,
				Delta: &d,
			}
			metrics = append(metrics, m)
			log.Print(m.ID)
		}
		// log.Print(metrics)
		return metrics, nil //s.MemStorage.Counter[m.ID]
	default:
		return nil, errMetricTypeError
	}
}

func (s *Storage) ReadMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {

	switch metric.MType {
	case "gauge":
		value := s.MemStorage.Gauge[metric.ID]
		metric.Value = &value
		return metric, nil
	case "counter":
		delta := s.MemStorage.Counter[metric.ID]
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errMetricTypeError
	}
}
