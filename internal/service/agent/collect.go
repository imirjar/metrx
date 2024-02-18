package agent

import (
	"math/rand"
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
		value := reflect.ValueOf(a.MemStats).FieldByName(ms)
		var metric = models.Metrics{
			ID:    ms,
			MType: "gauge",
		}

		if value.CanFloat() {
			val := value.Float()
			metric.Value = &val
			// a.Storage.AddGauge(ms, value.Float())
		} else {
			val := float64(value.Uint())
			metric.Value = &val
			// a.Storage.AddGauge(ms, float64(value.Uint()))
		}

		_, exists := a.Storage.ReadOne(metric)
		if exists {
			err := a.Storage.Update(metric)
			if err != nil {
				return
			}
		} else {
			err := a.Storage.Create(metric)
			if err != nil {
				return
			}
		}

		counter++
	}

	randV := rand.Float64()
	var randMetric = models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &randV,
	}
	_, exists := a.Storage.ReadOne(randMetric)
	if exists {
		err := a.Storage.Update(randMetric)
		if err != nil {
			return
		}
	} else {
		err := a.Storage.Create(randMetric)
		if err != nil {
			return
		}
	}
	counter++

	var cMetric = models.Metrics{
		ID:    "RandomValue",
		MType: "counter",
		Delta: &counter,
	}
	_, exists = a.Storage.ReadOne(cMetric)
	if exists {
		err := a.Storage.Update(cMetric)
		if err != nil {
			return
		}
	} else {
		err := a.Storage.Create(cMetric)
		if err != nil {
			return
		}
	}
}
