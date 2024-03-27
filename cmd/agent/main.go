package main

import (
	"github.com/imirjar/metrx/internal/app/agent"
)

func main() {
	agentApp := agent.NewAgentApp()
	agentApp.Run()
}
