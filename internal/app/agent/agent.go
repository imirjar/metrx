package agent

import (
	"time"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/service/agent"
)

type Servicer interface {
	CollectMetrix()
	SendMetrix() error
	SendBatch() error
}

type AgentApp struct {
	Service Servicer
}

func NewAgentApp(cfg config.AgentConfig) *AgentApp {
	return &AgentApp{
		Service: agent.NewAgentService(cfg),
	}
}

func (a *AgentApp) Run(path string, pollInterval, reportInterval time.Duration) error {

	poll := time.NewTicker(pollInterval)
	report := time.NewTicker(reportInterval)

	for {
		select {
		case <-poll.C:
			a.Service.CollectMetrix()
		case <-report.C:
			a.Service.SendBatch()
		}
	}
}
