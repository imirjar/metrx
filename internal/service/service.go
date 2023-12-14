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
	Storage storage.Storage
}

func (s *Service) Gauge(name, value string) []models.Gauge {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err)
	}
	if ok := s.Storage.GaugeRead(name); ok != nil {
		s.Storage.GaugeUpdate(name, floatValue)
	} else {
		gauge := models.Gauge{
			Name:  name,
			Value: floatValue,
		}
		s.Storage.GaugeCreate(&gauge)
	}
	return s.Storage.GaugeReadAll()

}

func (s *Service) Counter(name, value string) []models.Counter {
	intValue, err := strconv.ParseInt(value, 36, 64)
	if err != nil {
		fmt.Println(err)
	}
	if ok := s.Storage.CounterRead(name); ok != nil {
		s.Storage.CounterUpdate(name, intValue)
	} else {
		counter := models.Counter{
			Name:  name,
			Value: intValue,
		}

		s.Storage.CounterCreate(&counter)
	}
	return s.Storage.CounterReadAll()
}
