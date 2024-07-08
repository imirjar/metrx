package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	config "github.com/imirjar/metrx/config/server"
	gateway "github.com/imirjar/metrx/internal/server/controller/http"
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
	storage := storage.NewStorage(cfg.DBConn, cfg.FilePath, cfg.Interval.Duration, cfg.AutoImport)

	// Service layer
	service := service.NewServerService()
	service.MemStorager = storage

	//GATEWAY layer
	gw := gateway.NewGateway(cfg.Addr, cfg.CryptoKey, cfg.Secret, cfg.DBConn)
	gw.Service = service

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		go func() {
			<-serverCtx.Done()

			if serverCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := gw.Server.Shutdown(serverCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := gw.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
