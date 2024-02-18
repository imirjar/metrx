package agent

import (
	"runtime"

	"github.com/imirjar/metrx/internal/models"
)

func (a *AgentService) SendMetrix(url string) {

	gauges, err := a.Storage.ReadAll("gauge")
	if err != nil {
		panic(err)
	}
	counters, err := a.Storage.ReadAll("counter")
	if err != nil {
		panic(err)
	}

	for _, g := range gauges {
		a.MetricsClient.POSTMetric(&g)
	}

	for _, c := range counters {
		a.MetricsClient.POSTMetric(&c)
	}
}

func (a *AgentService) CollectMetrix() {
	var counter int64 = 0
	runtime.ReadMemStats(&a.MemStats)

	for _, ms := range a.GaugeList {
		var metric = models.Metrics{
			ID:    ms,
			MType: "gauge",
		}
		// value, ok := a.ReadMemValue(ms)
		// if ok {
		// 	metric.Value = &value
		// }
		metric.SetRandomValue()

		err := a.Storage.AddGauge(metric)
		if err != nil {
			return
		}

		counter++
	}

	var randMetric = models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
	}
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
