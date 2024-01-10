package agent

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

type config struct {
	url            string
	pollInterval   time.Duration
	reportInterval time.Duration
	memStats       runtime.MemStats
	gaugeMetrics   []string
}

func newConfig() *config {
	port := flag.String("a", "localhost:8080", "executable port")
	pollInterval := flag.Int("p", 2, "executable port")
	reportInterval := flag.Int("r", 10, "executable port")
	flag.Parse()

	return &config{
		url:            fmt.Sprint("http://", *port),
		pollInterval:   time.Duration(*pollInterval),
		reportInterval: time.Duration(*reportInterval),
		gaugeMetrics: []string{
			"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
			"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
			"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
		},
	}
}
