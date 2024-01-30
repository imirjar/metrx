package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/app/server/http/middleware"
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
	update.HandleFunc("/", s.UpdateJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	//read metric value
	value := mux.PathPrefix("/value").Subrouter()
	value.HandleFunc("/gauge/{name}", s.ValueGauge).Methods("GET")
	value.HandleFunc("/counter/{name}", s.ValueCounter).Methods("GET")
	value.HandleFunc("/{other}/{name}", s.BadParams).Methods("GET") //status 400
	value.HandleFunc("/", s.ValueJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	//all metric values as a html page
	mux.HandleFunc("/", s.MainPage).Methods("GET")
	mux.Use(middleware.Logger)
	mux.Use(middleware.Encoder)

	return http.ListenAndServe(url, mux)
}
