package config

import (
	"log"
	"time"
)

func NewServerConfig() *ServerConfig {
	cfg := ServerConfig{
		URL: "localhost:8080",
	}
	cfg.setEnv()
	cfg.setFlags()
	log.Print("start on ", cfg.URL)
	return &cfg
}

func NewAgentConfig() *AgentConfig {
	cfg := AgentConfig{
		URL:            "http://localhost:8080",
		PollInterval:   time.Duration(1_000_000_000 * 2),  //2s
		ReportInterval: time.Duration(1_000_000_000 * 10), //10s
	}
	cfg.setEnv()
	cfg.setFlags()

	return &cfg
}
