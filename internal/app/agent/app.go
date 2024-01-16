package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/imirjar/metrx/internal/models"
)

// var memStats runtime.MemStats
var gaugeMetrics []string = []string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
}

var gaugeStore []models.Gauge
var counter int = 0

func Run() error {

	conf := newConfig()
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for {
			collectMetrix(conf, &wg)
		}
	}()

	time.Sleep(conf.pollInterval)

	go func() {
		for {
			sendMetrix(conf, &wg)
		}
	}()

	wg.Wait()
	return nil
}

func collectMetrix(c *config, wg *sync.WaitGroup) {
	wg.Add(1)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for _, name := range gaugeMetrics {
		gauge := models.Gauge{Name: name}
		value := reflect.ValueOf(memStats).FieldByName(name)
		if value.CanFloat() {
			gauge.Value = value.Float()
		} else if value.CanInt() {
			gauge.Value = float64(value.Int())
		} else {
			return
		}
		// fmt.Println(gauge.Name, "#", gauge.Value)
		c.store.AddGauge(gauge)
		counter += 1

	}
	randomGauge := models.Gauge{
		Name:  "RandomValue",
		Value: rand.Float64(),
	}

	c.store.AddGauge(randomGauge)
	counter += 1

	time.Sleep(c.pollInterval)
}

func sendMetrix(c *config, wg *sync.WaitGroup) {
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	var counter int = 0

	for _, metric := range c.store.ReadAllGauge() {
		fmt.Println(metric.Name, metric.Value)
		// sendMetric("gauge", metric.Name, c.url, metric.Value)
		path := fmt.Sprintf(c.url+"/update/%s/%s/%v", "gauge", metric.Name, metric.Value)
		resp, err := client.Post(path, "text/plain", nil)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()
		counter += 1
	}

	// sendMetric("counter", "PollCount", c.url, counter)
	path := fmt.Sprintf(c.url+"/update/%s/%s/%v", "counter", "PollCount", counter)
	resp, err := client.Post(path, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()
	c.store.Drop()

	time.Sleep(c.reportInterval)
	wg.Done()
}
