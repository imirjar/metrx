package agent

import (
	"time"

	"github.com/imirjar/metrx/internal/service"
)

// var memStats runtime.MemStats

func Run() error {

	conf := newConfig()
	agent := service.NewAgent()

	poll := time.NewTicker(conf.pollInterval)
	report := time.NewTicker(conf.reportInterval)

	for {
		select {
		case <-poll.C:
			agent.CollectMetrix()
		case <-report.C:
			agent.SendMetrix(conf.url)
		}
	}
}
