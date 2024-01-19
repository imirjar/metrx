package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
)

type Agent struct {
	Metrics  []string
	Client   *http.Client
	MemStats runtime.MemStats
}

func (a *Agent) CollectMetrix() {
	runtime.ReadMemStats(&a.MemStats)
}

func (a *Agent) SendMetrix(path string) {
	counter := 0
	for _, mName := range a.Metrics {
		value := reflect.ValueOf(a.MemStats).FieldByName(mName)
		fullPath := fmt.Sprintf(path+"/update/%s/%s/%v", "gauge", mName, value)
		_, err := a.Client.Post(fullPath, "text/plain", nil)
		if err != nil {
			fmt.Println(err)
		}
		counter += 1
	}

	fullPath := fmt.Sprintf(path+"/update/%s/%s/%v", "gauge", "RandomValue", rand.Float64())
	_, err := a.Client.Post(fullPath, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	counter += 1

	counterPath := fmt.Sprintf(path+"/update/%s/%s/%v", "counter", "PollCount", counter)
	resp, err := a.Client.Post(counterPath, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()
}
