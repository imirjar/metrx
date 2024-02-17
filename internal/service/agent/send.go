package agent

import (
	"github.com/imirjar/metrx/internal/models"
)

func (a *AgentService) SendMetrix(url string) {
	var counter = 0
	gauges := a.Storage.ReadAllGauge()
	counters := a.Storage.ReadAllCounter()

	for v := range gauges {
		value := gauges[v]
		gm := models.Metrics{
			ID:    v,
			MType: "gauge",
			Value: &value,
		}
		a.MetricsClient.SendJSON(&gm)
		counter++
	}

	for v := range counters {
		delta := counters[v]
		cm := models.Metrics{
			ID:    v,
			MType: "counter",
			Delta: &delta,
		}
		a.MetricsClient.SendJSON(&cm)
	}
}
