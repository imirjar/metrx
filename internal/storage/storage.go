package storage

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func New() *MemStorage {
	return &MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}

func (m *MemStorage) AddGauge(mName string, mValue float64) {
	m.Gauge[mName] = mValue
}

func (m *MemStorage) AddCounter(mName string, mValue int64) {
	m.Counter[mName] = mValue
}

func (m *MemStorage) ReadAllGauge() map[string]float64 {
	return m.Gauge
}

func (m *MemStorage) ReadAllCounter() map[string]int64 {
	return m.Counter
}

func (m *MemStorage) ReadGauge(mName string) (float64, bool) {
	v, ok := m.Gauge[mName]
	return v, ok
}

func (m *MemStorage) ReadCounter(mName string) (int64, bool) {
	v, ok := m.Counter[mName]
	return v, ok
}
