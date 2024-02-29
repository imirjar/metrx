package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/models"
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
	BatchUpdate(metric []models.Metrics) error

	Update(metric models.Metrics) (models.Metrics, error)
	UpdatePath(name, mType, mValue string) (string, error)
	View(metric models.Metrics) (models.Metrics, error)
	ViewPath(name, mType string) (string, error)

	MetricPage() (string, error)
}

type HTTPGateway struct {
	Service Service
	cfg     config.ServerConfig
}

func (h *HTTPGateway) Run() error {

	router := chi.NewRouter()

	router.Use(compressor.Compressor)
	// router.Use(logger.Logger)

	router.Route("/update", func(update chi.Router) {
		update.Post("/{type}/{name}/{value}", h.UpdatePathHandler)
		update.Post("/", h.UpdateJSONHandler)
	})

	router.Route("/value", func(value chi.Router) {
		value.Get("/{type}/{name}", h.ValuePathHandler)
		value.Post("/", h.ValueJSONHandler)
	})

	router.Route("/updates", func(value chi.Router) {
		value.Post("/", h.BatchHandler)
	})

	router.Get("/ping", h.Ping)
	router.Get("/", h.MainPage)

	s := &http.Server{
		Addr:    h.cfg.URL,
		Handler: router,
		// ReadTimeout:  1 * time.Second,
		// WriteTimeout: 1 * time.Second,
	}

	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	// defer stop()

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()

	// 	<-ctx.Done()
	// 	s.Shutdown(ctx)
	// 	h.Service.Backup()
	// }()

	return s.ListenAndServe()
	// wg.Wait()
	// return nil
}
