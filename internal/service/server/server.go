package server

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/storage"
)

func NewServerService(cfg config.ServerConfig) *ServerService {
	storage := storage.NewStorage(cfg)
	serverService := ServerService{
		Storage: storage,
		Dump:    storage,
		cfg:     &cfg.ServiceConfig,
	}
	// run dump auto-exporter
	if serverService.cfg.Interval > 0 {
		go serverService.PeriodicBackup(serverService.cfg.Interval)
	}
	return &serverService
}

type ServerService struct {
	Storage Storager
	Dump    Backuper
	cfg     *config.ServiceConfig
}
