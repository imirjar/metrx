package mock

import (
	"github.com/imirjar/metrx/internal/models"
)

func (m *MemStorage) AddGauge(name string, value float64) (float64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Gauge[name] = value
	return value, nil
}

func (m *MemStorage) AddCounter(name string, delta int64) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	newDelta := m.Counter[name] + delta
	m.Counter[name] = newDelta
	return newDelta, nil
}

func (m *MemStorage) AddGauges(gauges map[string]float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i, v := range gauges {
		m.Gauge[i] = v
	}
	return nil
}

func (m *MemStorage) AddCounters(counters map[string]int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i, d := range counters {
		// fmt.Println(metric)
		m.Counter[i] = m.Counter[i] + d

	}
	return nil
}

func (m *MemStorage) ReadOne(metric models.Metrics) (models.Metrics, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	switch metric.MType {
	case "gauge":
		if value, ok := m.Gauge[metric.ID]; ok {
			metric.Value = &value
			return metric, true
		} else {
			return metric, false
		}

	case "counter":
		if delta, ok := m.Counter[metric.ID]; ok {
			metric.Delta = &delta
			return metric, true
		} else {
			return metric, false
		}
	default:
		return metric, false
	}
}

func (m *MemStorage) ReadAllGauges() (map[string]float64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Gauge, nil
}

func (m *MemStorage) ReadAllCounters() (map[string]int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Counter, nil
}
