package agent

import "runtime"

type MetricStore struct {
	MemStats  runtime.MemStats
	GaugeList []string
}

func (a *MetricStore) CollectMetrix() {
	runtime.ReadMemStats(&a.MemStats)
}
