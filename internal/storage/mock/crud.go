package mock

func (m *MemStorage) AddGauge(mName string, mValue float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Gauge[mName] = mValue
}

func (m *MemStorage) AddCounter(mName string, mValue int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// if s.cfg.Interval == 0 {
	// 	s.Backup()
	// }
	m.Counter[mName] = mValue
}

func (m *MemStorage) ReadAllGauge() map[string]float64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Gauge
}

func (m *MemStorage) ReadAllCounter() map[string]int64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Counter
}

func (m *MemStorage) ReadGauge(mName string) (float64, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, ok := m.Gauge[mName]
	return v, ok
}

func (m *MemStorage) ReadCounter(mName string) (int64, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, ok := m.Counter[mName]
	return v, ok
}
