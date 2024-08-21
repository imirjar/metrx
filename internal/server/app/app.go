package server

import (
	config "github.com/imirjar/metrx/config/server"
	"github.com/imirjar/metrx/internal/server/controller/grpc"
	HTTPServer "github.com/imirjar/metrx/internal/server/controller/http"
	"github.com/imirjar/metrx/internal/server/service"
	"github.com/imirjar/metrx/internal/server/storage"
)

type Server interface {
	Start(string) error
	Stop() error
}

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
	storage := storage.New(cfg.DBConn, cfg.FilePath, cfg.Interval.Duration, cfg.AutoImport)

	// Service layer
	service := service.New()
	service.MemStorager = storage

	// Controller layer
	HTTP := HTTPServer.New(cfg.CryptoKey, cfg.Secret, cfg.DBConn, cfg.TrustedSubnet)
	HTTP.Service = service

	GRPC := grpc.NewGRPCServer()
	GRPC.Service = service

	// Run the server
	HTTP.Start(cfg.Addr)
	GRPC.Start()
	// if err := server.Start(cfg.Addr); err != nil && err != http.ErrServerClosed {
	// 	log.Fatal(err)
	// }
}
