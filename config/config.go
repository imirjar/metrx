package config

import (
	"time"
)

func NewServerConfig() *ServerConfig {
	cfg := ServerConfig{
		AppConfig{
			URL: "localhost:8080",
		},
		ServiceConfig{
			Interval: time.Duration(1_000_000_000 * 300), //2s
		},
		StorageConfig{
			FilePath:   "/tmp/metrics-db.json",
			AutoImport: true,
		},
	}
	cfg.setEnv()
	cfg.setFlags()
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
