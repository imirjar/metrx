package agent

import (
	"reflect"
	"runtime"

	"github.com/imirjar/metrx/internal/models"
)

func (a *AgentService) CollectMetrix() {
	var counter int64 = 0
	runtime.ReadMemStats(&a.MemStats)
	var gaugeList = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	for _, ms := range gaugeList {
		var metric models.Metrics
		value := reflect.ValueOf(a.MemStats).FieldByName(ms)
		err := metric.SetMemValue(value, ms, "gauge")
		if err != nil {
			return
		}

		err = a.Storage.AddGauge(metric)
		if err != nil {
			return
		}

		counter++
	}

	var randMetric = models.Metrics{}
	randMetric.SetRandomValue()

	err := a.Storage.AddGauge(randMetric)
	if err != nil {
		return
	}

	counter++

	var cMetric = models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &counter,
	}
	err = a.Storage.AddGauge(cMetric)
	if err != nil {
		return
	}

}
