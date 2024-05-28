package mock

import (
	"time"
)

func NewMockStorage(filePath string) *Storage {

	storage := Storage{
		DumpPath:   filePath,
		AutoExport: false,
		MemStorage: MemStorage{

			Gauge:   map[string]float64{},
			Counter: map[string]int64{},
		},
	}

	return &storage
}

type Storage struct {
	DumpPath   string
	AutoExport bool
	MemStorage
}

func (s *Storage) Configure(filePath string, autoImport bool, interval time.Duration) {
	if autoImport {
		s.MemStorage.Import(filePath)
	}

	if interval == 0 {
		s.AutoExport = true
	}

	if interval > 0 {
		go func() {
			defer s.MemStorage.Export(filePath)
			for {
				time.Sleep(interval)
				s.MemStorage.Export(filePath)
			}
		}()
	}
}
