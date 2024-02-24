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

func (m *MemStorage) ReadGauge(metric models.Metrics) (float64, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if value, ok := m.Gauge[metric.ID]; ok {
		metric.Value = &value
		return value, true
	} else {
		return value, false
	}

}

func (m *MemStorage) ReadCounter(metric models.Metrics) (int64, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if delta, ok := m.Counter[metric.ID]; ok {
		metric.Delta = &delta
		return delta, true
	} else {
		return delta, false
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
