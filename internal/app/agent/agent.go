package agent

import (
	"time"

	"github.com/imirjar/metrx/internal/service/agent"
)

type AgentApp struct{}

func NewAgentApp() *AgentApp {
	return &AgentApp{}
}

func (a *AgentApp) Run(path string, pollInterval, reportInterval time.Duration) error {

	agent := agent.NewAgentService()

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
