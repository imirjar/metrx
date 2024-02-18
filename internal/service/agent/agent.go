package agent

import (
	"net/http"
	"runtime"
	"time"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/internal/storage/mock"
)

type Storager interface {
	// AddGauge(mName string, mValue float64) (float64, error)
	// AddCounter(mName string, mValue int64) (int64, error)
	// ReadAllGauge() map[string]float64
	// ReadAllCounter() map[string]int64
	// ReadGauge(mName string) (float64, bool)
	// ReadCounter(mName string) (int64, bool)
	AddGauge(metric models.Metrics) error
	AddCounter(name string, delta int64) error
	ReadOne(metric models.Metrics) (models.Metrics, bool)
	ReadAll(mType string) ([]models.Metrics, error)
	Delete(metric models.Metrics) error
}
type MetricsClient struct {
	Client http.Client
	Path   string
}

type AgentService struct {
	MetricsClient MetricsClient
	MemStats      runtime.MemStats
	Storage       Storager
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
		Storage: mock.NewMockStorage(config.Testcfg),
		GaugeList: []string{
			"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
			"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
			"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
		},
	}

	return agent
}
