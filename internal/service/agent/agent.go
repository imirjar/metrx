package agent

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/imirjar/metrx/internal/models"
)

type AgentService struct {
	Metrics  []models.Metrics
	MemStats runtime.MemStats
}

func NewAgentService() *AgentService {

	agent := &AgentService{}
	for _, m := range models.MemStats {
		metric := models.Metrics{
			ID:    m,
			MType: "gauge",
		}
		agent.Metrics = append(agent.Metrics, metric)
	}

	return agent
}

func (a *AgentService) CollectMetrix() {
	runtime.ReadMemStats(&a.MemStats)
}

func (a *AgentService) SendMetrix(URI string) {

	var counter int64 = 0
	path := URI + "/update/"

	//for metric list
	for _, metric := range a.Metrics {
		a.ReadValueFromMemStats(&metric)
		err := metric.SendJSONToPath(path)
		if err != nil {
			fmt.Printf("%s", err)
		}
		counter += 1
	}

	//for random metric
	randomMetric := models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
	}
	randomMetric.SetRandomValue()

	err := randomMetric.SendJSONToPath(path)
	if err != nil {
		fmt.Printf("%s", err)
	}

	counter += 1

	//for metric counter
	counterMetric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &counter,
	}
	counterMetric.SendJSONToPath(path)

}

func (a *AgentService) ReadValueFromMemStats(metric *models.Metrics) {
	value := reflect.ValueOf(a.MemStats).FieldByName(metric.ID)
	// fmt.Println(value)
	if value.CanFloat() {
		v := value.Float()
		metric.Value = &v
	} else {
		v := float64(value.Uint())
		metric.Value = &v

	}
}
