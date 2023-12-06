package handlers

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/imirjar/metrx/internal/service"
)

type Handler struct {
	Service service.Service
}

func New() *Handler {

	return &Handler{
		// Routes:  defineRoutes(),
		Service: *service.New(),
	}
}

func middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		next.ServeHTTP(w, r)
	})
}
func (h *Handler) GaugeHandle(w http.ResponseWriter, r *http.Request) {

	mamePath, value := path.Split(r.URL.Path)
	name := path.Base(mamePath)

	gauge := h.Service.Gauge(name, value)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gauge)
}

func (h *Handler) CounterHandle(w http.ResponseWriter, r *http.Request) {

	mamePath, value := path.Split(r.URL.Path)
	name := path.Base(mamePath)

	counter := h.Service.Gauge(name, value)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(counter)
}

func (h *Handler) DefineRoutes() *http.ServeMux {
	gauge := http.NewServeMux()
	gauge.Handle("/", middleware(http.HandlerFunc(h.GaugeHandle)))

	counter := http.NewServeMux()
	counter.Handle("/", middleware(http.HandlerFunc(h.CounterHandle)))

	update := http.NewServeMux()
	update.Handle("/gauge/", http.StripPrefix("/gauge", gauge))
	update.Handle("/counter/", http.StripPrefix("/counter", counter))

	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update", update))
	return mux
}
