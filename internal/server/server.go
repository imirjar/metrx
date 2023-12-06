package server

import (
	"github.com/imirjar/metrx/internal/config"
	"github.com/imirjar/metrx/internal/handlers"
)

type Server struct {
	Handler *handlers.Handler
	Config  *config.Config
}

func New() *Server {
	return &Server{
		Handler: handlers.New(),
		Config:  config.New(),
	}
}
