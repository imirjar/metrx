package service

import (
	"fmt"
	"log"

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

func (s *Service) Gauge(name string, value float64) error {
	log.Println(name)
	gauge := models.Gauge{
		Name:  name,
		Value: value,
	}
	_, err := s.Storage.AddGauge(gauge)
	if err != nil {
		return fmt.Errorf("Ошибочка какая-то %w", err)
	}

	return nil
}

func (s *Service) Counter(name string, value int64) error {
	fmt.Println(name)
	counter := models.Counter{
		Name:  name,
		Value: value,
	}
	_, err := s.Storage.AddCounter(counter)
	if err != nil {
		return fmt.Errorf("Ошибочка какая-то %w", err)
	}

	return nil
}
