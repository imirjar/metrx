package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
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
	gw := gateway.NewGateway(cfg.Addr, cfg.CryptoKey, cfg.Secret, cfg.DBConn, cfg.TrustedSubnet)
	gw.Service = service

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		<-ctx.Done()
		fmt.Println("Break the loop")
		if err := gw.Server.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}
		return
		// case <-time.After(1 * time.Second):
		// 	fmt.Println("Hello in a loop")

	}()

	// Run the server
	err := gw.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	wg.Wait()
	fmt.Println("Main done")
}
