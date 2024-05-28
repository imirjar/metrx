package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/imirjar/metrx/internal/controller/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/controller/http/middleware/encryptor"
	"github.com/imirjar/metrx/internal/controller/http/middleware/logger"
)

func NewGateway(secret string) *HTTPGateway {
	// service := service.NewServerService(cfg)
	gw := HTTPGateway{
		Secret: secret,
	}
	return &gw
}

type HTTPGateway struct {
	Service Service
	Secret  string
}

func (h *HTTPGateway) Start(path, conn string) error {

	router := chi.NewRouter()
	router.Use(middleware.NoCache)

	// compression is upper then encrypting its matter!
	router.Use(compressor.Compressing())
	router.Use(encryptor.Encrypting(h.Secret))
	router.Use(encryptor.EncWrite(h.Secret))
	router.Use(logger.Logger())

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

	router.Mount("/debug", middleware.Profiler())

	s := &http.Server{
		Addr:    path,
		Handler: router,
	}

	return s.ListenAndServe()
}
