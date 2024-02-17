package agent

import (
	"math/rand"
	"reflect"
	"runtime"
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
		value := reflect.ValueOf(a.MemStats).FieldByName(ms)

		if value.CanFloat() {
			a.Storage.AddGauge(ms, value.Float())
		} else {
			a.Storage.AddGauge(ms, float64(value.Uint()))
		}
		counter++
	}

	a.Storage.AddGauge("RandomValue", rand.Float64())
	counter++

	a.Storage.AddCounter("RandomValue", counter)
}
