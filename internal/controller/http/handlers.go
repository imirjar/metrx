package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/ping"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Provide methods of the service layer
type Service interface {
	UpdateMetrics(ctx context.Context, metrics []models.Metrics) error
	UpdateMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error)
	ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error)
	MetricPage(ctx context.Context) (string, error)
}

// Html page consist of the saved metrics
func (h *HTTPGateway) MainPage() http.HandlerFunc {
	log.Println("HANDLER MAIN PAGE")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		page, err := h.Service.MetricPage(ctx)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, page)
	}
}

// Update metric value by passing params in url path
func (h *HTTPGateway) UpdatePathHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HANDLER UpdatePathHandler PAGE")

		var metric models.Metrics
		metric.ID = chi.URLParam(r, "name")
		metric.MType = chi.URLParam(r, "type")
		if err := metric.SetVal(chi.URLParam(r, "value")); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newMetric, err := h.Service.UpdateMetric(r.Context(), metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := newMetric.GetVal()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint(result)))
	}
}

// Read metric value by passing params in url path
func (h *HTTPGateway) ValuePathHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HANDLER ValuePathHandler PAGE")

		var metric models.Metrics
		metric.ID = chi.URLParam(r, "name")
		metric.MType = chi.URLParam(r, "type")

		newMetric, err := h.Service.ViewMetric(r.Context(), metric)
		if err != nil {
			log.Println("ViewMetric ERROR", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		result, err := newMetric.GetVal()
		if err != nil {
			log.Println("GETVAL ERROR", err)
			http.Error(w, errParamsIncorrect.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

// Update metric value with application/json by passing json
func (h *HTTPGateway) UpdateJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("HANDLER UpdateJSONHandler PAGE")
		var metric models.Metrics

		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			log.Println("HANDLER UpdateJSONHandler Decode ERROR", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		newMetric, err := h.Service.UpdateMetric(r.Context(), metric)
		if err != nil {
			// log.Println("Что-то не обновляется", metric.ID, metric.MType, value)
			log.Println("HANDLER UpdateJSONHandler UpdatePath ERROR", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(newMetric); err != nil {
			log.Println("HANDLER UpdateJSONHandler Encode ERROR", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Read metric value with application/json by passing json
func (h *HTTPGateway) ValueJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Println("HANDLER ValueJSONHandler PAGE")

		var metric models.Metrics
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			log.Println("HANDLER ValueJSONHandler Decode ERROR", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		newMetric, err := h.Service.ViewMetric(r.Context(), metric)
		if err != nil {
			log.Println("HANDLER ValueJSONHandler ViewPath ERROR", err)
			http.Error(w, errMetricNameIncorrect.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(newMetric); err != nil {
			log.Println("HANDLER ValueJSONHandler Encode ERROR", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Update metrics with application/json by passing json list of metrics
func (h *HTTPGateway) BatchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HANDLER BatchHandler PAGE")
		var metrics []models.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			log.Println("HANDLER ValueJSONHandler Decode ERROR", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err := h.Service.UpdateMetrics(r.Context(), metrics)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// Check db connection
func (h *HTTPGateway) Ping(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HANDLER Ping PAGE")
		ctx := r.Context()

		db, err := ping.NewDBPool(ctx, path)
		if err != nil {
			log.Println("HANDLER NewDBPool ERROR", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = db.Ping(ctx); err != nil {
			log.Println("HANDLER DB Ping ERROR", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
