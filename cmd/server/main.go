package main

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server/http"
)

func main() {
	cfg := config.NewServerConfig()
	app := http.NewGateway(*cfg)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
