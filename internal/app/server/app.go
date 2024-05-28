package server

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/controller/http"
	"github.com/imirjar/metrx/internal/service"
	"github.com/imirjar/metrx/internal/storage"
)

type Gatewayer interface {
	Start(path, conn string) error
}

type ServerApp struct {
	Gateway Gatewayer
}

func Run() {
	cfg := config.NewServerConfig()

	// //STORAGE layer
	// log.Println("###", cfg.DBConn, "###")
	storage := storage.NewStorage(cfg.DBConn, cfg.FilePath, cfg.Interval, cfg.AutoImport)

	// //SERVICE layer
	service := service.NewServerService()
	service.MemStorager = storage

	//GATEWAY layer
	gateway := http.NewGateway(cfg.SECRET)
	gateway.Service = service

	s := ServerApp{
		Gateway: gateway,
	}
	if err := s.Gateway.Start(cfg.URL, cfg.DBConn); err != nil {
		panic(err)
	}
}
