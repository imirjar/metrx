package main

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app"
)

func main() {
	config := config.NewServerConfig()
	serverApp := app.NewServerApp()
	if err := serverApp.Run(config.URL); err != nil {
		panic(err)
	}
}
