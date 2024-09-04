package system

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/imirjar/metrx/internal/models"
)

func Collect(ctx context.Context) ([]models.Metrics, error) {
	var rms runtime.MemStats
	runtime.ReadMemStats(&rms)
	var counter int64

	var metrics []models.Metrics
	var gaugeList = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	for _, n := range gaugeList {

		metric := models.Metrics{
			ID:    n,
			MType: "gauge",
		}

		v := getMemStat(&rms, n)
		if err := metric.SetVal(v); err != nil {
			log.Print(err)
		}

		metrics = append(metrics, metric)
		counter += 1

		// log.Println("###", *metric.Value)
	}

	cMetric := models.Metrics{
		ID:    "PollCounter",
		MType: "counter",
		Delta: &counter,
	}
	metrics = append(metrics, cMetric)

	return metrics, nil
}

func getMemStat(ms *runtime.MemStats, n string) string {
	switch n {
	case "Alloc":
		return fmt.Sprint(ms.Alloc)
	case "BuckHashSys":
		return fmt.Sprint(ms.BuckHashSys)
	case "Frees":
		return fmt.Sprint(ms.Frees)
	case "GCCPUFraction":
		return fmt.Sprint(ms.GCCPUFraction)
	case "GCSys":
		return fmt.Sprint(ms.GCSys)
	case "HeapAlloc":
		return fmt.Sprint(ms.HeapAlloc)
	case "HeapIdle":
		return fmt.Sprint(ms.HeapIdle)
	case "HeapInuse":
		return fmt.Sprint(ms.HeapInuse)
	case "HeapObjects":
		return fmt.Sprint(ms.HeapObjects)
	case "HeapReleased":
		return fmt.Sprint(ms.HeapReleased)
	case "HeapSys":
		return fmt.Sprint(ms.HeapSys)
	case "LastGC":
		return fmt.Sprint(ms.LastGC)
	case "Lookups":
		return fmt.Sprint(ms.Lookups)
	case "MCacheInuse":
		return fmt.Sprint(ms.MCacheInuse)
	case "MCacheSys":
		return fmt.Sprint(ms.MCacheSys)
	case "MSpanInuse":
		return fmt.Sprint(ms.MSpanInuse)
	case "MSpanSys":
		return fmt.Sprint(ms.MSpanSys)
	case "Mallocs":
		return fmt.Sprint(ms.Mallocs)
	case "NextGC":
		return fmt.Sprint(ms.NextGC)
	case "NumForcedGC":
		return fmt.Sprint(ms.NumForcedGC)
	case "NumGC":
		return fmt.Sprint(ms.NumGC)
	case "OtherSys":
		return fmt.Sprint(ms.OtherSys)
	case "PauseTotalNs":
		return fmt.Sprint(ms.PauseTotalNs)
	case "StackInuse":
		return fmt.Sprint(ms.StackInuse)
	case "StackSys":
		return fmt.Sprint(ms.StackSys)
	case "Sys":
		return fmt.Sprint(ms.Sys)
	case "TotalAlloc":
		return fmt.Sprint(ms.TotalAlloc)
	default:
		return "0"
	}
}
