package service

import (
	"net/http"
	"time"

	"github.com/imirjar/metrx/internal/storage"
)

func NewServer() *Server {
	return &Server{
		Storage: storage.New(),
	}
}

func NewAgent() *Agent {
	agent := &Agent{
		Metrics: []string{
			"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
			"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
			"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
		},
		Client: &http.Client{
			Timeout: time.Second * 1,
		},
	}

	return agent
}

type Storager interface {
	AddGauge(mName string, mValue float64)
	AddCounter(mName string, mValue int64)
	ReadAll() *storage.MemStorage
	ReadGauge(mName string) (float64, bool)
	ReadCounter(mName string) (int64, bool)
}
