package server

import (
	config "github.com/imirjar/metrx/config/server"
	"github.com/imirjar/metrx/internal/server/controller/http"
	"github.com/imirjar/metrx/internal/server/service"
	"github.com/imirjar/metrx/internal/server/storage"
)

func Run() {
	// Application configuration variables
	cfg := config.NewConfig()

	//Storage layer
	// cfg.DBConn for db connection
	// if database doesn't exist we create mock storage
	// witch can:
	// place dump to cfg.FilePath
	// witch cfg.Interval periodicity
	// and can autorestore if сfg.AutoImport
	storage := storage.NewStorage(cfg.DBConn, cfg.FilePath, cfg.Interval, cfg.AutoImport)

	// Service layer
	service := service.NewServerService()
	service.MemStorager = storage

	//GATEWAY layer
	gateway := http.NewGateway(cfg.Secret)
	gateway.Service = service

	//Run app on cfg.URL, pass dbconn for /ping handler
	if err := gateway.Start(cfg.Addr, cfg.DBConn); err != nil {
		panic(err)
	}
}
