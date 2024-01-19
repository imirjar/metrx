package server

import (
	"github.com/imirjar/metrx/internal/service"
	"github.com/imirjar/metrx/internal/service/server"
)

type Server struct {
	Service *server.Server
}

func NewServer() *Server {
	server := &Server{
		Service: service.NewServer(),
	}
	return server
}
