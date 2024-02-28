package mock

import (
	"github.com/imirjar/metrx/config"
)

func NewMockStorage(cfg config.ServerConfig) *Storage {

	storage := Storage{
		DumpPath:   cfg.FilePath,
		AutoExport: false,
		MemStorage: MemStorage{

			Gauge:   map[string]float64{},
			Counter: map[string]int64{},
		},
	}

	storage.configure(cfg)

	return &storage
}
