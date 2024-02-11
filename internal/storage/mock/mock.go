package mock

import (
	"sync"

	"github.com/imirjar/metrx/config"
)

func NewMockStorage(cfg *config.StorageConfig) *MemStorage {
	storage := MemStorage{
		cfg:     cfg,
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}

	if cfg.AutoImport {
		storage.Import()
	}

	return &storage
}

type MemStorage struct {
	cfg     *config.StorageConfig
	mutex   sync.Mutex
	Gauge   map[string]float64
	Counter map[string]int64
}
