package agent

import (
	"math/rand"
	"runtime"

	"github.com/imirjar/metrx/internal/models"
)

func (a *AgentService) SendMetrix(url string) {
	var counter int64 = 0
	for _, ms := range a.GaugeList {
		value := a.ReadMemValue(ms)

		metric := models.Metrics{
			ID:    ms,
			MType: "gauge",
			Value: &value,
		}
		a.MetricsClient.POSTMetric(metric)
		counter++
	}

	randV := rand.Float64()
	randMetric := models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &randV,
	}
	a.MetricsClient.POSTMetric(randMetric)
	counter++

	cMetric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &counter,
	}
	a.MetricsClient.POSTMetric(cMetric)
}

func (a *AgentService) CollectMetrix() {
	runtime.ReadMemStats(&a.MemStats)
}
