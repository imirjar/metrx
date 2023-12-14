package agent

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sync"

	"github.com/imirjar/metrx/internal/models"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func SendMetrics(client *http.Client, metrics []models.Metrix) error {

	wg := sync.WaitGroup{}
	wg.Add(len(metrics))

	for _, m := range metrics {

		log.Printf("Sending Metrix %s -> %v", m.Name, m.Value)

		go func() {
			defer wg.Done()

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", m.Type, m.Name, m.Value), nil)
			if err != nil {
				log.Println(err)
				return
			}
			req.Header.Add("Content-Type", "text/plain")
			res, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			res.Body.Close()
		}()
	}

	wg.Wait()

	return nil
}

func CollectMetrics() []models.Metrix {

	var GaugeMetrics = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	var metrics []models.Metrix
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	msvalue := reflect.ValueOf(memStats)
	mstype := msvalue.Type()

	for _, metric := range GaugeMetrics {
		field, ok := mstype.FieldByName(metric)

		if !ok {
			fmt.Println(field)
			continue
		}
		value := msvalue.FieldByName(metric)
		metrics = append(metrics, models.Metrix{Name: field.Name, Value: value})

	}

	return metrics
}
