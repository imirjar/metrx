package agent

import (
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/models"
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
			Path: cfg.URL,
		},
		GaugeList: []string{
			"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
			"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
			"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
		},
	}

	return agent
}

func (a *AgentService) SendMetrix() error {

	var counter int64 = 0
	for _, ms := range a.GaugeList {
		value := a.ReadMemValue(ms)

		metric := models.Metrics{
			ID:    ms,
			MType: "gauge",
			Value: &value,
		}
		err := a.MetricsClient.POSTMetric(metric)
		if err != nil {
			log.Print(err)
			return err
		}
		counter++
	}

	randV := rand.Float64()
	randMetric := models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &randV,
	}
	err := a.MetricsClient.POSTMetric(randMetric)
	if err != nil {
		log.Print(err)
		return err
	}

	counter++

	cMetric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &counter,
	}
	err = a.MetricsClient.POSTMetric(cMetric)
	if err != nil {
		log.Print(err)
		return err
	}
	return err
}

func (a *AgentService) SendBatch() error {
	var metrics []models.Metrics
	var counter int64 = 0

	for _, ms := range a.GaugeList {
		value := a.ReadMemValue(ms)

		metric := models.Metrics{
			ID:    ms,
			MType: "gauge",
			Value: &value,
		}
		metrics = append(metrics, metric)
		counter++
	}

	randV := rand.Float64()
	randMetric := models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &randV,
	}
	metrics = append(metrics, randMetric)
	counter++

	cMetric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &counter,
	}
	metrics = append(metrics, cMetric)

	err := a.MetricsClient.POSTMetrics(metrics)
	return err
}

func (a *AgentService) CollectMetrix() {
	runtime.ReadMemStats(&a.MemStats)
}

func (a *AgentService) ReadMemValue(ms string) float64 {
	switch ms {
	case "Alloc":
		return float64(a.MemStats.Alloc)
	case "BuckHashSys":
		return float64(a.MemStats.BuckHashSys)
	case "Frees":
		return float64(a.MemStats.Frees)
	case "GCCPUFraction":
		return float64(a.MemStats.GCCPUFraction)
	case "GCSys":
		return float64(a.MemStats.GCSys)
	case "HeapAlloc":
		return float64(a.MemStats.HeapAlloc)
	case "HeapIdle":
		return float64(a.MemStats.HeapIdle)
	case "HeapInuse":
		return float64(a.MemStats.HeapInuse)
	case "HeapObjects":
		return float64(a.MemStats.HeapObjects)
	case "HeapReleased":
		return float64(a.MemStats.HeapReleased)
	case "HeapSys":
		return float64(a.MemStats.HeapSys)
	case "LastGC":
		return float64(a.MemStats.LastGC)
	case "Lookups":
		return float64(a.MemStats.Lookups)
	case "MCacheInuse":
		return float64(a.MemStats.MCacheInuse)
	case "MCacheSys":
		return float64(a.MemStats.MCacheSys)
	case "MSpanInuse":
		return float64(a.MemStats.MSpanInuse)
	case "MSpanSys":
		return float64(a.MemStats.MSpanSys)
	case "Mallocs":
		return float64(a.MemStats.Mallocs)
	case "NextGC":
		return float64(a.MemStats.NextGC)
	case "NumForcedGC":
		return float64(a.MemStats.NumForcedGC)
	case "NumGC":
		return float64(a.MemStats.NumGC)
	case "OtherSys":
		return float64(a.MemStats.OtherSys)
	case "PauseTotalNs":
		return float64(a.MemStats.PauseTotalNs)
	case "StackInuse":
		return float64(a.MemStats.StackInuse)
	case "StackSys":
		return float64(a.MemStats.StackSys)
	case "Sys":
		return float64(a.MemStats.Sys)
	case "TotalAlloc":
		return float64(a.MemStats.TotalAlloc)
	default:
		return 0
	}
}
