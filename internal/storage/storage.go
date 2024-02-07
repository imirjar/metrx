package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/imirjar/metrx/config"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
	cfg     *config.StorageConfig
	mutex   sync.Mutex
}

func NewStorage(cfg config.ServerConfig) *MemStorage {

	storage := MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
		cfg:     &cfg.StorageConfig,
	}

	if storage.cfg.AutoImport {
		storage.Import()
	}

	return &storage
}

func (m *MemStorage) Export() error {
	m.mutex.Lock()
	file, err := os.OpenFile(m.cfg.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	file.Write(data)
	m.mutex.Unlock()
	return nil

}

func (m *MemStorage) Import() error {
	file, err := os.ReadFile(m.cfg.FilePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, m); err != nil {
		return err
	}
	return nil
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
