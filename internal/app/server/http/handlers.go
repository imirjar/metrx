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
func (h *HTTPGateway) MainPage(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	page, err := h.Service.MetricPage(ctx)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Header().Set("content-type", "text/html")
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, page)
}

// UPDATE ...
func (h *HTTPGateway) UpdatePathHandler(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	mType := chi.URLParam(req, "type")
	mName := chi.URLParam(req, "name")
	mValue := chi.URLParam(req, "value")

	if mType == "" || mName == "" || mValue == "" {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.Service.UpdatePath(ctx, mName, mType, mValue)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprint(result)))
}

func (h *HTTPGateway) UpdateJSONHandler(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var metric models.Metrics

	err := json.NewDecoder(req.Body).Decode(&metric)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.MType == "" || metric.ID == "" {
		http.Error(resp, "MType and ID must not be empty", http.StatusBadRequest)
		return
	}

	if metric.MType == "gauge" && metric.Value == nil {
		http.Error(resp, "Value must not be empty for gauge metric", http.StatusBadRequest)
		return
	} else if metric.MType == "counter" && metric.Delta == nil {
		http.Error(resp, "Delta must not be empty for counter metric", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	newMetric, err := h.Service.Update(ctx, metric)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if newMetric.MType == "gauge" {
		// if newMetric.Value != nil {
		// 	log.Println("###UPDATE mValue--->", *newMetric.Value)
		// }
		jsonResponse, err := json.Marshal(newMetric)
		if err != nil {
			http.Error(resp, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}

		resp.Header().Set("Content-Type", "application/json")
		resp.Write(jsonResponse)
	} else if newMetric.MType == "counter" {
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(resp).Encode(newMetric); err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// resp.Header().Set("content-type", "application/json")
	// resp.WriteHeader(http.StatusOK)

	// if err = json.NewEncoder(resp).Encode(newMetric); err != nil {
	// 	http.Error(resp, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// jsonResponse, err := json.Marshal(newMetric)
	// if err != nil {
	// 	http.Error(resp, "Error marshaling JSON", http.StatusInternalServerError)
	// 	return
	// }

	// resp.Header().Set("Content-Type", "application/json")
	// resp.Write(jsonResponse)
}

// VALUE ...
func (h *HTTPGateway) ValuePathHandler(resp http.ResponseWriter, req *http.Request) {
	log.Print("###ValuePathHandler")
	ctx := req.Context()
	mType := chi.URLParam(req, "type")
	mName := chi.URLParam(req, "name")

	if mType == "" || mName == "" {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.Service.ViewPath(ctx, mName, mType)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusNotFound)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprint(result)))
}

func (h *HTTPGateway) ValueJSONHandler(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var metric models.Metrics

	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil || metric.ID == "" || metric.MType == "" {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	newMetric, err := h.Service.View(ctx, metric)
	if err != nil {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
		return
	}

	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(resp).Encode(newMetric); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Batch ...
func (h *HTTPGateway) BatchHandler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var metrics []models.Metrics

	body, err := io.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err = json.Unmarshal(body, &metrics)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Service.BatchUpdate(ctx, metrics)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok"))
}

// Check ...
func (h *HTTPGateway) Ping(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	db, err := ping.NewDBPool(ctx, h.cfg.DBConn)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = db.Ping(ctx); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok"))
}
