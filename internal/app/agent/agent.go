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
	cfg     config.AgentConfig
}

func NewAgentApp() *AgentApp {
	cfg := config.NewAgentConfig()
	return &AgentApp{
		Service: agent.NewAgentService(*cfg),
		cfg:     *cfg,
	}
}

func (a *AgentApp) Run() error {

	poll := time.NewTicker(a.cfg.PollInterval)
	report := time.NewTicker(a.cfg.ReportInterval)

	for {
		select {
		case <-poll.C:
			// log.Println("Collect")
			a.Service.CollectMetrix()
		case <-report.C:
			// log.Println("Send")
			a.Service.SendBatch()
		}
	}
}
