package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/models"
)

func (h *HTTPGateway) BadParams(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusBadRequest)
}

func (h *HTTPGateway) Ping(resp http.ResponseWriter, req *http.Request) {
	err := h.Service.PingDB()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
}

func (h *HTTPGateway) UpdateGauge(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	vv, ok := vars["value"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mValue, err := strconv.ParseFloat(vv, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var metric = models.Metrics{
		ID:    vn,
		MType: "gauge",
		Value: &mValue,
	}

	err = h.Service.Update(metric)
	fmt.Println(metric)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}

}

func (h *HTTPGateway) UpdateCounter(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	vv, ok := vars["value"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mDelta, err := strconv.ParseInt(vv, 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var metric = models.Metrics{
		ID:    vn,
		MType: "counter",
		Delta: &mDelta,
	}

	err = h.Service.Update(metric)
	fmt.Println(metric)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}

}

func (h *HTTPGateway) ValueGauge(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var metric = models.Metrics{
		ID:    vars["name"],
		MType: "gauge",
	}

	value := fmt.Sprintf("%d", metric.Value)
	resp.Write([]byte(value))

}
func (h *HTTPGateway) ValueCounter(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var metric = models.Metrics{
		ID:    vars["name"],
		MType: "counter",
	}

	delta := fmt.Sprintf("%d", metric.Delta)
	resp.Write([]byte(delta))

}

func (h *HTTPGateway) MainPage(resp http.ResponseWriter, req *http.Request) {
	page, err := h.Service.MetricPage()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}
	resp.Header().Set("content-type", "text/html")
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, page)
}

func (h *HTTPGateway) UpdateJSON(resp http.ResponseWriter, req *http.Request) {
	var metric models.Metrics

	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.Update(metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := h.Service.View(metric)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(metric)
}

func (h *HTTPGateway) ValueJSON(resp http.ResponseWriter, req *http.Request) {
	var metric models.Metrics

	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := h.Service.View(metric)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(metric)
}
