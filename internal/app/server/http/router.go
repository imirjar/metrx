package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server/http/middleware"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/internal/service"
)

func NewGateway(cfg config.ServerConfig) *HTTPGateway {
	service := service.NewServerService(cfg)
	middleware := middleware.New()
	app := HTTPGateway{
		Service:    service,
		Middleware: middleware,
	}
	return &app
}

type Service interface {
	BatchUpdate(ctx context.Context, metrics []models.Metrics) error
	UpdatePath(ctx context.Context, name, mType, mValue string) (string, error)
	ViewPath(ctx context.Context, name, mType string) (string, error)
	MetricPage(ctx context.Context) (string, error)
}

type Middleware interface {
	Encrypting(key string) func(next http.Handler) http.Handler
	Logging() func(next http.Handler) http.Handler
	Compressing() func(next http.Handler) http.Handler
}

type HTTPGateway struct {
	Service    Service
	Middleware Middleware
}

func (h *HTTPGateway) Start(path, conn, secret string) error {

	router := chi.NewRouter()

	router.Use(h.Middleware.Encrypting(secret))
	router.Use(h.Middleware.Compressing())
	router.Use(h.Middleware.Logging())

	router.Route("/update", func(update chi.Router) {
		update.Post("/{type}/{name}/{value}", h.UpdatePathHandler())
		update.Post("/", h.UpdateJSONHandler())
	})

	router.Route("/value", func(value chi.Router) {
		value.Get("/{type}/{name}", h.ValuePathHandler())
		value.Post("/", h.ValueJSONHandler())
	})

	router.Route("/updates", func(batch chi.Router) {
		batch.Post("/", h.BatchHandler())
	})

	router.Get("/ping", h.Ping(conn))
	router.Get("/", h.MainPage())

	s := &http.Server{
		Addr:    path,
		Handler: router,
	}

	return s.ListenAndServe()
}
