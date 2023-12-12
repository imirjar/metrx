package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (h *Handler) Handler404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Неверный запрос"))
}

func (h *Handler) Handler400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Неверный запрос"))
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if _, err := strconv.Atoi(vars["value"]); err == nil {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
		}
	})
}

func (h *Handler) GaugeHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if vars["name"] != "" && vars["value"] != "" {
		gauge := h.Service.Gauge(vars["name"], vars["value"])

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		// w.Write(nil)
		json.NewEncoder(w).Encode(gauge)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
	}

}

func (h *Handler) CounterHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if vars["name"] != "" && vars["value"] != "" {
		counter := h.Service.Counter(vars["name"], vars["value"])
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		// w.Write(nil)
		json.NewEncoder(w).Encode(counter)
	} else {
		w.WriteHeader(http.StatusBadGateway)
		w.Write(nil)
	}
}

func (h *Handler) DefineRoutes() *mux.Router {

	mux := mux.NewRouter()
	mux.Use(Middleware)
	mux.HandleFunc("/update/gauge/{name}/{value}", h.GaugeHandle).Methods("POST")
	mux.HandleFunc("/update/counter/{name}/{value}", h.CounterHandle).Methods("POST")
	mux.HandleFunc("/update/{another}/{name}/{value}", h.Handler400).Methods("POST")

	return mux
}
