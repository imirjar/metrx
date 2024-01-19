package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/service/server"
)

type ServerApp struct {
	Service *server.ServerService
}

func (a *ServerApp) Run(url string) error {
	mux := mux.NewRouter()

	//middleware is cheking if the name and value exist
	update := mux.PathPrefix("/update").Subrouter()
	update.HandleFunc("/gauge/{name}/{value:[0-9]+[.]{0,1}[0-9]*}", a.UpdateGauge).Methods("POST") //value:[0-9]+[.]{0,1}[0-9]*
	update.HandleFunc("/counter/{name}/{value:[0-9]+}", a.UpdateCounter).Methods("POST")
	update.HandleFunc("/{other}/{name}/{value}", a.BadParams).Methods("POST")
	// update.Use(a.UpdateMiddleware)

	//middleware is cheking if the name exists
	value := mux.PathPrefix("/value").Subrouter()
	value.HandleFunc("/gauge/{name}", a.ValueGauge).Methods("GET")
	value.HandleFunc("/counter/{name}", a.ValueCounter).Methods("GET")
	value.HandleFunc("/{other}/{name}", a.BadParams).Methods("GET")
	// value.Use(a.ValueMiddleware)

	//without middleware
	mux.HandleFunc("/", a.MainPage).Methods("GET")

	return http.ListenAndServe(url, mux)
}
