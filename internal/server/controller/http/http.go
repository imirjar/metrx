package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/imirjar/metrx/internal/server/controller/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/server/controller/http/middleware/hasher"
	"github.com/imirjar/metrx/internal/server/controller/http/middleware/logger"
)

// Http gateway using secret value for encoding
// it has few endpoints for handlers whitch use Service interface
type HTTPGateway struct {
	Service Service
	Server  *http.Server
}

func NewGateway(path, secret, conn string) *HTTPGateway {

	gtw := HTTPGateway{}

	router := chi.NewRouter()
	router.Use(middleware.NoCache)

	router.Use(compressor.Compressing())
	router.Use(hasher.HashRead(secret))
	router.Use(hasher.HashWrite(secret))
	router.Use(logger.Logger())

	// Save metric
	router.Route("/update", func(update chi.Router) {
		update.Post("/{type}/{name}/{value}", gtw.UpdatePathHandler())
		update.Post("/", gtw.UpdateJSONHandler())
	})

	// Read metric if exists
	router.Route("/value", func(value chi.Router) {
		value.Get("/{type}/{name}", gtw.ValuePathHandler())
		value.Post("/", gtw.ValueJSONHandler())
	})

	// Save list passed metrics
	router.Route("/updates", func(batch chi.Router) {
		batch.Post("/", gtw.BatchHandler())
	})

	// Check db connection
	router.Get("/ping", gtw.Ping(conn))

	// HTML page witch all of metrics
	router.Get("/", gtw.MainPage())

	// Pprof package routes
	router.Mount("/debug", middleware.Profiler())

	gtw.Server = &http.Server{
		Addr:    path,
		Handler: router,
	}

	return &gtw
}
