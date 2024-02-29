package mock

import (
	"context"
	"log"
	"time"

	"github.com/imirjar/metrx/config"
)

type Storage struct {
	DumpPath   string
	AutoExport bool
	MemStorage
}

func (s *Storage) configure(cfg config.ServerConfig) {
	if cfg.AutoImport {
		s.MemStorage.Import(cfg.FilePath)
	}

	if cfg.Interval == 0 {
		log.Println("AutoExport")
		s.AutoExport = true
	}

	if cfg.Interval > 0 {
		go func() {
			defer s.MemStorage.Export(cfg.FilePath)
			for {
				time.Sleep(cfg.Interval)
				s.MemStorage.Export(cfg.FilePath)
			}
		}()
	}

}

func (s *Storage) AddGauge(ctx context.Context, name string, value float64) (float64, error) {
	s.MemStorage.Gauge[name] = value

	if s.AutoExport {
		s.MemStorage.Export(s.DumpPath)
	}
	return value, nil
}

func (s *Storage) AddCounter(ctx context.Context, name string, delta int64) (int64, error) {
	log.Println("###Я ТОПОВОЕ ХРАНИЛИЩЕ МЕТРИК")
	newDelta := s.MemStorage.Counter[name] + delta
	s.MemStorage.Counter[name] = newDelta

	if s.AutoExport {
		s.MemStorage.Export(s.DumpPath)
	}
	return newDelta, nil
}

func (s *Storage) AddGauges(ctx context.Context, gauges map[string]float64) error {

	for i, v := range gauges {
		s.MemStorage.Gauge[i] = v
	}

	if s.AutoExport {
		s.MemStorage.Export(s.DumpPath)
	}
	return nil
}

func (s *Storage) AddCounters(ctx context.Context, counters map[string]int64) error {

	for i, d := range counters {
		// fmt.Println(metric)
		s.MemStorage.Counter[i] = s.MemStorage.Counter[i] + d
	}

	if s.AutoExport {
		s.MemStorage.Export(s.DumpPath)
	}

	return nil
}

func (s *Storage) ReadGauge(ctx context.Context, name string) (float64, bool) {
	if value, ok := s.MemStorage.Gauge[name]; ok {
		// metric.Value = &value
		return value, true
	} else {
		return 0, false
	}
}

func (s *Storage) ReadCounter(ctx context.Context, name string) (int64, bool) {
	if delta, ok := s.MemStorage.Counter[name]; ok {
		// metric.Delta = &delta
		return delta, true
	} else {
		return 0, false
	}
}

func (s *Storage) ReadAllGauges(ctx context.Context) (map[string]float64, error) {
	return s.MemStorage.Gauge, nil
}

func (s *Storage) ReadAllCounters(ctx context.Context) (map[string]int64, error) {
	return s.MemStorage.Counter, nil
}
