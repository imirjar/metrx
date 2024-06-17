package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/imirjar/metrx/internal/server/controller/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/server/controller/http/middleware/encryptor"
	"github.com/imirjar/metrx/internal/server/controller/http/middleware/logger"
)

func NewGateway(secret string) *HTTPGateway {
	// service := service.NewServerService(cfg)
	gw := HTTPGateway{
		Secret: secret,
	}
	return &gw
}

// Http gateway using secret value for encoding
// it has few endpoints for handlers whitch use Service interface
type HTTPGateway struct {
	Service Service
	Secret  string
}

func (h *HTTPGateway) Start(path, conn string) error {

	router := chi.NewRouter()
	router.Use(middleware.NoCache)

	router.Use(compressor.Compressing())
	router.Use(encryptor.Encrypting(h.Secret))
	router.Use(encryptor.EncWrite(h.Secret))
	router.Use(logger.Logger())

	// Save metric
	router.Route("/update", func(update chi.Router) {
		update.Post("/{type}/{name}/{value}", h.UpdatePathHandler())
		update.Post("/", h.UpdateJSONHandler())
	})

	// Read metric if exists
	router.Route("/value", func(value chi.Router) {
		value.Get("/{type}/{name}", h.ValuePathHandler())
		value.Post("/", h.ValueJSONHandler())
	})

	// Save list passed metrics
	router.Route("/updates", func(batch chi.Router) {
		batch.Post("/", h.BatchHandler())
	})

	// Check db connection
	router.Get("/ping", h.Ping(conn))

	// HTML page witch all of metrics
	router.Get("/", h.MainPage())

	// Pprof package routes
	router.Mount("/debug", middleware.Profiler())

	s := &http.Server{
		Addr:    path,
		Handler: router,
	}

	return s.ListenAndServe()
}
