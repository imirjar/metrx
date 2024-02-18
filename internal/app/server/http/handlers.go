package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/models"
)

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

func (h *HTTPGateway) Update(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var metric = models.Metrics{
		ID:    vars["name"],
		MType: vars["type"],
	}

	switch vars["type"] {
	case "gauge":
		value, err := strconv.ParseFloat(vars["value"], 64)
		if err != nil {
			http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
			return
		}
		metric.Value = &value
	case "counter":
		delta, err := strconv.ParseInt(vars["value"], 10, 64)
		if err != nil {
			http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
			return
		}
		metric.Delta = &delta
	default:
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
		return
	}

	err := h.Service.Update(metric)
	fmt.Println(metric)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}
}

func (h *HTTPGateway) View(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var metric = models.Metrics{
		ID:    vars["name"],
		MType: vars["type"],
	}

	// newMetric, err := h.Service.View(metric) //here
	// if err != nil {
	// 	resp.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	switch metric.MType {
	case "gauge":
		value := fmt.Sprintf("%d", metric.Value)
		resp.Write([]byte(value))
		return
	case "counter":
		delta := fmt.Sprintf("%d", metric.Delta)
		resp.Write([]byte(delta))
		return
	}

}

func (h *HTTPGateway) MainPage(resp http.ResponseWriter, req *http.Request) {
	page := h.Service.MetricPage()
	resp.Header().Set("content-type", "text/html")
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, page)
}

func (h *HTTPGateway) UpdateJSON(resp http.ResponseWriter, req *http.Request) {
	// var metric models.Metrics
	var buf bytes.Buffer //byte buffer

	//request body -> buf
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	bMetric, err := h.Service.ByteUpdate(buf.Bytes())
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(bMetric)
}

func (h *HTTPGateway) ValueJSON(resp http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer //byte buffer

	//request body -> buf
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	bMetric, err := h.Service.ByteRead(buf.Bytes())
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(bMetric)
}
