package server

import (
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server/http"
)

type Gatewayer interface {
	Start(path, conn string) error
}

type ServerApp struct {
	cfg     *config.ServerConfig
	Gateway Gatewayer
}

func (s *ServerApp) Run() {
	if err := s.Gateway.Start(s.cfg.URL, s.cfg.DBConn); err != nil {
		panic(err)
	}
}

func NewServerApp() *ServerApp {
	cfg := config.NewServerConfig()
	return &ServerApp{
		cfg:     cfg,
		Gateway: http.NewGateway(*cfg),
	}
}
