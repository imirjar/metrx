package main

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app"
)

func main() {
	config := config.NewAgentConfig()
	agentApp := app.NewAgentApp()
	if err := agentApp.Run(config.URL, config.PollInterval, config.ReportInterval); err != nil {
		panic(err)
	}
}
