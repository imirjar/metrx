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

	//set metric value
	update := mux.PathPrefix("/update").Subrouter()
	update.HandleFunc("/gauge/{name}/{value:[0-9]+[.]{0,1}[0-9]*}", s.UpdateGauge).Methods("POST")
	update.HandleFunc("/counter/{name}/{value:[0-9]+}", s.UpdateCounter).Methods("POST")
	update.HandleFunc("/{other}/{name}/{value}", s.BadParams).Methods("POST") //status 400

	//read metric value
	value := mux.PathPrefix("/value").Subrouter()
	value.HandleFunc("/gauge/{name}", s.ValueGauge).Methods("GET")
	value.HandleFunc("/counter/{name}", s.ValueCounter).Methods("GET")
	value.HandleFunc("/{other}/{name}", s.BadParams).Methods("GET") //status 400

	//all metric values as a html page
	mux.HandleFunc("/", s.MainPage).Methods("GET")

	return http.ListenAndServe(url, mux)
}
