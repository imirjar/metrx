package main

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server"
)

func main() {
	config := config.NewServerConfig()
	serverApp := server.NewServerApp()
	if err := serverApp.Run(config.URL); err != nil {
		panic(err)
	}
}
