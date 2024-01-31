package server

import (
	"github.com/imirjar/metrx/internal/app/server/http"
	"github.com/imirjar/metrx/internal/service/server"
)

func NewServerApp() *http.HttpApp {
	return &http.HttpApp{
		Service: server.NewServerService(),
	}
}
