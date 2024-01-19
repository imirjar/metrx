package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/service/server"
)

type ServerApp struct {
	Service *server.ServerService
}

func (s *ServerApp) Run(url string) error {
	mux := mux.NewRouter()

	//middleware is cheking if the name and value exist
	update := mux.PathPrefix("/update").Subrouter()
	update.HandleFunc("/gauge/{name}/{value:[0-9]+[.]{0,1}[0-9]*}", s.UpdateGauge).Methods("POST") //value:[0-9]+[.]{0,1}[0-9]*
	update.HandleFunc("/counter/{name}/{value:[0-9]+}", s.UpdateCounter).Methods("POST")
	update.HandleFunc("/{other}/{name}/{value}", s.UpdateBadParams).Methods("POST")

	//middleware is cheking if the name exists
	value := mux.PathPrefix("/value").Subrouter()
	value.HandleFunc("/gauge/{name}", s.ValueGauge).Methods("GET")
	value.HandleFunc("/counter/{name}", s.ValueCounter).Methods("GET")
	value.HandleFunc("/{other}/{name}", s.ValueBadParams).Methods("GET")

	//without middleware
	mux.HandleFunc("/", s.MainPage).Methods("GET")

	return http.ListenAndServe(url, mux)
}
