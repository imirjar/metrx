package main

//go:generate go run main.go

import (
	_ "net/http/pprof"

	app "github.com/imirjar/metrx/internal/server/app"
)

func main() {
	app.Run()
}
