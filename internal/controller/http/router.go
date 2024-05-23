package http

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/controller/http/middleware"
	"github.com/imirjar/metrx/internal/service"
)

func NewGateway(cfg config.ServerConfig) *HTTPGateway {
	service := service.NewServerService(cfg)
	middleware := middleware.New()
	app := HTTPGateway{
		Service:    service,
		Middleware: middleware,
		Secret:     cfg.SECRET,
	}
	return &app
}

type Middleware interface {
	Encrypting(key string) func(next http.Handler) http.Handler
	Logging() func(next http.Handler) http.Handler
	Compressing() func(next http.Handler) http.Handler
	EncWrite(key string) func(next http.Handler) http.Handler
}

type HTTPGateway struct {
	Service    Service
	Middleware Middleware
	Secret     string
}

func (h *HTTPGateway) Start(path, conn string) error {

	router := chi.NewRouter()

	// compression is upper then encrypting its matter!
	router.Use(h.Middleware.Compressing())
	router.Use(h.Middleware.Encrypting(h.Secret))
	router.Use(h.Middleware.EncWrite(h.Secret))
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

	router.Mount("/debug", Profiler())

	s := &http.Server{
		Addr:    path,
		Handler: router,
	}

	return s.ListenAndServe()
}

func Profiler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.Handle("/vars", expvar.Handler())

	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))

	return r
}
