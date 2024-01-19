package service

import (
	"net/http"
	"time"

	"github.com/imirjar/metrx/internal/service/agent"
	"github.com/imirjar/metrx/internal/service/server"
	"github.com/imirjar/metrx/internal/storage"
)

// лучше вынести в подпакет
func NewServer() *server.Server {
	return &server.Server{
		Storage: storage.New(),
	}
}

// лучше вынести в подпакет
func NewAgent() *agent.Agent {
	agent := &agent.Agent{
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
