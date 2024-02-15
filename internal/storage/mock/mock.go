package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/imirjar/metrx/config"
)

func NewMockStorage(cfg config.ServerConfig) *MemStorage {
	storage := MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}

	storage.configure(cfg)

	return &storage
}

type MemStorage struct {
	mutex   sync.Mutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func (m *MemStorage) configure(cfg config.ServerConfig) {
	if cfg.AutoImport {
		fmt.Println("AUTO IMPORT")
		fmt.Println(cfg.AutoImport)
		m.Import(cfg.FilePath)
	}

	if cfg.Interval > 0 {
		go func() {
			for {
				fmt.Println("AUTO EXPORT")
				time.Sleep(cfg.Interval)
				m.Export(cfg.FilePath)
			}
		}()
	}

}
