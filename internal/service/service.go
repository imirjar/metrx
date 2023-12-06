package service

import (
	"fmt"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/internal/storage"
)

type Service struct {
	// Routes  *http.ServeMux
	Storage *storage.MemStorage
}

func (s *Service) Gauge(name, value string) []models.Gauge {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err)
	}

	gauge := models.Gauge{
		Name:  name,
		Value: floatValue,
	}

	find := s.Storage.Read("Name", name)

	if find == nil {
		s.Storage.Create(&gauge)
	} else {

		s.Storage.Update(&gauge, floatValue)
	}
	return s.Storage.GaugeStorage

}

func (s *Service) Counter(name, value string) []models.Counter {
	intValue, err := strconv.ParseInt(value, 36, 64)
	if err != nil {
		fmt.Println(err)
	}
	counter := models.Counter{
		Name:  name,
		Value: intValue,
	}

	s.Storage.CounterStorage = append(s.Storage.CounterStorage, counter)
	return s.Storage.CounterStorage
}

func New() *Service {

	return &Service{
		// Routes:  defineRoutes(),
		Storage: storage.New(),
	}
}
