package http

import (
	"context"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/logger"
	"github.com/imirjar/metrx/internal/service/server"
)

func NewGateway(cfg config.ServerConfig) *HTTPGateway {
	service := server.NewServerService(cfg)
	app := HTTPGateway{
		Service: service,
		cfg:     cfg,
	}
	return &app
}

type Service interface {
	UpdateGauge(mName string, mValue float64) error
	UpdateCounter(mName string, mValue int64) error
	ViewGaugeByName(mName string) (float64, error)
	ViewCounterByName(mName string) (int64, error)
	MetricPage() string
	Backup() error
	Restore() error
	PingDB() error
}

type HTTPGateway struct {
	Service Service
	cfg     config.ServerConfig
}

func (h *HTTPGateway) Run() error {

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

	// updates
	updates := router.PathPrefix("/updates").Subrouter()
	updates.HandleFunc("/", h.UpdatesMetrics).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	// all metric values as a html page
	router.HandleFunc("/ping", h.Ping).Methods("GET")
	router.HandleFunc("/updates", h.MainPage).Methods("GET")

	router.Use(compressor.Compressor)
	router.Use(logger.Logger)

	s := &http.Server{
		Addr:    h.cfg.URL,
		Handler: router,
		// ReadTimeout:  1 * time.Second,
		// WriteTimeout: 1 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				s.Shutdown(ctx)
				h.Service.Backup()
				return
			}
		}
	}()

	s.ListenAndServe()
	wg.Wait()
	return nil
}
