package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/models"
)

type AgentApp struct {
	cfg config.AgentConfig
	Collector
	Client
}

func NewAgentApp() *AgentApp {
	return &AgentApp{
		cfg: *config.NewAgentConfig(),
		Collector: Collector{
			MemStats: runtime.MemStats{},
		},
		Client: Client{
			Client: http.Client{
				Timeout: 1 * time.Second,
			},
		},
	}
}

func (a *AgentApp) SendMetrics() error {
	var counter int64 = 0
	// var batch models.Batch
	var metrics []models.Metrics
	var gaugeList = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	for _, ms := range gaugeList {
		val := a.Collector.ReadMemStatsValue(ms)
		m := models.Metrics{
			ID:    ms,
			MType: "gauge",
			Value: &val,
		}

		log.Println("agent.go SendMetrics", ms, "-->", val)
		metrics = append(metrics, m)
		counter++
	}

	randV := rand.Float64()
	rm := models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &randV,
	}
	metrics = append(metrics, rm)
	counter++

	// batch.AddCounter("PollCount", counter)
	c := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &counter,
	}
	metrics = append(metrics, c)

	if len(metrics) == 0 {
		log.Print("NO METRICS NO METRICS NO METRICS NO METRICSNO METRICS")
		return errors.New("NO METRICS NO METRICS NO METRICS NO METRICSNO METRICS")
	}
	mm, err := json.Marshal(metrics)
	if err != nil {
		log.Print("agent.go MARSHALL ERR")
		log.Print(err)
		return errors.New("agent.go MARSHALL ERR")
	}

	//compress resp.body
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)
	if _, err := gz.Write(mm); err != nil {
		log.Print("agent.go GZIP ERR")
		log.Print(err)
		return errors.New("agent.go GZIP ERR")
	}
	gz.Close()

	return a.Client.POST(a.cfg.URL+"/updates/", a.cfg.SECRET, mm)
}

func (a *AgentApp) Run() error {

	poll := time.NewTicker(a.cfg.PollInterval)
	report := time.NewTicker(a.cfg.ReportInterval)

	for {
		select {
		case <-poll.C:
			// log.Println("Collect")
			a.Collector.CollectMemStats()
		case <-report.C:
			// log.Println("Send")
			err := a.SendMetrics()
			if err != nil {
				panic(err)
			}
		}
	}
}
