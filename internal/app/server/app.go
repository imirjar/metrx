package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Run(url string) error {
	server := NewServer()
	mux := mux.NewRouter()

	//middleware is cheking if the name and value exist
	update := mux.PathPrefix("/update").Subrouter()
	update.HandleFunc("/gauge/{name}/{value}", server.UpdateGauge).Methods("POST") //value:[0-9]+[.]{0,1}[0-9]*
	update.HandleFunc("/counter/{name}/{value}", server.UpdateCounter).Methods("POST")
	update.HandleFunc("/{other}/{name}/{value}", server.BadRequest).Methods("POST")
	update.Use(server.UpdateMiddleware)

	//middleware is cheking if the name exists
	value := mux.PathPrefix("/value").Subrouter()
	value.HandleFunc("/gauge/{name}", server.ValueGauge).Methods("GET")
	value.HandleFunc("/counter/{name}", server.ValueCounter).Methods("GET")
	value.Use(server.ValueMiddleware)

	//without middleware
	mux.HandleFunc("/", server.MainPage).Methods("GET")

	return http.ListenAndServe(url, mux)
}
