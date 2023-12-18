package storage

import (
	"github.com/imirjar/metrx/internal/models"
)

type MemStorage struct {
	Gauge   []models.Gauge
	Counter []models.Counter
	Metrics any
}

type Storager interface {
	AddGauge(gauge models.Gauge) (*models.Gauge, error)
	AddCounter(counter models.Counter) (*models.Counter, error)
}

func New() Storager {
	return &MemStorage{
		Gauge:   []models.Gauge{},
		Counter: []models.Counter{},
	}
}

func (m *MemStorage) AddGauge(gauge models.Gauge) (*models.Gauge, error) {
	for i, v := range m.Gauge {
		if v.Name == gauge.Name {
			m.Gauge[i] = gauge
			// fmt.Println(m.Gauge)
			return &m.Gauge[i], nil
		}
	}
	m.Gauge = append(m.Gauge, gauge)
	mewElementRef := &m.Gauge[len(m.Gauge)-1]
	// fmt.Println(m.Gauge)
	return mewElementRef, nil
}

func (m *MemStorage) AddCounter(counter models.Counter) (*models.Counter, error) {
	for i, v := range m.Counter {
		if v.Name == counter.Name {
			m.Counter[i].Value = v.Sum(counter.Value)
			// fmt.Println(m.Counter)
			return &m.Counter[i], nil
		}
	}
	m.Counter = append(m.Counter, counter)
	mewElementRef := &m.Counter[len(m.Counter)-1]
	// fmt.Println(m.Counter)
	return mewElementRef, nil
}
