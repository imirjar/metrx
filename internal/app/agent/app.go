package agent

import (
	"log"
	"net/http"
	"time"
)

func Run() {
	client := &http.Client{
		Timeout: time.Minute,
	}
	metrics := CollectMetrics()

	for {
		log.Println("Agent is started")
		SendMetrics(client, metrics)
		time.Sleep(10 * time.Second)
	}

}
