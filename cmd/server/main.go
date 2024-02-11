package main

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server/http"
)

func main() {
	cfg := config.NewServerConfig()
	app := http.NewServerApp(*cfg)

	// gateway := http.NewServerApp(*cfg)
	// service := metrics.NewMetricsService(*cfg)
	// // gateway.Service = service
	// storage := mock.NewMockStorage(*cfg)
	// gateway.Service =
	// service.Storage = storage

	//export dump when app stoped
	defer app.Service.Backup()

	//run app
	if err := app.Run(); err != nil {
		panic(err)
	}

}
