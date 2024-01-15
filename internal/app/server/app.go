package server

import (
	"net/http"
)

func Run() error {
	server := newServer()
	server.router.HandleFunc("/update/{mType}/{name}/{value}", server.Update).Methods("POST")
	server.router.HandleFunc("/value/{mType}/{name}", server.View).Methods("GET")
	server.router.HandleFunc("/", server.MainPage).Methods("GET")

	return http.ListenAndServe(server.config.url, server.router)
}
