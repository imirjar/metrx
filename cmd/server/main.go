package main

import (
	"github.com/imirjar/metrx/internal/app/server"
)

func main() {
	app := server.NewServerApp()
	app.Run()
}
