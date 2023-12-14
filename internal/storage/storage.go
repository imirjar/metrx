package storage

import (
	"fmt"

	"github.com/imirjar/metrx/internal/models"
)

type MemStorage struct {
	CounterStorage []models.Counter
	GaugeStorage   []models.Gauge
}

type Storage interface {
	CounterUpdate(name string, value int64) error
	CounterRead(name string) *models.Counter
	CounterReadAll() []models.Counter
	CounterCreate(obj *models.Counter)
	GaugeUpdate(name string, value float64) error
	GaugeRead(name string) *models.Gauge
	GaugeReadAll() []models.Gauge
	GaugeCreate(obj *models.Gauge)
}

func New() Storage {
	return &MemStorage{
		CounterStorage: []models.Counter{},
		GaugeStorage:   []models.Gauge{},
	}
}

func (m *MemStorage) GaugeCreate(obj *models.Gauge) {
	m.GaugeStorage = append(m.GaugeStorage, *obj)
}

func (m *MemStorage) GaugeRead(name string) *models.Gauge {
	for _, v := range m.GaugeStorage {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func (m *MemStorage) GaugeReadAll() []models.Gauge {
	return m.GaugeStorage
}

func (m *MemStorage) GaugeUpdate(name string, value float64) error {

	for i, v := range m.GaugeStorage {
		if v.Name == name {
			m.GaugeStorage[i].Value = value
			return nil
		}
	}
	return fmt.Errorf("Указанная запись не существует")
}

func (m *MemStorage) CounterCreate(obj *models.Counter) {
	m.CounterStorage = append(m.CounterStorage, *obj)
}

func (m *MemStorage) CounterRead(name string) *models.Counter {
	for _, v := range m.CounterStorage {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func (m *MemStorage) CounterReadAll() []models.Counter {
	return m.CounterStorage
}

func (m *MemStorage) CounterUpdate(name string, value int64) error {

	for i, v := range m.CounterStorage {
		if v.Name == name {
			m.CounterStorage[i].Value += value
			return nil
		}
	}
	return fmt.Errorf("Указанная запись не существует")
}
