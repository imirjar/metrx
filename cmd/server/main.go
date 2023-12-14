package main

import (
	"github.com/imirjar/metrx/internal/app/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
