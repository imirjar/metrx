package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/logger"
	"github.com/imirjar/metrx/internal/service/server"
)

func NewServerApp(cfg config.ServerConfig) *HTTPApp {
	app := HTTPApp{
		Service: server.NewServerService(cfg),
		cfg:     &cfg.AppConfig,
	}
	return &app
}

type HTTPApp struct {
	Service *server.ServerService
	cfg     *config.AppConfig
}

func (h *HTTPApp) Run() error {

	router := mux.NewRouter()
	// set metric value
	update := router.PathPrefix("/update").Subrouter()
	update.HandleFunc("/gauge/{name}/{value:[0-9]+[.]{0,1}[0-9]*}", h.UpdateGauge).Methods("POST")
	update.HandleFunc("/counter/{name}/{value:[0-9]+}", h.UpdateCounter).Methods("POST")
	update.HandleFunc("/{other}/{name}/{value}", h.BadParams).Methods("POST") //status 400
	update.HandleFunc("/", h.UpdateJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	// read metric value
	value := router.PathPrefix("/value").Subrouter()
	value.HandleFunc("/gauge/{name}", h.ValueGauge).Methods("GET")
	value.HandleFunc("/counter/{name}", h.ValueCounter).Methods("GET")
	value.HandleFunc("/{other}/{name}", h.BadParams).Methods("GET") //status 400
	value.HandleFunc("/", h.ValueJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	router.Use(compressor.Compressor)
	router.Use(logger.Logger)

	// all metric values as a html page
	router.HandleFunc("/", h.MainPage).Methods("GET")

	// DB connection test
	router.HandleFunc("/ping", h.Ping).Methods("GET")

	s := &http.Server{
		Addr:         h.cfg.URL,
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	return s.ListenAndServe()
}
