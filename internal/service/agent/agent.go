package agent

import (
	"net/http"
	"runtime"
	"time"

	"github.com/imirjar/metrx/config"
)

type MetricsClient struct {
	Client http.Client
	Path   string
}

type AgentService struct {
	MetricsClient MetricsClient
	MemStats      runtime.MemStats
	GaugeList     []string
}

func NewAgentService(cfg config.AgentConfig) *AgentService {

	agent := &AgentService{
		MetricsClient: MetricsClient{
			Client: http.Client{
				Timeout: 1 * time.Second,
			},
			Path: cfg.URL + "/update/",
		},
		GaugeList: []string{
			"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
			"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
			"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
		},
	}

	return agent
}
