package server

import (
	"github.com/imirjar/metrx/internal/app/server/http"
	"github.com/imirjar/metrx/internal/service/server"
)

func NewServerApp() *http.HTTPApp {
	return &http.HTTPApp{
		Service: server.NewServerService(),
	}
}
