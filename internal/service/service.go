package service

import (
	"net/http"
	"time"

	"github.com/imirjar/metrx/internal/service/agent"
	"github.com/imirjar/metrx/internal/service/server"
	"github.com/imirjar/metrx/internal/storage"
)

// лучше вынести в подпакет
func NewServerService() *server.ServerService {
	return &server.ServerService{
		Storage: storage.New(),
	}
}

// лучше вынести в подпакет
func NewAgentService() *agent.AgentService {
	agent := &agent.AgentService{
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
