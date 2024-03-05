package main

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/agent"
)

func main() {
	config := config.NewAgentConfig()
	agentApp := agent.NewAgentApp(*config)
	if err := agentApp.Run(config.URL, config.PollInterval, config.ReportInterval); err != nil {
		panic(err)
	}
}
