package agent

import (
	"runtime"
)

type Collector struct {
	MemStats runtime.MemStats
}

func (c *Collector) CollectMemStats() {
	runtime.ReadMemStats(&c.MemStats)
}

func (c *Collector) ReadMemStatsValue(ms string) float64 {
	switch ms {
	case "Alloc":
		return float64(c.MemStats.Alloc)
	case "BuckHashSys":
		return float64(c.MemStats.BuckHashSys)
	case "Frees":
		return float64(c.MemStats.Frees)
	case "GCCPUFraction":
		return float64(c.MemStats.GCCPUFraction)
	case "GCSys":
		return float64(c.MemStats.GCSys)
	case "HeapAlloc":
		return float64(c.MemStats.HeapAlloc)
	case "HeapIdle":
		return float64(c.MemStats.HeapIdle)
	case "HeapInuse":
		return float64(c.MemStats.HeapInuse)
	case "HeapObjects":
		return float64(c.MemStats.HeapObjects)
	case "HeapReleased":
		return float64(c.MemStats.HeapReleased)
	case "HeapSys":
		return float64(c.MemStats.HeapSys)
	case "LastGC":
		return float64(c.MemStats.LastGC)
	case "Lookups":
		return float64(c.MemStats.Lookups)
	case "MCacheInuse":
		return float64(c.MemStats.MCacheInuse)
	case "MCacheSys":
		return float64(c.MemStats.MCacheSys)
	case "MSpanInuse":
		return float64(c.MemStats.MSpanInuse)
	case "MSpanSys":
		return float64(c.MemStats.MSpanSys)
	case "Mallocs":
		return float64(c.MemStats.Mallocs)
	case "NextGC":
		return float64(c.MemStats.NextGC)
	case "NumForcedGC":
		return float64(c.MemStats.NumForcedGC)
	case "NumGC":
		return float64(c.MemStats.NumGC)
	case "OtherSys":
		return float64(c.MemStats.OtherSys)
	case "PauseTotalNs":
		return float64(c.MemStats.PauseTotalNs)
	case "StackInuse":
		return float64(c.MemStats.StackInuse)
	case "StackSys":
		return float64(c.MemStats.StackSys)
	case "Sys":
		return float64(c.MemStats.Sys)
	case "TotalAlloc":
		return float64(c.MemStats.TotalAlloc)
	default:
		return 0
	}
}
