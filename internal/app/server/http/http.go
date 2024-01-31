package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/logger"
	"github.com/imirjar/metrx/internal/service/server"
)

type HTTPApp struct {
	Service  *server.ServerService
	DumpPath string
}

func (h *HTTPApp) Run(url string, useDump bool, dumpPath string, backupInterval time.Duration) error {
	h.DumpPath = dumpPath

	if useDump {
		h.Service.Storage.Import(dumpPath)
	}

	router := mux.NewRouter()
	// set metric value
	update := router.PathPrefix("/update").Subrouter()
	update.HandleFunc("/gauge/{name}/{value:[0-9]+[.]{0,1}[0-9]*}", h.UpdateGauge).Methods("POST")
	update.HandleFunc("/counter/{name}/{value:[0-9]+}", h.UpdateCounter).Methods("POST")
	update.HandleFunc("/{other}/{name}/{value}", h.BadParams).Methods("POST") //status 400
	update.HandleFunc("/", h.UpdateJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")
	// backup storage dump when metrics have updated if ->
	if backupInterval == 0 {
		update.Use(h.Backuper)
	}

	// read metric value
	value := router.PathPrefix("/value").Subrouter()
	value.HandleFunc("/gauge/{name}", h.ValueGauge).Methods("GET")
	value.HandleFunc("/counter/{name}", h.ValueCounter).Methods("GET")
	value.HandleFunc("/{other}/{name}", h.BadParams).Methods("GET") //status 400
	value.HandleFunc("/", h.ValueJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	// all metric values as a html page
	router.HandleFunc("/", h.MainPage).Methods("GET")
	router.Use(compressor.Compressor)
	router.Use(logger.Logger)

	return http.ListenAndServe(url, router)
}