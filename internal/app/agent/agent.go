package agent

import (
	"time"

	"github.com/imirjar/metrx/internal/service"
)

type AgentApp struct{}

func (a *AgentApp) Run(path string, pollInterval, reportInterval time.Duration) error {

	agent := service.NewAgentService()

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
