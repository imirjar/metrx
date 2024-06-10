package main

import (
	_ "net/http/pprof"

	app "github.com/imirjar/metrx/internal/app/server"
)

func main() {
	app.Run()
}
