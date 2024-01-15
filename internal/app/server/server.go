package server

import "github.com/gorilla/mux"

type server struct {
	router  *mux.Router
	handler *Handler
	config  *config
}

func newServer() *server {
	return &server{
		config:  newConfig(),
		router:  newRouter(),
		handler: newHandler(),
	}
}
