package mock

import (
	"github.com/imirjar/metrx/internal/models"
)

func (m *MemStorage) Create(metric models.Metrics) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	switch metric.MType {
	case "gauge":
		m.Gauge[metric.ID] = *metric.Value
		return nil
	case "counter":
		m.Counter[metric.ID] = *metric.Delta
		return nil
	default:
		return errMetricStructureError
	}
}

func (m *MemStorage) ReadOne(metric models.Metrics) (models.Metrics, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	switch metric.MType {
	case "gauge":
		value := m.Gauge[metric.ID]
		metric.Value = &value
		return metric, false
	case "counter":
		delta := m.Counter[metric.ID]
		metric.Delta = &delta
		return metric, false
	default:
		return metric, false
	}
}

func (m *MemStorage) ReadAll(mType string) ([]models.Metrics, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var metrics = []models.Metrics{}
	// for i:= range m[mType]{

	// }
	switch mType {
	case "gauge":
		gMap := m.Gauge
		for i, v := range gMap {
			metrics = append(metrics, models.Metrics{ID: i, MType: mType, Value: &v})
		}
		return metrics, nil
	case "counter":
		cMap := m.Counter
		for i, v := range cMap {
			metrics = append(metrics, models.Metrics{ID: i, MType: mType, Delta: &v})
		}
		return metrics, nil
	default:
		return metrics, errMetricStructureError
	}
}

func (m *MemStorage) Update(metric models.Metrics) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	switch metric.MType {
	case "gauge":
		m.Gauge[metric.ID] = *metric.Value
		return nil
	case "counter":
		m.Counter[metric.ID] = *metric.Delta
		return nil
	default:
		return errMetricStructureError
	}
}

func (m *MemStorage) Delete(metric models.Metrics) error {
	return nil
}

// func (m *MemStorage) AddGauge(mName string, mValue float64) (float64, error) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	m.Gauge[mName] = mValue
// 	return mValue, nil
// }

// func (m *MemStorage) AddCounter(mName string, mValue int64) (int64, error) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	// if s.cfg.Interval == 0 {
// 	// 	s.Backup()
// 	// }
// 	cValue, ok := m.Counter[mName]
// 	if !ok {
// 		m.Counter[mName] = mValue
// 		return mValue, nil
// 	} else {
// 		newValue := cValue + mValue
// 		m.Counter[mName] = newValue
// 		return newValue, nil
// 	}

// }

// func (m *MemStorage) ReadAllGauge() map[string]float64 {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	return m.Gauge
// }

// func (m *MemStorage) ReadAllCounter() map[string]int64 {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	return m.Counter
// }

// func (m *MemStorage) ReadGauge(mName string) (float64, bool) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	v, ok := m.Gauge[mName]
// 	return v, ok
// }

// func (m *MemStorage) ReadCounter(mName string) (int64, bool) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	v, ok := m.Counter[mName]
// 	return v, ok
// }
