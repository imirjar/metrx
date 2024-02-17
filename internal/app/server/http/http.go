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
	ByteUpdate(bMetric []byte) ([]byte, error)
	ByteRead(bMetric []byte) ([]byte, error)
	Update(mName, mType, mValue string) error
	View(mName, mType string) (string, error)
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
	update.HandleFunc("/{type}/{name}/{value}", h.Update).Methods("POST")
	update.HandleFunc("/", h.UpdateJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	// read metric value
	value := router.PathPrefix("/value").Subrouter()
	value.HandleFunc("/{type}/{name}", h.View).Methods("GET")
	value.HandleFunc("/", h.ValueJSON).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	// all metric values as a html page
	router.HandleFunc("/ping", h.Ping).Methods("GET")
	router.HandleFunc("/", h.MainPage).Methods("GET")

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
