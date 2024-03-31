package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
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

func (a *AgentApp) SendMetrics(ctx context.Context) {
	var counter int64 = 0
	var batch models.Batch
	var gaugeList = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	for _, ms := range gaugeList {
		value := a.Collector.ReadMemStatsValue(ms)
		if ms == "Mallocs" {
			log.Println("SendMetrics", ms, "-->", value)
		}
		//
		batch.AddGauge(ms, value)
		counter++
	}

	randV := rand.Float64()
	batch.AddGauge("RandomValue", randV)
	counter++

	batch.AddCounter("PollCount", counter)

	mm, err := json.Marshal(batch.Metrics)
	if err != nil {
		log.Fatal(err)
	}

	//compress resp.body
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)
	gz.Write(mm)
	gz.Close()

	err = a.Client.POST(ctx, a.cfg.URL+"/updates/", a.cfg.SECRET, mm)
	if err != nil {
		log.Println("POST ERROR", err)
		return
	}
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
			a.SendMetrics(context.Background())
		}
	}
}
