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

	var metric = models.Metrics{
		ID: mName,
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
		resp.WriteHeader(http.StatusBadRequest)
	}

	if _, err := h.Service.Update(metric); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (h *HTTPGateway) UpdateJSONHandler(resp http.ResponseWriter, req *http.Request) {
	var metric models.Metrics

	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	newMetric, err := h.Service.Update(metric)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(newMetric)
}

// VALUE ...
func (h *HTTPGateway) ValuePathHandler(resp http.ResponseWriter, req *http.Request) {

	mType := chi.URLParam(req, "type")
	mName := chi.URLParam(req, "name")

	var metric = models.Metrics{
		ID:    mName,
		MType: mType,
	}

	metric, err := h.Service.View(metric)
	if err != nil {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
		return
	}

	switch mType {
	case "gauge":
		// r := strconv.FormatFloat(float64(*metric.Value), 'f', -1, 64)
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprint(*metric.Value)))
	case "counter":
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
	default:
		http.Error(resp, errMetricTypeUnexpected.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HTTPGateway) ValueJSONHandler(resp http.ResponseWriter, req *http.Request) {
	var metric models.Metrics

	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := h.Service.View(metric)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusNotFound)
		return
	}

	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(metric)
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
