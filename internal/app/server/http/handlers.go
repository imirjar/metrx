package http

import (
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

// MainPage ...
func (h *HTTPGateway) MainPage() http.HandlerFunc {
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

// UPDATE ...
func (h *HTTPGateway) UpdatePathHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		mType := chi.URLParam(r, "type")
		mName := chi.URLParam(r, "name")
		mValue := chi.URLParam(r, "value")

		if mType == "" || mName == "" || mValue == "" {
			http.Error(w, errMetricNameIncorrect.Error(), http.StatusBadRequest)
			return
		}

		result, err := h.Service.UpdatePath(ctx, mName, mType, mValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint(result)))
	}
}

func (h *HTTPGateway) UpdateJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		var metric models.Metrics

		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			log.Print("Что-то с докодером")
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		value, err := metric.GetVal()
		if err != nil {
			log.Print("Что-то не вычленилось значение")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := h.Service.UpdatePath(ctx, metric.ID, metric.MType, value)
		if err != nil {
			log.Print("Что-то не обновляется")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = metric.SetVal(result)
		if err != nil {
			log.Print("Что-то не присваевается новое значение")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// response, err := metric.Marshal()
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(metric); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// w.Write(response)
	}
}

// VALUE ...
func (h *HTTPGateway) ValuePathHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		mType := chi.URLParam(r, "type")
		mName := chi.URLParam(r, "name")

		result, err := h.Service.ViewPath(ctx, mName, mType)
		if err != nil {
			http.Error(w, errParamsIncorrect.Error(), http.StatusNotFound)
			return
		}

		log.Println("RESULT VALUE PATH HANDLER--->", result)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

func (h *HTTPGateway) ValueJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		var metric models.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		val, err := h.Service.ViewPath(ctx, metric.ID, metric.MType)
		metric.SetVal(val)

		if err != nil {
			http.Error(w, errMetricNameIncorrect.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(metric); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Batch ...
func (h *HTTPGateway) BatchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var metrics []models.Metrics

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, m := range metrics {
			if m.MType == "gauge" {
				log.Println("handlers.go batch metric ===>", m.ID, *m.Value)
			} else {
				log.Println("handlers.go batch metric ===>", m.ID, *m.Delta)
			}
		}

		err = h.Service.BatchUpdate(ctx, metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// Check ...
func (h *HTTPGateway) Ping(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		db, err := ping.NewDBPool(ctx, path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
