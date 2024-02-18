package agent

func (a *AgentService) ReadMemValue(ms string) float64 {

	switch ms {
	case "Alloc":
		return float64(a.MemStats.Alloc)
	case "BuckHashSys":
		return float64(a.MemStats.BuckHashSys)
	case "Frees":
		return float64(a.MemStats.Frees)
	case "GCCPUFraction":
		return float64(a.MemStats.GCCPUFraction)
	case "GCSys":
		return float64(a.MemStats.GCSys)
	case "HeapAlloc":
		return float64(a.MemStats.HeapAlloc)
	case "HeapIdle":
		return float64(a.MemStats.HeapIdle)
	case "HeapInuse":
		return float64(a.MemStats.HeapInuse)
	case "HeapObjects":
		return float64(a.MemStats.HeapObjects)
	case "HeapReleased":
		return float64(a.MemStats.HeapReleased)
	case "HeapSys":
		return float64(a.MemStats.HeapSys)
	case "LastGC":
		return float64(a.MemStats.LastGC)
	case "Lookups":
		return float64(a.MemStats.Lookups)
	case "MCacheInuse":
		return float64(a.MemStats.MCacheInuse)
	case "MCacheSys":
		return float64(a.MemStats.MCacheSys)
	case "MSpanInuse":
		return float64(a.MemStats.MSpanInuse)
	case "MSpanSys":
		return float64(a.MemStats.MSpanSys)
	case "Mallocs":
		return float64(a.MemStats.Mallocs)
	case "NextGC":
		return float64(a.MemStats.NextGC)
	case "NumForcedGC":
		return float64(a.MemStats.NumForcedGC)
	case "NumGC":
		return float64(a.MemStats.NumGC)
	case "OtherSys":
		return float64(a.MemStats.OtherSys)
	case "PauseTotalNs":
		return float64(a.MemStats.PauseTotalNs)
	case "StackInuse":
		return float64(a.MemStats.StackInuse)
	case "StackSys":
		return float64(a.MemStats.StackSys)
	case "Sys":
		return float64(a.MemStats.Sys)
	case "TotalAlloc":
		return float64(a.MemStats.TotalAlloc)
	default:
		return 0
	}

}
