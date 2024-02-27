package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/ping"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// MainPage ...
func (h *HTTPGateway) MainPage(resp http.ResponseWriter, req *http.Request) {
	page, err := h.Service.MetricPage()
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
	mType := chi.URLParam(req, "type")
	mName := chi.URLParam(req, "name")
	mValue := chi.URLParam(req, "value")

	if mType == "" || mName == "" || mValue == "" {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
		return
	}

	var metric = models.Metrics{
		ID:    mName,
		MType: mType,
	}

	switch mType {

	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		metric.MType = "gauge"
		metric.Value = &value

	case "counter":
		delta, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		metric.MType = "counter"
		metric.Delta = &delta

	default:
		http.Error(resp, errMetricTypeIncorrect.Error(), http.StatusBadRequest)
		return
	}
	if _, err := h.Service.Update(metric); err != nil {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok"))
}

func (h *HTTPGateway) UpdateJSONHandler(resp http.ResponseWriter, req *http.Request) {
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

	newMetric, err := h.Service.Update(metric)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(resp).Encode(newMetric); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

// VALUE ...
func (h *HTTPGateway) ValuePathHandler(resp http.ResponseWriter, req *http.Request) {

	mType := chi.URLParam(req, "type")
	mName := chi.URLParam(req, "name")

	if mType == "" || mName == "" {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
		return
	}

	var metric = models.Metrics{
		ID:    mName,
		MType: mType,
	}

	switch mType {
	case "gauge":
		metric, err := h.Service.View(metric)
		if err != nil {
			http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
			return
		}
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprint(*metric.Value)))

	case "counter":
		metric, err := h.Service.View(metric)
		if err != nil {
			http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
			return
		}
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
	default:
		http.Error(resp, errMetricTypeIncorrect.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HTTPGateway) ValueJSONHandler(resp http.ResponseWriter, req *http.Request) {
	var metric models.Metrics

	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil || metric.ID == "" || metric.MType == "" {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	newMetric, err := h.Service.View(metric)
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

	err = h.Service.BatchUpdate(metrics)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok"))
}

// Check ...
func (h *HTTPGateway) Ping(resp http.ResponseWriter, req *http.Request) {
	connIsValid := ping.PingPgx(req.Context(), h.cfg.DBConn)

	if !connIsValid {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok"))
}
