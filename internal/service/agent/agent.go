package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type AgentService struct {
	Metrics  []string
	Client   *http.Client
	MemStats runtime.MemStats
}

func NewAgentService() *AgentService {
	agent := &AgentService{
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

func (a *AgentService) CollectMetrix() {
	runtime.ReadMemStats(&a.MemStats)
}

func (a *AgentService) SendMetrix(path string) {
	counter := 0
	for _, mName := range a.Metrics {
		value := reflect.ValueOf(a.MemStats).FieldByName(mName)
		fullPath := fmt.Sprintf(path+"/update/%s/%s/%v", "gauge", mName, value)
		resp, err := a.Client.Post(fullPath, "text/plain", nil)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()
		counter += 1
	}

	fullPath := fmt.Sprintf(path+"/update/%s/%s/%v", "gauge", "RandomValue", rand.Float64())
	resp, err := a.Client.Post(fullPath, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()
	counter += 1

	counterPath := fmt.Sprintf(path+"/update/%s/%s/%v", "counter", "PollCount", counter)
	resp, err = a.Client.Post(counterPath, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()
	defer a.Client.CloseIdleConnections()
}
