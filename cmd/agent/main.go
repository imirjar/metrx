package main

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/agent"
)

func main() {
	config := config.NewAgentConfig()
	if err := agent.Run(config.URL, config.PollInterval, config.ReportInterval); err != nil {
		panic(err)
	}
}
