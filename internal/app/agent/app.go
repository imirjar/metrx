package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

var GaugeMetrics = []string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
}

func sendMetric(metricType string, metric string, value any) {
	_, err := http.Post(fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metricType, metric, value), "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}

}

func Run() {
	var memStats runtime.MemStats
	counter := 0
	go func() {
		for {
			// fmt.Println("Updating memstats")
			runtime.ReadMemStats(&memStats)
			// fmt.Println(memStats.Alloc)
			counter += 1
			time.Sleep(pollInterval)
		}
	}()

	for {
		for _, metric := range GaugeMetrics {
			sendMetric("gauge", metric, reflect.ValueOf(memStats).FieldByName(metric))
		}
		sendMetric("gauge", "RandomValue", rand.Intn(100))
		sendMetric("counter", "PollCount", counter)
		// fmt.Println(counter)
		counter = 0
		time.Sleep(reportInterval)
	}
}
