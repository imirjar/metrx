package server

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/storage/mock"
)

func NewServerService(cfg config.ServerConfig) *ServerService {
	store := mock.NewMockStorage(&cfg.StorageConfig)
	serverService := ServerService{
		MemStorager: store,
		Backuper:    store,
		cfg:         &cfg.ServiceConfig,
	}
	// run dump auto-exporter
	if serverService.cfg.Interval > 0 {
		go serverService.PeriodicBackup(serverService.cfg.Interval)
	}
	return &serverService
}

type ServerService struct {
	MemStorager MemStorager
	Backuper    Backuper
	cfg         *config.ServiceConfig
}
