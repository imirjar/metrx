package main

import (
	"net/http"

	"github.com/imirjar/metrx/internal/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	server := server.New()
	return http.ListenAndServe(":8080", server.Handler.DefineRoutes())
}
