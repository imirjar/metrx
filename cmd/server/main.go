package main

import (
	_ "net/http/pprof"

	"github.com/imirjar/metrx/internal/app/server"
)

func main() {
	server.Run()
}
