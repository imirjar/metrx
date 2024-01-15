package agent

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

var memStats runtime.MemStats
var gaugeMetrics []string = []string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
}

func Run() error {
	// var metricStorage []models.Gauge
	conf := newConfig()
	// fmt.Printf("Client issue on %s", conf.url)

	counter := 0
	go func() {
		for {
			runtime.ReadMemStats(&memStats)
			counter += 1
			time.Sleep(conf.pollInterval)
		}
	}()

	for {
		for _, metric := range gaugeMetrics {
			sendMetric("gauge", metric, conf.url, reflect.ValueOf(memStats).FieldByName(metric))
		}
		sendMetric("gauge", "RandomValue", conf.url, rand.Intn(100))
		sendMetric("counter", "PollCount", conf.url, counter)
		counter = 0
		time.Sleep(conf.reportInterval)
	}

}
