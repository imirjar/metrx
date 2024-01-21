package app

import (
	"github.com/imirjar/metrx/internal/app/agent"
	"github.com/imirjar/metrx/internal/app/server/http"
	"github.com/imirjar/metrx/internal/service/server"
)

func NewServerApp() *http.ServerApp {
	return &http.ServerApp{
		Service: server.NewServerService(),
	}
}

func NewAgentApp() *agent.AgentApp {
	return &agent.AgentApp{}
}
