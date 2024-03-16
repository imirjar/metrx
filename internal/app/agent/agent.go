package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/encrypt"
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

func (a *AgentApp) SendMetrics() {
	var counter int64 = 0
	var batch models.Batch
	var gaugeList = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	for _, ms := range gaugeList {
		value := a.Collector.ReadMemStatsValue(ms)
		// log.Println("#####agent.go MemValue-->", value)
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

	if a.cfg.SECRET != "" {
		hash, err := encrypt.EncryptSHA256(mm, []byte(a.cfg.SECRET))
		if err != nil {
			log.Fatal(err)
		}
		a.Client.POST(a.cfg.URL+"/updates/", mm, hash)
	} else {
		a.Client.POST(a.cfg.URL+"/updates/", mm)
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
			a.SendMetrics()
		}
	}
}
