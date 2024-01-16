package storage

import (
	"github.com/imirjar/metrx/internal/models"
)

type MemStorage struct {
	Gauge   []models.Gauge
	Counter []models.Counter
}

type Storager interface {
	AddGauge(gauge models.Gauge) (*models.Gauge, error)
	AddCounter(counter models.Counter) (*models.Counter, error)
	ReadAllGauge() []models.Gauge
	ReadAllCounter() []models.Counter
	ReadGauge(name string) *models.Gauge
	ReadCounter(name string) *models.Counter
	Drop()
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

func (m *MemStorage) ReadAllGauge() []models.Gauge {
	return m.Gauge
}

func (m *MemStorage) ReadAllCounter() []models.Counter {
	return m.Counter
}

func (m *MemStorage) ReadGauge(name string) *models.Gauge {
	for _, g := range m.Gauge {
		if name == g.Name {
			return &g
		}
	}
	return nil
}

func (m *MemStorage) ReadCounter(name string) *models.Counter {
	for _, c := range m.Counter {
		if name == c.Name {
			return &c
		}
	}
	return nil
}

func (m *MemStorage) Drop() {
	m = nil
}
