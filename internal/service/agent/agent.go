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
	ReadOne(metric models.Metrics) (models.Metrics, bool)
	ReadAll(mType string) ([]models.Metrics, error)
	Delete(metric models.Metrics) error
}

type AgentService struct {
	MetricsClient models.MetricsClient
	MemStats      runtime.MemStats
	Storage       Storager
}

func NewAgentService(cfg config.AgentConfig) *AgentService {

	agent := &AgentService{
		MetricsClient: models.MetricsClient{
			Client: http.Client{
				Timeout: 1 * time.Second,
			},
			Path: cfg.URL + "/update/",
		},
		Storage: mock.NewMockStorage(config.Testcfg),
	}

	return agent
}
