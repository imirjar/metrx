package agent

import (
	"time"

	"github.com/imirjar/metrx/internal/service"
)

// var memStats runtime.MemStats

func Run(path string, pollInterval, reportInterval time.Duration) error {

	agent := service.NewAgent()

	poll := time.NewTicker(pollInterval)
	report := time.NewTicker(reportInterval)

	for {
		select {
		case <-poll.C:
			agent.CollectMetrix()
		case <-report.C:
			agent.SendMetrix(path)
		}
	}
}
