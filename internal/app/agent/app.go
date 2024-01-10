package agent

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

func Run() error {
	conf := newConfig()
	// fmt.Printf("Client issue on %s", conf.url)

	counter := 0
	go func() {
		for {
			runtime.ReadMemStats(&conf.memStats)
			counter += 1
			time.Sleep(conf.pollInterval)
		}
	}()

	for {
		for _, metric := range conf.gaugeMetrics {
			sendMetric("gauge", metric, conf.url, reflect.ValueOf(conf.memStats).FieldByName(metric))
		}
		sendMetric("gauge", "RandomValue", conf.url, rand.Intn(100))
		sendMetric("counter", "PollCount", conf.url, counter)
		counter = 0
		time.Sleep(conf.reportInterval)
	}

}
