package server

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/service"
)

type Server struct {
	Router  *mux.Router
	Service *service.Server
	Config  *config
}

func NewServer() *Server {
	server := &Server{
		Config:  newConfig(),
		Router:  newRouter(),
		Service: service.NewServer(),
	}
	server.Router.HandleFunc("/update/gauge/{name}/{value}", server.UpdateGauge).Methods("POST")
	server.Router.HandleFunc("/update/counter/{name}/{value:[0-9]+}", server.UpdateCounter).Methods("POST")
	server.Router.HandleFunc("/value/gauge/{name}", server.ValueGauge).Methods("GET")
	server.Router.HandleFunc("/value/counter/{name}", server.ValueCounter).Methods("GET")
	server.Router.HandleFunc("/", server.MainPage).Methods("GET")
	return server
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}

func newConfig() *config {
	cfg := config{
		url: "localhost:8080",
	}
	cfg.setEnv()
	cfg.setFlags()
	log.Print("start on ", cfg.url)
	return &cfg
}
