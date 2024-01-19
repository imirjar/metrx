package app

import (
	"github.com/imirjar/metrx/internal/app/agent"
	"github.com/imirjar/metrx/internal/app/server/http"
	"github.com/imirjar/metrx/internal/service"
)

func NewServerApp() *http.ServerApp {
	return &http.ServerApp{
		Service: service.NewServerService(),
	}
}

func NewAgentApp() *agent.AgentApp {
	return &agent.AgentApp{}
}
