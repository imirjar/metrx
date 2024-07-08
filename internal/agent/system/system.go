package system

import (
	"context"
	"fmt"
	"runtime"

	"github.com/imirjar/metrx/internal/models"
)

func NewSystem() *Collector {
	return &Collector{}
}

type Collector struct {
	Ms runtime.MemStats
}

func (c *Collector) Collect(ctx context.Context) ([]models.Metrics, error) {
	runtime.ReadMemStats(&c.Ms)
	var counter int64

	var metrics []models.Metrics
	var gaugeList = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	for _, n := range gaugeList {
		v, err := getMemStat(&c.Ms, n)
		if err != nil {
			return metrics, err
		}
		metric := models.Metrics{
			ID:    n,
			MType: "gauge",
		}

		counter += 1

		metric.SetVal(v)
		metrics = append(metrics, metric)
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

func getMemStat(ms *runtime.MemStats, n string) (string, error) {
	switch n {
	case "Alloc":
		return fmt.Sprint(ms.Alloc), nil
	case "BuckHashSys":
		return fmt.Sprint(ms.BuckHashSys), nil
	case "Frees":
		return fmt.Sprint(ms.Frees), nil
	case "GCCPUFraction":
		return fmt.Sprint(ms.GCCPUFraction), nil
	case "GCSys":
		return fmt.Sprint(ms.GCSys), nil
	case "HeapAlloc":
		return fmt.Sprint(ms.HeapAlloc), nil
	case "HeapIdle":
		return fmt.Sprint(ms.HeapIdle), nil
	case "HeapInuse":
		return fmt.Sprint(ms.HeapInuse), nil
	case "HeapObjects":
		return fmt.Sprint(ms.HeapObjects), nil
	case "HeapReleased":
		return fmt.Sprint(ms.HeapReleased), nil
	case "HeapSys":
		return fmt.Sprint(ms.HeapSys), nil
	case "LastGC":
		return fmt.Sprint(ms.LastGC), nil
	case "Lookups":
		return fmt.Sprint(ms.Lookups), nil
	case "MCacheInuse":
		return fmt.Sprint(ms.MCacheInuse), nil
	case "MCacheSys":
		return fmt.Sprint(ms.MCacheSys), nil
	case "MSpanInuse":
		return fmt.Sprint(ms.MSpanInuse), nil
	case "MSpanSys":
		return fmt.Sprint(ms.MSpanSys), nil
	case "Mallocs":
		return fmt.Sprint(ms.Mallocs), nil
	case "NextGC":
		return fmt.Sprint(ms.NextGC), nil
	case "NumForcedGC":
		return fmt.Sprint(ms.NumForcedGC), nil
	case "NumGC":
		return fmt.Sprint(ms.NumGC), nil
	case "OtherSys":
		return fmt.Sprint(ms.OtherSys), nil
	case "PauseTotalNs":
		return fmt.Sprint(ms.PauseTotalNs), nil
	case "StackInuse":
		return fmt.Sprint(ms.StackInuse), nil
	case "StackSys":
		return fmt.Sprint(ms.StackSys), nil
	case "Sys":
		return fmt.Sprint(ms.Sys), nil
	case "TotalAlloc":
		return fmt.Sprint(ms.TotalAlloc), nil
	default:
		return "", fmt.Errorf("there is no metric named %s", n)
	}
}
